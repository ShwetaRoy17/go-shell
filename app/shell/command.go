package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)
var builtIns = []string{"type", "echo", "exit"}


func TypFun(argv []string) {

	if len(argv) == 1 {
		return
	}

	val := argv[1]

	if slices.Contains(builtIns, val) {
		fmt.Printf("%s is a shell builtin\n",val)
		return
	}
	if file, exists := FindInPath(val); exists == true {
		fmt.Printf("%s is %s\n", val, file)
		return
	}
	fmt.Printf("%s: not found\n", val)

}

func FindInPath(bin string) (string, bool) {
	paths := os.Getenv("PATH")
	arr := strings.Split(paths, ":")

	for _, path := range arr{
		file:=filepath.Join(path, bin)
		if _, err := os.Stat(file); err == nil {
			return file, true
		}
	}
	return "", false
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
	output := strings.Join(argv[1:], " ")
	fmt.Println(output)
}


func isExectutable(filePath string) (string,bool) {
	
	// this will tell if the command exists in the path or not
    path, err := exec.LookPath(filePath)
    if err != nil {
        return "", false
    }
    return path, true
}

func ExtProg(argv []string) {
	path,exists := isExectutable(argv[0])
	if exists  {
		cmd := exec.Command(path,argv[1:]...)
		// cmd.Args = argv.Args // Set argv to use original command name as argv[0]
		cmd.Args = argv
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Error executing command:", err)
		}

	} else {
		fmt.Printf("%s: command not found\n", argv[0])
	}

}
