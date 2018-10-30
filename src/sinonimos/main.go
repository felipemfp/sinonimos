package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"

	cli "github.com/jawher/mow.cli"
)

var name = "sinonimos"
var version = "v0.2.0"
var s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)

func main() {
	app := cli.App(name, "Encontre sinônimos")

	app.Spec = "PALAVRA"

	var (
		word = app.StringArg("PALAVRA", "", "Palavra para encontrar sinônimos")
	)

	app.Action = func() {
		if *word == "" {
			fmt.Fprintln(os.Stderr, "Error: incorrect usage")
			app.PrintHelp()
			os.Exit(1)
		}
		fmt.Printf("Buscando sinônimos para \"%s\":\n", *word)
		err := find(*word)
		if err != nil {
			fmt.Printf("\n... falhou (%s)\n", err.Error())
		}
	}

	app.Version("v version", fmt.Sprintf("%s %s", name, version))

	app.Run(os.Args)
}

func find(word string) error {
	s.Start()
	resp, err := http.Get(fmt.Sprintf("https://www.sinonimos.com.br/%s/", word))
	if err != nil {
		return err
	}
	s.Stop()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("  %s\n", aurora.Red("Desculpa, mas não encontramos nenhum sinônimo"))
		os.Exit(1)
	}

	body := charmap.ISO8859_1.NewDecoder().Reader(resp.Body)
	root, err := html.Parse(body)
	if err != nil {
		return err
	}

	meaningSections := scrape.FindAll(root, scrape.ByClass("s-wrapper"))
	for j, meaningSection := range meaningSections {
		if meaning, ok := scrape.Find(meaningSection, scrape.ByClass("sentido")); ok {
			fmt.Printf("\n> %s\n", aurora.Colorize(scrape.Text(meaning), getColors(j)).Bold())

			synonyms := scrape.FindAll(meaningSection, scrape.ByClass("sinonimo"))
			fmt.Print("  ")
			for i, synonym := range synonyms {
				fmt.Print(scrape.Text(synonym))
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

func getColors(index int) aurora.Color {
	colors := []aurora.Color{
		aurora.GreenFg,
		aurora.BrownFg,
		aurora.BlueFg,
		aurora.MagentaFg,
		aurora.CyanFg,
	}
	if index != 0 {
		index = index % len(colors)
	}
	return colors[index]
}
