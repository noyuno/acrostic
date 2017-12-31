package acrostic

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/noyuno/lgo/algo"
	"github.com/noyuno/lgo/color"
	log "github.com/sirupsen/logrus"
)

type Arrange struct {
	Options     *Options
	Instance    *Instance
	Sentences   []Sentence
	Text        []rune
	Keyword     []rune
	Number      int
	Results     [][]ArrangeMatrixResult
	Surfaces    [][]rune
	Width       int
	Color       bool
	Count       []int
	Writer      *ArrangeWriter
	WipedLength []int
}

type BasicPhraseArrange struct {
	*BasicPhrase
	UseSynonyms    bool
	SynonymsID     int
	ArrangeSurface []rune
	ArrangeKana    bool
}

func NewArrange(o *Options, i *Instance, s []Sentence, t []rune, kn int, k []rune, width int) (*Arrange, error) {
	ret := new(Arrange)
	ret.Options = o
	ret.Instance = i
	ret.Sentences = s
	ret.Text = t
	ret.Number = kn
	ret.Keyword = k
	ret.Width = width
	ret.Surfaces = make([][]rune, 0)
	ret.Writer = NewArrangeWriter(ret.Options, ret.Number, ret.Keyword)
	if ret.Number == 0 {
		// ファイルを初期化
		err := ret.Writer.Truncate()
		if err != nil {
			return nil, err
		}
	}
	ret.Writer.OutputKeyword()
	return ret, nil
}

func factorial(x int64) int64 {
	if x == 0 {
		return 1
	}
	return x * factorial(x-1)
}

func sentencePatternLength(array [][][]BasicPhrase, swap bool) int64 {
	ret := int64(1)
	if swap {
		ret = factorial(int64(len(array)))
	}
	for i := range array {
		ret *= int64(len(array[i]))
	}
	return ret
}

//func bpPatternLength(array [][][]BasicPhrase, swap bool, lim int) int64 {
//	ret := int64(1)
//	for i := range array {
//		for k := range array[i][0] {
//			if k >= lim {
//				break
//			}
//			o := len(array[i][0][k].Pattern)
//			if o != 0 {
//				ret *= int64(o)
//				log.Debugf("k=%v o=%v, ret=%v", k, o, ret)
//			}
//		}
//	}
//	return ret
//}

func (a *Arrange) makeSentences() [][][]BasicPhrase {
	begin := 0
	sentences := make([][][]BasicPhrase, len(a.Sentences))
	for si, s := range a.Sentences {
		sentences[si] = make([][]BasicPhrase, 0)
		found := false
		for pi, p := range s.CaseAnalysisPhrases {
			if s.CaseAnalysisResults[pi] == false {
				continue
			}
			t := make([]BasicPhrase, 0)
			for _, phrase := range p {
				t = append(t, phrase.BasicPhrases...)
			}
			//a.Array = append(a.Array, t)
			match := true
			if found == false {
				// 入力文の順序か？
				for i := range t {
					if t[i].ID != i+begin {
						match = false
						break
					}
				}
				if match {
					found = true
				}
			} else {
				match = false
			}

			if match {
				//log.Debugf("matched in %v", si)
				// 入力文は頭に挿入したほうが見栄えが良い
				sentences[si] = append(sentences[si], []BasicPhrase{})
				for i := len(sentences[si]) - 1; i >= 1; i-- {
					sentences[si][i] = sentences[si][i-1]
				}
				sentences[si][0] = t
			} else {
				sentences[si] = append(sentences[si], t)
			}
		}

		// 入力文の順序が格解析結果にないときは，入力文の順序を挿入する
		if found == false {
			//log.Debugf("input sentence order not found in anaphora")
			bps := make([]BasicPhrase, len(s.BasicPhrases))
			keys := []int{}
			for k := range s.BasicPhrases {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			bpsi := 0
			for _, k := range keys {
				//log.Debugf("len(bps): %v, len(s.BasicPhrase): %v, bpsi: %v, k: %v", len(bps), len(s.BasicPhrases), bpsi, k)
				bps[bpsi] = s.BasicPhrases[k]
				bpsi++
			}
			//a.Array = append(a.Array, bps)
			//sentences[si] = append(sentences[si], bps)
			sentences[si] = append(sentences[si], nil)
			for i := len(sentences[si]) - 1; i >= 1; i-- {
				sentences[si][i] = sentences[si][i-1]
			}
			sentences[si][0] = bps
		}
		begin += len(s.BasicPhrases)
	}
	return sentences
}

func bpPatternLength(sentences []Sentence, swap bool, lim int) int64 {
	ret := int64(1)
	for i := range sentences {
		for _, v := range sentences[i].BasicPhrases {
			if len(v.Pattern) != 0 {
				ret *= int64(len(v.Pattern))
				log.Debugf("len=%v", len(v.Pattern))
			}
		}
	}
	return ret
}

func sentencePattern(array [][][]BasicPhrase, swap bool) chan []BasicPhrase {
	c := make(chan []BasicPhrase)
	go func(c chan []BasicPhrase) {
		defer close(c)

		if swap {
			// 入れ替え（レシート）
			perm := make([]int, len(array))
			for i := range perm {
				perm[i] = i
			}
			for e := range algo.Permutations(perm) {
				comb := make([][]int, len(array))
				for i := range array {
					comb[i] = make([]int, len(array[i]))
					for k := range array[i] {
						comb[i][k] = k
					}
				}
				for v := range algo.Combinations(comb) {
					bp := make([]BasicPhrase, 0)
					for i := range v {
						newline := false
						for bpi := range array[e[i]][v[i]] {
							if array[e[i]][v[i]][bpi].NewLine {
								newline = true
								array[e[i]][v[i]][bpi].NewLine = false
							}
						}
						if newline {
							array[e[i]][v[i]][0].NewLine = true
						}
						bp = append(bp, array[e[i]][v[i]]...)
					}
					c <- bp
				}
			}
		} else {
			comb := make([][]int, len(array))
			for i := range array {
				comb[i] = make([]int, len(array[i]))
				for k := range array[i] {
					comb[i][k] = k
				}
			}
			for v := range algo.Combinations(comb) {
				bp := make([]BasicPhrase, 0)
				for i := range v {
					newline := false
					for bpi := range array[i][v[i]] {
						if array[i][v[i]][bpi].NewLine {
							newline = true
							array[i][v[i]][bpi].NewLine = false
						}
					}
					if newline {
						array[i][v[i]][0].NewLine = true
					}
					bp = append(bp, array[i][v[i]]...)
				}
				c <- bp
			}
		}
	}(c)
	return c
}

func (a *Arrange) Arrange() (bool, error) {
	//T, _ := i18n.Tfunc(a.Options.Language)

	// 下準備
	sentences := a.makeSentences()
	// 実行
	// OutputEachPattern のときは，結果が帰ってきたらすぐに出力する
	a.Results = make([][]ArrangeMatrixResult, 0)
	a.Count = make([]int, 0)
	a.WipedLength = make([]int, 0)

	bpai := 0
	for bpa := range sentencePattern(sentences, a.Options.SwapSentences) {
		if a.Options.OutputEachPattern {
			a.WipedLength = append(a.WipedLength, 0)
			a.Count = append(a.Count, 0)
			mret, err := a.ArrangePattern(bpai, bpa)
			if err != nil {
				return false, err
			}
			surface := []rune("")
			for i := range bpa {
				if bpa[i].NewLine && i != 0 {
					surface = append(surface, []rune("\\n")...)
				}
				surface = append(surface, bpa[i].Surface...)
			}
			a.Surfaces = append(a.Surfaces, surface)
			_, err = a.Writer.OutputPattern(
				a.Keyword, surface, bpai, mret, true, a.WipedLength[bpai], a.Width)
			if err != nil {
				return false, err
			}
			a.Count[bpai] = len(mret)
			if a.Options.One && len(mret) > 0 {
				break
			}

			//log.Debugf("Count=%v WipedLength=%v", a.Count, a.WipedLength)
			//if a.Count[bpai] == 0 && a.WipedLength[bpai] == 0 {
			//	fmt.Print("no pattern found\n")
			//} else {
			//	fmt.Printf("found %v patterns\n", a.Count[bpai]+a.WipedLength[bpai])
			//}
		} else {
			mret, err := a.ArrangePattern(bpai, bpa)
			if err != nil {
				return false, err
			}
			a.Results = append(a.Results, mret)
			a.Count = append(a.Count, len(mret))
			if a.Options.One && len(mret) > 0 {
				break
			}
		}
		//if a.Options.Verbose || a.Options.Verbosely {
		//	log.Debugf("before CG: " + MemoryInfo())
		//}
		if a.Options.GCHeapSize <= HeapAlloc() {
			log.Debugf("Arrange: GC")
			runtime.GC()
		}
		//if a.Options.Verbose || a.Options.Verbosely {
		//	log.Debugf("after  CG: " + MemoryInfo())
		//}
		bpai++
	}

	for i := range a.WipedLength {
		if a.WipedLength[i] > 0 {
			return true, nil
		}
	}
	for i := range a.Count {
		if a.Count[i] > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (a *Arrange) ArrangePattern(bpai int, bpa []BasicPhrase) ([]ArrangeMatrixResult, error) {
	o := fmt.Sprintf("%2v-%2v: ", a.Number, bpai)
	for bpi, bp := range bpa {
		if bp.NewLine && bpi != 0 {
			o += color.FRed + "\\n" + color.Reset
		}
		o += string(bp.Surface)
	}
	fmt.Println(o)
	//bplength := bpPatternLength(sentences, a.Options.SwapSentences,
	//	a.Options.ProgressDepth)
	progress := NewArrangeProgress(a.Options, bpa)
	progressid := progress.Add("main")
	progress.Start()

	//mat := make([][]rune, 0)
	//mat = append(mat, make([]rune, a.Width))
	// 最大文字数
	textlength := 0
	for i := range bpa {
		textlength += bpa[i].PatternMaxLength
	}
	am, err := NewArrangeMatrix(a.Options, a.Keyword, a.Number, 0, bpa, 0, 0,
		bpai, nil, []int{0, 0}, a.Writer, progress, progressid,
		[]int{}, []int{}, a.Width)
	if err != nil {
		return nil, err
	}
	err = am.Search([]int{}, 0)
	progress.Stop()
	if err != nil {
		return nil, err
	}
	a.WipedLength[bpai] = am.WipedLength
	return am.MatrixResult, nil
}

func (a *Arrange) Output() error {
	ocount, err := a.Writer.Output(
		a.Keyword, a.Surfaces, a.Results, a.WipedLength, a.Width)
	total := 0
	count := make([]int, len(a.Count))
	for i := range a.Count {
		count[i] += a.Count[i]
		count[i] += a.WipedLength[i]
		if len(ocount) > i {
			count[i] += ocount[i]
		}
		total += count[i]
	}
	fmt.Printf("%v results found, %v\n", total, count)
	if a.Options.Verbose || a.Options.Verbosely {
		log.Debug(MemoryInfo())
	}

	return err
}
