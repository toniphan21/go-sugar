package sugar

type API interface {
	SugarToGo(line, column int) (int, int, error)

	GoToSugar(line, column int) (int, int, error)
}

type Snapshot struct {
	content []byte
	lex     []Lexeme
}

func (s *Snapshot) Scan() []Lexeme {
	if s.lex == nil {
		s.lex = Lex(s.content)
	}
	return s.lex
}

func (s *Snapshot) SugarToGo(line, column int) (int, int, error) {
	return line, column, nil
}

func (s *Snapshot) GoToSugar(line, column int) (int, int, error) {
	return line, column, nil
}

var _ API = (*Snapshot)(nil)
