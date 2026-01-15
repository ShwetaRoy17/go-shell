package internal

import (
	"os"
	"errors"
)


func RedirectInput(args []string) (clean []string, writeOutput bool, writeError bool, outputFile string, errorFile string, mode rune) {
	
	mode = 'a'
	clean = []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case ">", "1>":
			if i+1 < len(args) {
				outputFile = args[i+1]
				writeOutput = true
				mode = 'w'
				i++
				continue
			}
			return clean, writeOutput, writeError, outputFile, errorFile, mode

		case ">>", "1>>":
			if i+1 < len(args) {
				outputFile = args[i+1]
				writeOutput = true
				mode = 'a'
				i++
				continue
			}
			return clean, writeOutput, writeError, outputFile, errorFile, mode

		case "2>":
			if i+1 < len(args) {
				errorFile = args[i+1]
				writeError = true
				mode = 'w'
				i++
				continue
			}
			return clean, writeOutput, writeError, outputFile, errorFile, mode

		case "2>>":
			if i+1 < len(args) {
				errorFile = args[i+1]
				writeError = true
				mode = 'a'
				i++
				continue
			}
			return clean, writeOutput, writeError, outputFile, errorFile, mode
		}

		clean = append(clean,args[i])
	}
	return clean, writeOutput, writeError, outputFile, errorFile, mode
}


// file creation based on mode
func CreateFile(filepath string, mode rune) (*os.File, error) {
	var file *os.File
	var err error
	if mode == 'a' {
		file, err = os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(filepath)
		if err != nil {
			return nil, errors.New("Error creating file " + filepath)
		}
	}
	return file, err

}


