package vocabulary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/saeub/hubullu/pkg/translation"
)

// Vocabulary is a collection of vocabulary items.
// Can be serialized to JSON.
type Vocabulary struct {
	filename string
	Items    []Item `json:"items"`
}

// Item represents a single translation and all of its backtranslations for a word.
// Can be serialized to JSON.
type Item struct {
	Word             string                    `json:"word"`
	Translation      translation.Translation   `json:"translation"`
	Backtranslations []translation.Translation `json:"backtranslations"`
}

// Add adds an item to the vocabulary.
func (v *Vocabulary) Add(item Item) {
	v.Items = append(v.Items, item)
}

// NewVocabulary returns a new, empty vocabulary without an associated filename.
func NewVocabulary() *Vocabulary {
	return &Vocabulary{}
}

func getArbitraryFilename() string {
	base := time.Now().Format("vocab_060102-150405")
	suffix := 0
	extension := ".json"
	var fn string
	for {
		if suffix == 0 {
			fn = fmt.Sprintf("%s", base) + extension
		} else {
			fn = fmt.Sprintf("%s_%d", base, suffix) + extension
		}
		_, err := os.Stat(fn)
		if os.IsNotExist(err) {
			break
		}
		suffix++
	}
	return fn
}

// LoadVocabulary loads a vocabulary from a JSON file at the specified location.
// If the file does not exist, it will create it and return an empty vocabulary.
func LoadVocabulary(filename string) (*Vocabulary, error) {
	vocab := Vocabulary{
		filename: filename,
	}
	data, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		// try to create the file so the name won't get stolen
		_, err = os.Create(filename)
		return &vocab, err
	}
	if err != nil {
		return &vocab, err
	}
	err = json.Unmarshal(data, &vocab)
	return &vocab, err
}

// Save saves the vocabulary to a JSON file.
// If the vocabulary does not already have a filename associated with it,
// an arbitrary non-existing one in the working directory will be picked.
func (v *Vocabulary) Save() error {
	if v.filename == "" {
		v.filename = getArbitraryFilename()
	}
	data, err := json.Marshal(*v)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(v.filename, data, 0644)
	return err
}

// CreateItems fetches translations and backtranslations
// using the given translator and compiles them in vocabulary items.
func CreateItems(word string, t translation.Translator) (items []Item, err error) {
	trls, err := t.Translate(word)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	for i, trl := range trls {
		item := Item{
			Word:        word,
			Translation: trl,
		}
		items = append(items, item)
		wg.Add(1)
		go func(i int) {
			text := items[i].Translation.Text
			if text != "" {
				items[i].Backtranslations, _ = t.Translate(text)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return items, nil
}
