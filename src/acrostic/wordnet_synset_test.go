package acrostic

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func tWordSynset(t *testing.T, r []SynsetResult, synset string, astep int, bstep int, depth int) {
	found := false
	if r == nil {
		t.Errorf("no results given")
	}
	for i := range r {
		if r[i].Synset == synset {
			if r[i].AStep != astep {
				t.Errorf("same synset %v found, but AStep mismatch (want %v, returned %v)", synset, astep, r[i].AStep)
			}
			if r[i].BStep != bstep {
				t.Errorf("same synset %v found, but BStep mismatch (want %v, returned %v)", synset, bstep, r[i].BStep)
			}
			if r[i].Depth != depth {
				t.Errorf("same synset %v found, but Depth mismatch (want %v, returned %v)", synset, depth, r[i].Depth)
			}
			found = true
			break
		}
	}
	if !found {
		t.Errorf("unable to find synset %v", synset)
	}
}

func TestWordnetSynset(t *testing.T) {
	f, err := os.Create("/tmp/log")
	if err != nil {
		t.Errorf("unable to open log file")
	}
	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{DisableSorting: false, QuoteEmptyFields: true})

	o := &Options{
		JumanCommand:         "jumanpp",
		KnpCommand:           "knp -tab -anaphora",
		JumanDirectory:       "/usr/local/share/juman",
		WordNetDatabase:      "../../third-party/wnjpn/wnjpn.db",
		SynonymsJapaneseOnly: true,
		UseSynsetList:        false,
		Interactive:          false,
	}
	i := &Instance{}
	var jk *JumanKnp
	jk, err = NewJumanKnp(o, i)
	i.JumanKnp = jk

	if err != nil {
		t.Errorf("cannot initialize JumanKnp, error: %v", err.Error())
		return
	}

	i.WordNet, err = NewWordNet(o, i)
	if err != nil {
		t.Errorf("constructor error: %v", err.Error())
		return
	}

	ws := NewWordNetSynset(o, i)
	var ret []SynsetResult
	ret, err = ws.NearestSynset([]rune("バナナ"), WNNounPart, []rune("りんご"), WNNounPart)
	tWordSynset(t, ret, "07705931-n", 0, 0, 7)                                     //edible_fruit, 0
	ret, err = ws.NearestSynset([]rune("猫"), WNNounPart, []rune("寝る"), WNNounPart) // ->nil
	ret, err = ws.NearestSynset([]rune("マグロ"), WNNounPart, []rune("寿司"), WNNounPart)
	tWordSynset(t, ret, "00020827-n", 4, 4, 3) // matter, 0.75
	ret, err = ws.NearestSynset([]rune("マグロ"), WNNounPart, []rune("タイ"), WNNounPart)
	tWordSynset(t, ret, "07775905-n", 0, 0, 7)  // saltwater_fish, 0
	tWordSynset(t, ret, "02554730-n", 1, 0, 15) // percoid, 0

	for i := range ret {
		t.Logf("%v: %v, %v[%v %v], %v\n",
			i, ret[i].Synset, ret[i].Depth,
			ret[i].AStep, ret[i].BStep, ret[i].Approximation)
	}
}
