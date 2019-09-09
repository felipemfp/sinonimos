package sinonimos

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gosimple/slug"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/text/encoding/charmap"
)

var (
	// ErrNotFound is returned when an expression is not found on sinonimos.com.br.
	ErrNotFound = errors.New("expression not found")
)

// Meaning contains information about an meaning.
//
// See Also
//
// Find
type Meaning struct {
	Description string
	Synonyms    []string
	Examples    []string
}

// FindInput contains the input data require to Find.
//
// See Also
//
// Find
type FindInput struct {
	Expression string
}

// FindOutput contains the output payload from Find.
//
// See Also
//
// Find
type FindOutput struct {
	Meanings []Meaning
}

// Find try to find meanings for an expression on sinonimos.com.br.
func Find(input *FindInput) (*FindOutput, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.sinonimos.com.br/%s/", slug.Make(input.Expression)))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	body := charmap.ISO8859_1.NewDecoder().Reader(resp.Body)
	root, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	meaningSections := scrape.FindAll(root, scrape.ByClass("s-wrapper"))
	meanings := make([]Meaning, len(meaningSections))

	for j, meaningSection := range meaningSections {
		if meaning, ok := scrape.Find(meaningSection, scrape.ByClass("sentido")); ok {
			meanings[j].Description = scrape.Text(meaning)
		}

		synonyms := scrape.FindAll(meaningSection, synonymMatcher)
		meanings[j].Synonyms = make([]string, len(synonyms))
		for i, synonym := range synonyms {
			meanings[j].Synonyms[i] = scrape.Text(synonym)
		}

		examples := scrape.FindAll(meaningSection, scrape.ByClass("exemplo"))
		meanings[j].Examples = make([]string, len(examples))
		for i, example := range examples {
			meanings[j].Examples[i] = scrape.Text(example)
		}
	}

	return &FindOutput{
		Meanings: meanings,
	}, nil
}

func synonymMatcher(n *html.Node) bool {
	if n.DataAtom == atom.A || n.DataAtom == atom.Span {
		if n.Parent != nil {
			return scrape.Attr(n.Parent, "class") == "sinonimos" && scrape.Attr(n, "class") != "exemplo"
		}
	}
	return false
}
