package acrostic

import "github.com/noyuno/lgo/runes"

// CaseAnalysisType : 格解析のタイプ
type CaseAnalysisType int

const (
	// NoneSide : なし
	NoneSide = iota
	// PredicateSide : 述語側
	PredicateSide
	// CaseElementSide : 格要素側
	CaseElementSide
)

// GetCaseAnalysisType : 格解析のタイプを取得する
// t: テキスト
// return: 格解析のタイプ
func GetCaseAnalysisType(t []rune) CaseAnalysisType {
	pred := runes.Index(t, []rune("<格解析結果:"), 0)
	if pred != -1 {
		return PredicateSide
	}
	ce := runes.Index(t, []rune("<解析格:"), 0)
	if ce != -1 {
		return CaseElementSide
	}
	ce = runes.Index(t, []rune("解析連絡:"), 0)
	if ce != -1 {
		return CaseElementSide
	} else {
		return NoneSide
	}
}

func (cat CaseAnalysisType) String() string {
	switch cat {
	case PredicateSide:
		return "述語側"
	case CaseElementSide:
		return "格要素側"
	default:
		return "なし"
	}
}
