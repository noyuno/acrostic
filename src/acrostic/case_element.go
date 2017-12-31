package acrostic

import (
	"errors"

	"github.com/noyuno/lgo/runes"
)

// CaseElement : 格要素側(PredicateのCaseElementGroup（格要素群）ではない)
type CaseElement struct {
	// knpの基本句の出力行
	BasicPhrase []rune

	// AnalysisCase : 解析格(被連体修飾詞以外)
	AnalysisCase []rune
	// HasAnalysisCase : 解析格を持つかどうか
	HasAnalysisCase bool

	// AnalysisConnection : 解析連絡(被連体修飾詞)
	AnalysisConnection []rune
	// HasAnalysisConnection : 解析連絡を持つかどうか
	HasAnalysisConnection bool
}

// NewCaseElement : constructor
func NewCaseElement(bp []rune) *CaseElement {
	ret := new(CaseElement)
	ret.BasicPhrase = bp
	return ret
}

func (c *CaseElement) Copy() *CaseElement {
	r := NewCaseElement(c.BasicPhrase)
	r.AnalysisCase = runes.Copy(c.AnalysisCase)
	r.HasAnalysisCase = c.HasAnalysisCase
	r.AnalysisConnection = runes.Copy(c.AnalysisConnection)
	r.HasAnalysisCase = c.HasAnalysisCase
	return r
}

// Analyze : 格要素側の解析
func (ce *CaseElement) Analyze() error {
	begincase := []rune("<解析格:")
	beginconnection := []rune("<解析連絡:")
	end := []rune(">")
	acb := runes.Index(ce.BasicPhrase, begincase, 0)
	if acb != -1 {
		ace := runes.Index(ce.BasicPhrase, end, acb)
		if ace == -1 {
			return errors.New("CaseElement.Analyze found begin of AnalysisCase, but not found end of.")
		}
		ce.AnalysisCase = ce.BasicPhrase[acb+len(begincase) : ace]
		ce.HasAnalysisCase = true
		//return errors.New("CaseElement.Analyze cannot find begin of tag")
	}
	acb = runes.Index(ce.BasicPhrase, beginconnection, 0)
	if acb != -1 {
		ace := runes.Index(ce.BasicPhrase, end, acb)
		if ace == -1 {
			return errors.New("CaseElement.Analyze found begin of AnalysisConnection, but not found end of.")
		}
		ce.AnalysisConnection = ce.BasicPhrase[acb+len(beginconnection) : ace]
		ce.HasAnalysisConnection = true
	}
	if !ce.HasAnalysisCase && !ce.HasAnalysisConnection {
		return errors.New("CaseElement.Text has not contain any of case element tags")
	}

	//log.WithFields(log.Fields{
	//	"AnalysisCase":       string(ce.AnalysisCase),
	//	"AnalysisConnection": string(ce.AnalysisConnection)}).Debug("[CaseElement]")
	return nil
}
