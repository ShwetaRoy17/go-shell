package shell

import (
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
			break
		}
	

		if input == "" {
			continue
		}
		input = strings.Trim(input,"\n")
		if internal.IsPipeline(input){
			s.ExecutePipeline(input)
		}else {
		s.Execute(input)
		}
		

	}

}

