package acrostic

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func tNewBasicPhrases(t ...string) []BasicPhrase {
	ret := make([]BasicPhrase, 0)
	for i := range t {
		ret = append(ret, BasicPhrase{Surface: []rune(t[i])})
	}
	return ret
}

func TestSentencePattern(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		DisableSorting:   false,
		QuoteEmptyFields: true,
		ForceColors:      true,
		FullTimestamp:    false,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	in := [][][]BasicPhrase{
		[][]BasicPhrase{
			tNewBasicPhrases("こんにちは，", "私の", "名前は", "綾地", "寧々", "です．"),
			tNewBasicPhrases("こんにちは，", "私の", "名前は", "明日原", "ユウキ", "です．"),
		},
		[][]BasicPhrase{
			tNewBasicPhrases("今年の", "みかんは", "酸味が", "あって", "おいしい．"),
			tNewBasicPhrases("今年の", "オレンジは", "酸味が", "あって", "おいしい．"),
		},
	}
	i := 0
	log.Debugf("swap = false")
	for v := range sentencePattern(in, false) {
		o := ""
		for _, bp := range v {
			o += string(bp.Surface)
		}
		log.Debugf("%v: %v", i, o)
		i++
	}
	i = 0
	log.Debugf("swap = true")
	for v := range sentencePattern(in, true) {
		o := ""
		for _, bp := range v {
			o += string(bp.Surface)
		}
		log.Debugf("%v: %v", i, o)
		i++
	}
}
