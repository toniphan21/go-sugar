package sugar

type API interface {
	SugarToGo(line, column int) (int, int, error)

	GoToSugar(line, column int) (int, int, error)
}

type Snapshot struct {
	content []byte
	lex     []Lexeme
}

func (s *Snapshot) SugarToGo(line, column int) (int, int, error) {
	return line, column, nil
}

func (s *Snapshot) GoToSugar(line, column int) (int, int, error) {
	return line, column, nil
}

var _ API = (*Snapshot)(nil)
