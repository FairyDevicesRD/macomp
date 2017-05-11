package macomp

import (
	"strings"

	"github.com/fatih/color"
	runewidth "github.com/mattn/go-runewidth"
)

//Blank is an expression for a non segment
const Blank = " "

//Segment is an expression for a valid segment
const Segment = "|"

//IllegalSegment is an expression for an illegal segment
const IllegalSegment = "/"

//IllegalUnsegment is an expression for an illegal unsegmented position
const IllegalUnsegment = "_"

func getSurfsCapacity(surfs []string) (surflen int) {
	for _, surf := range surfs {
		surflen += len(surf)
	}
	return surflen*2 + 2
}

//DecorateSurfaces returns a decorated line
func DecorateSurfaces(text string, surfs []string) (string, bool) {
	goldSureSeps, _, goldSeps, start, end := GetConstraints(text)
	isValid := true
	if len(goldSeps) == 0 {
		return PrettySurfaces(surfs), isValid
	}

	ret := make([]string, 0, getSurfsCapacity(surfs))
	posit := 0
	for _, surf := range surfs {
		if posit < start || posit > end {
			ret = append(ret, Segment)
		} else {
			if _, ok := goldSeps[posit]; ok {
				ret = append(ret, Segment)
			} else {
				ret = append(ret, color.RedString((IllegalSegment)))
				isValid = false
			}
		}

		rsurf := []rune(surf)
		for idx, char := range rsurf {
			ret = append(ret, string(char))
			posit++
			if idx != len(rsurf)-1 {
				if posit < start || posit > end {
					ret = append(ret, Blank)
				} else {
					if _, ok := goldSureSeps[posit]; ok {
						ret = append(ret, color.RedString(IllegalUnsegment))
						isValid = false
					} else {
						ret = append(ret, Blank)
					}
				}
			}
		}
	}
	ret = append(ret, Segment)
	return strings.Join(ret, ""), isValid
}

//PrettySurfaces returns a pretty line
func PrettySurfaces(surfs []string) string {
	ret := make([]string, 0, getSurfsCapacity(surfs))
	posit := 0
	for _, surf := range surfs {
		ret = append(ret, Segment)

		rsurf := []rune(surf)
		for idx, char := range rsurf {
			ret = append(ret, string(char))
			if idx != len(rsurf)-1 {
				ret = append(ret, Blank)
			}
			posit++
		}
	}
	ret = append(ret, Segment)
	return strings.Join(ret, "")
}

//PrettyFeatures returns a pretty line
func PrettyFeatures(surfs []string, features []string) string {
	if len(surfs) != len(features) {
		return ""
	}

	ret := make([]string, 0, len(surfs)*2+2)
	for i, surf := range surfs {
		rsurf := []rune(surf)
		width := runewidth.StringWidth(surf) + len(rsurf) - 1
		var pos string
		rf := []rune(features[i])
		if len(rf) > 0 {
			pos = string(rf[0])
		}

		surfRWidth := runewidth.StringWidth(surf)
		if surfRWidth >= 2 {
			ret = append(ret, Segment)
		}
		ret = append(ret, runewidth.FillRight(pos, width))
	}
	ret = append(ret, Segment)
	return strings.Join(ret, "")
}
