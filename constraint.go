package macomp

//GetConstraints returns segmentation information
func GetConstraints(query string) (IntSet, IntSet, IntSet, int, int) {
	goldSureSeps := IntSet{}
	goldAmbSeps := IntSet{}

	goldSeps := IntSet{}

	idx := 0
	start, end := -1, -1
	isFirst := true
	for _, r := range query {
		posit := idx - len(goldSeps)
		idx++
		if r == '|' {
			goldSeps.Add(posit)
			goldSureSeps.Add(posit)
			if isFirst {
				start = posit
				end = posit
				isFirst = false
			} else {
				end = posit
			}
		} else if r == '?' {
			goldSeps.Add(posit)
			goldAmbSeps.Add(posit)
		}
	}
	return goldSureSeps, goldAmbSeps, goldSeps, start, end
}

//GetSegments returns segment positions
func GetSegments(surfs []string) IntSet {
	seps := make(IntSet, len(surfs)+1)
	seps[0] = struct{}{}

	last := 0
	for _, m := range surfs {
		last += len([]rune(m))
		seps.Add(last)
	}
	return seps
}
