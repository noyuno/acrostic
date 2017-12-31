package acrostic

// Part : 品詞
// 接尾辞は品詞ではないが，処理の都合上，分けるのは面倒なので，品詞の一つとみなす
type Part int

const (
	// UnknownPart : 不明な品詞
	UnknownPart Part = iota
	// VerbPart : 動詞
	VerbPart
	// AdjectivePart : 形容詞
	AdjectivePart
	// AdjectiveVerbPart : 形容動詞
	AdjectiveVerbPart
	// NounPart : 名詞
	NounPart
	// AdnominalPart : 連体詞
	AdnominalPart
	// AdverbPart : 副詞
	AdverbPart
	// ConjunctionPart : 接続詞
	ConjunctionPart
	// EmotiveVerbPart : 感動詞
	EmotiveVerbPart
	// AuxiliaryVerbPart : 助動詞
	AuxiliaryVerbPart
	// ParticlePart : 助詞
	ParticlePart
	// SuffixPart : 接尾辞（品詞ではない）
	SuffixPart
	// PrefixPart : 接頭辞（品詞ではない）
	PrefixPart
	// SpecialPart : 特殊（句読点等，当然品詞ではない）
	SpecialPart
	// DeterminePart : 判定詞（だ）
	DeterminePart
	// DemonstrativePart : 指示詞
	DemonstrativePart
)

func (p Part) String() string {
	switch p {
	case VerbPart:
		return "動詞"
	case AdjectivePart:
		return "形容詞"
	case AdjectiveVerbPart:
		return "形容動詞"
	case NounPart:
		return "名詞"
	case AdnominalPart:
		return "連体詞"
	case AdverbPart:
		return "副詞"
	case ConjunctionPart:
		return "接続詞"
	case EmotiveVerbPart:
		return "感動詞"
	case AuxiliaryVerbPart:
		return "助動詞"
	case ParticlePart:
		return "助詞"
	case SuffixPart:
		return "接尾辞"
	case PrefixPart:
		return "接頭辞"
	case SpecialPart:
		return "特殊"
	case DeterminePart:
		return "判定詞"
	case DemonstrativePart:
		return "指示詞"
	default:
		return "不明"
	}
}

func (p Part) Rune() []rune {
	return []rune(p.String())
}

func NewPart(s []rune) Part {
	t := string(s)
	switch t {
	case "動詞":
		return VerbPart
	case "形容詞":
		return AdjectivePart
	case "形容動詞":
		return AdjectiveVerbPart
	case "名詞":
		return NounPart
	case "連体詞":
		return AdnominalPart
	case "副詞":
		return AdverbPart
	case "接続詞":
		return ConjunctionPart
	case "感動詞":
		return EmotiveVerbPart
	case "助動詞":
		return AuxiliaryVerbPart
	case "助詞":
		return ParticlePart
	case "接尾辞":
		return SuffixPart
	case "接頭辞":
		return PrefixPart
	case "特殊":
		return SpecialPart
	case "判定詞":
		return DeterminePart
	case "指示詞":
		return DemonstrativePart
	default:
		return UnknownPart
	}
}

// IsIndependent : 自立語かどうかを判定
// return : trueならばその品詞は自立語，falseならばその品詞は付属語または不明な品詞
func (p Part) IsIndependent() bool {
	return !p.IsAdjunct() && !p.IsSuffix() && !p.IsSpecial() && p != UnknownPart
}

// IsAdjunct : 付属語かどうかを判定
// return trueならば付属語
func (p Part) IsAdjunct() bool {
	return (p == AuxiliaryVerbPart || p == ParticlePart)
}

// IsSuffix : 接尾辞かどうかを判定
// return trueならば接尾辞
func (p Part) IsSuffix() bool {
	return p == SuffixPart
}

// IsSpecial : 特殊かどうかを判定
// return trueならば特殊
func (p Part) IsSpecial() bool {
	return p == SpecialPart
}

// Inflection : 語形変化するかどうかを判定
// return : trueならば語形変化する
func (p Part) IsFlection() bool {
	return p == VerbPart || p == AuxiliaryVerbPart
}
