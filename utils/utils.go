package utils

import (
	"regexp"
	"strings"

	"github.com/rivo/uniseg"
)

func splitGraphemeClusters(text string, lenCnt func(string) int, lenLim int, parts []string) []string {
	parts = append(parts, "")
	gr := uniseg.NewGraphemes(text)

	for gr.Next() {
		partLen := lenCnt(parts[len(parts)-1])
		runes := string(gr.Runes())

		if lenLim > partLen+len(runes) {
			parts[len(parts)-1] += runes
		} else {
			parts = append(parts, runes)
		}
	}

	return parts
}

// Splits a message into submessages. Tries to not split words unless it
// absolutely has to in which case it splits based on grapheme clusters.
func Split(text string, lenCnt func(string) int, lenLim int) []string {
	parts := []string{""}
	r := regexp.MustCompile(`[^\s]+|\s+`)

	for _, word := range r.FindAllString(text, -1) {
		wordLen := lenCnt(word)
		partLen := lenCnt(parts[len(parts)-1])

		if lenLim > partLen+wordLen {
			parts[len(parts)-1] += word
		} else if lenLim > wordLen {
			parts = append(parts, word)
		} else {
			parts = splitGraphemeClusters(word, lenCnt, lenLim, parts)
		}
	}

	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}

	return parts
}

func GetCommandArgsPrefix(s, prefix string) (string, []string, bool) {
	args := strings.Fields(s)

	var cmd string
	if len(args) > 0 {
		cmd = args[0]
		args = args[1:]
	}

	hasPrefix := strings.HasPrefix(cmd, prefix)
	cmd = strings.TrimPrefix(cmd, prefix)

	return strings.ToLower(cmd), args, hasPrefix
}
