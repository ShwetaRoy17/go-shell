package internal

import (
	"strings"
)

func IsPipeline(input string) bool {
	return strings.Contains(input, " | ")
}

func ParsePipeline(input string) []string {
	var segments []string
	var curr strings.Builder
	singleQ, inDoubleQ := false, false
	for i := 0; i < len(input); i++ {
		ch := input[i]

		if ch == '\'' && !inDoubleQ {
			singleQ = !singleQ
		} else if ch == '"' && !singleQ {

		} else if ch == '|' && !singleQ && !inDoubleQ {
			segments = append(segments, curr.String())
			curr.Reset()
		}
		curr.WriteByte(ch)

	}
	if curr.Len() > 0 {
		segments = append(segments, curr.String())
	}
	return segments
}


