package main


func isHdr(idx int, t *token, tokens []*token) bool {
	if len(tokens) < (idx+1) || idx == 0 {
		return false
	}
	return tokens[idx-1].kind == DASHED && t.kind == UPPER && tokens[idx+1].kind == DASHED
}

func getHdrIndices(tokens []*token)[]int{
	var hdrIndices []int
	for idx, t := range tokens {
		if isHdr(idx, t, tokens){
			hdrIndices = append(hdrIndices, idx)
		}
	}
	return hdrIndices
}


func getSections(tokens []*token) map[string][]*token{
	sections := make(map[string][]*token)
	hdrIndices := getHdrIndices(tokens)
	last := hdrIndices[len(hdrIndices)-1]
	for i, hdrIdx := range hdrIndices {
		secName := tokens[hdrIdx].literal
		if last == hdrIdx {
			sections[secName] = tokens[hdrIdx:]
			continue
		}
		end := hdrIndices[i+1] - 1
		sections[secName] = tokens[hdrIdx+1:end]
	}
	return sections
}