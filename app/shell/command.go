package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ShwetaRoy17/go-shell/app/utility"
)

func TypFun(argv []string) {

	if len(argv) == 0 {
		return
	}

	val := argv[0]
	outputString := ""
	if utility.BuiltIns[val] {
		outputString = fmt.Sprintf("%s is a shell builtin\n", val)

	} else if file, exists := findInPath(val); exists {
		outputString = fmt.Sprintf("%s is %s\n", val, file)

	} else {
		outputString = fmt.Sprintf("%s: not found\n", val)
	}

	fmt.Printf(outputString)

}

func ExitCmd(argv []string) {
	code := 0
	if len(argv) > 0 {
		val, err := strconv.Atoi(argv[0])
		if err == nil {
			code = val
		}
	}
	os.Exit(code)
}

func EchoCmd(argv []string) {
	output := strings.Join(argv, " ")
	fmt.Println(output)
}

func ExtProg(command string, argv []string, oFile, eFile *os.File) {
	path, exists := isExecutable(command)

	if exists {
		cmd := exec.Command(path, argv...)
		cmd.Args[0] = command

		cmd.Stdin = os.Stdin
		cmd.Stdout = oFile
		cmd.Stderr = eFile

		if err := cmd.Run(); err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				fmt.Fprintf(eFile, "%s: %v\n", command, err)
			}
		}
	} else {
		fmt.Printf("%s: command not found\n", command)
	}
}

func isExecutable(filePath string) (string, bool) {
	path, err := exec.LookPath(filePath)
	if err != nil {
		return "", false
	}
	return path, true
}

func Pwd() {
	dir, err := os.Getwd()
	if err == nil {
		if strings.HasPrefix(dir, "//") {
			dir = dir[1:]
		}

		fmt.Println(dir)
	}
}

func Cd(argv []string) {

	path := ""

	if len(argv) < 1 {
		path = "~"
	} else {
		path = argv[0]
	}

	if path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {

			fmt.Fprintln(os.Stderr, "cd: could not find home directory")
			return
		}
		path = homeDir
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", path)
	}
}

func SplitCmd(command string) []string {

	s := []string{}
	singleQ, doubleQ, esc := false, false, false
	curr := ""

	n := len(command)
	for i := 0; i < n; i++ {
		ch := command[i]
		if esc && doubleQ {
			if !utility.EscCh[ch] {
				curr += "\\"

			}
			curr += string(ch)
			esc = false
		} else if esc {
			esc = false
			curr += string(ch)
		} else if ch == '\'' && !doubleQ {
			singleQ = !singleQ
		} else if ch == '"' && !singleQ {
			doubleQ = !doubleQ
		} else if ch == '\\' && !singleQ {
			esc = true
		} else if ch == ' ' && !singleQ && !doubleQ {
			if curr != "" {
				s = append(s, curr)
				curr = ""
			}
		} else {
			curr += (string)(ch)
		}

	}

	if curr != "" {
		s = append(s, curr)
	}

	return s
}

func findInPath(bin string) (string, bool) {
	if file, exec := isExecutable(bin); exec {
		return file, true
	}
	paths := os.Getenv("PATH")
	arr := strings.Split(paths, ":")
	for _, path := range arr {
		fullpath := filepath.Join(path, bin)
		if file, err := isExecutable(fullpath); err {
			return file, true
		}
		if _, err := os.Stat(fullpath); err == nil {
			return fullpath, true
		}
	}
	return "", false
}

func HistoryCmd(argv []string, historyList *[]string, appendHistoryList *[]string) {
	n := len(argv)
	if n > 2 {
		return
	}
	hisLen := len(*historyList)
	if n == 2 {
		filepath := argv[1]
		switch argv[0] {
		case "-r":
			readHistory(historyList, filepath)
		case "-a":
			appendHistory(historyList, filepath)
		case "-w":
			writeHistory(historyList, filepath)
		default:
		}
	}
	if n == 1 {
		var err error
		n, err = strconv.Atoi(argv[0])
		if err == nil {
			n = min(n, hisLen)
		}
	} else {
		n = hisLen
	}

	for i := hisLen - n; i < hisLen; i++ {
		fmt.Printf("    %d  %s\n", i+1, (*historyList)[i])
	}
}

func readHistory(historyList *[]string, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			*historyList = append(*historyList, line)
		}
	}
	return scanner.Err()
}

func writeHistory(historyList *[]string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range *historyList {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func appendHistory(historyList *[]string, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, line := range *historyList {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()
	*historyList = []string{}
	return nil
}
