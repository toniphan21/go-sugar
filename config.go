package sugar

const DefaultGoModFileName = "go.mod"
const DefaultGoFileExt = ".go"
const DefaultFileExt = ".gos"

type Env struct {
	GoModFileName string
	GoFileExt     string
	SugarFileExt  string
}

func (e *Env) IsSugarFile(ext string) bool {
	return ext == e.SugarFileExt
}

func (e *Env) GoFilePath(path string) string {
	return path + e.GoFileExt
}

type Config struct {
	Env Env
}

func (c *Config) env() Env {
	e := c.Env
	if e.GoModFileName == "" {
		e.GoModFileName = DefaultGoModFileName
	}

	if e.GoFileExt == "" {
		e.GoFileExt = DefaultGoFileExt
	}

	if e.SugarFileExt == "" {
		e.SugarFileExt = DefaultFileExt
	}
	return e
}

func DefaultConfig() Config {
	return Config{
		Env: Env{
			GoModFileName: DefaultGoModFileName,
			GoFileExt:     DefaultGoFileExt,
			SugarFileExt:  DefaultFileExt,
		},
	}
}
