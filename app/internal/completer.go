package internal

import (
	"fmt"
	"github.com/chzyer/readline"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"github.com/ShwetaRoy17/go-shell/app/utility"

)

type MyCompleter struct {
	trie       *utility.Trie
	lastPrefix string
	tabcnt     int
}

func NewCompleter() readline.AutoCompleter {
	t := &utility.Trie{Root: utility.NewTrieNode()}
	t.Insert("echo")
	t.Insert("exit")
	return &MyCompleter{
		t,
		"",
		0,
	}
}
func (m *MyCompleter) reset() {
	m.lastPrefix = ""
	m.tabcnt = 0
}

func (m *MyCompleter) resetTabCnt() {
	m.tabcnt = 0
}

func (m *MyCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	strline := string(line[:pos])
	strline = strings.TrimLeft(strline, " ")
	l := strings.Fields(strline)
	
	if len(l) == 1 && strings.HasPrefix(strline, " ") {
		m.reset()
		return nil, 0
	}
	
	prefix := l[0]
	if prefix != m.lastPrefix {
		m.reset()
		m.lastPrefix = prefix
	}
	m.tabcnt++

	completions := m.findAllCompletion(prefix)
	if len(completions) == 0 {
		fmt.Fprintf(os.Stdout, "\x07")
		return nil, 0
	} else if len(completions) == 1 {
		completion := completions[0]
		suffix := completion[len(prefix):]
		suffix += " "
		newLine = [][]rune{[]rune(suffix)}
		return newLine, len(prefix)

	} else {
		if m.tabcnt == 1 {
			fmt.Fprintf(os.Stdout, "\x07")
			return nil, 0
		} else if m.tabcnt >= 2 {
			m.resetTabCnt()
			sort.Strings(completions)
			fmt.Fprintf(os.Stdout, "\n")
			for ind, comp := range completions {
				if ind > 0 {
					fmt.Fprintf(os.Stdout, "  ")
				}
				fmt.Fprintf(os.Stdout, "%s", comp)
			}
			fmt.Fprintf(os.Stdout, "\n")
			commonPrefix := findCommonPrefix(completions)
			if len(commonPrefix) > len(prefix) {
				suffix := commonPrefix[len(prefix):]
				newLine = [][]rune{[]rune(suffix)}
				length = len(prefix)
				return newLine, length
			}
			newline := [][]rune {[]rune("")}
			return newline,0
		}
	}
	return nil, 0
}

func findCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	if len(strs) == 1 {
		return strs[0]
	}

	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		for !strings.HasPrefix(strs[i], prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}

	}
	return prefix
}

func (m *MyCompleter) findAllCompletion(prefix string) []string {
	seen := make(map[string]bool)
	var completions []string
	triecompletions := m.trie.FindCompletion(prefix)
	for _, comp := range triecompletions {
		if !seen[comp] {
			seen[comp] = true
			completions = append(completions, comp)
		}
	}

	pathCompletions := findExecutableInPath(prefix)
	for _, comp := range pathCompletions {
		if !seen[comp] {
			seen[comp] = true
			completions = append(completions, comp)
		}
	}
	return completions
}

func findExecutableInPath(prefix string) []string {
	var executables []string
	paths := os.Getenv("PATH")
	seen := make(map[string]bool)
	if paths == "" {
		return executables
	}
	dirs := filepath.SplitList(paths)

	for _, dir := range dirs {
		dir = filepath.Clean(dir)
		if dir == "" || dir == "." {
			continue
		}
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			filename := file.Name()
			if strings.HasPrefix(filename, prefix) {

				filePath := filepath.Join(dir, filename)
				if !seen[filename] && isExecutable(filePath) {
					seen[filename] = true
					executables = append(executables, filename)
				}
			}
		}
	}
	return executables
}

func isExecutable(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if info.Mode()&0111 != 0 {
		return true
	}
	return false
}

