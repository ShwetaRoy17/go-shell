package shell

import (
	"fmt"
	"os"
	"strings"
	"github.com/chzyer/readline"
	"github.com/ShwetaRoy17/go-shell/app/internal"
)

const (
	prompt = "$ "
)

type Shell struct {
	autocompleter readline.AutoCompleter
}

func NewShell() *Shell {
	completer := internal.NewCompleter()
	return &Shell{autocompleter: completer}

}

func (s *Shell) Run() {
	
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       prompt,
		HistoryFile:  "/tmp/my-shell.history",
		AutoComplete: s.autocompleter,
		Stdout:       os.Stdout,
	})

	if err != nil {
		panic(err)
	}

	defer rl.Close()

	for true {
		input, err :=rl.Readline()
		if err != nil {
			// panic(err)
			break
		}
	

		if input == "" {
			continue
		}
		input = strings.Trim(input,"\n")
		s.Execute(input)

	}

}

func (s *Shell) Execute(input string) {
	cmd, args := ParseCmd(input)
	// fmt.Println(cmd,args)
	if cmd == "" {
		return
	}
	
	oFile,eFile := os.Stdout,os.Stderr
	origStdout := os.Stdout
	origStderr := os.Stderr
	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()
	
	var err error

	argv, writeOutput, writeError, outputFile, errorFile, mode := redirectInput(args)

	if writeOutput {
		oFile, err = CreateFile(outputFile, mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "redirect error : %v\n", err)
		}
		os.Stdout = oFile
		defer oFile.Close()
	}

	if writeError {
		eFile, err = CreateFile(errorFile, mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "redirect error : %v\n", err)

		}
		os.Stderr = eFile
		defer eFile.Close()
	}

	switch cmd {
	case "type":
		TypFun(argv)
	case "echo":
		EchoCmd(argv)
	case "exit":
		ExitCmd(argv)
	case "pwd":
		Pwd()
	case "cd":
		Cd(argv)
	default:
		ExtProg(cmd, argv,oFile,eFile)

	}

}
