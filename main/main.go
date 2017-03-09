package main

import (
	"github.com/markbates/pop/soda/cmd"
    "fmt"
    "os"
)


func main() {
	fmt.Print("--==> HERE\n")
	os.Args = make([]string, 8, 8)
	os.Args[1] = "create"
	os.Args[2] = "-e"
	os.Args[3] = "sqlserver"
	os.Args[4] = "-c"
	os.Args[5] = "/home/wellsj/workspace/go/src/github.com/markbates/pop/database.yml"
	cmd.Execute()
}
