package versioncmd

import (
	"encoding/json"
	"fmt"
	"os"

	"nhatp.com/go/sugar"
)

type Format int

const (
	FormatNormal Format = iota
	FormatJSON
	FormatSemver
)

type Arguments struct {
	Sugar  sugar.Sugar
	Format Format
}

type Version struct {
	ID      string `json:"id"`
	Binary  string `json:"binary"`
	Version string `json:"version"`
}

func Run(stdin, stdout, stderr *os.File, args Arguments) error {
	binary := args.Sugar.Binary()

	switch args.Format {
	case FormatJSON:
		result := Version{
			ID:      args.Sugar.ID(),
			Binary:  binary.Name,
			Version: binary.Version,
		}
		out, err := json.Marshal(result)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(stdout, string(out))

	case FormatSemver:
		_, _ = fmt.Fprintln(stdout, binary.Version)

	default:
		v := fmt.Sprintf("%s - v%s", binary.Name, binary.Version)
		_, _ = fmt.Fprintln(stdout, v)
	}
	return nil
}
