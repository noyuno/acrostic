package acrostic

import (
	"strings"

	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

const HiraganaList = `
ぁあぃいぅうぇえぉおかがきぎくぐ
けげこごさざしじすずせぜそぞただ
ちぢっつづてでとどなにぬねのはば
ぱひびぴふぶぷへべぺほぼぽまみむ
めもゃやゅゆょよらりるれろゎわゐ
ゑをんー・ゔヵヶ`

const KatakanaList = `
ァアィイゥウェエォオカガキギクグ
ケゲコゴサザシジスズセゼソゾタダ
チヂッツヅテデトドナニヌネノハバ
パヒビピフブプヘベペホボポマミム
メモャヤュユョヨラリルレロヮワヰ
ヱヲンー・ヴヵヶ`

const NumberList = `0123456789-`
const NumberZenkakuList = `０１２３４５６７８９`

const AlphabetList = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz., %`
const AlphabetZenkakuList = `
ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰ
ＱＲＳＴＵＶＷＸＹＺａｂｃｄｅｆ
ｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖ
ｗｘｙｚ．，％　`

func HasOnlyKana(text []rune) bool {
	for _, c := range text {
		if strings.Index(HiraganaList, string(c)) == -1 &&
			strings.Index(KatakanaList, string(c)) == -1 &&
			strings.Index(NumberList, string(c)) == -1 &&
			strings.Index(NumberZenkakuList, string(c)) == -1 &&
			strings.Index(AlphabetList, string(c)) == -1 &&
			strings.Index(AlphabetZenkakuList, string(c)) == -1 {
			return false
		}
	}
	return true
}

func KatakanaToHiragana(text []rune) []rune {
	ret := make([]rune, len(text))
	kl := []rune(KatakanaList)
	hl := []rune(HiraganaList)
	for i := range text {
		if p := runes.Index(kl, text[i:i+1], 0); p != -1 {
			ret[i] = hl[p]
		} else {
			ret[i] = text[i]
		}
	}
	return ret
}

type Kana struct {
	Options   *Options
	Instance  *Instance
	kanaCache map[string][]rune
}

func NewKana(o *Options, i *Instance) *Kana {
	ret := new(Kana)
	ret.Options = o
	ret.Instance = i
	ret.kanaCache = map[string][]rune{}
	return ret
}

func (k *Kana) Get(text []rune) ([]rune, bool) {
	if v, o := k.kanaCache[string(text)]; o {
		return v, true
	}
	ret := []rune("")
	for _, mode := range k.Options.KanaModeOrder {
		if mode == "juman" {
			ret = k.Instance.JumanKnp.GetKana(text)
		}
		if mode == "mecab" {
			ret = KatakanaToHiragana(k.Instance.MeCab.GetKana(text))

		}
		if mode == "kakasi" {
			ret = k.Instance.Kakasi.GetKana(text)
		}
		if HasOnlyKana(ret) {
			k.kanaCache[string(text)] = ret
			return ret, true
		}
		log.WithFields(log.Fields{
			"Kana": string(ret),
			"Text": string(text),
		}).Debugf("failed to get Kana using %v.", mode)
	}
	//log.WithFields(log.Fields{
	//	"Kana": string(ret),
	//	"Text": string(text),
	//}).Warningf("could not get Kana in any mode: %v", k.Options.KanaModeOrder)
	return nil, false
}

func Wide(c rune) string {
	cc := string(c)
	if strings.Contains(NumberList, cc) || strings.Contains(AlphabetList, cc) {
		return cc + " "
	}
	return cc
}
