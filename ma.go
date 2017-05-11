package macomp

import (
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	jk "github.com/shirayu/go-jk"
	mecab "github.com/shogo82148/go-mecab"
)

//MaResource contains resources related the server
type MaResource struct {
	Mecabs map[string]*mecab.MeCab
	Jumans map[string]*jk.Juman
}

//NewMaResource returns MaResource
func NewMaResource(settings map[string]MaSetting) (*MaResource, error) {
	res := new(MaResource)
	res.Mecabs = map[string]*mecab.MeCab{}
	res.Jumans = map[string]*jk.Juman{}

	for name, dicinfo := range settings {
		if dicinfo.Disable {
			continue
		}
		matype := strings.ToLower(dicinfo.MaType)
		switch matype {
		case "mecab":
			if len(dicinfo.Path) != 0 {
				return res, errors.New("Path designation is not supported for MeCab")
			}
			tagger, err := mecab.New(dicinfo.Options)
			if err != nil {
				return res, err
			}
			// avoid GC problem with MeCab 0.996 (see https://github.com/taku910/mecab/pull/24)
			tagger.Parse("")

			res.Mecabs[name] = &tagger
		case "juman", "jumanpp":
			if len(dicinfo.Path) == 0 {
				var err error
				dicinfo.Path, err = exec.LookPath(matype)
				if err != nil {
					return res, err
				}
			}
			opts := make([]string, 0, len(dicinfo.Options)*2)
			for k, v := range dicinfo.Options {
				if len(k) == 0 {
					continue
				}
				opts = append(opts, k)
				if len(v) != 0 {
					opts = append(opts, v)
				}
			}
			jm, err := jk.NewJuman(dicinfo.Path, opts...)
			if err != nil {
				return res, err
			}
			res.Jumans[name] = jm

		default:
			return res, fmt.Errorf("Unknown type %s", dicinfo.MaType)
		}
	}
	return res, nil
}

//Destroy all objects
func (mr *MaResource) Destroy() {
	for _, tagger := range mr.Mecabs {
		tagger.Destroy()
	}
}

//Parse returns results
func (mr *MaResource) Parse(text string) []MaResult {
	results := []MaResult{}
	c := make(chan *MaResult)

	waits := len(mr.Mecabs) + len(mr.Jumans)
	for name, tagger := range mr.Mecabs {
		go func(name string, tagger *mecab.MeCab) {
			node, err := tagger.ParseToNode(text)
			mr := NewMaResult(name, &node, err)
			c <- mr
		}(name, tagger)
	}

	for name, tagger := range mr.Jumans {
		go func(name string, tagger *jk.Juman) {
			sent, err := tagger.Parse(text)
			mr := new(MaResult)
			mr.Name = name
			mr.Error = err
			if err == nil {
				mr.Surfaces = make([]string, len(sent.Morphemes))
				mr.Features = make([]string, len(sent.Morphemes))

				for i, mrph := range sent.Morphemes {
					mr.Surfaces[i] = mrph.Surface
					mr.Features[i] = mrph.Pos0 + "," + mrph.Pos1 + "," + mrph.RootForm + "," + mrph.Pronunciation + "," + mrph.CType + "," + mrph.CForm + `"` + mrph.Seminfo + `"`
				}
			}
			c <- mr
		}(name, tagger)
	}

	for i := 0; i < waits; i++ {
		result := <-c
		results = append(results, *result)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results
}
