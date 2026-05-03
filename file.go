package sugar

type File struct {
}

func (f *File) Update(content []byte) *File {
	return f
}

func (f *File) SugarToGo(line, column int) (int, int, error) {
	return line, column, nil
}

func (f *File) GoToSugar(line, column int) (int, int, error) {
	return line, column, nil
}

var _ API = (*File)(nil)

func NewFile(content []byte) *File {
	return &File{}
}
