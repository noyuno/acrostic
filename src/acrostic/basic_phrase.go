package acrostic

import (
	"errors"
	"strings"

	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

// BasicPhrase : 自立語の基本句(+から始まる)および形態素(*+以外から始まる)
type BasicPhrase struct {
	// Options : オプション
	Options *Options
	// Instance : Instance
	Instance *Instance

	// Text : テキスト
	Text [][]rune

	Keywords [][]rune

	NewLine bool

	// Part : 基本句の品詞
	Part Part
	// BasicPhrase : 基本句(+から始まる)
	BasicPhrase []rune

	// Independent : 自立語の形態素(*+以外から始まる)
	Independent [][]rune
	// HasIndependent : 自立語があるかどうか
	HasIndependent bool
	// IndependentSurface : 自立語の表層
	IndependentSurface [][]rune

	AllIndependentSurface []rune

	// Suffix : 接尾辞(「渋/み」の「み」「開催/さ/れる」の「れる」)
	//Suffix []rune
	// HasSuffix : 接尾辞があるかどうか
	HasSuffix bool
	// SuffixSurface : 接尾辞の表層
	SuffixSurface [][]rune
	SuffixKana    [][]rune
	// 接尾辞の形
	SuffixForm [][]rune

	AllSuffixSurface []rune

	// 接頭辞があるかどうか
	HasPrefix bool
	// 接頭辞の表層
	PrefixSurface [][]rune
	PrefixKana    [][]rune

	//AllPrefixSurface []rune

	// Adjunct : 付属語の形態素(*+以外から始まる)
	//Adjunct []rune
	// HasAdjunct : 付属語があるかどうか
	//HasAdjunct bool

	// 助動詞
	HasAuxiliaryVerb     bool
	AuxiliaryVerbSurface []rune

	// 助詞
	HasParticle        bool
	ParticleSurface    [][]rune
	AllParticleSurface []rune

	// Special : 特殊文字（句読点など）
	Special []rune
	// HasSpecial : 特殊文字があるかどうか
	HasSpecial bool
	// SpecialSurface : 特殊文字の表層
	SpecialSurface []rune

	// Determine : 判定詞（だ）
	Determine []rune
	// HasDetermine : 判定詞があるかどうか
	HasDetermine bool
	// DetermineSurface : 判定詞の表層
	DetermineSurface []rune
	// 判定詞の原形
	DetermineOrigin []rune

	// IndependentSurfaceは活用する語で，原文は丁寧語であるかどうか
	HasInflectionPolite bool
	// IndependentSurfaceは活用する語で，丁寧語にした語
	InflectionPolite []rune

	// Surface : 表層．自立語，接尾辞，付属語をまとめた表層
	Surface []rune

	// SurfaceOrder : 表層に現れる品詞の順番
	// 例えばIndependentSurfaceのような複数許容しているものはその数だけ追加されるべきである
	SurfaceOrder []Part

	// Origin : 原形
	Origin []rune

	// Kana : IndependentSurfaceのかな
	Kana []rune

	// CaseAnalysisType : 格解析のタイプ
	CaseAnalysisType CaseAnalysisType

	// HasPredicate : 述語側であるかどうか（<格解析結果:...>を持つ基本句かどうか）
	//HasPredicate bool
	// Predicate : 述語側（<格解析結果:...>）
	Predicate Predicate

	// HasCaseElement : 格要素側であるかどうか（<解析格:...>, <解析連格:...>を持つ基本句かどうか）
	//HasCaseElement bool
	// CaseElement : 格要素側（<解析格:...>, <解析連格:...>を持つ基本句かどうか）
	CaseElement CaseElement

	// PredicateTerm : 格解析
	PredicateTerm PredicateTerm

	// PhraseNumber : この基本句が含まれている文節の番号
	PhraseNumber int

	// Number : Phraseのはじめからつけられた固有の番号
	Number int

	// ID : Sentenceのはじめからつけられた固有の番号
	ID int

	// Synonyms : 類語
	Synonyms []WordNetResult

	// 活用形の型
	InflectionType []rune

	// 活用形の形
	InflectionForm []rune

	// ドメイン
	Domain []rune

	// 類義語およびそのかなのをぜんぶまとめたもの
	Pattern [][]rune

	PatternLengthMap map[int]bool

	// Patternの最大文字数
	PatternMaxLength int

	// Pattern中に含まれるキーワードの位置
	// キーワード番号, キーワード文字列番目, 出現文字数(map)->パターン番号
	PatternKeywordPos [][]map[int][]int

	// 過去形かどうか
	Past bool
}

// NewBasicPhrase : constructor
// bw: knpの基本句の出力行
func NewBasicPhrase(
	o *Options,
	i *Instance,
	t [][]rune,
	pn int,
	n int,
	begin int,
	newline bool,
	keywords [][]rune,
) *BasicPhrase {
	ret := new(BasicPhrase)
	ret.Text = t
	ret.Options = o
	ret.Instance = i
	ret.PhraseNumber = pn
	ret.Number = n
	ret.ID = begin
	//ret.ID = i.Variables.BasicPhrase
	ret.SurfaceOrder = make([]Part, 0)
	ret.NewLine = newline
	ret.Keywords = keywords
	ret.PatternLengthMap = map[int]bool{}
	return ret
}

func (b *BasicPhrase) Copy() *BasicPhrase {
	var text [][]rune
	if b.Options.EnableDeepCopy {
		text = runes.CopyArray(b.Text)
	} else {
		text = b.Text
	}
	r := NewBasicPhrase(b.Options, b.Instance, text, b.PhraseNumber, b.Number, b.ID, b.NewLine, b.Keywords)
	r.Surface = runes.Copy(b.Surface)
	//r.AdjunctSurface = runes.Copy(b.AdjunctSurface)
	r.AuxiliaryVerbSurface = runes.Copy(b.AuxiliaryVerbSurface)
	r.ParticleSurface = runes.CopyArray(b.ParticleSurface)
	r.AllParticleSurface = runes.Copy(b.AllParticleSurface)
	r.NewLine = b.NewLine
	r.PatternKeywordPos = b.PatternKeywordPos
	r.PatternLengthMap = b.PatternLengthMap
	if b.Options.EnableDeepCopy {
		r.Part = b.Part
		r.BasicPhrase = runes.Copy(b.BasicPhrase)
		r.Independent = runes.CopyArray(b.Independent)
		r.HasIndependent = b.HasIndependent
		r.IndependentSurface = runes.CopyArray(b.IndependentSurface)
		r.AllIndependentSurface = runes.Copy(b.AllIndependentSurface)
		//r.Suffix = runes.Copy(b.Suffix)
		r.HasSuffix = b.HasSuffix
		r.SuffixSurface = runes.CopyArray(b.SuffixSurface)
		r.SuffixKana = runes.CopyArray(b.SuffixKana)
		r.SuffixForm = runes.CopyArray(b.SuffixForm)
		r.AllSuffixSurface = runes.Copy(b.AllSuffixSurface)
		r.HasPrefix = b.HasPrefix
		r.PrefixSurface = runes.CopyArray(b.PrefixSurface)
		r.PrefixKana = runes.CopyArray(b.PrefixKana)
		//r.AllPrefixSurface = runes.Copy(b.AllPrefixSurface)
		//r.Adjunct = runes.Copy(b.Adjunct)
		//r.HasAdjunct = b.HasAdjunct
		r.HasAuxiliaryVerb = b.HasAuxiliaryVerb
		r.HasParticle = b.HasParticle
		r.Special = runes.Copy(b.Special)
		r.HasSpecial = b.HasSpecial
		r.SpecialSurface = runes.Copy(b.SpecialSurface)
		r.Determine = runes.Copy(b.Determine)
		r.HasDetermine = b.HasDetermine
		r.DetermineSurface = runes.Copy(b.DetermineSurface)
		r.DetermineOrigin = runes.Copy(b.DetermineOrigin)
		r.HasInflectionPolite = b.HasInflectionPolite
		r.InflectionPolite = runes.Copy(b.InflectionPolite)
		r.SurfaceOrder = make([]Part, len(b.SurfaceOrder))
		for i := range b.SurfaceOrder {
			r.SurfaceOrder[i] = b.SurfaceOrder[i]
		}
		r.Origin = runes.Copy(b.Origin)
		r.Kana = runes.Copy(b.Kana)
		r.CaseAnalysisType = b.CaseAnalysisType
		r.Predicate = *b.Predicate.Copy()
		r.CaseElement = *b.CaseElement.Copy()
		r.PredicateTerm = *b.PredicateTerm.Copy()
		r.Synonyms = make([]WordNetResult, len(b.Synonyms))
		for i := range b.Synonyms {
			r.Synonyms[i] = *b.Synonyms[i].Copy()
		}
		r.PhraseNumber = b.PhraseNumber
		r.Number = b.Number
		r.ID = b.ID
		r.InflectionType = runes.Copy(b.InflectionType)
		r.InflectionForm = runes.Copy(b.InflectionForm)
		r.Domain = runes.Copy(b.Domain)
		r.Pattern = runes.CopyArray(b.Pattern)
	} else {
		r.Part = b.Part
		r.BasicPhrase = b.BasicPhrase
		r.Independent = b.Independent
		r.HasIndependent = b.HasIndependent
		r.IndependentSurface = b.IndependentSurface
		r.AllIndependentSurface = b.AllIndependentSurface
		//r.Suffix = b.Suffix
		r.HasSuffix = b.HasSuffix
		r.SuffixSurface = b.SuffixSurface
		r.SuffixKana = b.SuffixKana
		r.SuffixForm = b.SuffixForm
		r.AllSuffixSurface = b.AllSuffixSurface
		r.HasPrefix = b.HasPrefix
		r.PrefixSurface = b.PrefixSurface
		r.PrefixKana = b.PrefixKana
		//r.AllPrefixSurface = b.AllPrefixSurface
		//r.Adjunct = b.Adjunct
		//r.HasAdjunct = b.HasAdjunct
		r.HasAuxiliaryVerb = b.HasAuxiliaryVerb
		r.HasParticle = b.HasParticle
		r.Special = b.Special
		r.HasSpecial = b.HasSpecial
		r.SpecialSurface = b.SpecialSurface
		r.Determine = b.Determine
		r.HasDetermine = b.HasDetermine
		r.DetermineSurface = b.DetermineSurface
		r.DetermineOrigin = b.DetermineOrigin
		r.HasInflectionPolite = b.HasInflectionPolite
		r.InflectionPolite = b.InflectionPolite
		r.SurfaceOrder = b.SurfaceOrder
		r.Origin = b.Origin
		r.Kana = b.Kana
		r.CaseAnalysisType = b.CaseAnalysisType
		r.Predicate = b.Predicate
		r.CaseElement = b.CaseElement
		r.PredicateTerm = b.PredicateTerm
		r.Synonyms = b.Synonyms
		r.PhraseNumber = b.PhraseNumber
		r.Number = b.Number
		r.ID = b.ID
		r.InflectionType = b.InflectionType
		r.InflectionForm = b.InflectionForm
		r.Domain = b.Domain
		r.Pattern = b.Pattern
	}
	return r
}

// 言い換え
func (bp *BasicPhrase) AppendParaphrase(s []rune) {
	para := bp.Instance.Paraphrase.Replace(s)
	if !runes.Compare(s, para) {
		log.Debugf("%v can paraphrase into %v", string(s), string(para))
		bp.Pattern = append(bp.Pattern, para)
	}
}

// sにAdjunct, Suffix, Special, Determineを付加してPatternsに追加する
func (bp *BasicPhrase) AppendPattern(s []rune, suffix bool, onlykeywords bool) {
	//log.Debugf("AppendPattern: %v", string(s))
	bp.AppendPatternBase(s, bp.DetermineSurface, suffix, onlykeywords)
	if !runes.Compare(bp.DetermineSurface, bp.DetermineOrigin) {
		bp.AppendPatternBase(s, bp.DetermineOrigin, suffix, onlykeywords)
	}
}

func (bp *BasicPhrase) AppendPatternBase(s []rune, determine []rune, suffix bool, onlykeywords bool) {
	ret := make([]rune, 0)
	particle := 0
	suf := 0
	if bp.HasPrefix {
		//for i = range bp.AllPrefixSurface {
		//	ret[i] = bp.AllPrefixSurface[i]
		//}
		if bp.Options.UseKanji {
			for i := range bp.PrefixSurface {
				ret = append(ret, bp.PrefixSurface[i]...)
			}
		} else {
			for i := range bp.PrefixKana {
				ret = append(ret, bp.PrefixKana[i]...)
			}
		}
	}
	ret = append(ret, s...)
	for _, part := range bp.SurfaceOrder {
		if part == AuxiliaryVerbPart {
			ret = append(ret, bp.AuxiliaryVerbSurface...)
		} else if part == ParticlePart {
			ret = append(ret, bp.ParticleSurface[particle]...)
			particle++
		} else if part.IsSuffix() {
			if suffix {
				if bp.Options.UseKanji {
					ret = append(ret, bp.SuffixSurface[suf]...)
				} else {
					ret = append(ret, bp.SuffixKana[suf]...)
				}
				suf++
			}
		} else if part.IsSpecial() {
			ret = append(ret, bp.SpecialSurface...)
		} else if part == DeterminePart {
			ret = append(ret, determine...)
		} else if part == NounPart ||
			part == VerbPart ||
			part == AdjectivePart ||
			part == PrefixPart ||
			part == AdverbPart ||
			part == DemonstrativePart {
			// do not anything
		} else {
			log.Fatalf("BasicPhrase.AppendPatternBase: unknown part: %v", part.String())
		}
	}
	// 同じ文字列長の単語が入っているか
	if bp.Options.AllWordLength {
		if _, ok := bp.PatternLengthMap[len(ret)]; !ok {
			// force
			// すでに入っているわけがない
			bp.Pattern = append(bp.Pattern, ret)
			bp.PatternLengthMap[len(ret)] = true
			return
		}
	}
	if onlykeywords && bp.Options.OnlyKeywords {
		if bp.HasKeyword(string(ret)) == false {
			return
		}
	}
	// キーワードの文字が入っているか
	if onlykeywords && bp.Options.OnlyKeywords {
		found := false
		for i := range bp.Keywords {
			for k := range bp.Keywords[i] {
				if strings.Contains(string(ret), string(bp.Keywords[i][k])) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found == false {
			return
		}
	}
	// すでに入っているか
	for i := range bp.Pattern {
		if runes.Compare(bp.Pattern[i], ret) {
			return
		}
	}
	bp.Pattern = append(bp.Pattern, ret)
	bp.PatternLengthMap[len(ret)] = true
}

func (bp *BasicPhrase) Analyze() error {
	basicPhraseToken := []rune("+")[0]
	phraseToken := []rune("*")[0]
	spaceToken := []rune(" ")
	domainToken := []rune("ドメイン:")
	dquoteToken := []rune("\"") // "
	// check input text
	if len(bp.Text) < 2 && len(bp.Text) > 3 {
		return errors.New("The number of BasicPhrase.Text lines is wrong. Require 2 or 3 lines text.")
	}
	if bp.Text[0][0] != basicPhraseToken {
		return errors.New("BasicPhrase.Text[0] must be start with +(plus), but given " + string(bp.Text[0]))
	}
	bp.BasicPhrase = bp.Text[0]
	for k := 1; k < len(bp.Text); k++ {
		//log.WithFields(log.Fields{"len(bp.Text)": len(bp.Text), "k": k}).Debug()
		if bp.Text[k][0] == phraseToken || bp.Text[k][0] == basicPhraseToken {
			return errors.New("BasicPhrase.Text[" + string(2) +
				"] must NOT be start with +(plus) and *(asterisk)")
		}
		a := runes.Split(bp.Text[k], spaceToken)
		part := NewPart(a[3])
		bp.SurfaceOrder = append(bp.SurfaceOrder, part)
		//log.WithFields(log.Fields{"part": part.String(), "surface": string(a[0])}).Debug()
		if part.IsIndependent() {
			if bp.HasIndependent {
				//return errors.New("BasicPhrase.Part.Independent has already exists")
				//log.Debug("append BasicPhrase.Part.Independent")
			}
			bp.HasIndependent = true
			bp.Independent = append(bp.Independent, bp.Text[k])
			arr := runes.Split(bp.Independent[len(bp.Independent)-1], spaceToken)
			bp.IndependentSurface = append(bp.IndependentSurface, arr[0])
			bp.AllIndependentSurface = append(bp.AllIndependentSurface, arr[0]...)
			bp.Part = part
			bp.InflectionType = arr[7]
			bp.InflectionForm = arr[9]
			for i := 11; i < len(arr); i++ {
				//log.WithFields(log.Fields{"arr[i]": string(arr[i]), "i": i}).Debug()
				if runes.Compare(arr[i][0:len(domainToken)], domainToken) {
					bp.Domain = arr[i][len(domainToken):]
				}
				if runes.Compare(arr[i][len(arr[i])-1:], []rune(dquoteToken)) {
					break
				}
			}
			bp.Kana = append(bp.Kana, a[1]...)
			if HasOnlyKana(bp.Kana) == false {
				ok := false
				bp.Kana, ok = bp.Instance.Kana.Get(bp.Kana)
				if ok == false {
					log.Warnf("could not get kana: %v", string(bp.Kana))
				}
			}
			if part.IsFlection() {
				if bp.HasInflectionPolite == false {
					bp.HasInflectionPolite = bp.Instance.JumanKnp.IsInflectionPolite(arr[9])
				}
				bp.Past = bp.Instance.JumanKnp.IsPast(arr[9])
			}
		} else if part.IsSuffix() {
			//if bp.HasSuffix {
			//	return fmt.Errorf("BasicPhrase.Part.Suffix has already exists: %v", string(bp.SuffixSurface))
			//}
			bp.HasSuffix = true
			//bp.Suffix = bp.Text[k]
			s := runes.Split(bp.Text[k], spaceToken)
			bp.SuffixSurface = append(bp.SuffixSurface, s[0])
			k, o := bp.Instance.Kana.Get(s[0])
			if o {
				bp.SuffixKana = append(bp.SuffixKana, k)
			} else {
				log.Warnf("could not get Kana of Suffix: %v", string(s[0]))
			}
			bp.SuffixForm = append(bp.SuffixForm, s[9])
			bp.AllSuffixSurface = append(bp.AllSuffixSurface, s[0]...)
		} else if part == PrefixPart {
			bp.HasPrefix = true
			s := runes.Split(bp.Text[k], spaceToken)
			bp.PrefixSurface = append(bp.PrefixSurface, s[0])
			k, o := bp.Instance.Kana.Get(s[0])
			if o {
				bp.PrefixKana = append(bp.PrefixKana, k)
			} else {
				log.Warnf("could not get Kana of Prefix: %v", string(s[0]))
			}
			//bp.AllPrefixSurface = append(bp.AllPrefixSurface, s[0]...)
		} else if part == AuxiliaryVerbPart {
			if bp.HasAuxiliaryVerb {
				return errors.New("BasicPhrase.Part.AuxiliaryVerbPart has already exists")
			}
			bp.HasAuxiliaryVerb = true
			bp.AuxiliaryVerbSurface = runes.Split(bp.Text[k], spaceToken)[0]
		} else if part == ParticlePart {
			//if bp.HasParticle {
			//	return errors.New("BasicPhrase.Part.ParticlePart has already exists")
			//}
			bp.HasParticle = true
			s := runes.Split(bp.Text[k], spaceToken)[0]
			bp.ParticleSurface = append(bp.ParticleSurface, s)
			bp.AllParticleSurface = append(bp.AllParticleSurface, s...)
		} else if part.IsSpecial() {
			bp.HasSpecial = true
			bp.Special = bp.Text[k]
			bp.SpecialSurface = runes.Split(bp.Special, spaceToken)[0]
		} else if part == DeterminePart {
			bp.HasDetermine = true
			bp.Determine = bp.Text[k]
			det := runes.Split(bp.Determine, spaceToken)
			bp.DetermineSurface = det[0]
			bp.DetermineOrigin = det[2]
			bp.Past = bp.Instance.JumanKnp.IsPast(det[9])
		} else {
			return errors.New("unknown BasicPhrase.Part: " + string(a[3]))
		}
		bp.Surface = append(bp.Surface, a[0]...)
	}
	bp.Origin = runes.Split(bp.Text[1], spaceToken)[2]
	bp.Kana = KatakanaToHiragana(bp.Kana)

	bp.CaseAnalysisType = GetCaseAnalysisType(bp.Text[0])

	// 格解析の種類: 述語側か格解析側かを判定
	var err error
	if bp.Part.IsSpecial() == false {
		switch bp.CaseAnalysisType {
		case PredicateSide:
			bp.Predicate = *NewPredicate(bp.BasicPhrase)
			err = bp.Predicate.Analyze()
		case CaseElementSide:
			bp.CaseElement = *NewCaseElement(bp.BasicPhrase)
			err = bp.CaseElement.Analyze()
		case NoneSide:
			//return errors.New("unable to find case analysis type")
			//log.WithFields(log.Fields{"Surface": string(bp.Surface)}).Info(
			//	"BasicPhrase: unable to find case analysis type")
		}
		if err != nil {
			return err
		}

		// 格解析をする
		if bp.Options.CaseAnalysis {
			bp.PredicateTerm = *NewPredicateTerm(bp.Text[0])
			err = bp.PredicateTerm.Analyze()
			if err != nil {
				return err
			}
		}
	}

	if bp.NewLine {
		log.Debugf("BasicPhrase: %v '%v' NewLine = true", bp.ID, string(bp.Surface))
	}

	// パターン更新
	err = bp.UpdatePattern()
	if err != nil {
		return err
	}
	//log.Debugf("SurfaceOrder: %v", bp.SurfaceOrder)

	// debug
	//log.WithFields(log.Fields{
	//	"Part":             bp.Part.String(),
	//	"BasicPhrase":      string(bp.BasicPhrase),
	//	"Independent":      string(bp.Independent),
	//	"Suffix":           string(bp.Suffix),
	//	"Adjunct":          string(bp.Adjunct),
	//	"CaseAnalysisType": bp.CaseAnalysisType.String(),
	//	"Surface":          string(bp.Surface),
	//	"Synonyms":         len(bp.Synonyms)}).Debug("BasicPhrase.Analyze")
	return nil
}

func (bp *BasicPhrase) HasKeyword(s string) bool {
	// キーワードの文字が入っているか
	found := false
	for i := range bp.Keywords {
		for k := range bp.Keywords[i] {
			if strings.Contains(s, string(bp.Keywords[i][k])) {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	return found
}

func (bp *BasicPhrase) Given(s string) bool {
	return string(bp.Surface) == s ||
		string(bp.Kana) == s
}

func (bp *BasicPhrase) RemoveSameLengthPattern() {
	le := map[int]bool{}
	array := make([][]rune, 0)
	for i := range bp.Pattern {
		p := string(bp.Pattern[i])
		if bp.HasKeyword(p) || bp.Given(p) {
			//log.Debugf("%v contains keyword", string(bp.Pattern[i]))
			array = append(array, bp.Pattern[i])
			le[len(bp.Pattern[i])] = true
		} else {
			//log.Debugf("%v does not contain keyword", string(bp.Pattern[i]))
		}
	}
	for i := range bp.Pattern {
		if _, ok := le[len(bp.Pattern[i])]; ok == false {
			//log.Debugf("length %v is not in Pattern, insert %v",
			//	len(bp.Pattern[i]), string(bp.Pattern[i]))
			array = append(array, bp.Pattern[i])
			le[len(bp.Pattern[i])] = true
		}
	}
	bp.Pattern = bp.Pattern[:0]
	bp.Pattern = array
}

func (bp *BasicPhrase) UpdateSurface() {
	independent := 0
	suffix := 0
	particle := 0
	log.Debugf("BasicPhrase.UpdateSurface: %v-%v", bp.ID, string(bp.Surface))
	oldsurface := make([]rune, len(bp.Surface))
	copy(oldsurface, bp.Surface)
	bp.Surface = bp.Surface[:0]
	for _, part := range bp.SurfaceOrder {
		log.Debugf("part = %v", part)
		if part.IsIndependent() {
			bp.Surface = append(bp.Surface, bp.IndependentSurface[independent]...)
			independent++
		} else if part == AuxiliaryVerbPart {
			bp.Surface = append(bp.Surface, bp.AuxiliaryVerbSurface...)
		} else if part == ParticlePart {
			log.Debugf("len(bp.ParticleSurface)=%v, particle=%v", len(bp.ParticleSurface), particle)
			bp.Surface = append(bp.Surface, bp.ParticleSurface[particle]...)
			particle++
		} else if part.IsSpecial() {
			bp.Surface = append(bp.Surface, bp.SpecialSurface...)
		} else if part.IsSuffix() {
			if bp.Options.UseKanji {
				bp.Surface = append(bp.Surface, bp.SuffixSurface[suffix]...)
			} else {
				bp.Surface = append(bp.Surface, bp.SuffixKana[suffix]...)
			}
			suffix++
		} else if part == PrefixPart {
			//bp.Surface = append(bp.Surface, bp.AllPrefixSurface...)
			if bp.Options.UseKanji {
				for i := range bp.PrefixSurface {
					bp.Surface = append(bp.Surface, bp.PrefixSurface[i]...)
				}
			} else {
				for i := range bp.PrefixKana {
					bp.Surface = append(bp.Surface, bp.PrefixKana[i]...)
				}
			}
		} else {
			log.Fatalf("unknown Part: %v", part.String())
		}
	}
	log.Debugf("BasicPhrase.UpdateSurface %v -> %v",
		string(oldsurface), string(bp.Surface))
}

func (bp *BasicPhrase) UpdatePattern() error {
	var err error
	bp.Pattern = make([][]rune, 0)

	if bp.Options.UseKanji {
		bp.Pattern = append(bp.Pattern, bp.Surface)
		bp.AppendParaphrase(bp.Surface)
	}
	if bp.Options.UseKana {
		bp.AppendPattern(bp.Kana, true, false)
	}

	// 丁寧語
	if bp.Part.IsFlection() && bp.Options.UsePolite {
		if bp.HasInflectionPolite {
			// この語形変化する語はすでに丁寧
			// 丁寧でない形にするために，
			// 接尾辞を取って語形変化する語の形を接尾辞の形にする
			_, i, f := bp.Instance.JumanKnp.Inflection(
				bp.IndependentSurface[len(bp.IndependentSurface)-1], bp.SuffixForm[len(bp.SuffixForm)-1])
			if f {
				a := make([]rune, 0)
				for i := 0; i < len(bp.IndependentSurface)-1; i++ {
					a = append(a, bp.IndependentSurface[i]...)
				}
				a = append(a, i...)
				if bp.Options.UseKanji {
					bp.AppendPattern(a, false, false)
				}
				if bp.Options.UseKana {
					k, f := bp.Instance.Kana.Get(a)
					if f {
						bp.AppendPattern(k, true, false)
					} else {
						log.Warnf("could not get kana: %v", string(a))
					}
				}
			}
		} else {
			// この語形変化する語は丁寧でない
			// 丁寧な形にする
			//log.Warnf("not polite: %v", string(bp.Origin))
			p, f := bp.Instance.JumanKnp.InflectionPolite(bp.Origin, bp.InflectionForm)
			if f {
				//log.Warnf("got InflectionPolite: %v", string(p))
				a := make([]rune, 0)
				for i := 0; i < len(bp.IndependentSurface)-1; i++ {
					a = append(a, bp.IndependentSurface[i]...)
				}
				a = append(a, p...)
				if bp.Options.UseKanji {
					bp.AppendPattern(a, false, false)
				}
				if bp.Options.UseKana {
					k, f := bp.Instance.Kana.Get(a)
					if f {
						bp.AppendPattern(k, true, false)
					} else {
						log.Warnf("could not get kana: %v", string(a))
					}
				}
			} else {
				log.Warnf("polite was not created: %v", string(bp.Surface))
			}
		}
	}
	//bp.UpdatePatternMaxLength()

	// 類語検索をするよ
	if bp.Options.Synonyms && bp.HasIndependent &&
		((bp.Part == VerbPart && bp.Options.SynonymsVerb) || bp.Part != VerbPart) {
		var wnr []WordNetResult
		wnr, err = bp.Instance.WordNet.GetSynonyms(bp, []WordNetLink{WNSynonym, WNHype})
		if err != nil {
			return err
		}
		bp.Synonyms = wnr
		for _, s := range bp.Synonyms {
			if s.HasInflection {
				if bp.Options.UseKanji {
					bp.AppendPattern(s.InflectionSurface, true, true)
				}
				if s.HasPolite {
					log.Fatalf("polite found: %v", string(s.PoliteSurface))
					if bp.Options.UseKanji {
						bp.AppendPattern(s.PoliteSurface, true, true)
					}
					if bp.Options.UseKana {
						bp.AppendPattern(s.PoliteKana, true, true)
					}
				}
			} else {
				if bp.Options.UseKanji {
					bp.AppendPattern(s.Surface, true, true)
				}
			}
			if s.HasKana && bp.Options.UseKana {
				bp.AppendPattern(s.Kana, true, true)
			}
		}
		// test whether bp.Synonyms contains other part
		//for i := range bp.Synonyms {
		//	if ToWordNetPart(bp.Part) != bp.Synonyms[i].Part {
		//		log.WithFields(log.Fields{"Synonym": bp.Synonyms[i], "i": i}).Fatal("BasicPhrase: contains other part")
		//	}
		//}
	}

	// 同じ文字数で，キーワードの文字を含んでいなければ削除
	if bp.Options.SkipSameLength {
		bp.RemoveSameLengthPattern()
	}

	// パターン数を制限する
	// 組み合わせ問題であるから，演算はパターン数で放物線状に増加する．
	// 枝刈りでも防げないときは，それぞれのパターン数を制限するしかない
	if len(bp.Pattern) > bp.Options.WordPatternLimit {
		bp.Pattern = bp.Pattern[:bp.Options.WordPatternLimit]
	}
	bp.UpdatePatternMaxLength()

	bp.MarkKeywordPos(bp.Keywords)

	return nil
}

// Patternの最大文字数を取得する
func (bp *BasicPhrase) UpdatePatternMaxLength() {
	o := 0
	for i := range bp.Pattern {
		if o < len(bp.Pattern[i]) {
			o = len(bp.Pattern[i])
		}
	}
	bp.PatternMaxLength = o
}

func (bp *BasicPhrase) MarkKeywordPos(keywords [][]rune) {
	bp.PatternKeywordPos = make([][]map[int][]int, len(keywords))
	for i := range keywords {
		bp.PatternKeywordPos[i] = make([]map[int][]int, len(keywords[i]))
		for k := range keywords[i] {
			key := []rune(string(keywords[i][k]))
			bp.PatternKeywordPos[i][k] = map[int][]int{}
			//log.Debugf("len(bp.Pattern)=%v", len(bp.Pattern))
			for m := range bp.Pattern {
				r := runes.Index(bp.Pattern[m], key, 0)
				for r != -1 {
					if _, ok := bp.PatternKeywordPos[i][k][r]; !ok {
						bp.PatternKeywordPos[i][k][r] = make([]int, 0)
					}
					bp.PatternKeywordPos[i][k][r] = append(bp.PatternKeywordPos[i][k][r], m)
					r = runes.Index(bp.Pattern[m], key, r+1)
				}
				//if strings.Contains(string(bp.Pattern[m]), string(keywords[i][k])) {
				//	bp.PatternKeywordPos[i][k] = append(bp.PatternKeywordPos[i][k], m)
				//	//log.Debugf("pkp[%v][%v] = %v(%v)", i, k, m, string(bp.Pattern[m]))
				//}
			}
		}
	}
}
