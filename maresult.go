package macomp

import mecab "github.com/shogo82148/go-mecab"

//MaResult is a result store
type MaResult struct {
	Name     string
	Surfaces []string
	Features []string
	Error    error
}

//NewMaResult returns MaResult
func NewMaResult(name string, node *mecab.Node, err error) *MaResult {
	mr := new(MaResult)
	mr.Name = name

	mr.Surfaces = []string{}
	mr.Features = []string{}

	if err != nil {
		mr.Error = err
		return mr
	}

	for n := node.Next(); n.Stat() != mecab.EOSNode; n = n.Next() {
		mr.Surfaces = append(mr.Surfaces, n.Surface())
		mr.Features = append(mr.Features, n.Feature())
	}
	return mr
}
