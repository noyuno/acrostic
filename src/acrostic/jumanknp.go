package acrostic

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/kr/pty"
	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

type JumanKnp struct {
	JumanPipe *os.File
	JKPipe    *os.File
	//Imis      *os.File

	InflectionDB        map[string]map[string]string
	inflectionTypeCache map[string][]rune
	Options             *Options
	Instance            *Instance
}

func NewJumanKnp(o *Options, i *Instance) (*JumanKnp, error) {
	var err error
	ret := new(JumanKnp)
	ret.Options = o
	ret.Instance = i
	c := strings.Split(ret.Options.JumanCommand, " ")[0]
	err = exec.Command("which", c).Run()
	if err != nil {
		return nil, errors.New("command not found: " + c)
	}
	err = exec.Command("which", strings.Split(ret.Options.KnpCommand, " ")[0]).Run()
	if err != nil {
		return nil, errors.New("JumanKnp.KnpCommand not found: " + ret.Options.KnpCommand)
	}

	j := exec.Command("sh", "-c", ret.Options.JumanCommand)
	ret.JumanPipe, err = pty.Start(j)
	if err != nil {
		return nil, err
	}
	//log.Debug("sh -c " + ret.Options.JumanCommand + "|" + ret.Options.KnpCommand)
	jk := exec.Command("sh", "-c", ret.Options.JumanCommand+"|"+ret.Options.KnpCommand)
	ret.JKPipe, err = pty.Start(jk)
	if err != nil {
		return nil, err
	}

	//ret.Imis, err = os.Open(ret.Options.JumanPPDirectory + "/dic.imis")
	//if err != nil {
	//	return nil, err
	//}

	err = ret.ReadInflection()
	if err != nil {
		return nil, err
	}

	ret.inflectionTypeCache = map[string][]rune{}
	return ret, nil
}

func (jk *JumanKnp) Execute(text []rune, knp bool) []rune {
	var f *os.File
	if knp {
		f = jk.JKPipe
	} else {
		f = jk.JumanPipe
	}
	//log.Debugf("JumanKnp.Execute: %v", string(text))
	// 入力チェック
	if strings.TrimSpace(string(text)) == "" {
		log.Fatalf("JumanKnp.Execute: empty input")
		return nil
	}

	f.Write([]byte(string(text) + "\n"))
	s := bufio.NewScanner(f)
	ret := make([]rune, 0, 100)
	i := 0
	for s.Scan() {
		t := s.Text()
		if i == 0 {
			i++
			continue
		}
		ret = append(ret, []rune(t)...)
		ret = append(ret, []rune("\n")...)
		if t == "EOS" {
			break
		}
		i++
	}
	out := ""
	for i := range ret {
		out += string(ret[i])
	}
	//log.Debugf("JumanKnp.Execute: %v;", out)
	return ret
}

type JumanKnpVerb struct {
	Surface []rune
	Kana    []rune
}

//func (jk *JumanKnp) GetVerbs(domain []rune) []JumanKnpVerb {
//	ret := make([]JumanKnpVerb, 0)
//	scanner := bufio.NewScanner(jk.Imis)
//	for scanner.Scan() {
//		t := scanner.Text()
//		s := strings.Split(t, " ")
//		if strings.HasPrefix(s[0], "ドメイン：") {
//			d := strings.Split(strings.Split(s[0], "：")[1], ";")
//			matched := false
//			for i := range d {
//				if string(domain) == d[i] {
//					// match domain
//					matched = true
//					break
//				}
//			}
//			if matched {
//				c := strings.Split(s[1], ":")
//				sl := strings.Split(c[len(c)-1], "/")
//				ret = append(ret, JumanKnpVerb{Surface: []rune(sl[0]), Kana: []rune(sl[1])})
//			}
//		}
//	}
//	return ret
//}

// Origin : 辞書形（原形）を取得する
// text: 入力
// return: 辞書形
func (jk *JumanKnp) Origin(text []rune) []rune {
	origin := []rune("")
	if len(text) >= 2 && runes.Index([]rune("する;ます"), text[len(text)-2:], 0) >= 0 {
		// 末尾が「する」であれば2文字削る
		origin = append(origin, text[:len(text)-2]...)
		//log.Debugf("%v suffix is する, origin is %v", string(text), string(origin))
	} else {
		origin = append(origin, text[:len(text)-1]...)
		//log.Debugf("%v suffix is not する, origin is %v", string(text), string(origin))
	}
	return origin
}

// GetInflection : 指定された単語から希望する語形変化を取得する
// text: 動詞の表層の原形
// form: 活用形の名前
// return: 活用型の名前，語形変化した文字列，可否
func (jk *JumanKnp) Inflection(text []rune, form []rune) ([]rune, []rune, bool) {
	var o bool
	var f map[string]string
	spacet := []rune(" ")
	att := []rune("@")
	eost := []rune("EOS")
	lft := []rune("\n")
	itype := []rune("")
	if itype, o = jk.inflectionTypeCache[string(text)]; !o {
		j := jk.Execute(text, false)
		for _, line := range runes.Split(j, lft) {
			if runes.Compare(line[0:1], att) || runes.Compare(line[0:3], eost) {
				continue
			}

			// 「押さえ込む」 = 「押さえ」「込む」で、活用は「込む」なので上書き
			out := runes.Split(line, spacet)
			if NewPart(out[3]).IsFlection() {

			}
			itype = out[7]
		}
		jk.inflectionTypeCache[string(text)] = itype
		if runes.Compare(itype, []rune("")) {
			return nil, nil, false
		}
	}

	if f, o = jk.InflectionDB[string(itype)]; o {
		if a, k := f[string(form)]; k {
			if a == "*" {
				// アスタリスクは，データベースにないことを表すようだ
				return itype, nil, false
			}
			inf := make([]rune, 0)
			inf = append(jk.Origin(text), []rune(a)...)
			//log.WithFields(log.Fields{
			//	"Text":       string(text),
			//	"Type":       string(itype),
			//	"Form":       string(form),
			//	"Inflection": string(inf),
			//	"Suffix":     a}).Debug("found inflection")
			return itype, inf, true
		}
		return itype, nil, false
	}
	return itype, nil, false
}

// InflectionPolite : 活用する語(動詞，形容詞，形容動詞)を丁寧にした語を取得する
// text : 活用する語の表層の原形
// past: true: 過去, false: 現在
// return : 丁寧語にした語, 可否
func (jk *JumanKnp) InflectionPolite(text []rune, form []rune) ([]rune, bool) {
	log.Debugf("text: %v, form: %v", string(text), string(form))
	_, infl, flag := jk.Inflection(text, []rune("基本連用形"))
	if !flag {
		infl = jk.Origin(text)
	}
	var polinfl []rune
	_, polinfl, flag = jk.Inflection([]rune("ます"), form)
	if flag {
		return append(infl, polinfl...), true
	} else {
		log.Warnf("JumanKnp.InflectionPolite: could not make polite: %v (%v)",
			string(text), string(form))
	}
	//} else {
	//log.Warnf("JumanKnp.InflectionPolite: could not get inflection: %v (%v)",
	//	string(text), string(form))

	//if flag2 {
	//	log.Fatalf("%v", string(infl2))
	//} else {
	//	log.Fatalf("JumanKnp.InflectionPolite: could not found any inflection of polite: %v (%v)",
	//		string(text), string(form))
	//}
	//}
	return nil, false
}

// IsPast : 過去形かどうか取得する
// form : 活用形の名前
func (jk *JumanKnp) IsPast(form []rune) bool {
	return runes.Compare(form[len(form)-2:], []rune("タ形"))
}

// IsInflectionPolite : 活用する語が丁寧語かどうか取得する
func (jk *JumanKnp) IsInflectionPolite(form []rune) bool {
	return runes.Compare(form[:3], []rune("デス列"))
}

func (jk *JumanKnp) GetKana(text []rune) []rune {
	commentt := []rune("#")
	phraset := []rune("*")
	bphraset := []rune("+")
	spacet := []rune(" ")
	att := []rune("@")
	lft := []rune("\n")
	out := jk.Execute(text, false)
	t := runes.Split(out, lft)
	var ret []rune
	for i := range t {
		//fmt.Println(string(t[i]))
		if !(runes.Compare(t[i][0:1], commentt) ||
			runes.Compare(t[i][0:1], phraset) ||
			runes.Compare(t[i][0:1], bphraset) ||
			runes.Compare(t[i][0:1], att)) {
			p := runes.Split(t[i], spacet)
			if len(p) > 1 {
				ret = append(ret, p[1]...)
			}
		}
	}
	return ret
}

func (jk *JumanKnp) ReadInflection() error {
	file, err := os.Open(jk.Options.JumanDirectory + "/dic/JUMAN.katuyou")
	if err != nil {
		return err
	}
	jk.InflectionDB = map[string]map[string]string{}
	spacet := []rune(" ")
	tabt := []rune("\t")
	semict := []rune(";")
	braot := []rune("(")
	bract := []rune(")")
	scanner := bufio.NewScanner(file)
	bracket := 0
	cinf := []rune("")
	inf := 0
	cform := []rune("")
	form := 0
	surface := 0
	waitsurface := false
	accept := false
	bt := []rune("")
	for scanner.Scan() {
		s := scanner.Text()
		t := []rune(s)
		i := 0
		//log.WithFields(log.Fields{"inf": inf, "i": i, "bt": string(bt)}).Debug()
		if i == 0 {
			if inf != 0 {
				cinf = bt[inf:]
				inf = 0
				//log.WithFields(log.Fields{"cinf": string(cinf)}).Debug()
				jk.InflectionDB[string(cinf)] = map[string]string{}
			}
			if form != 0 {
				return errors.New("form must close immediately")
				//cform = bt[form:]
				//form = 0
			}
		}
		for _, c := range t {

			if c == semict[0] {
				// comment
				//log.Debug("comment")
				break
			}
			if c == braot[0] {
				bracket++
				if bracket == 1 {
					inf = i + 1
				} else if bracket == 3 {
					form = i + 1
					surface = 0
				}
			}
			if waitsurface {
				surface = i
				if c == spacet[0] || c == tabt[0] {
					//log.WithFields(log.Fields{"i": i, "c": string(c)}).Debug("wait surface")
					i++
					continue
				} else {
					//log.WithFields(log.Fields{"i": i, "c": string(c)}).Debug("read surface")
					waitsurface = false
				}
			}
			if (c == spacet[0] || c == tabt[0] || c == bract[0]) && !accept {
				if bracket == 1 {
					if inf != 0 {
						//log.WithFields(log.Fields{"c": string(c), "inf": inf, "i": i}).Debug()
						// 「母音動詞」を受理
						cinf = t[inf:i]
						inf = 0
						jk.InflectionDB[string(cinf)] = map[string]string{}
					}
				} else if bracket == 3 {
					if surface == 0 {
						if waitsurface == false {
							if form == 0 {
								form = i
							} else {
								cform = t[form:i]
								form = 0
								waitsurface = true
								// 「語幹」を受理
								//log.WithFields(log.Fields{"cform": string(cform)}).Debug()
							}
						}
					} else {
						// 「だろう」を受理
						jk.InflectionDB[string(cinf)][string(cform)] = string(t[surface:i])
						//log.WithFields(log.Fields{
						//	"cinf":    string(cinf),
						//	"cform":   string(cform),
						//	"surface": string(t[surface:i])}).Debug()
						surface = 0
						waitsurface = false
						accept = true
					}
				}
			}

			if c == bract[0] {
				bracket--
				accept = false
			}
			i++
		}
		bt = make([]rune, len(t))
		copy(bt, t)
		//for i := range t {
		//	bt[i] = t[i]
		//}
		//log.WithFields(log.Fields{"bt": string(bt), "t": string(t), "len(bt)": len(bt), "len(t)": len(t)}).Debug("next")
	}
	return nil
}

// 表層形 読み 見出し語 品詞大分類 品詞大分類ID 品詞細分類 品詞細分類ID 活用型 活用型ID 活用形 活用形ID 意味情報
