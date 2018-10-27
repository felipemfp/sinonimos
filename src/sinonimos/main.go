package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"

	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("sinonimos", "Encontre sinônimos")

	app.Spec = "PALAVRA"

	var (
		word = app.StringArg("PALAVRA", "", "Palavra para encontrar sinônimos")
	)

	app.Action = func() {
		fmt.Printf("Buscando sinônimos para \"%s\":\n", *word)
		err := find(*word)
		if err != nil {
			fmt.Printf("... falhou (%s)\n", err.Error())
		}
	}

	app.Run(os.Args)
}

func find(word string) error {
	resp, err := http.Get(fmt.Sprintf("https://www.sinonimos.com.br/%s/", word))
	if err != nil {
		return err
	}

	body := charmap.ISO8859_1.NewDecoder().Reader(resp.Body)
	root, err := html.Parse(body)
	if err != nil {
		return err
	}

	meaningSections := scrape.FindAll(root, scrape.ByClass("s-wrapper"))
	for _, meaningSection := range meaningSections {
		if meaning, ok := scrape.Find(meaningSection, scrape.ByClass("sentido")); ok {
			fmt.Printf("\n> %s\n", scrape.Text(meaning))

			synonyms := scrape.FindAll(meaningSection, scrape.ByClass("sinonimo"))
			fmt.Print("  ")
			for i, synonym := range synonyms {
				fmt.Printf("%s", scrape.Text(synonym))
				if i == (len(synonyms) - 1) {
					fmt.Print("\n")
				} else {
					fmt.Print(", ")
				}
			}
		}
	}

	return nil
}
