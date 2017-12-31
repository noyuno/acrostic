package acrostic

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type WordNetSynset struct {
	Options  *Options
	Instance *Instance
}

func NewWordNetSynset(o *Options, i *Instance) *WordNetSynset {
	ret := new(WordNetSynset)
	ret.Options = o
	ret.Instance = i
	return ret
}

func (w *WordNetSynset) WordID(a []rune, part WordNetPart) (string, error) {
	rows, err := w.Instance.WordNet.DB.Query("select wordid,pos from word where lemma=?", string(a))
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			wordid string
			pos    string
		)
		if err = rows.Scan(&wordid, &pos); err != nil {
			return "", err
		}
		if part.String() == pos {
			return wordid, nil
		}
	}
	return "", nil
}

func (w *WordNetSynset) Synset(id string) ([]string, error) {
	var ret []string
	rows, err := w.Instance.WordNet.DB.Query("select synset from sense where wordid=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var synset string
		if err = rows.Scan(&synset); err != nil {
			return nil, err
		}
		ret = append(ret, synset)
	}
	return ret, nil
}

func (w *WordNetSynset) Name(syns string) (string, error) {
	rows, err := w.Instance.WordNet.DB.Query("select name from synset where synset=?", syns)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return "", err
		}
		return name, nil
	}
	return "", nil
}

func (w *WordNetSynset) Hype(syns string) (string, error) {
	rows, err := w.Instance.WordNet.DB.Query("select synset2 from synlink where synset1=? and link=\"hype\"",
		syns)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var synset string
		if err = rows.Scan(&synset); err != nil {
			return "", err
		}
		return synset, nil
	}
	return "", nil
}

// Nearest: aについて，bに共通で最も近いsynset（概念）のIDを取得する．
// a: 概念を取得するための文字列
// b: aの概念のうち，最も近いものとして挙げられる文字列
func (w *WordNetSynset) NearestSynset(
	a []rune,
	apart WordNetPart,
	b []rune,
	bpart WordNetPart) ([]SynsetResult, error) {
	// wordidを取得
	var (
		aid string
		bid string
		err error
	)
	aid, err = w.WordID(a, apart)
	if err != nil {
		return nil, err
	}
	bid, err = w.WordID(b, bpart)
	if err != nil {
		return nil, err
	}
	if len(aid) == 0 {
		return nil, errors.New("string a has not id")
	}
	if len(bid) == 0 {
		return nil, errors.New("string b has not id")
	}
	log.Debugf("aid: %v, bid: %v", aid, bid)
	// synsetを取得
	asyn, err := w.Synset(aid)
	if err != nil {
		return nil, err
	}
	bsyn, err := w.Synset(bid)
	if err != nil {
		return nil, err
	}
	log.Debugf("asyn: %v, bsyn: %v", asyn, bsyn)

	// 探索
	sr := make([]SynsetResult, 0)
	astep := make([]int, len(asyn))
	bstep := make([]int, len(bsyn))
	for i := range astep {
		astep[i] = 0
	}
	for i := range bstep {
		bstep[i] = 0
	}
	sr, err = w.Search(asyn, astep, 0, bsyn, bstep, 0, sr)
	if err != nil {
		return nil, err
	}
	log.Debugf("finished Search(), len = %v", len(sr))

	// ルートまで行く
	for r := range sr {
		v := sr[r].Synset
		for v != "" {
			v, err = w.Hype(v)
			if err != nil {
				return nil, err
			}
			sr[r].Depth++
		}
		if sr[r].AStep == 0 && sr[r].BStep == 0 {
			sr[r].Approximation = 0
		} else {
			sr[r].Approximation = (2.0 * float32(sr[r].Depth) /
				(float32(sr[r].AStep) + float32(sr[r].BStep)))
		}
		log.Debugf("sr[%v]: %v, %v[%v %v] = %v",
			r, sr[r].Synset, sr[r].Depth, sr[r].AStep, sr[r].BStep, sr[r].Approximation)
	}
	return sr, nil
}

type SynsetResult struct {
	Synset        string
	AStep         int
	BStep         int
	Depth         int
	Approximation float32
}

// 検索をする
//
func (w *WordNetSynset) Search(
	asyn []string,
	astep []int,
	an int,
	bsyn []string,
	bstep []int,
	bn int,
	sr []SynsetResult) ([]SynsetResult, error) {

	log.Debugf("Search: %v-%v(%v), %v-%v-(%v)", astep, len(asyn), asyn, bstep, len(bsyn), bsyn)

	// 探索
	for ai, as := range asyn {
		if as == "" {
			continue
		}
		for bi, bs := range bsyn {
			if bs == "" {
				continue
			}
			if as == bs {
				// 同じsynset
				found := false
				for i := range sr {
					if sr[i].Synset == as {
						found = true
						break
					}
				}
				if found == false {
					log.Debugf("Search: same synset %v [%v %v]", as, astep, bstep)
					sr = append(sr, SynsetResult{Synset: as, AStep: astep[ai], BStep: bstep[bi]})
					asyn[ai] = ""
					bsyn[bi] = ""
				}
			}
		}
	}
	// 上
	new := false
	for _, as := range asyn {
		h, err := w.Hype(as)
		if err != nil {
			return nil, err
		}
		if h == "" {
			// root
			continue
		}
		found := false
		for s := range sr {
			if sr[s].Synset == h {
				found = true
				break
			}
		}
		if found == false {
			for s := range asyn {
				if asyn[s] == h {
					found = true
					break
				}
			}
			if found == false {
				new = true
				asyn = append(asyn, h)
				astep = append(astep, an)
			}
		}
	}
	for _, bs := range bsyn {
		h, err := w.Hype(bs)
		if err != nil {
			return nil, err
		}
		if h == "" {
			// root
			continue
		}
		found := false
		for s := range sr {
			if sr[s].Synset == h {
				found = true
				break
			}
		}
		if found == false {
			for s := range bsyn {
				if bsyn[s] == h {
					found = true
					break
				}
			}
			if found == false {
				new = true
				bsyn = append(bsyn, h)
				bstep = append(bstep, bn)
			}
		}
	}
	if new {
		return w.Search(asyn, astep, an+1, bsyn, bstep, bn+1, sr)
	} else {
		return sr, nil
	}
}
