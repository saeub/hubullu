package ui

import (
	"github.com/saeub/hubullu/internal/vocabulary"
	"github.com/saeub/hubullu/pkg/translation"
)

// UserInterface is the interface all UI types must satisfy.
type UserInterface interface {
	Run(*vocabulary.Vocabulary, []translation.Translator)
}
