package sort

import (
	"bytes"
	"fmt"
	"github.com/neovim/go-client/nvim"
	"regexp"
	"sort"
	"strings"
)

func getPattern(v *nvim.Nvim, filetype string, args []string) string {
	var pattern string
	if len(args) > 0 {
		pattern = args[0]
	} else {
		patterns := map[string]string{}
		v.Var("sort_patterns", patterns)
		if val, ok := patterns[filetype]; ok {
			pattern = val
		} else {
			pattern = ".*"
		}
	}
	return pattern
}

// ByPattern ...
type ByPattern struct {
	input   [][][]byte
	pattern string
}

func (s ByPattern) Len() int {
	return len(s.input)
}

func (s ByPattern) Swap(i, j int) {
	s.input[i], s.input[j] = s.input[j], s.input[i]
}

func (s ByPattern) Less(i, j int) bool {
	r, _ := regexp.Compile(s.pattern)
	a := r.Find(s.input[i][0])
	b := r.Find(s.input[j][0])
	return bytes.Compare(a, b) == -1
}

// ToBlocks converts lines to array of lines
func ToBlocks(lines [][]byte) [][][]byte {
	result := [][][]byte{}

	r, _ := regexp.Compile(`\S`)
	indent := r.FindIndex(lines[0])[0]

	temp := [][]byte{}

	for _, line := range lines {
		if indent == r.FindIndex(line)[0] {
			fmt.Printf("temp length %d\n", len(temp))
			if len(temp) > 0 {
				fmt.Println("temp")
				temp = append(temp, line)
				result = append(result, temp)
				temp = [][]byte{}
			} else {
				fmt.Println("not temp")
				result = append(result, [][]byte{line})
			}
		} else {
			fmt.Println("else")
			if len(temp) > 0 {
				temp = append(temp, line)
			} else {
				temp = append(temp, result[len(result)-1][0])
				result = result[0 : len(result)-1]
				temp = append(temp, line)
			}
		}
	}

	return result
}

// ToLines flattens blocks
func ToLines(blocks [][][]byte) [][]byte {
	results := [][]byte{}
	for _, block := range blocks {
		for _, line := range block {
			results = append(results, line)
		}
	}
	return results
}

// Sort ...
func Sort(v *nvim.Nvim, args []string, r [2]int, eval *struct {
	Filetype string `eval:"&filetype"`
}) error {
	pattern := getPattern(v, strings.Split(eval.Filetype, ".")[0], args)
	buffer, _ := v.CurrentBuffer()
	lineStart := r[0] - 1
	lineEnd := r[1]
	lines, _ := v.BufferLines(buffer, lineStart, lineEnd, false)
	blocks := ToBlocks(lines)
	sort.Sort(ByPattern{blocks, pattern})
	sorted := ToLines(blocks)
	return v.SetBufferLines(buffer, lineStart, lineEnd, false, [][]byte(sorted))
}
