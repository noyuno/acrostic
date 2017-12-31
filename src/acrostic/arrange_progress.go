package acrostic

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type ArrangeProgress struct {
	Options *Options
	Max     []int
	Current []ArrangeProgressItem
	ticker  *time.Ticker
	stop    chan bool
}

type ArrangeProgressItem struct {
	Name    string
	Current []int
	Enable  bool
}

func NewArrangeProgress(
	o *Options,
	bpa []BasicPhrase,
) *ArrangeProgress {
	ret := new(ArrangeProgress)
	ret.Options = o
	ret.Max = make([]int, len(bpa))
	for i := range bpa {
		ret.Max[i] = len(bpa[i].Pattern)
	}
	p := int64(1)
	bpnormal := 0
	bpmax := 0
	// 数え上げは直積である
	i := 0
	for ; i < len(bpa); i++ {
		p *= int64(len(bpa[i].Pattern))
		bpnormal += len(bpa[i].Surface)
		bpmax += bpa[i].PatternMaxLength
	}
	log.Infof("max=%v len=%v p=%v surface=%v-%v",
		ret.Max, len(ret.Max), p, bpnormal, bpmax)
	return ret
}

func (a *ArrangeProgress) Add(name string) int {
	a.Current = append(a.Current, ArrangeProgressItem{
		Name:    name,
		Current: make([]int, len(a.Max)),
		Enable:  true,
	})
	fmt.Printf("\n")
	return len(a.Current) - 1
}

func (a *ArrangeProgress) Remove(id int) {
	a.Current[id].Enable = false
}

func (a *ArrangeProgress) Start() {
	a.ticker = time.NewTicker(500 * time.Millisecond)
	a.stop = make(chan bool)
	go func() {
	loop:
		for {
			select {
			case <-a.ticker.C:
				a.Print()
			case <-a.stop:
				break loop
			}
		}
	}()
}

func (a *ArrangeProgress) Stop() {
	a.ticker.Stop()
	close(a.stop)
}

func (a *ArrangeProgress) Set(id int, stack []int) {
	if a.Options.Progress {
		a.Current[id].Current = stack
	}
}

func (a *ArrangeProgress) Print() {
	//if len(a.Current) == 1 {
	//	fmt.Printf("\x1b[1K\x1b[0G")
	//} else {
	//	fmt.Printf("\x1b[0G", len(a.Current)-1)
	//	// \x1b[%vF
	//}
	fmt.Printf("\x1b[%vA", len(a.Current))
	for i := range a.Current {
		if a.Current[i].Enable {
			fmt.Printf("\x1b[2K%v %v\n", a.Current[i].Name, a.Current[i].Current)
			//if i+1 < len(a.Current) {
			//} else {

			//	fmt.Printf("%v %v", a.Current[i].Name, a.Current[i].Current)
			//}
		}
	}
}
