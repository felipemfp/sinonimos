package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/felipemfp/sinonimos"
	cli "github.com/jawher/mow.cli"
	"github.com/logrusorgru/aurora"
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

		// if expression == "" {
		// 	fmt.Fprintln(os.Stderr, "Error: incorrect usage")
		// 	app.PrintHelp()
		// 	os.Exit(1)
		// }

		log.SetFlags(0)

		log.Printf("Buscando sinônimos para \"%s\":\n", expression)
		s.Start()
		output, err := sinonimos.Find(&sinonimos.FindInput{
			Expression: expression,
		})
		s.Stop()
		if err != nil {
			switch err {
			case sinonimos.ErrNotFound:
				log.Fatalf("  %s\n", aurora.Red("Desculpa, mas não encontramos nenhum sinônimo"))
			default:
				log.Fatal(aurora.Red(fmt.Sprintf("\n... falhou (%s)\n", err.Error())))
			}
		}

		for j, meaning := range output.Meanings {
			if meaning.Description != "" {
				log.Printf("\n> %s\n", aurora.Colorize(meaning.Description, getColors(j)).Bold())
			} else {
				log.Print("\n> -\n")
			}

			log.Printf("  %s\n", strings.Join(meaning.Synonyms, ", "))

			for _, example := range meaning.Examples {
				log.Print(aurora.Gray(fmt.Sprintf("  + %s\n", example)))
			}
		}
	}

	app.Version("v version", fmt.Sprintf("%s %s", name, version))

	app.Run(os.Args)
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
