package acrostic

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/noyuno/lgo/color"
	log "github.com/sirupsen/logrus"
)

type ArrangeWriter struct {
	Options   *Options
	StdWriter *bufio.Writer
	Color     bool
	Keyword   []rune
	Number    int
	Mutex     sync.RWMutex
}

func NewArrangeWriter(o *Options, kn int, keyword []rune) *ArrangeWriter {
	ret := new(ArrangeWriter)
	ret.Options = o
	if ret.Options.OutFileName == "" {
		ret.StdWriter = bufio.NewWriter(os.Stdout)
		ret.Color = true
	}
	ret.Mutex = sync.RWMutex{}
	ret.Number = kn
	ret.Keyword = keyword
	return ret
}

func (a *ArrangeWriter) Truncate() error {
	a.Mutex.Lock()
	// ファイルを初期化
	if a.Options.OutFileName != "" {
		f, err := os.Create(a.Options.OutFileName)
		if err != nil {
			a.Mutex.Unlock()
			return err
		}
		defer f.Close()
	}
	a.Mutex.Unlock()
	return nil
}

func (a *ArrangeWriter) OutputKeyword() error {
	a.Mutex.Lock()
	if a.Options.OutFileName == "" {

	} else {
		f, err := os.OpenFile(a.Options.OutFileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			a.Mutex.Unlock()
			return err
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		_, err = w.WriteString(fmt.Sprintf("# %v: keyword: %v\n", a.Number, string(a.Keyword)))
		if err != nil {
			a.Mutex.Unlock()
			return err
		}
		err = w.Flush()
		if err != nil {
			a.Mutex.Unlock()
			return err
		}
	}
	a.Mutex.Unlock()
	return nil
}

// Output : まとめて出力します
// begin: 今まで書き込んだ数
func (a *ArrangeWriter) Output(
	keyword []rune,
	surfaces [][]rune,
	mats [][]ArrangeMatrixResult,
	begin []int,
	width int) ([]int, error) {

	var count []int
	var err error

	a.Mutex.Lock()
	if a.Options.OutputEachPattern == false {
		var w *bufio.Writer
		if a.Options.OutFileName == "" {
			w = a.StdWriter
		} else {
			var f *os.File
			f, err = os.OpenFile(a.Options.OutFileName, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				a.Mutex.Unlock()
				return nil, err
			}
			defer f.Close()
			w = bufio.NewWriter(f)
		}
		count, err = a.write(w, keyword, surfaces, mats, true, begin, width)
		if err != nil {
			return nil, err
		}
		err = w.Flush()
		if err != nil {
			a.Mutex.Unlock()
			return nil, err
		}
	} else {
		// do not anything
		//count = make([]int, len(array))
	}
	a.Mutex.Unlock()
	return count, nil

}

func (a *ArrangeWriter) OutputPattern(
	keyword []rune,
	surface []rune,
	mreti int,
	mret []ArrangeMatrixResult,
	writecount bool,
	begin int,
	width int) (int, error) {

	var w *bufio.Writer

	a.Mutex.Lock()
	if a.Options.OutFileName == "" {
		w = a.StdWriter
	} else {
		f, err := os.OpenFile(a.Options.OutFileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			a.Mutex.Unlock()
			return 0, err
		}
		defer f.Close()
		w = bufio.NewWriter(f)
	}
	count, err := a.writePattern(w, keyword, surface, mreti, mret, writecount, begin, width)
	if err != nil {
		return 0, err
	}
	err = w.Flush()
	if err != nil {
		a.Mutex.Unlock()
		return 0, err
	}

	//fmt.Printf("%v results found\n", count)
	//log.Debug(MemoryInfo())
	if a.Options.Verbose || a.Options.Verbosely {
		//fmt.Printf("arrange process time: %v\n", a.ElapsedTime)
	}
	a.Mutex.Unlock()
	return count, nil
}

func (a *ArrangeWriter) write(
	w *bufio.Writer,
	keyword []rune,
	surfaces [][]rune,
	mats [][]ArrangeMatrixResult,
	writecount bool,
	begin []int,
	width int) ([]int, error) {
	T, _ := i18n.Tfunc(a.Options.Language)

	total := 0
	r := make([]int, len(mats))
	for i := range r {
		r[i] = 0
	}
	var err error
	for reti, ret := range mats {
		r[reti], err = a.writePattern(
			w, keyword, surfaces[reti], reti, ret, writecount, begin[reti], width)
		if err != nil {
			return nil, err
		}
		total += r[reti]
	}
	if total == 0 {
		_, err = w.WriteString("# " + T("not found any patterns") + "\n")
	} else {
		_, err = w.WriteString(fmt.Sprintf("# "+T("total %v results found")+"\n", total))
	}
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (a *ArrangeWriter) writePattern(
	w *bufio.Writer,
	keyword []rune,
	surface []rune,
	reti int,
	ret []ArrangeMatrixResult,
	writecount bool,
	begin int,
	width int) (int, error) {
	T, _ := i18n.Tfunc(a.Options.Language)

	out := fmt.Sprintf("# %v-%v: %v (%v)\n",
		a.Number, reti, string(surface), width)
	//for _, r := range ret {
	//	out += (string(r.Surface))
	//}
	if keyword == nil {
		log.Fatal("writePattern: keyword == nil")
	}

	for ti, t := range ret {
		out += (fmt.Sprintf("# %v-%v-%v stack:%v bstack:%v\n",
			a.Number, reti, ti+begin, t.PatternStack, t.BranchStack))
		//w.WriteString(fmt.Sprintf("KeywordEnd = %v\n", t.KeywordEnd))
		if t.KeywordEnd == nil {
			log.Warnf("writePattern: ArrangeMatrixResult[%v].KeywordEnd == nil", ti)
			continue
		}
		if len(t.KeywordEnd) != 2 {
			log.Warnf("writePattern: ArrangeMatrixResult[%v].KeywordEnd length is %v", len(t.KeywordEnd))
			continue
		}
		startrow := t.KeywordEnd[0] - len(keyword) + 1
		for ri, r := range t.Matrix {
			//end := false
			for ci, c := range r {
				if string(c) == "" || c == 0 {
					//end = true
					break
				}
				if a.Color &&
					startrow <= ri && ri <= t.KeywordEnd[0] &&
					t.KeywordEnd[1] == ci {
					out += (color.FGreen + Wide(c) + color.Reset)
				} else {
					out += Wide(c)
				}
			}
			out += ("\n")
			//if end {
			//	break
			//}
		}
	}
	if writecount {
		if len(ret)+begin == 0 {
			out += "# " + T("not found") + "\n"
		} else {
			out += (fmt.Sprintf("# "+T("%v results found")+"\n", len(ret)+begin))
		}
	}
	_, err := w.WriteString(out)
	if err != nil {
		return 0, err
	}

	return len(ret), nil
}
