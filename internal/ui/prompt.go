package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/saeub/hubullu/internal/vocabulary"
	"github.com/saeub/hubullu/pkg/translation"
)

// Prompt is a REPL-style user interface type.
type Prompt struct {
	vocab       *vocabulary.Vocabulary
	translators []translation.Translator
}

// NewPrompt initializes a new prompt-type user interface.
func NewPrompt(vocab *vocabulary.Vocabulary, trlrs []translation.Translator) *Prompt {
	return &Prompt{vocab, trlrs}
}

// Run starts the interface.
// While running, the interface will take care of displaying and saving translations.
func (p *Prompt) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	var items []vocabulary.Item

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			fmt.Print("\n")
			break
		}
		word := scanner.Text()
		if number, err := strconv.Atoi(word); err == nil {
			if number >= 1 && number <= len(items) {
				item := items[number-1]
				p.vocab.Add(item)
				err = p.vocab.Save()
				if err != nil {
					fmt.Printf("saving failed: %v\n\n", err)
				} else {
					fmt.Print("saved\n\n")
				}
			} else {
				fmt.Print("out of range\n\n")
			}
		} else {
			i := 1
			for _, trltr := range p.translators {
				fmt.Printf("%s:\n", trltr.Name())
				items, err = vocabulary.CreateItems(word, trltr)
				for _, item := range items {
					trlStr := renderTranslation(item.Translation)
					btrlStr := ""
					for i, btrl := range item.Backtranslations {
						if i >= 5 {
							break
						}
						if i > 0 {
							btrlStr += ", "
						}

						b := aurora.Blue(btrl.Text)
						if strings.ToLower(btrl.Text) == strings.ToLower(word) {
							b = aurora.Bold(b)
						}
						btrlStr += b.String()
					}
					fmt.Printf("[%d] %s → %s\n", i, trlStr, btrlStr)
					i++
				}
				fmt.Print("\n")
			}
		}
	}
}

func renderTranslation(trl translation.Translation) string {
	result := aurora.Bold(trl.Text).String()
	if trl.Grammar != "" {
		result += fmt.Sprintf(" <%s>", trl.Grammar)
	}
	if trl.Context != "" {
		result += fmt.Sprintf(" (%s)", trl.Context)
	}
	return result
}
