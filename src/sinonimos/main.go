package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"net/http"
	"os"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"

	cli "github.com/jawher/mow.cli"
)

var s = spinner.New(spinner.CharSets[12], 100 * time.Millisecond)

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

	app.Version("v version", "sinonimos v0.1.0")

	app.Run(os.Args)
}

func find(word string) error {
	s.Start()
	resp, err := http.Get(fmt.Sprintf("https://www.sinonimos.com.br/%s/", word))
	if err != nil {
		return err
	}
	s.Stop()

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
