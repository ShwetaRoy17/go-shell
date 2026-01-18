package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ShwetaRoy17/go-shell/app/internal"
	"github.com/ShwetaRoy17/go-shell/app/utility"
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

var externalCmds []*exec.Cmd

func (s *Shell) ExecutePipeline(input string) {
	segments := internal.ParsePipeline(input)

	if len(segments) == 0 {
		return
	}

	if len(segments) == 1 {
		s.Execute(segments[0])
		return
	}

	commands := make([]utility.PipelineCommand, len(segments))
	for i, seg := range segments {
		cmd, args := internal.ParseCmd(seg)
		commands[i] = utility.PipelineCommand{
			Name: cmd,
			Args: args,
		}
	}

	var nxtReader *os.File

	for ind, cmd := range commands {
		var stdin io.Reader = os.Stdin
		var stdout *os.File = os.Stdout
		var currWriter *os.File

		if ind > 0 {
			stdin = nxtReader
		}

		if ind < len(commands)-1 {
			r, w, err := os.Pipe()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Pipe error: %v\n", err)
				return
			}
			stdout = w
			currWriter = w
			nxtReader = r

		}

		isBuiltIn := utility.BuiltIns[cmd.Name]

		if isBuiltIn {
			args, writeOutput, writeError, outputFile, errorFile, mode := internal.RedirectInput(cmd.Args)
			var oriStderr io.Writer = os.Stderr 

			if writeOutput {
				f, err := internal.CreateFile(outputFile, mode)
				if err == nil {
					stdout = f
					defer f.Close()
				}
			}

			if writeError {
				f, err := internal.CreateFile(errorFile, mode)
				if err == nil {
					oriStderr = f 
					defer f.Close()
				}
			}

			var stdinData []byte
			if ind > 0 {
				stdinData, _ = io.ReadAll(stdin)
			}
			ExecuteBuiltInWithIO(s, cmd.Name, args, bytes.NewReader(stdinData), stdout, oriStderr)

		} else {
			path, exists := isExecutable(cmd.Name)
			if !exists {
				fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd.Name)

				if ind < len(commands)-1 {
					currWriter.Close()
				}
				if ind > 0 {
					if r, ok := stdin.(*os.File); ok && r != os.Stdin {
						r.Close()
					}
				}
				continue 

			}
			args, writeOutput, writeError, outputFile, errorFile, mode := internal.RedirectInput(cmd.Args)
			currCmd := exec.Command(path, args...)
			currCmd.Args[0] = cmd.Name

			
			currCmd.Stdin = stdin

			if writeOutput {
				f, err := internal.CreateFile(outputFile, mode)
				if err == nil {
					currCmd.Stdout = f
					defer f.Close()
				}
			} else if ind < len(commands)-1 {
				currCmd.Stdout = stdout
			} else {
				currCmd.Stdout = os.Stdout 
			}

			if writeError {
				f, err := internal.CreateFile(errorFile, mode)
				if err == nil {
					currCmd.Stderr = f
					defer f.Close()
				}
			} else {
				currCmd.Stderr = os.Stderr
			}

			if err := currCmd.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", cmd.Name, err)
			} else {
				externalCmds = append(externalCmds, currCmd)
			}

		}
		if currWriter != nil {
			currWriter.Close()
		}
		if r, ok := stdin.(*os.File); ok && r != os.Stdin {
			r.Close()
		}
	}
	for _, cmd := range externalCmds {
		cmd.Wait()
	}

}

func ExecuteBuiltInWithIO(s *Shell, cmd string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	switch cmd {
	case "echo":
		output := strings.Join(args, " ")
		fmt.Fprintln(stdout, output)
		return nil
	case "type":
		if len(args) == 0 {
			return nil
		}
		val := args[0]
		if utility.BuiltIns[val] {
			fmt.Fprintf(stdout, "%s is a shell builtin\n", val)
		} else if file, exists := findInPath(val); exists {
			fmt.Fprintf(stdout, "%s is a %s\n", val, file)
		} else {
			fmt.Fprintf(stdout, "%s: not found\n", val)
		}
		return nil
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			if strings.HasPrefix(dir, "//") {
				dir = dir[1:]
			}
		}
		fmt.Fprintln(stdout, dir)
		return nil
	case "exit":
		ExitCmd(args)
		return nil
	case "cd":
		Cd(args)
		return nil
	default:
		return fmt.Errorf("unknown buildin command")
	}

}
