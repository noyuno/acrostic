package acrostic

import (
	"fmt"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/noyuno/lgo/color"
	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

// Paragraph : 文章
type Paragraph struct {
	SentencesText    [][]rune
	Sentences        []Sentence
	Text             []rune
	Keywords         [][]rune
	Options          *Options
	Instance         *Instance
	Arrange          *Arrange
	FoundBasicPhrase bool
}

// NewParagraph : constructor
func NewParagraph(o *Options, i *Instance, t []rune, keywords [][]rune) *Paragraph {
	ret := new(Paragraph)
	ret.Options = o
	ret.Instance = i
	ret.Text = t
	ret.Keywords = keywords
	return ret
}

// GetPunctuation : 句点を返す
func GetPunctuation() []rune {
	return []rune("．。！？.!?\n")
}

// Analyze : 文章を解析する
func (p *Paragraph) Analyze() error {
	begin := 0
	array, newline := SplitSentence(p.Text, p.Options.SwapSentences)
	for i := range array {
		log.Debugf("Paragraph: %v, newline=%v", string(array[i]), newline[i])
		sentence := NewSentence(p.Options, p.Instance, array[i], newline[i], p.Keywords)
		if p.Options.KnpOnly {
			s := p.Instance.JumanKnp.Execute(sentence.Text, true)
			fmt.Println(string(s))
		} else {
			var err error
			begin, err = sentence.Analyze(begin)
			if err != nil {
				return err
			}
			if begin > 0 {
				p.FoundBasicPhrase = true
			}
			p.Sentences = append(p.Sentences, *sentence)
		}
	}
	return nil
}

// ファイル全体を改行または句点で分割
// return: 分割した文字列, 直前に改行があるかどうか
func SplitSentence(t []rune, swap bool) ([][]rune, []bool) {
	ret := make([][]rune, 0)
	newline := make([]bool, 0)
	currentnewline := swap
	foundnewline := true
	buf := []rune("")
	punc := string(GetPunctuation())
	found := false
	for i := 0; i < len(t); i++ {
		//log.Debugf("pos=%v, char='%v'", i, string(t[i]))
		if strings.Index(punc, string(t[i])) != -1 {
			//log.Debug("found punc")
			found = true
		} else {
			if found {
				// 「文．．．．」や「文！？」に対応
				if strings.Trim(string(buf), punc) != "" {
					bufc := make([]rune, len(buf))
					copy(bufc, buf)
					ret = append(ret, bufc)
					newline = append(newline, currentnewline)
					//log.Debugf("append '%v', nl=%v", string(bufc), currentnewline)
					buf = buf[:0]
					currentnewline = false
				} else {
					// do not anything
				}
			} else {
			}
			found = false
		}
		if string(t[i]) == "\n" {
			//log.Debug("newline")
			foundnewline = true
		} else {
			if foundnewline {
				//log.Debug("currentnewline=true")
				currentnewline = true
			}
			foundnewline = false
		}
		if string(t[i]) != "\n" {
			buf = append(buf, t[i])
		}
	}

	if strings.Trim(string(buf), punc) != "" {
		bufc := make([]rune, len(buf))
		copy(bufc, buf)
		ret = append(ret, bufc)
		newline = append(newline, currentnewline)
		//log.Debugf("append remained %v, nl=%v", string(bufc), currentnewline)
	} else {
		// do not anything
	}

	return ret, newline
}

func (p *Paragraph) PrintAnalyzeResult() {
	T, _ := i18n.Tfunc(p.Options.Language)
	fmt.Println(color.FBlue + T("analyze_result") + ":" + color.Reset)
	p.PrintSynonyms()
	p.PrintBPPatterns()
	if !p.Options.UseSynsetList && p.Options.Interactive {
		p.Instance.WordNet.PrintAnswer()
	}
	p.PrintCaseAnalysisPatterns()
	if p.Options.Verbose || p.Options.Verbosely {
		log.Debug(MemoryInfo())
	}
	//p.PrintNumberOfPattern()
	p.PrintProgressWarn()
}

func (p *Paragraph) PrintProgressWarn() {
	//T, _ := i18n.Tfunc(p.Options.Language)
	bplen := 0
	for i := range p.Sentences {
		bplen += len(p.Sentences[i].BasicPhrases)
	}
	//if p.Options.Progress && p.Options.ProgressDepth < bplen {
	//	log.Warnf(T("exceeds-progress-bar", map[string]interface{}{
	//		"P": p.Options.ProgressDepth,
	//		"L": bplen,
	//	}))
	//}
}

//func (p *Paragraph) PrintNumberOfPattern() {
//	anaphora := make([]uint64, len(p.Sentences))
//	synonyms := make([][]uint64, len(p.Sentences))
//	kana := make([][]uint64, len(p.Sentences))
//	total := uint64(1)
//	for i, s := range p.Sentences {
//		synonyms[i] = make([]uint64, len(s.BasicPhrases))
//		kana[i] = make([]uint64, len(s.BasicPhrases))
//		anaphora[i] = uint64(len(s.AnaphoraPatterns))
//		total *= anaphora[i]
//		for k, bp := range s.BasicPhrases {
//			synonyms[i][k-s.BasicPhraseBegin] = uint64(len(bp.Synonyms))
//			total *= synonyms[i][k-s.BasicPhraseBegin] + 1
//			if !runes.Compare(bp.AllIndependentSurface, bp.Kana) {
//				kana[i][k-s.BasicPhraseBegin] = 1
//			}
//			total *= kana[i][k-s.BasicPhraseBegin] + 1
//		}
//	}
//
//	fmt.Printf(color.FGreen+
//		"statistics: anaphora: %v, synonyms: %v, kana: %v, total: %v patterns\n"+
//		color.Reset,
//		anaphora, synonyms, kana, total)
//}

func (p *Paragraph) PrintCaseAnalysisPatterns() {
	T, _ := i18n.Tfunc(p.Options.Language)
	fmt.Println(color.FGreen + T("case analysis patterns") + color.Reset)
	for si, s := range p.Sentences {
		for rowi, row := range s.CaseAnalysisPhrases {
			o := ""
			if s.CaseAnalysisResults[rowi] {
				o += fmt.Sprintf("%2v-%2v %vok %v ", si, rowi, color.FGreen, color.Reset)
			} else {
				o += fmt.Sprintf("%2v-%2v %vbad%v ", si, rowi, color.FYellow, color.Reset)
			}
			r := ""
			for i := range row {
				r += fmt.Sprintf("%v", row[i].Number)
				if i+1 < len(row) {
					r += "->"
				}
			}
			o += fmt.Sprintf("%-26v ", r)
			for i := range row {
				o += fmt.Sprintf("%s", string(s.CaseAnalysisPhrases[rowi][i].Surface()))
			}
			fmt.Println(o)
		}
	}
}

func (p *Paragraph) PrintSynonyms() {
	T, _ := i18n.Tfunc(p.Options.Language)
	fmt.Println(color.FGreen + T("synonyms") + color.Reset)
	for _, s := range p.Sentences {
		for _, h := range s.Phrases {
			for _, bp := range h.BasicPhrases {
				a := string(bp.AllIndependentSurface)
				d := ""
				for _, syn := range bp.Synonyms {
					if syn.HasInflection {
						d += string(syn.InflectionSurface) // + "[" + string(syn.Surface) + "]"
						if syn.HasPolite {
							d += "[" + string(syn.PoliteSurface) + "]"
						}
					} else {
						d += string(syn.Surface)
					}
					if p.Options.PrintKana && syn.HasKana {
						d += "(" + string(syn.Kana) + ")"
					}
					d += " "
				}
				fmt.Printf("%v: %v\n", a, d)
				//max[i] = len(bp.Synonyms)
				//pnum *= max[i] + 1
				//i++
			}
		}
	}
	//pat := len(pattern)
	//fmt.Printf("max: %v, (%v synsets * %v patterns) = %v\n",
	//	max, pnum, patterns, pnum*patterns)
}

func (p *Paragraph) PrintBPPatterns() {
	T, _ := i18n.Tfunc(p.Options.Language)
	fmt.Println(color.FGreen + T("Phrase patterns") + color.Reset)
	for _, s := range p.Sentences {
		for _, h := range s.Phrases {
			for _, bp := range h.BasicPhrases {
				a := string(bp.AllIndependentSurface)
				d := ""
				for _, a := range bp.Pattern {
					d += string(a) + " "
				}
				fmt.Printf("%v: %v\n", a, d)
			}
		}
	}
}

func (p *Paragraph) PrintPatterns() {
	fmt.Println(color.FGreen + "patterns" + color.Reset)
	for _, s := range p.Sentences {
		for _, h := range s.Phrases {
			for _, bp := range h.BasicPhrases {
				d := ""
				for _, para := range bp.Pattern {
					d += string(para) + " "
				}
				fmt.Printf("%v: %v\n", string(bp.Surface), d)
			}
		}
	}
}

func (p *Paragraph) CheckContainsKeyword(keyword []rune, keywordi int) bool {
	T, _ := i18n.Tfunc(p.Options.Language)
	//log.Debugf("Paragraph.CheckContainsKeyword")
	for _, k := range keyword {
		kk := []rune(string(k))
		//log.Debugf("k=%v", string(k))
		found := false
		for _, s := range p.Sentences {
			for _, p := range s.Phrases {
				for _, b := range p.BasicPhrases {
					for _, a := range b.Pattern {
						if runes.Index(a, kk, 0) != -1 {
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			// キーワード中少なくとも1文字が見つからなかった
			log.Warnf(T("keyword-not-found", map[string]interface{}{
				"C":    string(k),
				"K":    string(keyword),
				"Kpos": keywordi,
			}))
			return false
		}
	}
	return true
}

func (p *Paragraph) Generate(k []rune, n int, width int) (bool, error) {
	r := false
	arrange, err := NewArrange(p.Options, p.Instance, p.Sentences, p.Text, n, k, width)
	if err != nil {
		return false, err
	}
	r, err = arrange.Arrange()
	if err != nil {
		return false, err
	}
	arrange.Output()
	return r, nil
}
