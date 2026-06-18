package sugar

import (
	"cmp"
	"slices"
	"sort"
)

type Position struct {
	Offset int `json:"offset"` // byte offset from start of file, starting at 0
	Line   int `json:"line"`   // line number, starting at 1
	Column int `json:"column"` // column number, starting at 1 (byte count)
}

type Region struct {
	Pos Position `json:"pos"` // inclusive
	End Position `json:"end"` // exclusive
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
	Sugar Region `json:"sugar"`
	Go    Region `json:"go"`
	Kind  Kind   `json:"kind"`
}

type Header struct {
	SugarFile string   `json:"sugarFile"`
	GoFile    string   `json:"goFile"`
	SugarHash [32]byte `json:"sugarHash"` // SHA256 of sugar at generation time
	GoHash    [32]byte `json:"goHash"`    // SHA256 of go at generation time
}

type SourceMap struct {
	Header  Header  `json:"header"`
	Entries []Entry `json:"entries"` // insertion order, unsorted
	ByGo    []int   `json:"byGo"`    // indices into Entries, sorted by Go.Pos.Offset
	BySugar []int   `json:"bySugar"` // indices into Entries, sorted by Sugar.Pos.Offset
}

func (sm *SourceMap) buildBySugar() {
	sm.BySugar = make([]int, len(sm.Entries))
	for i := range sm.BySugar {
		sm.BySugar[i] = i
	}
	slices.SortFunc(sm.BySugar, func(a, b int) int {
		return cmp.Compare(sm.Entries[a].Sugar.Pos.Offset, sm.Entries[b].Sugar.Pos.Offset)
	})
}

func (sm *SourceMap) buildByGo() {
	sm.ByGo = make([]int, len(sm.Entries))
	for i := range sm.ByGo {
		sm.ByGo[i] = i
	}
	slices.SortFunc(sm.ByGo, func(a, b int) int {
		return cmp.Compare(sm.Entries[a].Go.Pos.Offset, sm.Entries[b].Go.Pos.Offset)
	})
}

func (sm *SourceMap) GoToSugarByOffset(offset int) (int, bool) {
	idx, found := sort.Find(len(sm.ByGo), func(i int) int {
		entry := sm.Entries[sm.ByGo[i]]
		if offset < entry.Go.Pos.Offset {
			return -1
		}
		if offset >= entry.Go.End.Offset {
			return 1
		}
		return 0
	})
	if !found {
		return -1, false
	}
	return sm.Entries[sm.ByGo[idx]].Sugar.Pos.Offset, true
}
