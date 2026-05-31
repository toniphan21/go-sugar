package generatecmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/gen-lib/cli/color"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/concurrent"
	"nhatp.com/go/sugar/watch"
)

func runWatch(stdin io.Reader, stdout, stderr io.Writer, args Arguments, log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	mod, err := sugar.NewModule(wd, sugar.DefaultConfig(), sugar.WithBinary(sugar.BinaryFullName), sugar.WithVersion(sugar.BinaryVersion))
	if err != nil {
		return err
	}

	targets, err := mod.Resolve(args.inputs()...)
	if err != nil {
		log.Error(cli.ColorRed(err.Error()), slog.Any("error", err))
		return err
	}

	watchers := concurrent.Init(ctx, targets, func(target sugar.Target) <-chan watch.Event {
		return watch.New(ctx, target, watch.WithDefaultLogger(log))
	})

	event, trace := concurrent.Tee(ctx, watchers)
	go func() {
		for evt := range trace {
			log.Info("\thandling  " + color.Source(evt.FilePath.DisplayPath) + "...")
		}
	}()

	generateCode := func() <-chan generateResult {
		return generate(ctx, mod, event)
	}
	for result := range concurrent.Parallelize(ctx, runtime.NumCPU(), generateCode) {
		handleGenerateResult(stdout, args, log, result)
	}
	return nil
}

type generateResult struct {
	Target           sugar.Target
	FilePath         sugar.FilePath
	Err              error
	Content          []byte
	GeneratedRelPath string
	Generated        []byte
}

func (r *generateResult) toOutputPair() OutputPair {
	return OutputPair{}
}

func handleGenerateResult(stdout io.Writer, args Arguments, log *slog.Logger, result generateResult) {
	switch {
	case result.Err != nil:
		log.Error(cli.ColorRed(result.Err.Error()))

	case args.JSON:
		output := Output{
			Argument:   result.Target.Input,
			WorkingDir: result.Target.WorkingDir,
			ModuleRoot: result.Target.Root,
		}

		source := OutputFile{
			DisplayPath: result.FilePath.DisplayPath,
			RelPath:     result.FilePath.RelPath,
			Content:     string(result.Content),
		}

		displayPath, err := filepath.Rel(result.Target.WorkingDir, result.GeneratedRelPath)
		if err != nil {
			displayPath = result.GeneratedRelPath
		}
		generated := OutputFile{
			DisplayPath: displayPath,
			RelPath:     result.GeneratedRelPath,
			Content:     string(result.Generated),
		}
		output.Files = append(output.Files, OutputPair{Source: source, Generated: generated})

		out, err := json.Marshal(output)
		if err != nil {
			log.Error(cli.ColorRed(err.Error()))
		}
		_, _ = fmt.Fprintln(stdout, string(out))

	case args.DryRun:
		generatedDisplayPath, err := filepath.Rel(result.Target.WorkingDir, result.GeneratedRelPath)
		if err != nil {
			generatedDisplayPath = result.GeneratedRelPath
		}

		_, _ = fmt.Fprintf(stdout, "// === go-sugar: %v    ===\n", result.FilePath.DisplayPath)
		_, _ = fmt.Fprintf(stdout, "// ---  desugar: %v ---\n", generatedDisplayPath)
		_, _ = fmt.Fprint(stdout, string(result.Generated))

	default:
		absPath := filepath.Join(result.Target.Root, result.GeneratedRelPath)
		displayPath, err := filepath.Rel(result.Target.WorkingDir, absPath)
		if err != nil {
			displayPath = result.GeneratedRelPath
		}
		log.Info("\tgenerated " + color.Generated(displayPath))
		if err := os.WriteFile(absPath, result.Generated, 0644); err != nil {
			log.Error(cli.ColorRed(err.Error()))
		}
	}
}

func generate(ctx context.Context, mod *sugar.Module, event <-chan watch.Event) <-chan generateResult {
	stream := make(chan generateResult)
	go func() {
		defer close(stream)

		select {
		case <-ctx.Done():
			return
		case e, ok := <-event:
			if !ok {
				return
			}

			content, err := os.ReadFile(e.FilePath.AbsPath)
			if err != nil {
				stream <- generateResult{Err: err}
			}
			generatedRelPath, generated, err := mod.GenerateOnDemand(e.FilePath.RelPath, content)
			if err != nil {
				stream <- generateResult{Err: err}
			}

			stream <- generateResult{
				Target:           e.Target,
				FilePath:         e.FilePath,
				Content:          content,
				GeneratedRelPath: generatedRelPath,
				Generated:        generated,
			}
		}
	}()
	return stream
}
