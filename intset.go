package macomp

//IntSet is a set of int
type IntSet map[int]struct{}

//Add addes a value
func (is IntSet) Add(v int) {
	is[v] = struct{}{}
}
