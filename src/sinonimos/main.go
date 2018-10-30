package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	. "github.com/logrusorgru/aurora"

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
				output := fmt.Sprintf("%s", scrape.Text(synonym))
				fmt.Print(Colorize(output, getColors(i)))
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

func getColors(index int) Color {
	colors := [6]Color{RedFg, GreenFg, BrownFg, BlueFg, MagentaFg, CyanFg}
	if index != 0 {
		index = index % len(colors)
	}
	return colors[index]
}
