package acrostic

import (
	"fmt"
	"os"
	"testing"

	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

func TestJumanKnpExecute(t *testing.T) {

	o := &Options{
		JumanCommand:   "jumanpp",
		KnpCommand:     "knp -tab -anaphora",
		JumanDirectory: "/usr/local/share/juman",
	}
	i := &Instance{}
	jk, err := NewJumanKnp(o, i)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	t.Logf("execute")
	jk.Execute([]rune("2丁目の花子さんは日曜日に一郎さんとピクニックに行った．"), true)
	//t.Logf("%v\n", string(r))
	jk.Execute([]rune("帽子を被った田中さんと横山さんはゲームセンターに行くようだ．"), false)
	//t.Logf("%v\n", string(r))
}

func TestReadInflection(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{DisableSorting: false, QuoteEmptyFields: true})
	o := &Options{
		JumanCommand:         "jumanpp",
		KnpCommand:           "knp -tab -anaphora",
		JumanDirectory:       "/usr/local/share/juman",
		WordNetDatabase:      "third-party/wnjpn/wnjpn.db",
		SynonymsJapaneseOnly: true,
		UseSynsetList:        false,
		Interactive:          false,
	}
	i := &Instance{}
	jk, err := NewJumanKnp(o, i)
	if err != nil {
		t.Errorf("cannot initialize JumanKnp")
	}

	err = jk.ReadInflection()
	if err != nil {
		t.Errorf(err.Error())
	}
	if v, o := jk.InflectionDB["子音動詞カ行"]; o {
		if _, ok := v["意志形"]; ok {
			fmt.Println("found")
		} else {
			t.Errorf("form not found")
			for i := range v {
				t.Errorf(v[i])
			}
			if len(v) == 0 {
				t.Errorf("form is null in the type")
			}
		}
	} else {
		t.Errorf("type not found")
	}
	tyexpected := []rune("子音動詞カ行")
	rexpected := []rune("歩こう")
	var ty []rune
	var r []rune
	f := false
	ty, r, f = jk.Inflection([]rune("歩く"), []rune("意志形"))
	if f == false {
		t.Errorf("not found")
	}
	if !runes.Compare(ty, tyexpected) {
		t.Errorf("want %v, but returned %v", tyexpected, ty)
	}
	if !runes.Compare(r, rexpected) {
		t.Errorf("want %v, but returned %v", rexpected, r)
	}
}
