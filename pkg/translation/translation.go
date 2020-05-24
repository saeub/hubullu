package translation

// Translation contains all the information returned by a translator.
// Can be serialized to JSON.
type Translation struct {
	Text    string `json:"text"`
	Context string `json:"context"`
	Grammar string `json:"grammar"`
}

// Translator is a source for translations.
type Translator interface {
	Translate(string) ([]Translation, error)
	Name() string
}
