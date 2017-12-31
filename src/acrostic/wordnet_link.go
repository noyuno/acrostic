package acrostic

import log "github.com/sirupsen/logrus"

// http://compling.hss.ntu.edu.sg/wnja/
type WordNetLink int

func NewWordNetLink(t string) WordNetLink {
	switch t {
	case "synonyms":
		return WNSynonym
	case "hype":
		return WNHype
	case "hypo":
		return WNHypo
	}
	log.Fatalf("not implemented: value of %v", t)
	return WNEnd
}

const (
	// WNSynonym : 同義語
	WNSynonym WordNetLink = iota
	// WNHype : 上位語
	WNHype
	// WNHypo : 下位語
	WNHypo
	// WNMprt : 被構成要素(部分)
	WNMprt
	// WNHprt : 構成要素(部分)
	WNHprt
	// WNHmem : 構成要素(構成員)
	WNHmem
	// WNMmem : 被構成要素(構成員)
	WNMmem
	// WNMsub : 被構成要素(物質・材料)
	WNMsub
	// WNHsub : 構成要素(物質・材料)
	WNHsub
	// WNDmnc : 被包含領域(カテゴリ)
	WNDmnc
	// WNDmtc : 包含領域(カテゴリ)
	WNDmtc
	// WNDmnu : 被包含領域(語法)
	WNDmnu
	// WNDmtu : 包含領域(語法)
	WNDmtu
	// WNDmnr : 被包含領域(地域)
	WNDmnr
	// WNDmtr : 包含領域(地域)
	WNDmtr
	// WNInst : 例
	WNInst
	// WNHasi : 例あり
	WNHasi
	// WNEnta : 含意
	WNEnta
	// WNCaus : 引き起こし
	WNCaus
	// WNAlso : 関連
	WNAlso
	// WNAttr : 属性
	WNAttr
	// WNSim : 近似
	WNSim
	// WNEnd : 終わり
	WNEnd
)
