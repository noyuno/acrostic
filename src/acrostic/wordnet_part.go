package acrostic

type WordNetPart int

const (
	WNUnknownPart WordNetPart = iota
	WNVerbPart
	WNAdjectivePart
	WNNounPart
	WNAdverbPart
	WNPartEnd
)

func (w WordNetPart) String() string {
	switch w {
	case WNUnknownPart:
		return "unknown"
	case WNVerbPart:
		return "v"
	case WNAdjectivePart:
		return "a"
	case WNNounPart:
		return "n"
	case WNAdverbPart:
		return "r"
	default:
		return "unknown"
	}
}
func NewWordNetPart(s string) WordNetPart {
	switch s {
	case "v":
		return WNVerbPart
	case "a":
		return WNAdjectivePart
	case "n":
		return WNNounPart
	case "r":
		return WNAdverbPart
	default:
		return WNUnknownPart
	}
}

func ToWordNetPart(p Part) WordNetPart {
	switch p {
	case VerbPart:
		return WNVerbPart
	case AdjectivePart:
		return WNAdjectivePart
	case AdjectiveVerbPart:
		return WNAdjectivePart
	case NounPart:
		return WNNounPart
	case AdnominalPart:
		return WNAdjectivePart
	case AdverbPart:
		return WNAdverbPart
	case ConjunctionPart:
		return WNUnknownPart
	case EmotiveVerbPart:
		return WNUnknownPart
	case AuxiliaryVerbPart:
		return WNUnknownPart
	case ParticlePart:
		return WNUnknownPart
	case SuffixPart:
		return WNUnknownPart
	case SpecialPart:
		return WNUnknownPart
	default:
		return WNUnknownPart
	}
}
