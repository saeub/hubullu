package vocabulary

import (
	"os"
	"reflect"
	"testing"

	"github.com/saeub/hubullu/pkg/translation"
)

func TestSaveLoad(t *testing.T) {
	vocab1 := NewVocabulary()
	vocab1.Add(Item{
		Word: "Test",
		Translation: translation.Translation{
			Text:    "koe",
			Context: "ark",
		},
		Backtranslations: []translation.Translation{
			translation.Translation{
				Text:    "Test",
				Grammar: "m",
			},
			translation.Translation{
				Text:    "Test",
				Grammar: "f",
			},
		},
	})
	err := vocab1.Save()
	if err != nil {
		t.Errorf("saving vocabulary to %s failed: %v", vocab1.filename, err)
	}
	vocab2, err := LoadVocabulary(vocab1.filename)
	if err != nil {
		t.Errorf("loading vocabulary from %s failed: %v", vocab1.filename, err)
	}
	if !reflect.DeepEqual(vocab1, vocab2) {
		t.Errorf("vocabulary not the same after loading: %v != %v",
			vocab1.Items[0], vocab2.Items[0])
	}
	err = os.Remove(vocab1.filename)
	if err != nil {
		t.Errorf("deleting %s failed: %v", vocab1.filename, err)
	}
}
