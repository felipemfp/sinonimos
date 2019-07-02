package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html/atom"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"
	"github.com/metal3d/go-slugify"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"

	cli "github.com/jawher/mow.cli"
)

var name = "sinonimos"
var version = "dev"
var s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)

func main() {
	app := cli.App(name, "Encontre sinônimos")

	app.Spec = "EXPRESSAO..."

	var (
		expressionArr = app.StringsArg("EXPRESSAO", nil, "EXPRESSAO para encontrar sinônimos")
	)

	app.Action = func() {
		expression := strings.Join(*expressionArr, " ")

		if expression == "" {
			fmt.Fprintln(os.Stderr, "Error: incorrect usage")
			app.PrintHelp()
			os.Exit(1)
		}

		fmt.Printf("Buscando sinônimos para \"%s\":\n", expression)
		err := find(expression)
		if err != nil {
			fmt.Print(aurora.Red(fmt.Sprintf("\n... falhou (%s)\n", err.Error())))
		}
	}

	app.Version("v version", fmt.Sprintf("%s %s", name, version))

	app.Run(os.Args)
}

func find(expression string) error {
	s.Start()
	resp, err := http.Get(fmt.Sprintf("https://www.sinonimos.com.br/%s/", slugify.Marshal(expression)))
	s.Stop()
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("  %s\n", aurora.Red("Desculpa, mas não encontramos nenhum sinônimo"))
		os.Exit(1)
	}

	body := charmap.ISO8859_1.NewDecoder().Reader(resp.Body)
	root, err := html.Parse(body)
	if err != nil {
		return err
	}

	synonymMatcher := func(n *html.Node) bool {
		if n.DataAtom == atom.A || n.DataAtom == atom.Span {
			if n.Parent != nil {
				return scrape.Attr(n.Parent, "class") == "sinonimos" && scrape.Attr(n, "class") != "exemplo"
			}
		}
		return false
	}

	meaningSections := scrape.FindAll(root, scrape.ByClass("s-wrapper"))
	for j, meaningSection := range meaningSections {
		if meaning, ok := scrape.Find(meaningSection, scrape.ByClass("sentido")); ok {
			fmt.Printf("\n> %s\n", aurora.Colorize(scrape.Text(meaning), getColors(j)).Bold())
		} else {
			fmt.Print("\n> -\n")
		}

		synonyms := scrape.FindAll(meaningSection, synonymMatcher)
		fmt.Print("  ")
		for i, synonym := range synonyms {
			fmt.Print(scrape.Text(synonym))
			if i == (len(synonyms) - 1) {
				fmt.Print("\n")
			} else {
				fmt.Print(", ")
			}
		}

		examples := scrape.FindAll(meaningSection, scrape.ByClass("exemplo"))
		for _, example := range examples {
			fmt.Print(aurora.Gray(fmt.Sprintf("  + %s\n", scrape.Text(example))))
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
