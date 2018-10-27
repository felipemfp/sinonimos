package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("sinonimos", "Encontre sinônimos")

	app.Spec = "PALAVRA"

	var (
		word = app.StringsArg("PALAVRA", nil, "Palavra para encontrar sinônimos")
	)

	app.Action = func() {
		fmt.Printf("Buscando sinônimos para %s\n", *word)
	}

	app.Run(os.Args)
}
