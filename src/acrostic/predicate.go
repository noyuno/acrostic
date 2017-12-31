package acrostic

import (
	"errors"
	"strconv"
	"unicode"

	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

const PredicateHeaderLength = 7

// PredicateCaseElementGroup : 格要素群
type PredicateCaseElementGroup struct {
	Text []rune

	// 格
	CaseType []rune
	// フラグ
	Flag []rune
	// 表記
	Notation []rune

	// 基本句番号を持つかどうか
	HasBasicPhraseNumber bool
	// 基本句番号
	BasicPhraseNumber int

	// N文前を持つかどうか
	HasBeforeSentence bool
	// N文前
	BeforeSentence int

	// 文IDを持つかどうか
	HasSentenceId bool
	// 文ID
	SentenceId int
}

func caseElementInt(text []rune) (bool, int) {
	v, e := strconv.Atoi(string(text))
	if e != nil {
		return false, 0
	}
	return true, v
}

// NewPredicateCaseElementGroup : constructor
// text : スラッシュで6分割できなければならない
func NewPredicateCaseElementGroup(text []rune) *PredicateCaseElementGroup {
	ret := new(PredicateCaseElementGroup)
	ret.Text = text
	return ret
}

func (p *PredicateCaseElementGroup) Copy() *PredicateCaseElementGroup {
	ret := NewPredicateCaseElementGroup(p.Text)
	ret.CaseType = runes.Copy(p.CaseType)
	ret.Flag = runes.Copy(p.Flag)
	ret.Notation = runes.Copy(p.Notation)
	ret.HasBasicPhraseNumber = p.HasBasicPhraseNumber
	ret.BasicPhraseNumber = p.BasicPhraseNumber
	ret.HasBeforeSentence = p.HasBeforeSentence
	ret.BeforeSentence = p.BeforeSentence
	ret.HasSentenceId = p.HasSentenceId
	ret.SentenceId = p.SentenceId
	return ret
}

func (ce *PredicateCaseElementGroup) Analyze() error {
	if ce.Text == nil {
		return errors.New("PredicateCaseElementGroup.Text is null")
	}
	v := runes.Split(ce.Text, []rune("/"))
	if len(v) != 6 {
		return errors.New("cannot split PredicateCaseElementGroup.Text by 6")
	}
	ce.CaseType = v[0]
	ce.Flag = v[1]
	ce.Notation = v[2]
	ce.HasBasicPhraseNumber, ce.BasicPhraseNumber = caseElementInt(v[3])
	ce.HasBeforeSentence, ce.BeforeSentence = caseElementInt(v[4])
	ce.HasSentenceId, ce.SentenceId = caseElementInt(v[5])
	return nil
}

// Predicate : 述語側
type Predicate struct {
	// RepresentativeNotation : 格フレームIDの代表表記
	RepresentativeNotation []rune
	// ProverbsType : 格フレームIDの用言の種類
	ProverbsType []rune
	// CaseFrameNumber : 格フレームIDの格フレーム番号
	CaseFrameNumber []rune
	// PredicateCaseElementGroups : 格要素群
	PredicateCaseElementGroups []PredicateCaseElementGroup
	// BasicPhrase : knpの基本句の出力行
	BasicPhrase []rune
}

// NewPredicate : constructor
func NewPredicate(bp []rune) *Predicate {
	ret := new(Predicate)
	ret.BasicPhrase = bp
	return ret
}

func (p *Predicate) Copy() *Predicate {
	r := NewPredicate(p.BasicPhrase)
	r.RepresentativeNotation = runes.Copy(p.RepresentativeNotation)
	r.ProverbsType = runes.Copy(p.ProverbsType)
	r.CaseFrameNumber = runes.Copy(p.CaseFrameNumber)
	r.PredicateCaseElementGroups = make([]PredicateCaseElementGroup, len(p.PredicateCaseElementGroups))
	for i := range p.PredicateCaseElementGroups {
		r.PredicateCaseElementGroups[i] = p.PredicateCaseElementGroups[i]
	}
	return r
}

func (p *Predicate) Analyze() error {
	// <よりも前の文字列は読み飛ばす
	begintag := []rune("<格解析結果::")
	begin := runes.Index(p.BasicPhrase, begintag, 0)
	if begin == -1 {
		return errors.New("Predicate.Analyze cannot find begin of tag")
	}
	tags := make([][]rune, 0, 10)
	for {
		end := runes.Index(p.BasicPhrase, []rune(">"), begin)
		if end == -1 {
			return errors.New("cannot find end of tag")
		}
		tags = append(tags, p.BasicPhrase[begin:end])
		begin := runes.Index(p.BasicPhrase, begintag, end)
		if begin == -1 {
			break
		}
	}

	for _, t := range tags {
		beginc := 0
		rnum := beginc
		for foundnum := false; foundnum == false || unicode.IsDigit(t[rnum]); rnum++ {
			if unicode.IsDigit(t[rnum]) {
				foundnum = true
			}
		}
		//log.WithFields(log.Fields{
		//	"len(begintag)": len(begintag),
		//	"rnum":          rnum}).Debug("[Predicate.Analyze]")
		caseFrameId := t[len(begintag):rnum]
		slash := runes.Index(caseFrameId, []rune("/"), 0)
		if slash == -1 {
			continue
		}
		p.RepresentativeNotation = caseFrameId[0:slash]
		colon2 := runes.Index(caseFrameId, []rune(":"), 0)
		p.ProverbsType = caseFrameId[slash+1 : colon2]
		p.CaseFrameNumber = caseFrameId[colon2+1:]
		caseElems := t[rnum:]
		for _, ce := range runes.Split(caseElems, []rune(";")) {
			caseElement := NewPredicateCaseElementGroup(ce)
			err := caseElement.Analyze()
			if err != nil {
				return err
			}
			p.PredicateCaseElementGroups = append(p.PredicateCaseElementGroups, *caseElement)
		}
	}

	//p.DebugPrint()
	return nil
}

func caseElementIntToString(b bool, i int) string {
	if b {
		return strconv.Itoa(i)
	} else {
		return "-"
	}
}

func (c *PredicateCaseElementGroup) DebugPrint() {
	sBasicPhraseNumber := caseElementIntToString(c.HasBasicPhraseNumber, c.BasicPhraseNumber)
	sBeforeSentence := caseElementIntToString(c.HasBeforeSentence, c.BeforeSentence)
	sSentenceId := caseElementIntToString(c.HasSentenceId, c.SentenceId)

	log.WithFields(log.Fields{
		"CaseType":          string(c.CaseType),
		"Flag":              string(c.Flag),
		"Notation":          string(c.Notation),
		"BasicPhraseNumber": sBasicPhraseNumber,
		"BeforeSentence":    sBeforeSentence,
		"SentenceId":        sSentenceId}).Debug("[PredicateCaseElementGroup]")
}

func (p *Predicate) DebugPrint() {
	log.WithFields(log.Fields{
		"RepresentativeNotation": string(p.RepresentativeNotation),
		"ProverbsType":           string(p.ProverbsType),
		"CaseFrameNumber":        string(p.CaseFrameNumber)}).Debug("[Predicate]")
	for _, c := range p.PredicateCaseElementGroups {
		c.DebugPrint()
	}
}
