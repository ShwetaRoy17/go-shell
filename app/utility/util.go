package utility

var EscCh = map[byte]bool{'"': true, '\\': true, '$': true, '`': true}
var BuiltIns = map[string]bool{"type": true, "echo": true, "exit": true, "pwd": true}
var ExtCmd = map[string]bool{"cat": true, "ls": true, "date": true, "touch": true, "rm": true, "mkdir": true, "rmdir": true}


type PipelineCommand struct {
	Name string
	Args []string
}