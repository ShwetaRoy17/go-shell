package shell

import (
	"bufio"
	"fmt"
	"os"
)

type Shell struct {
}

func New() *Shell {
	return &Shell{}
}

func (s *Shell) Run() {
	for true {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		// fmt.Println("input:",input)
		if err != nil {
			panic(err)
		}
		// input = strings.Trim(input,"\n")
		// input = strings.Trim(input," ")

		if input == "" {
			continue
		}

		s.Execute(input)

	}

}

func (s *Shell) Execute(input string) {
	cmd, args := ParseCmd(input)
	

	argv, writeOutput, writeError, outputFile, errorFile, mode := redirectInput(args)
	
	if writeOutput {
		oFile, err := CreateFile(outputFile, mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "redirect error : %v\n", err)
		}
		os.Stdout = oFile
		defer oFile.Close()
	}

	if writeError {
		eFile, err := CreateFile(errorFile, mode)
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
		Pwd(argv[len(argv)-1])
	case "cd":
		Cd(argv)
	default:
		ExtProg(cmd,argv)

	}

}
