package macomp

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
)

func TestDecorateSurfaces(t *testing.T) {
	tests := []struct {
		text  string
		surfs []string
		gold  string
		valid bool
	}{
		{
			text:  "|ジャパリ?まん|は|美味しい",
			surfs: []string{"ジャ", "パリ", "まんは", "美味しい"},
			gold:  Segment + "ジ ャ" + IllegalSegment + "パ リ" + Segment + "ま ん" + IllegalUnsegment + "は" + Segment + "美 味 し い" + Segment,
			valid: false,
		},
	}
	for _, test := range tests {
		sysout, isValid := DecorateSurfaces(test.text, test.surfs)
		if sysout != test.gold {
			t.Errorf("Expected %v but got %v", test.gold, sysout)
		}
		if isValid != test.valid {
			t.Errorf("Expected %v but got %v", test.valid, isValid)
		}
	}
}

func BenchmarkDecorateSurfaces(b *testing.B) {
	for i := 0; i < 10000; i++ {
		DecorateSurfaces("|ジャパリ?まん|は|美味しい", []string{"ジャ", "パリ", "まんは", "美味しい"})
		DecorateSurfaces("ジ|ャパリ?まん|は|美味しい", []string{"ジャ", "パリ", "まんは", "美味しい"})
		DecorateSurfaces("ジャパ|リ?まん|は|美味しい", []string{"ジャ", "パリ", "まんは", "美味しい"})
	}
}

func TestPrettySurfaces(t *testing.T) {
	tests := []struct {
		surfs []string
		gold  string
	}{
		{
			surfs: []string{"ジャパリ", "まん", "は", "美味しい"},
			gold:  fmt.Sprintf("%sジ ャ パ リ%sま ん%sは%s美 味 し い%s", Segment, Segment, Segment, Segment, Segment),
		},
	}
	for _, test := range tests {
		sysout := PrettySurfaces(test.surfs)
		if sysout != test.gold {
			t.Errorf("Expected %v but got %v", test.gold, sysout)
		}
	}
}

func TestPrettyFeatures(t *testing.T) {
	tests := []struct {
		surfs    []string
		features []string
		gold     string
	}{
		{
			surfs:    []string{"ジャパリ", "まん", "は", "美味しい"},
			features: []string{"名", "名", "助", "形"},
			gold:     fmt.Sprintf("%s名         %s名   %s助%s形         %s", Segment, Segment, Segment, Segment, Segment),
		},
	}
	for _, test := range tests {
		sysout := PrettyFeatures(test.surfs, test.features)
		if sysout != test.gold {
			t.Errorf("Expected %v but got %v", test.gold, sysout)
		}
	}
}

func init() {
	color.NoColor = true
}
