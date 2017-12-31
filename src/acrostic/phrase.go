package acrostic

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"

	"github.com/noyuno/lgo/runes"
)

// Phrase : 文節
type Phrase struct {
	//Options : Options
	Options *Options
	//Instance : Instance
	Instance *Instance
	Text     [][]rune
	Keywords [][]rune
	NewLine  bool
	// Destination : 係り先の文節番号
	Destination int
	// DependencyType : 係り受けの種類(DとかPとか出るやつ)
	DependencyType []rune
	// BasicPhrases : 自立語(接尾辞を含む)および自立語にぶら下がった付属語
	BasicPhrases []BasicPhrase
	// Surface : 表層
	//Surface []rune
	// Kana : かな
	//Kana []rune
	// Number : 文節番号
	Number int
	// 並列する語かどうか
	Parallel bool
}

// NewPhrase : constructor
// tはknpの文節の出力を行ごとにvectorに格納されたもの（*から始まり，*の手前で終わる複数行．
func NewPhrase(o *Options, i *Instance, t [][]rune, n int, newline bool, keywords [][]rune) *Phrase {
	ret := new(Phrase)
	ret.Text = t
	ret.Options = o
	ret.Instance = i
	ret.Number = n
	ret.NewLine = newline
	ret.Keywords = keywords
	return ret
}

func (p *Phrase) Copy() *Phrase {
	ret := NewPhrase(p.Options, p.Instance, p.Text, p.Number, p.NewLine, p.Keywords)
	ret.Destination = p.Destination
	ret.BasicPhrases = make([]BasicPhrase, len(p.BasicPhrases))
	for i := range p.BasicPhrases {
		ret.BasicPhrases[i] = *p.BasicPhrases[i].Copy()
	}
	ret.Parallel = p.Parallel
	if p.Options.EnableDeepCopy {
		ret.DependencyType = runes.Copy(p.DependencyType)
	} else {
		ret.DependencyType = p.DependencyType
	}
	return ret
}

// Analyze : 解析
func (p *Phrase) Analyze(begin int) (int, error) {
	out := ""
	for i := range p.Text {
		out += string(p.Text[i]) + "\n"
	}
	//log.Warnf("Phrase.Analyze: %v;", out)
	if len(p.Text) <= 2 {
		return 0, fmt.Errorf("Phrase.Text: want len(p.Text)>2, but given %v", len(p.Text))
	}
	if !(string(p.Text[0][0]) == "*") {
		return 0, fmt.Errorf("Phrase.Text: want start at *(asterisk), but given %v", string(p.Text[0][0]))
	}

	space := runes.Index(p.Text[0], []rune(" "), 0)
	dt := space
	foundalphabet := false
	for foundalphabet == false || unicode.IsDigit(p.Text[0][dt]) ||
		string(p.Text[0][dt]) == "-" {
		if !(unicode.IsDigit(p.Text[0][dt]) || string(p.Text[0][dt]) == "-") {
			foundalphabet = true
		}
		dt++
	}
	var err error
	//log.WithFields(log.Fields{"Destination": string(p.Text[0][space+1 : dt])}).Debug("")
	p.Destination, err = strconv.Atoi(string(p.Text[0][space+1 : dt]))
	if err != nil {
		return 0, errors.New("cannot convert p.Text[0][" + strconv.Itoa(space) + ":" +
			strconv.Itoa(dt) + "]:" + string(p.Text[0][space:dt]) + " to int")
	}
	p.DependencyType = p.Text[0][dt : dt+1]
	if string(p.DependencyType) == "P" {
		p.Parallel = true
	}

	// BasicPhrase
	start := 1
	var i int
	// 付属語は，必ずBasicPhraseにぶら下げる．
	//log.WithFields(log.Fields{"len(p.Text)": len(p.Text)}).Debug("")
	for i = 1; i < len(p.Text); i++ {
		//log.WithFields(log.Fields{"i": i}).Debug("")
		if i != start && string(p.Text[i][0]) == "+" {
			//for _, d := range p.Text[start:i] {
			//	log.Debug(string(d))
			//}
			newline := len(p.BasicPhrases) == 0 && p.NewLine
			bp := NewBasicPhrase(p.Options, p.Instance, p.Text[start:i],
				p.Number, len(p.BasicPhrases), begin+len(p.BasicPhrases), newline, p.Keywords)
			err := bp.Analyze()
			if err != nil {
				return 0, err
			}
			p.BasicPhrases = append(p.BasicPhrases, *bp)
			start = i
		}
	}
	if start+1 == i {
		return 0, errors.New("remained only one line")
	} else if start < i {
		//for _, d := range p.Text[start:] {
		//	log.Debug(string(d))
		//}
		newline := len(p.BasicPhrases) == 0 && p.NewLine
		bp := NewBasicPhrase(p.Options, p.Instance, p.Text[start:],
			p.Number, len(p.BasicPhrases), begin+len(p.BasicPhrases), newline, p.Keywords)
		err := bp.Analyze()
		if err != nil {
			return 0, err
		}
		p.BasicPhrases = append(p.BasicPhrases, *bp)
	}

	// debug
	//log.WithFields(log.Fields{"Destination": strconv.Itoa(p.Destination),
	//	"DependencyType": string(p.DependencyType),
	//	"BasicPhrases":   len(p.BasicPhrases)}).Debug("[Phrase]")
	return begin + len(p.BasicPhrases), nil
}

func (p *Phrase) Surface() []rune {
	// append ... で配列内をコピーしてくれる(検証済み)
	ret := []rune("")
	for i := range p.BasicPhrases {
		ret = append(ret, p.BasicPhrases[i].Surface...)
	}
	return ret
}

func (p *Phrase) Kana() []rune {
	// append ... で配列内をコピーしてくれる(検証済み)
	ret := []rune("")
	for i := range p.BasicPhrases {
		ret = append(ret, p.BasicPhrases[i].Kana...)
	}
	return ret
}
