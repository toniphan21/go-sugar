package watch

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"nhatp.com/go/sugar"
)

var errNonSugarFile = errors.New("target is non sugar file")

type Event struct {
	Target   sugar.Target
	FilePath sugar.FilePath
}

func New(ctx context.Context, target sugar.Target, options ...Option) <-chan Event {
	o := &opts{
		log: &defaultLogger{log: slog.New(slog.NewTextHandler(io.Discard, nil))},
	}
	for _, f := range options {
		f.apply(o)
	}
	return makeStream(ctx, o.log, target)
}

type Option interface {
	apply(*opts)
}

type optionFunc func(*opts)

func (f optionFunc) apply(opts *opts) { f(opts) }

type opts struct {
	log Logger
}

func WithDefaultLogger(v *slog.Logger) Option {
	return optionFunc(func(m *opts) {
		m.log = &defaultLogger{log: v}
	})
}

func makeStream(ctx context.Context, log Logger, target sugar.Target) <-chan Event {
	stream := make(chan Event)
	go func() {
		defer close(stream)

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error("failed to create watcher", target, err)
			return
		}
		defer func() {
			if err := watcher.Close(); err != nil {
				log.Error("failed to close watcher", target, err)
			}
		}()

		// setup watcher
		handleEvent, fps, err := setupWatcher(watcher, target)
		if errors.Is(err, errNonSugarFile) {
			log.Skip(target)
			return
		}
		if err != nil {
			log.Error("failed to setup watcher", target, err)
			return
		}

		// stream initial FilePaths in the setup phase
		for _, fp := range fps {
			stream <- Event{Target: target, FilePath: fp}
		}

		// start watcher
		log.Watching(target)
		for {
			select {
			case <-ctx.Done():
				log.Stop(target)
				return

			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if fp, ok := handleEvent(event, watcher, target); ok {
					stream <- Event{Target: target, FilePath: *fp}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					log.Error("watcher error", target, err)
					return
				}
			}
		}
	}()
	return stream
}

func setupWatcher(watcher *fsnotify.Watcher, target sugar.Target) (watcherEventHandler, []sugar.FilePath, error) {
	switch {
	case target.IsDir && target.Recursive:
		dirs, err := target.Dirs()
		if err != nil {
			return nil, nil, err
		}

		for _, dir := range dirs {
			err := watcher.Add(dir)
			if err != nil {
				return nil, nil, err
			}
		}
		return recursiveDirTargetEventHandler, nil, nil

	case target.IsDir:
		if err := watcher.Add(target.Path); err != nil {
			return nil, nil, err
		}
		fps, err := target.Resolve()
		if err != nil {
			return nil, nil, err
		}
		return dirTargetEventHandler, fps, nil

	default:
		fp, ok := target.IsSugarFilePath(target.Path)
		if !ok {
			return nil, nil, errNonSugarFile
		}

		dir := filepath.Dir(target.Path)
		if err := watcher.Add(dir); err != nil {
			return nil, nil, err
		}
		return fileTargetEventHandler, []sugar.FilePath{fp}, nil
	}
}

type watcherEventHandler func(event fsnotify.Event, watcher *fsnotify.Watcher, target sugar.Target) (*sugar.FilePath, bool)

func fileTargetEventHandler(event fsnotify.Event, _ *fsnotify.Watcher, target sugar.Target) (*sugar.FilePath, bool) {
	if !event.Has(fsnotify.Create) && !event.Has(fsnotify.Write) {
		return nil, false
	}

	if target.Path != event.Name {
		return nil, false
	}

	if fp, ok := target.IsSugarFilePath(target.Path); ok {
		return &fp, true
	}
	return nil, false
}

func dirTargetEventHandler(event fsnotify.Event, watcher *fsnotify.Watcher, target sugar.Target) (*sugar.FilePath, bool) {
	if !event.Has(fsnotify.Create) && !event.Has(fsnotify.Write) {
		return nil, false
	}

	p := event.Name
	stat, err := os.Stat(p)
	if err != nil || stat.IsDir() {
		return nil, false
	}

	if fp, ok := target.IsSugarFilePath(p); ok {
		return &fp, true
	}
	return nil, false
}

func recursiveDirTargetEventHandler(event fsnotify.Event, watcher *fsnotify.Watcher, target sugar.Target) (*sugar.FilePath, bool) {
	if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
		_ = watcher.Remove(event.Name)
		return nil, false
	}

	if !event.Has(fsnotify.Create) && !event.Has(fsnotify.Write) && !event.Has(fsnotify.Remove) {
		return nil, false
	}

	p := event.Name
	stat, err := os.Stat(p)
	if err != nil {
		return nil, false
	}

	if stat.IsDir() && event.Has(fsnotify.Create) {
		_ = watcher.Add(p)
		return nil, false
	}

	if fp, ok := target.IsSugarFilePath(p); ok {
		return &fp, true
	}
	return nil, false
}
