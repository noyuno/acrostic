package acrostic

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Paraphrase struct {
	Options  *Options
	database map[string]string
}

func NewParaphrase(o *Options) (*Paraphrase, error) {
	ret := new(Paraphrase)
	ret.Options = o
	ret.database = map[string]string{}
	fp, err := os.Open(ret.Options.ParaphraseDatabase)
	if err != nil {
		return nil, errors.New("unable to open paraphrase database: " + ret.Options.ParaphraseDatabase)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	trim := " 　	"
	i := 0
	for scanner.Scan() {
		t := scanner.Text()
		if strings.Trim(t, trim) == "" || strings.HasPrefix(t, "#") {
			continue
		}
		s := strings.Split(t, ",")
		if len(s) == 1 {
			// カンマで分割できないときは，削除とみなす
			a := strings.Trim(s[0], trim)
			ret.database[a] = ""
		} else if len(s) == 2 {
			a := strings.Trim(s[0], trim)
			b := strings.Trim(s[1], trim)
			ret.database[a] = b
		}
		i++
	}
	//log.Debugf("paraphrase %v items loaded", i)
	//for k, v := range ret.database {
	//	log.Debugf("%v -> %v", k, v)
	//}
	return ret, nil
}

func (p *Paraphrase) GetS(t string) (string, bool) {
	v, ok := p.database[t]
	return v, ok
}

func (p *Paraphrase) Get(t []rune) ([]rune, bool) {
	v, ok := p.database[string(t)]
	return []rune(v), ok
}

func (p *Paraphrase) Replace(t []rune) []rune {
	s := string(t)
	for k, v := range p.database {
		s = strings.Replace(s, k, v, -1)
	}
	return []rune(s)
}
