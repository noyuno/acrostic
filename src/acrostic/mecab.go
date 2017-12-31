package acrostic

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/kr/pty"
)

type MeCab struct {
	MeCabPipe *os.File

	Options *Options
}

func NewMeCab(o *Options) (*MeCab, error) {
	var err error
	ret := new(MeCab)
	ret.Options = o
	c := strings.Split(ret.Options.MeCabCommand, " ")[0]
	err = exec.Command("which", c).Run()
	if err != nil {
		return nil, errors.New("command not found: " + c)
	}

	j := exec.Command("sh", "-c", ret.Options.MeCabCommand)
	ret.MeCabPipe, err = pty.Start(j)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (m *MeCab) GetKana(text []rune) []rune {
	m.MeCabPipe.Write([]byte(string(text) + "\n"))
	s := bufio.NewScanner(m.MeCabPipe)
	o := ""
	i := 0
	for s.Scan() {
		if i == 0 {
			i++
			continue
		}
		o += s.Text()
		break
	}
	return []rune(o)
}
