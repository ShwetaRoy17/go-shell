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
	historyList       []string
	appendHistoryList []string
	shouldExit        bool
	exitCode          int
}

func NewShell() *Shell {
	completer := internal.NewCompleter()
	return &Shell{autocompleter: completer}

}

func (s *Shell) Run() int {
	
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

	historyFile := os.Getenv("HISTFILE")
	if historyFile != ""{
		err := s.readHistory(historyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
		defer func(){
			s.dumpHistory(historyFile)
		}()
	}
	
	for !s.shouldExit {
		input, err :=rl.Readline()
		if err != nil {
			break
		}
	

		if input == "" {
			continue
		}
		input = strings.Trim(input,"\n")
		s.historyList = append(s.historyList,input)
		s.appendHistoryList = append(s.appendHistoryList,input)
		if internal.IsPipeline(input){
			s.ExecutePipeline(input)
		}else {
		s.Execute(input)
		}
		

	}
return s.exitCode
}

