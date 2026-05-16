package versioncmd

import (
	"encoding/json"
	"fmt"
	"io"

	"nhatp.com/go/gen-lib/cli/color"
	"nhatp.com/go/sugar"
)

type Format int

const (
	FormatNormal Format = iota
	FormatJSON
	FormatSemver
)

type Arguments struct {
	Format Format
}

type Binary struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Version struct {
	Module  string `json:"module"`
	Binary  Binary `json:"binary"`
	Version string `json:"version"`
}

func Run(stdin io.Reader, stdout, stderr io.Writer, arg Arguments) error {
	switch arg.Format {
	case FormatJSON:
		result := Version{
			Module: sugar.BinaryPackagePath,
			Binary: Binary{
				Name: sugar.BinaryName,
				Path: sugar.BinaryPath,
			},
			Version: sugar.Version,
		}
		out, err := json.Marshal(result)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(stdout, string(out))

	case FormatSemver:
		_, _ = fmt.Fprintln(stdout, sugar.Version)
	default:
		v := fmt.Sprintf("%s%s - %s", sugar.BinaryPath, color.Binary(sugar.BinaryName), color.Version(sugar.BinaryVersion))
		_, _ = fmt.Fprintln(stdout, v)
	}
	return nil
}
