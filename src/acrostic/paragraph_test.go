package acrostic

import (
	"os"
	"testing"

	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

func TSS1(t *testing.T) {
	in := []rune(`
浜松沖墜落の空自ヘリ、海底から１遺体引き揚げ

　航空自衛隊の救難ヘリコプターが静岡県浜松市沖に墜落し、乗員４人が行方不明となっている事故で、空自は３日、現場周辺の海底から乗員１人の遺体を引き揚げたと発表した。
　空自は、発見されていない３人の捜索を続けている。
　発表によると、遺体は先月２９日、水深約７００メートルの海底で発見された。３日午後、ワイヤを使って海上へ引き揚げた。飛行服などから身元を特定したが、空自は、「遺族の了解が得られていない」として、隊員の氏名を明らかにしていない。
　事故は１０月１７日に発生。夜間訓練のため浜松基地を離陸した「ＵＨ６０Ｊ」が、同基地南約３０キロの太平洋上で墜落した。
`)

	expectedsent := [][]rune{
		[]rune(`浜松沖墜落の空自ヘリ、海底から１遺体引き揚げ`),
		[]rune(`　航空自衛隊の救難ヘリコプターが静岡県浜松市沖に墜落し、乗員４人が行方不明となっている事故で、空自は３日、現場周辺の海底から乗員１人の遺体を引き揚げたと発表した。`),
		[]rune(`　空自は、発見されていない３人の捜索を続けている。`),
		[]rune(`　発表によると、遺体は先月２９日、水深約７００メートルの海底で発見された。`),
		[]rune(`３日午後、ワイヤを使って海上へ引き揚げた。`),
		[]rune(`飛行服などから身元を特定したが、空自は、「遺族の了解が得られていない」として、隊員の氏名を明らかにしていない。`),
		[]rune(`　事故は１０月１７日に発生。`),
		[]rune(`夜間訓練のため浜松基地を離陸した「ＵＨ６０Ｊ」が、同基地南約３０キロの太平洋上で墜落した。`),
	}
	expectednl := []bool{true, true, true, true, false, false, true, false}

	sent, newline := SplitSentence(in, false)

	for i := range expectedsent {
		if !runes.Compare(expectedsent[i], sent[i]) {
			t.Errorf("want sent[%v]=%v, but returned %v", i, string(expectedsent[i]), string(sent[i]))
		}
		if expectednl[i] != newline[i] {
			t.Errorf("want newline[%v]=%v, but returned %v", i, expectednl[i], newline[i])
		}
	}
}

func TSS2(t *testing.T) {
	in := []rune(`おいしい牛乳　５００ｍｌ
たまごＭサイズ　６コ
ブタコマ
塩鮭切り身
カットヤサイ
あらびきソーセージ　１個

`)

	expectedsent := [][]rune{
		[]rune(`おいしい牛乳　５００ｍｌ`),
		[]rune(`たまごＭサイズ　６コ`),
		[]rune(`ブタコマ`),
		[]rune(`塩鮭切り身`),
		[]rune(`カットヤサイ`),
		[]rune(`あらびきソーセージ　１個`),
	}
	expectednl := []bool{true, true, true, true, true, true}

	sent, newline := SplitSentence(in, true)

	for i := range expectedsent {
		if !runes.Compare(expectedsent[i], sent[i]) {
			t.Errorf("want sent[%v]=%v, but returned %v", i, string(expectedsent[i]), string(sent[i]))
		}
		if expectednl[i] != newline[i] {
			t.Errorf("want newline[%v]=%v, but returned %v", i, expectednl[i], newline[i])
		}
	}

}

func TestSplitSentence(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		DisableSorting:   false,
		QuoteEmptyFields: true,
		ForceColors:      true,
		FullTimestamp:    false,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	TSS1(t)
	TSS2(t)
}
