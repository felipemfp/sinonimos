package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/felipemfp/sinonimos"
	cli "github.com/jawher/mow.cli"
)

var name = "sinonimos"
var version = "v0.7.0"
var s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)

func main() {
	app := cli.App(name, "Encontre sinônimos")

	app.Spec = "EXPRESSAO..."

	var (
		expressionArr = app.StringsArg("EXPRESSAO", nil, "EXPRESSAO para encontrar sinônimos")
	)

	log.SetFlags(0)
	log.SetOutput(color.Output)

	app.Action = func() {
		expression := strings.Join(*expressionArr, " ")

		log.Printf("Buscando sinônimos para \"%s\":\n", expression)
		s.Start()
		output, err := sinonimos.Find(&sinonimos.FindInput{
			Expression: expression,
		})
		s.Stop()
		if err != nil {
			switch err {
			case sinonimos.ErrHTTPLayer:
				log.Fatalf("  %s\n", color.RedString("Desculpe, não foi possível obter resposta do " +
																	"servidor do sinonimos...a internet está acessível?"))
			case sinonimos.ErrInvalidFormatBody:
				log.Fatalf("  %s\n", color.RedString("Desculpa, obtivemos uma resposta em um formato " +
					"												inválido do servidor do sinonimos.com.br"))
			case sinonimos.ErrNotFound:
				log.Fatalf("  %s\n", color.RedString("Desculpa, mas não encontramos nenhum sinônimo"))
			default:
				log.Fatal(color.RedString(fmt.Sprintf("\n... falhou (%s)\n", err.Error())))
			}
		}

		for j, meaning := range output.Meanings {
			if meaning.Description != "" {
				log.Printf("\n> %s\n", getColorized(j, meaning.Description))
			} else {
				log.Print("\n> -\n")
			}

			log.Printf("  %s\n", strings.Join(meaning.Synonyms, ", "))
			for _, example := range meaning.Examples {
				log.Print(color.HiWhiteString(fmt.Sprintf("  + %s", example)))
			}
		}
	}

	app.Version("v version", fmt.Sprintf("%s %s", name, version))

	app.Run(os.Args)
}

var colors = []color.Attribute{
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
}

func getColorized(index int, str string) string {
	c := color.New(colors[index%len(colors)], color.Bold)
	return c.Sprint(str)
}
