package puff

import (
	"os"
	"strings"
)

func readDescription(file string, lineNumber int, ok bool) string {
	if !ok {
		return ""
	}
	srcfile, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(srcfile), "\n")
	comments := []string{}
	// read file in reverse
	for i := lineNumber - 2; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" { // empty line (ignore)
			continue
		}
		if strings.HasPrefix(line, "//") { // we have reached comment
			comment := strings.TrimSpace(line[2:])
			comments = append([]string{comment}, comments...) // since we're reading in reverse, we have to prepend
			continue
		}
		break // since it wasn't an empty line or a comment, we have reached actual code. we are no longer in a comment
	}
	return strings.Join(comments, " ") // trim space removed the spaces
}
