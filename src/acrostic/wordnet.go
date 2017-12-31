package acrostic

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/noyuno/lgo/color"
	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

type WordNetAnswer int

const (
	WNANone = iota
	WNAAll
	WNASynonyms
)

func (w WordNetAnswer) String() string {
	switch w {
	case WNANone:
		return ""
	case WNAAll:
		return "all"
	case WNASynonyms:
		return "synonyms"
	}
	log.Fatalf("WordNetAnswer: unknown value")
	return ""
}

func NewWordNetAnswer(s string) WordNetAnswer {
	switch s {
	case "":
		return WNANone
	case "all":
		return WNAAll
	case "synonyms":
		return WNASynonyms
	}
	return WNANone
}

type WordNet struct {
	DB       *sql.DB
	Instance *Instance
	Options  *Options
	Answer   map[int]map[string]WordNetAnswer
}

type WordNetResult struct {
	Surface           []rune
	Link              WordNetLink
	Language          string
	Part              WordNetPart
	Kana              []rune
	HasKana           bool
	HasInflection     bool
	InflectionSurface []rune
	InflectionType    []rune
	InflectionForm    []rune
	HasPolite         bool
	PoliteSurface     []rune
	HasPoliteKana     bool
	PoliteKana        []rune
}

func (w *WordNetResult) Copy() *WordNetResult {
	r := new(WordNetResult)
	r.Surface = runes.Copy(w.Surface)
	r.Link = w.Link
	r.Language = w.Language
	r.Part = w.Part
	r.Kana = runes.Copy(w.Kana)
	r.HasKana = w.HasKana
	r.HasInflection = w.HasInflection
	r.InflectionSurface = runes.Copy(w.InflectionSurface)
	r.InflectionType = runes.Copy(w.InflectionType)
	r.InflectionForm = runes.Copy(w.InflectionForm)
	return r
}

func NewWordNet(o *Options, i *Instance) (*WordNet, error) {
	var err error
	ret := new(WordNet)
	ret.Options = o
	ret.Instance = i

	if o.WordNetDatabase == "" {
		return nil, errors.New("require WordNetDatabase as WordNet database filename")
	}
	ret.DB, err = sql.Open("sqlite3", ret.Options.WordNetDatabase)
	if err != nil {
		return nil, err
	}
	if err = ret.DB.Ping(); err != nil {
		log.Fatalf("ping failure to sqlite(%v): %v", ret.Options.WordNetDatabase, err.Error())
	}

	if o.JumanDirectory == "" {
		return nil, errors.New("require JumanDirectory")
	}
	ret.Answer = map[int]map[string]WordNetAnswer{}
	return ret, nil
}

// getSynset : 指定された単語の指定された関係である単語を取得する．
// bphrase: 基本句
// link: 基本句に対するリンク
// synset: グループ記号
func (w *WordNet) getSynset(
	bphrase []rune, bppart WordNetPart, link WordNetLink, synset string) ([]WordNetResult, error) {
	var ret []WordNetResult
	var lang string
	var rows *sql.Rows
	var err error

	if w.Options.SynonymsJapaneseOnly {
		lang = "and word.lang='jpn'"
	} else {
		lang = ""
	}
	if link == WNSynonym {
		rows, err = w.DB.Query(`select lemma,word.lang,pos from sense, word
			where synset=? and sense.wordid=word.wordid `+lang, synset)
	} else {
		rows, err = w.DB.Query(`select lemma,word.lang,pos from synlink, sense, word
			where link=? and synset1=? and synset2=synset
			and sense.wordid=word.wordid `+lang, link.DBString(), synset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		var l string
		var p string
		if err = rows.Scan(&s, &l, &p); err != nil {
			return nil, err
		}
		// check part
		//part := NewWordNetPart(p)
		//log.WithFields(log.Fields{"p": part.String(), "bppart": bppart.String()}).Debug()
		//if bppart != part {
		//	log.WithFields(log.Fields{"BPSurface": string(bphrase),
		//		"BPPart":  bppart.String(),
		//		"Link":    link.DBString(),
		//		"synset":  synset,
		//		"Surface": string(s),
		//		"Part":    part.String()}).Warning("mismatch part")
		//	continue
		//}
		ri := WordNetResult{
			Surface:  []rune(s),
			Link:     link,
			Language: l,
			Part:     NewWordNetPart(p)}
		//if !runes.Compare(inflectionForm, []rune("")) {
		//	ri.InflectionType, ri.InflectionSurface, err =
		//		w.Instance.JumanKnp.GetInflection(ri.Surface, inflectionForm)
		//	if err != nil {
		//		return nil, err
		//	}
		//}

		ret = append(ret, ri)
	}

	return ret, nil
}

func (w *WordNet) GetSynonyms(bp *BasicPhrase, link []WordNetLink) ([]WordNetResult, error) {
	T, _ := i18n.Tfunc(w.Options.Language)
	part := ToWordNetPart(bp.Part)

	rows, err := w.DB.Query("select wordid,pos from word where lemma=?",
		string(bp.Origin))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []WordNetResult
	var ids []int
	for rows.Next() {
		var (
			wordid int
			pos    string
		)
		if err = rows.Scan(&wordid, &pos); err != nil {
			return nil, err
		}
		if part.String() == pos {
			ids = append(ids, wordid)
		}
	}
	//log.WithFields(log.Fields{"ids": ids, "bp.Surface": string(bp.Surface)}).Debug("WordNet: step 1 ok")
	if len(ids) == 0 {
		log.Debugf("unable to get wordid: %v", string(bp.Origin))
	}
	for _, id := range ids {
		rows, err = w.DB.Query("select synset from sense where wordid=?", strconv.Itoa(id))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var synset string
			if err = rows.Scan(&synset); err != nil {
				return nil, err
			}
			//log.WithFields(log.Fields{"synset": synset}).Debug("WordNet: step 2 ok")

			var answer WordNetAnswer
			answer = WNANone
			if w.Options.UseSynsetList {
				if v, ok := w.Options.SynsetList[bp.ID]; ok {
					for k, vv := range v {
						if synset == k {
							answer = vv
							log.Debugf("GetSynonyms selected %v by SynsetList[%v] = %v", vv.String(), bp.ID, synset)
							break
						}
					}
					if answer == WNANone {
						log.Debugf("GetSynonyms skipped by SynsetList[%v] = %v", bp.ID, synset)
						continue
					}
				}

			}

			h := make([]WordNetResult, 0)
			for _, s := range link {
				if answer == WNASynonyms && s != WNSynonym {
					// synonymsだけ
					continue
				}
				var g []WordNetResult
				g, err = w.getSynset(bp.Origin, part, s, synset)
				if err != nil {
					return nil, err
				}
				for gi := range g {
					if part == WNVerbPart || part == WNAdverbPart {
						//g[gi].InflectionType = inflectionType
						// Surfaceの末尾がひらがなでかつう段でなければ，
						// 名詞と判断し，末尾に「する」を付けて動詞っぽくする
						last := string(g[gi].Surface[len(g[gi].Surface)-1])
						if !strings.Contains("うくすつぬふむゆる", last) {
							g[gi].Surface = append(g[gi].Surface, []rune("する")...)
						}
						g[gi].InflectionForm = bp.InflectionForm
						g[gi].InflectionType, g[gi].InflectionSurface,
							g[gi].HasInflection =
							w.Instance.JumanKnp.Inflection(
								g[gi].Surface, g[gi].InflectionForm)
						if w.Options.UsePolite {
							g[gi].PoliteSurface, g[gi].HasPolite =
								w.Instance.JumanKnp.InflectionPolite(
									g[gi].Surface, g[gi].InflectionForm)
						}
						if g[gi].HasInflection {
							g[gi].Kana, g[gi].HasKana =
								w.Instance.Kana.Get(g[gi].InflectionSurface)
							if g[gi].HasPolite {
								g[gi].PoliteKana, g[gi].HasPoliteKana =
									w.Instance.Kana.Get(g[gi].PoliteSurface)
							}
							h = append(h, g[gi])
						} else {
							log.WithFields(log.Fields{
								"Surface":        string(g[gi].Surface),
								"InflectionForm": string(g[gi].InflectionForm),
								"InflectionType": string(g[gi].InflectionType)}).Info(
								"failed to get inflection")
						}
					} else {
						g[gi].Kana, g[gi].HasKana =
							w.Instance.Kana.Get(g[gi].Surface)
						h = append(h, g[gi])
					}
				}
			}
			if (answer == WNANone) && w.Options.Interactive {
				var a string
				fmt.Printf(T("number")+": %v, "+T("origin")+": "+color.FGreen+"%v"+
					color.Reset+", "+T("synset")+": %v\n",
					bp.ID, string(bp.Origin), synset)
				outmap := map[string][]string{}
				for _, i := range h {
					dbs := i.Link.DBString()
					if _, ok := outmap[i.Link.DBString()]; !ok {
						outmap[dbs] = make([]string, 0)
					}
					if i.HasInflection {
						outmap[dbs] = append(outmap[dbs], string(i.InflectionSurface))
					} else {
						outmap[dbs] = append(outmap[dbs], string(i.Surface))
					}
				}
				for k, v := range outmap {
					fmt.Printf("%v: %v\n", k, v)
				}
				fmt.Printf("\n"+T("select above %v items?")+" [Yns]:", len(h))
				fmt.Scanln(&a)
				if MaybeYes(a) {
					ret = append(ret, h...)
					if _, o := w.Answer[bp.ID]; !o {
						w.Answer[bp.ID] = map[string]WordNetAnswer{}
					}
					w.Answer[bp.ID][synset] = WNAAll
				} else if a == "s" {
					// 類義語のみ
					for _, i := range h {
						if i.Link == WNSynonym {
							ret = append(ret, i)
						}
					}
					if _, o := w.Answer[bp.ID]; !o {
						w.Answer[bp.ID] = map[string]WordNetAnswer{}
					}
					w.Answer[bp.ID][synset] = WNASynonyms
				}
			} else {
				ret = append(ret, h...)
			}
		}
	}

	//log.WithFields(log.Fields{"len(ret)": len(ret)}).Debug("WordNet: step 3 ok")
	//log.WithFields(log.Fields{"ID": bp.ID, "len(ret)": len(ret)}).Debug("Wordnet: GetSynonyms finished")
	return ret, nil
}

func (w *WordNet) PrintAnswer() {
	T, _ := i18n.Tfunc(w.Options.Language)
	out := "\n"
	for k, v := range w.Answer {
		out += strconv.Itoa(k) + ":"
		a := []string{}
		for kk, vv := range v {
			a = append(a, kk+"="+vv.String())
		}
		out += strings.Join(a, ",") + ";"
	}
	fmt.Println(color.FGreen + T("synonyms-description") + out + color.Reset)
	//fmt.Println(color.FGreen + `If you specify the following text as argument
	//'-symonyms-string TEXT' or '-s FILE', you can omit next inputs.
	//` + out + color.Reset)
	fmt.Printf(T("Do you want to save it?") + " [Yn]:")
	a := ""
	fmt.Scanln(&a)
	if MaybeYes(a) {
		// 名前をつけて保存
		ext := filepath.Ext(w.Options.TextFileName)
		filename := w.Options.TextFileName[:len(w.Options.TextFileName)-len(ext)] + ".synonyms"
		fmt.Printf(T("filename [default %v]")+": ", filename)
		fmt.Scanln(&a)
		if a != "" {
			filename = a
		}
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf(T("unable to create file")+": %v", filename)
		}
		defer f.Close()
		f.WriteString(out)
		fmt.Printf(T("saved") + "\n")
		TypeToContinue()
	}

	//TypeToContinue()
}

func GetWordNetLinkList() []WordNetLink {
	ret := make([]WordNetLink, int(WNEnd))
	for i := 0; i < int(WNEnd); i++ {
		ret[i] = WordNetLink(i)
	}
	return ret
}

func (l WordNetLink) DBString() string {
	return strings.ToLower(l.String()[2:])
}
