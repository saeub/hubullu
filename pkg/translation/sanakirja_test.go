package translation

import (
	"testing"

	"golang.org/x/text/language"
)

func TestLanguages(t *testing.T) {
	srcLang := language.Finnish
	for trgLang := range languageIDs {
		if srcLang == trgLang {
			continue
		}
		sk, err := NewSanakirjaOrg(srcLang, trgLang)
		if err != nil {
			t.Errorf("creating %v→%v translator failed (err: %v)",
				srcLang, trgLang, err)
		}
		trls, err := sk.Translate("hei")
		if len(trls) == 0 {
			t.Errorf("translation %v→%v failed (no translations returned)",
				srcLang, trgLang)
		}
	}
}

func TestTranslate(t *testing.T) {
	sk, err := NewSanakirjaOrg(language.German, language.Finnish)
	if err != nil {
		t.Errorf("creating translator failed (err: %v)", err)
	}
	trls, err := sk.Translate("Baum")
	if err != nil {
		t.Errorf("translation failed (err: %v)", err)
	}
	if len(trls) == 0 {
		t.Errorf("translation failed (no translations returned)")
	} else {
		firstTrl := trls[0]
		if firstTrl.Text == "" {
			t.Errorf("translation failed (text empty)")
		}
	}
}
