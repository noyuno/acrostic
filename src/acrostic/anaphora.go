package acrostic

import (
	"errors"
	"strconv"
	"unicode"

	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

// AnaphoraCaseElementGroup : 格要素群
type AnaphoraCaseElementGroup struct {
	Text []rune

	// 格
	CaseType []rune
	// フラグ
	Flag []rune
	// 表記
	Notation []rune
	// EID
	EntityID int
}

// NewAnaphoraCaseElementGroup : constructor
// text : スラッシュで4分割できなければならない
func NewAnaphoraCaseElementGroup(text []rune) *AnaphoraCaseElementGroup {
	ret := new(AnaphoraCaseElementGroup)
	ret.Text = text
	return ret
}

func (a *AnaphoraCaseElementGroup) Copy() *AnaphoraCaseElementGroup {
	r := NewAnaphoraCaseElementGroup(a.Text)
	r.CaseType = runes.Copy(a.CaseType)
	r.Flag = runes.Copy(a.Flag)
	r.Notation = runes.Copy(a.Notation)
	r.EntityID = a.EntityID
	return r
}

func (ce *AnaphoraCaseElementGroup) Analyze() error {
	if ce.Text == nil {
		return errors.New("AnaphoraCaseElementGroup.Text is null")
	}
	v := runes.Split(ce.Text, []rune("/"))
	if len(v) != 4 {
		return errors.New("cannot split AnaphoraCaseElementGroup.Text by 4")
	}
	ce.CaseType = v[0]
	ce.Flag = v[1]
	ce.Notation = v[2]
	eid, err := strconv.Atoi(string(v[3]))
	if err != nil {
		return errors.New("cannot convert AnaphoraCaseElementGroup.EntityID []rune to int")
	}
	ce.EntityID = eid
	return nil
}

type PredicateTerm struct {
	Text []rune
	// EntityID : Entity ID (常に取得できる)
	EntityID int

	// IsAvailable : 項構造が参照できるかどうか
	IsAvailable bool

	// 以下の変数は，IsAvailable == trueのときのみ取得できる
	// RepresentativeNotation : 格フレームIDの代表表記
	RepresentativeNotation []rune
	// ProverbsType : 格フレームIDの用言の種類
	ProverbsType []rune
	// CaseFrameNumber : 格フレームIDの格フレーム番号
	CaseFrameNumber int
	// AnaphoraCaseElementGroups : 格要素群
	AnaphoraCaseElementGroups []AnaphoraCaseElementGroup
}

func NewPredicateTerm(t []rune) *PredicateTerm {
	ret := new(PredicateTerm)
	ret.Text = t
	return ret
}

func (a *PredicateTerm) Copy() *PredicateTerm {
	r := NewPredicateTerm(a.Text)
	r.EntityID = a.EntityID
	r.IsAvailable = a.IsAvailable
	r.RepresentativeNotation = runes.Copy(a.RepresentativeNotation)
	r.ProverbsType = runes.Copy(a.ProverbsType)
	r.CaseFrameNumber = a.CaseFrameNumber
	r.AnaphoraCaseElementGroups = make([]AnaphoraCaseElementGroup, len(a.AnaphoraCaseElementGroups))
	for i := range a.AnaphoraCaseElementGroups {
		r.AnaphoraCaseElementGroups[i] = *a.AnaphoraCaseElementGroups[i].Copy()
	}
	return r
}

func (a *PredicateTerm) Analyze() error {
	var err error
	eidbegintag := []rune("<EID:")
	structbegintag := []rune("<項構造:")
	endtag := []rune(">")
	splittag := []rune(":")
	semicolon := []rune(";")

	// eid
	ebegin := runes.Index(a.Text, eidbegintag, 0)
	if ebegin == -1 {
		return errors.New("Anaphora.Analyze cannot find begin of eid tag")
	}
	eend := runes.Index(a.Text, endtag, ebegin)
	if eend == -1 {
		return errors.New("Anaphora.Analyze cannot find end of eid tag")
	}
	eid := a.Text[ebegin+len(eidbegintag) : eend]
	a.EntityID, err = strconv.Atoi(string(eid))
	//log.Debugf("Anaphora.EntityID: %v", a.EntityID)
	if err != nil {
		return errors.New("cannot convert Anaphora.EntityID from []rune to int")
	}

	// struct
	sbegin := runes.Index(a.Text, structbegintag, 0)
	a.IsAvailable = sbegin != -1
	if sbegin == -1 {
		//return errors.New("Anaphora.Analyze cannot find begin of struct tag")
		log.Debug("Anaphora.Analyze cannot find begin of struct tag")
	} else {
		send := runes.Index(a.Text, endtag, sbegin)
		if send == -1 {
			return errors.New("Anaphora.Analyze cannot find end of struct tag")
		}
		stru := a.Text[sbegin+len(structbegintag) : send]
		colon := runes.Index(stru, splittag, 0)
		if colon == -1 {
			return errors.New("Anaphora.Analyze cannot find colon in struct tag")
		}
		a.RepresentativeNotation = stru[:colon]

		i := colon + 1
		for ; i < len(stru); i++ {
			if unicode.IsDigit(stru[i]) {
				break
			}
		}
		a.ProverbsType = stru[colon+1 : i]
		//log.WithFields(log.Fields{"stru": string(stru), "RepresentativeNotation": string(a.RepresentativeNotation), "ProverbsType": string(a.ProverbsType)}).Debug("[Anaphora]")
		fid := i
		for ; i < len(stru); i++ {
			//log.WithFields(log.Fields{"stru[i]": string(stru[i])}).Debug()
			if !unicode.IsDigit(stru[i]) {
				break
			}
		}
		a.CaseFrameNumber, err = strconv.Atoi(string(stru[fid:i]))
		if err != nil {
			return errors.New("cannot convert CaseFrameNumber from []rune to int")
		}
		//log.WithFields(log.Fields{"stru": string(stru), "i": i, "len(stru)": len(stru)}).Debug()
		if i+1 < len(stru) {
			ceg := runes.Split(stru[i+1:], semicolon)
			for _, c := range ceg {
				ce := NewAnaphoraCaseElementGroup(c)
				ce.Analyze()
				a.AnaphoraCaseElementGroups = append(a.AnaphoraCaseElementGroups, *ce)
			}
		}
	}

	//a.DebugPrint()
	return nil
}
func (c *AnaphoraCaseElementGroup) DebugPrint() {
	log.WithFields(log.Fields{
		"CaseType": string(c.CaseType),
		"Flag":     string(c.Flag),
		"Notation": string(c.Notation),
		"EntityID": strconv.Itoa(c.EntityID)}).Debug("[AnaphoraCaseElementGroup]")
}

func (p *PredicateTerm) DebugPrint() {
	log.WithFields(log.Fields{
		"IsAvailable":            p.IsAvailable,
		"EntityID":               strconv.Itoa(p.EntityID),
		"RepresentativeNotation": string(p.RepresentativeNotation),
		"ProverbsType":           string(p.ProverbsType),
		"CaseFrameNumber":        strconv.Itoa(p.CaseFrameNumber)}).Debug("[Anaphora]")
	for _, c := range p.AnaphoraCaseElementGroups {
		c.DebugPrint()
	}
}
