package main

import (
	"fmt"
	"os"
"github.com/ShwetaRoy17/go-shell/app/shell" 
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func main() {
	s := &shell.Shell{}
	s.Run()

}
