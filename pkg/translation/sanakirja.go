package translation

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/text/language"
)

var languageIDs = map[language.Tag]int{
	language.Bulgarian: 1,
	language.Estonian:  2,
	language.English:   3,
	language.Spanish:   4,
	// language.Esperanto: 5, // supported by Sanakirja.org, but not by language package
	language.Italian: 6,
	language.Greek:   7,
	// language.Latin: 8, // supported by Sanakirja.org, but not by language package
	language.Latvian:    9,
	language.Lithuanian: 10,
	language.Norwegian:  11,
	language.Portuguese: 12,
	language.Polish:     13,
	language.French:     14,
	language.Swedish:    15,
	language.German:     16,
	language.Finnish:    17,
	language.Danish:     18,
	language.Czech:      19,
	language.Turkish:    20,
	language.Hungarian:  21,
	language.Russian:    22,
	language.Dutch:      23,
	language.Japanese:   24,
}

var grammarRegex = regexp.MustCompile(` \{([^}]+)\}`)

// SanakirjaOrg is a translator type which scrapes Sanakirja.org.
type SanakirjaOrg struct {
	urlPattern string
}

// NewSanakirjaOrg initializes a new translator for Sanakirja.org.
// If the given languages are not supported, this returns a non-nil error.
// Note: Sanakirja.org automatically searches for translations in the opposite
// direction when there are no results.
func NewSanakirjaOrg(sourceLanguage, targetLanguage language.Tag) (*SanakirjaOrg, error) {
	sourceLanguageID, sourceOk := languageIDs[sourceLanguage]
	targetLanguageID, targetOk := languageIDs[targetLanguage]
	if !sourceOk {
		return &SanakirjaOrg{}, fmt.Errorf("language %v not supported", sourceLanguage)
	}
	if !targetOk {
		return &SanakirjaOrg{}, fmt.Errorf("language %v not supported", targetLanguage)
	}
	urlPattern := fmt.Sprintf("https://www.sanakirja.org/search.php?l=%d&l2=%d&q=%%s",
		sourceLanguageID, targetLanguageID)
	return &SanakirjaOrg{urlPattern}, nil
}

// Name returns the source of the translations.
func (t *SanakirjaOrg) Name() string {
	return "Sanakirja.org"
}

// Translate scrapes Sanakirja.org for translations of the given word.
func (t *SanakirjaOrg) Translate(word string) ([]Translation, error) {
	query := url.PathEscape(word)
	url := fmt.Sprintf(t.urlPattern, query)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("%s responded with status code %d",
			url, response.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	var trls []Translation

	textColumn := -1
	contextColumn := -1
	doc.Find("table#translations tbody tr").Each(func(_ int, row *goquery.Selection) {
		if row.HasClass("th_class") {
			row.Find("th").Each(func(i int, header *goquery.Selection) {
				switch header.Text() {
				case "Käännös":
					textColumn = i + 1
				case "Konteksti":
					contextColumn = i + 1
				}
			})
		} else if !row.HasClass("group_name") {
			cells := row.Find("td")
			trl := Translation{}
			if textColumn != -1 {
				text := cells.Eq(1).Text()
				if grammarRegex.MatchString(text) {
					trl.Grammar = grammarRegex.FindStringSubmatch(text)[1]
					text = grammarRegex.ReplaceAllString(text, "")
				}
				trl.Text = text
			}
			if contextColumn != -1 {
				trl.Context = cells.Eq(contextColumn).Text()
			}
			trls = append(trls, trl)
		}
	})
	return trls, nil
}
