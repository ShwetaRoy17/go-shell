package shell

import (
	"fmt"
	"io"
	"os"

	"github.com/ShwetaRoy17/go-shell/app/internal"
)

func (s *Shell) Execute(input string) {
	cmd, args := internal.ParseCmd(input)
	// fmt.Println(cmd,args)
	if cmd == "" {
		return
	}

	oFile, eFile := os.Stdout, os.Stderr
	origStdout := os.Stdout
	origStderr := os.Stderr
	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	var err error

	argv, writeOutput, writeError, outputFile, errorFile, mode := internal.RedirectInput(args)

	if writeOutput {
		oFile, err = internal.CreateFile(outputFile, mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "redirect error : %v\n", err)
		}
		os.Stdout = oFile
		defer oFile.Close()
	}

	if writeError {
		eFile, err = internal.CreateFile(errorFile, mode)
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
		ExtProg(cmd, argv, oFile, eFile)

	}

}

func (s *Shell) ExecutePipeline(input string) {
	

}

func ExecuteBuiltInWithIO(s *Shell, cmd string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	return nil
}
