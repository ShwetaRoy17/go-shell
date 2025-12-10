package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var builtIns = map[string]bool{"type": true, "echo": true, "exit": true, "pwd": true}
var extCmd = map[string]bool{"cat": true, "ls": true, "date": true, "touch": true, "rm": true, "mkdir": true, "rmdir": true}

func TypFun(argv []string) {

	if len(argv) == 0 {
		return
	}

	val := argv[0]
	outputString := ""
	if builtIns[val] {
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
	if len(argv) > 1 {
		val, err := strconv.Atoi(argv[1])
		if err != nil {
			code = val
		}
	}
	os.Exit(code)
}

func EchoCmd(argv []string) {
	output := strings.Join(argv, " ")
	fmt.Println(output)
}

func ExtProg(command string,argv []string) {
	if extCmd[command] {
		var cmd *exec.Cmd
		if len(argv) < 2 {
			cmd = exec.Command(command)
		} else {
			cmd = exec.Command(command, argv...)
		}
		// cmd.Args = argv.Args // Set argv to use original command name as argv[0]

		cmd.Args = argv
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
		}
		return
	}

	path, exists := isExecutable(command)
	if exists || builtIns[path] {
		cmd := exec.Command(path, argv...)
		// cmd.Args = argv.Args // Set argv to use original command name as argv[0]
		cmd.Args = argv
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {

			errorOp := "Error executing command:" + err.Error()

			fmt.Fprintln(os.Stderr, errorOp)
		}

	} else {
		output := fmt.Sprintf("%s: command not found\n", command)

		fmt.Printf(output)

	}

}

func isExecutable(filePath string) (string, bool) {

	// this will tell if the command exists in the path or not
	path, err := exec.LookPath(filePath)
	if err != nil {
		return "", false
	}
	return path, true
}

func Pwd(filename string) {
	dir, err := filepath.Abs(".")

	if err == nil {
		fmt.Println(dir)
	}

}

func Cd(argv []string) {
	if len(argv) < 2 {
		return
	}
	if argv[1] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {

			fmt.Fprintln(os.Stderr, "cd: could not find home directory")
			return
		}
		err = os.Chdir(homeDir)
		if err != nil {
			op := fmt.Sprintf("cd: %s: No such file or directory\n", homeDir)

			fmt.Fprintf(os.Stderr, op)

		}
		return
	}
	path := argv[1]
	err := os.Chdir(path)
	if err != nil {
		op := fmt.Sprintf("cd: %s: No such file or directory\n", path)

		fmt.Fprintf(os.Stderr, op)

	}
}

var escCh = map[byte]bool{'"': true, '\\': true, '$': true, '`': true}

func SplitCmd(command string) []string {

	s := []string{}
	singleQ, doubleQ, esc := false, false, false
	curr := ""

	n := len(command)
	for i := 0; i < n-1; i++ {
		ch := command[i]
		if esc && doubleQ {
			if !escCh[ch] {
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

