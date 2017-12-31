package acrostic

import (
	"errors"
	"reflect"
	"strings"

	"github.com/noyuno/lgo/algo"
	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

// Sentence : 文
type Sentence struct {
	// Phrases : 文節
	Phrases []Phrase

	// PhraseLength : 文節の数
	//PhraseLength int

	// BasicPhrases : 基本句
	BasicPhrases map[int]BasicPhrase

	// 基本句番号の初め
	BasicPhraseBegin int

	// Text : テキスト
	Text []rune

	Keywords [][]rune

	// NewLine : 直前に改行があるかどうか
	NewLine bool

	// Options : オプション
	Options *Options

	// Instance : Instance
	Instance *Instance

	// CaseAnalysisPatterns : 格解析で通過したパターン
	CaseAnalysisPatterns [][]int

	// すべての格解析結果
	CaseAnalysisResults []bool

	Pattern *Pattern

	ParallelIndex   map[int][]int
	ParallelPhrases [][][]*Phrase

	// CaseAnalysisPhrases : 格解析により入れ替えたものの助詞を調整したもの
	CaseAnalysisPhrases [][]Phrase

	ZenkakuSpaceFirst bool
}

// NewSentence : constructor
func NewSentence(o *Options, i *Instance, t []rune, newline bool, keywords [][]rune) *Sentence {
	ret := new(Sentence)
	ret.Options = o
	ret.Instance = i
	ret.Text = t
	ret.ParallelIndex = map[int][]int{}
	ret.BasicPhrases = map[int]BasicPhrase{}
	ret.NewLine = newline
	ret.Keywords = keywords
	return ret
}

// Analyze : 解析する
func (s *Sentence) Analyze(begin int) (int, error) {
	if s.Options.Mode != "" {
		return 0, errors.New("not implemented")
	}
	return s.AnalyzeJumanKnp(begin)
}

func (s *Sentence) Split(v []rune) [][][]rune {
	lfToken := []rune("\n")
	phraseToken := []rune("*")
	eosToken := []rune("EOS")
	t := runes.Split(v, lfToken)

	ret := make([][][]rune, 0)
	phrase := make([][]rune, 0)
	enable := false
	for i := 0; i < len(t); i++ {
		if t[i][0] == phraseToken[0] {
			if len(phrase) > 0 {
				phrasec := make([][]rune, len(phrase))
				for i := range phrase {
					phrasec[i] = make([]rune, len(phrase[i]))
					copy(phrasec[i], phrase[i])
				}
				phrase = phrase[:0]
				ret = append(ret, phrasec)
			}
			enable = true
		}
		if enable {
			if !runes.Compare(eosToken, t[i][0:3]) {
				phrase = append(phrase, t[i])
			}
		}
	}
	if len(phrase) > 0 {
		phrasec := make([][]rune, len(phrase))
		for i := range phrase {
			phrasec[i] = make([]rune, len(phrase[i]))
			copy(phrasec[i], phrase[i])
		}
		ret = append(ret, phrasec)
	}
	return ret
}

// AnalyzeJumanKnp : Juman and KNP で解析する
func (s *Sentence) AnalyzeJumanKnp(begin int) (int, error) {
	//log.Debugf("Sentence.AnalyzeJumanKnp: %v;", string(s.Text))
	out := s.Instance.JumanKnp.Execute(s.Text, true)
	array := s.Split(out)
	for i := range array {
		newline := i == 0 && s.NewLine
		phrase := NewPhrase(s.Options, s.Instance, array[i], len(s.Phrases), newline, s.Keywords)
		var err error
		s.BasicPhraseBegin = begin
		if begin, err = phrase.Analyze(begin); err != nil {
			return 0, err
		}
		s.Phrases = append(s.Phrases, *phrase)
		for i := range phrase.BasicPhrases {
			s.BasicPhrases[s.BasicPhraseBegin+i] = phrase.BasicPhrases[i]
		}
		if phrase.Parallel {
			/*
				並列する語であれば，出力"* NP"のNをキーとして，
				phrase.Numberを値としてマップに登録する．
				こうすることで，並列する語はその最後の語が出現したとき
				受理するようにコーディングすることができるようになるだろうね
				P
				└P
				 └D
				KNPの係先番号は，上のように，最後の並列する語ではなく，
				次の並列する語に向かってつけられているので，
				このままではこれら並列する語を同格とみなしてシャッフル
				できないため，並列する語の係先番号を最後の並列する語に
				向かうように付け替えなければならない．
				そのため，すでに，phraseに向いているParallelIndexがあれば，
				あらかじめそれを削除して，phrase.Destinationに付け替える．
			*/
			//log.Debugf("phrase %v is parallel destination to %v",
			//	phrase.Number, phrase.Destination)
			removedpp := []int{}
			if _, ok := s.ParallelIndex[phrase.Number]; ok {
				// 外す
				log.Debugf("%v is already exists in ParallelIndex, replace it to %v",
					phrase.Number, phrase.Destination)
				for _, v := range s.ParallelIndex[phrase.Number] {
					removedpp = append(removedpp, v)
				}
				delete(s.ParallelIndex, phrase.Number)
			}
			// 付ける
			if _, ok := s.ParallelIndex[phrase.Destination]; !ok {
				s.ParallelIndex[phrase.Destination] = []int{}
			}
			for i := range removedpp {
				s.ParallelIndex[phrase.Destination] =
					append(s.ParallelIndex[phrase.Destination], removedpp[i])
			}
			// 削除
			removedpp = removedpp[:0]
			s.ParallelIndex[phrase.Destination] =
				append(s.ParallelIndex[phrase.Destination], phrase.Number)
			//log.Debugf("current ParallelIndex[%v] = %v",
			//	phrase.Destination, s.ParallelIndex[phrase.Destination])
		}

	}
	// debug
	//for i := range s.BasicPhrases {
	//	log.Debugf("%v: %v", i, string(s.BasicPhrases[i].Surface))
	//}

	mat := make([][]bool, len(s.Phrases))
	for i := range mat {
		mat[i] = make([]bool, len(s.Phrases))
	}
	init := make([]int, len(s.Phrases))
	for i, ph := range s.Phrases {
		if ph.Destination == -1 {
			// 終端なので，宛先は存在しない
		} else {
			if ph.Parallel == false {
				mat[i][ph.Destination] = true
			} else {
				// do nothing
			}
		}
		init[i] = ph.Destination
		// 全角スペース
		if strings.Compare(string(ph.Surface()), "　") == 0 && i == 0 {
			if len(mat[i]) > 0 {
				// 先頭
				//for m := 0; m < len(mat[0]) && m < 3; m++ {
				//	mat[0][m] = true
				//}
				mat[0][1] = true
				init[0] = 1
				s.ZenkakuSpaceFirst = true
			}
		}
	}

	subindex := make([][]int, len(s.ParallelIndex))
	subindexi := 0
	for k, _ := range s.ParallelIndex {
		subindex[subindexi] = make([]int, len(s.ParallelIndex[k])+1)
		for c := range s.ParallelIndex[k] {
			log.WithFields(log.Fields{
				"subindexi": subindexi,
				"c":         c,
				"k":         k,
				"[k][c]":    s.ParallelIndex[k][c],
			}).Debug()
			subindex[subindexi][c] = s.ParallelIndex[k][c]
		}
		subindex[subindexi][len(subindex[subindexi])-1] = k
		subindexi++
	}
	// debug
	log.Debug("mat")
	for _, row := range mat {
		matout := ""
		for _, i := range row {
			if i {
				matout += "1 "
			} else {
				matout += "0 "
			}
		}
		log.Debug(matout)
	}
	log.Debug("subindex")
	for mi, m := range subindex {
		log.Debugf("%v: %v", mi, m)
	}
	s.Pattern = NewPattern(s.Options, len(s.Phrases), init, mat, subindex)
	err := s.Pattern.Shuffle()
	if err != nil {
		return 0, err
	}

	log.Debugf("%v patterns found", len(s.Pattern.Orders))
	log.Debugf("Sentence: change particle part order and append to CaseAnalysisPhrases")

	s.MakeParallelPhrase()
	s.ParticleOrder()

	caseAnalysisInvalidPatterns := make([][]int, 0, 10)
	log.Debug("case analysis")
	for _, p := range s.Phrases {
		for _, bp := range p.BasicPhrases {
			//log.WithFields(log.Fields{
			//	"surface":   string(bp.Surface),
			//	"available": bp.Anaphora.IsAvailable}).Debug()
			if bp.PredicateTerm.IsAvailable {
				for _, ace := range bp.PredicateTerm.AnaphoraCaseElementGroups {
					//log.Debugf("Sentence: %v (%v)", ace.EntityID, len(s.BasicPhrases))
					f := s.BasicPhrases[ace.EntityID].PhraseNumber
					//log.Debugf("Sentence: %v (%v)", bp.Anaphora.EntityID, len(s.BasicPhrases))
					t := s.BasicPhrases[bp.PredicateTerm.EntityID].PhraseNumber
					log.Debugf("%v(%v) -> %v(%v)",
						string(s.Phrases[f].Surface()), f,
						string(s.Phrases[t].Surface()), t)
					caseAnalysisInvalidPatterns = append(caseAnalysisInvalidPatterns, []int{t, f})
				}
			}
		}
	}
	s.CaseAnalysisResults = make([]bool, len(s.Pattern.Orders))
	for rowi, row := range s.Pattern.Orders {
		if algo.ValidateOrder(row, caseAnalysisInvalidPatterns) {
			log.Debugf("%v", row)
			if len(row) > 0 && s.ZenkakuSpaceFirst {
				if row[0] != 0 {
					log.Debugf("row start with zenkaku space does not start at 0, actual: %v", row[0])
					continue
				}
			}
			s.CaseAnalysisPatterns = append(s.CaseAnalysisPatterns, row)
			s.CaseAnalysisResults[rowi] = true
		}
	}

	if len(s.BasicPhrases) > 0 {
		return begin, nil
	}
	return 0, nil
}

func (s *Sentence) MakeParallelPhrase() {
	log.Debugf("MakeParallelPhrase")
	s.ParallelPhrases = make([][][]*Phrase, len(s.Pattern.SubPatterns))
	for pati, pat := range s.Pattern.SubPatterns {
		s.ParallelPhrases[pati] = make([][]*Phrase, len(pat))
		for pi, p := range pat {
			s.ParallelPhrases[pati][pi] = make([]*Phrase, len(p))
			//log.Debugf("pi=%v, p=%v", pi, p)
			for i, v := range p {
				// copy
				//log.Debugf("i=%v, v=%v", i, v)
				orig := s.Phrases[s.Pattern.SubIndex[pati][i]]
				//log.Debugf("s.Pattern.SubIndex[%v][%v]=%v, orig: %v",
				//	pati, i, s.Pattern.SubIndex[pati][i], string(orig.Surface()))
				t := s.Phrases[v].Copy()
				for bpi, _ := range t.BasicPhrases {
					//log.Debugf("t=%v len=%v, orig=%v len=%v",
					//	string(t.Surface()), len(t.BasicPhrases), string(orig.Surface()), len(orig.BasicPhrases))
					if len(orig.BasicPhrases) <= bpi {
						break
					}
					if t.BasicPhrases[bpi].HasParticle &&
						orig.BasicPhrases[bpi].HasParticle {
						// Order内で助詞の数が違うとout of rangeになるので
						// ParticleSurface に合わせる
						// 追記：数が違うのは，KNPが並立でないのに誤って判定している
						// 例：「帝京大と2位までの」
						tcount := 0
						for _, so := range t.BasicPhrases[bpi].SurfaceOrder {
							if so == ParticlePart {
								tcount++
							}
						}
						origcount := 0
						for _, so := range orig.BasicPhrases[bpi].SurfaceOrder {
							if so == ParticlePart {
								origcount++
							}
						}
						if origcount == tcount {
							t.BasicPhrases[bpi].ParticleSurface =
								runes.CopyArray(orig.BasicPhrases[bpi].ParticleSurface)
							t.BasicPhrases[bpi].AllParticleSurface =
								runes.Copy(orig.BasicPhrases[bpi].AllParticleSurface)
							t.BasicPhrases[bpi].UpdateSurface()
							t.BasicPhrases[bpi].UpdatePattern()
						}
						//if tcount > origcount {
						//	// 削除
						//	neworder := make([]Part, 0)
						//	tc := 0
						//	for _, so := range t.BasicPhrases[bpi].SurfaceOrder {
						//		if so == ParticlePart {
						//			if tc < origcount {
						//				neworder = append(neworder, so)
						//			}
						//			tc++
						//		} else {
						//			neworder = append(neworder, so)
						//		}
						//	}
						//	t.BasicPhrases[bpi].SurfaceOrder = neworder
						//} else if tcount < origcount {
						//	// 追加
						//	for i := tcount; tcount < origcount; i++ {
						//		t.BasicPhrases[bpi].SurfaceOrder =
						//			append(t.BasicPhrases[bpi].SurfaceOrder, ParticlePart)
						//	}
						//}
					}
				}
				s.ParallelPhrases[pati][pi][i] = t
			}
		}
	}
	//log.Debugf("created ParallelPhrases(%v)", len(s.ParallelPhrases))
	//for i := range s.ParallelPhrases {
	//	for k := range s.ParallelPhrases[i] {
	//		o := []rune("")
	//		for m := range s.ParallelPhrases[i][k] {
	//			o = append(o, s.ParallelPhrases[i][k][m].Surface()...)
	//		}
	//		log.Debugf("%v: %v", i, string(o))
	//	}
	//}
}

func (s *Sentence) ParticleOrder() {
	log.Debugf("ParticleOrder")
	s.CaseAnalysisPhrases = make([][]Phrase, len(s.Pattern.Orders))
	for o := range s.Pattern.Orders {
		//log.Debugf("o[%v]=%v", o, s.Pattern.Orders[o])
		s.CaseAnalysisPhrases[o] = make([]Phrase, len(s.Pattern.Orders[o]))
		// 並立する語の順番．0 1 4 2 3 5 6 7 -> [[4 2 3] [6 7]]
		found := make([][]int, len(s.Pattern.SubIndex))
		for p := range s.CaseAnalysisPhrases[o] {
			order := s.Pattern.Orders[o][p]

			// 並列チェック & 巡回チェック
			for i := range s.Pattern.SubIndex {
				if algo.Contains(s.Pattern.SubIndex[i], order) {
					// 何番目？
					found[i] = append(found[i], order)
				}
			}
		}

		//log.Debugf("found = %v", found)

		// SubPatternsの場所特定．SubPatterns[0]の[4 2 3], [6 7]がある場所
		// 4: [0 4 0]
		// 2: [0 4 1]
		// 3: [0 4 2]
		// 6: [1 2 0]
		// 7: [1 2 1]
		ppi := map[int][]int{}
		for i := range s.Pattern.SubPatterns {
			for k := range s.Pattern.SubPatterns[i] {
				if reflect.DeepEqual(found[i], s.Pattern.SubPatterns[i][k]) {
					// 一致
					for f := range found[i] {
						ppi[found[i][f]] = []int{i, k, f}
					}
					break
				}
			}
		}

		//log.Debugf("ppi = %v", ppi)

		// 作成
		for p := range s.CaseAnalysisPhrases[o] {
			order := s.Pattern.Orders[o][p]
			f := false
			for k := range ppi {
				if k == order {
					// 並立する
					//log.Debugf("found at %v", ppi[k])
					//log.Debugf("len of ParallelPhrases: %v", len(s.ParallelPhrases))
					f = true
					s.CaseAnalysisPhrases[o][p] =
						*s.ParallelPhrases[ppi[k][0]][ppi[k][1]][ppi[k][2]]
					break
				}
			}
			if f == false {
				// 並立しない
				//log.Debugf("not found o=%v, p=%v", p, o)
				s.CaseAnalysisPhrases[o][p] = s.Phrases[order]
			}
		}
	}
}
