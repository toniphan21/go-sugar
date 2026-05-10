package sugar

type Position struct {
	Offset int // byte offset from start of file, starting at 0
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (byte count)
}

type Region struct {
	Pos Position // inclusive
	End Position // exclusive
}

type Kind int

const (
	// KindExpand sugar grows into larger generated block
	KindExpand Kind = iota

	//Shrink    // sugar shrinks (future)
	//Split     // one sugar → multiple generated regions (future)
	//Phantom   // generated region with no sugar counterpart (future)
)

type Entry struct {
	Sugar Region
	Go    Region
	Kind  Kind
}

type Header struct {
	SugarFile string
	GoFile    string
	SugarHash [32]byte // SHA256 of sugar at generation time
	GoHash    [32]byte // SHA256 of go at generation time
}

type SourceMap struct {
	Header  Header
	Entries []Entry // insertion order, unsorted
	ByGo    []int   // indices into Entries, sorted by Go.Pos.Offset
	BySugar []int   // indices into Entries, sorted by Sugar.Pos.Offset
}
