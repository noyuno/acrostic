package acrostic

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/kr/pty"
	log "github.com/sirupsen/logrus"
)

type Kakasi struct {
	KakasiPipe *os.File

	Options *Options
}

func NewKakasi(o *Options) (*Kakasi, error) {
	var err error
	ret := new(Kakasi)
	ret.Options = o
	c := strings.Split(ret.Options.KakasiCommand, " ")[0]
	err = exec.Command("which", c).Run()
	if err != nil {
		return nil, errors.New("command not found: " + c)
	}

	j := exec.Command("sh", "-c", ret.Options.KakasiCommand)
	ret.KakasiPipe, err = pty.Start(j)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (k *Kakasi) GetKana(text []rune) []rune {
	k.KakasiPipe.Write([]byte(string(text) + "\n"))
	s := bufio.NewScanner(k.KakasiPipe)
	ret := make([]rune, 0, 10)
	i := 0
	for s.Scan() {
		if i == 0 {
			i++
			continue
		}
		log.Debug(s.Text())
		ret = append(ret, []rune(s.Text())...)
		break
	}
	return ret
}
