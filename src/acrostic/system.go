package acrostic

import (
	"fmt"
	"runtime"

	"github.com/nicksnyder/go-i18n/i18n"
)

var SystemLanguage string

func TypeToContinue() {
	T, _ := i18n.Tfunc(SystemLanguage)
	fmt.Print(T("Type Enter key to continue") + ": ")
	dummy := ""
	fmt.Scanln(&dummy)
}

func Indent(i int) string {
	ret := ""
	for ; i > 0; i-- {
		ret += "  "
	}
	return ret
}

func MemoryInfo() string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return fmt.Sprintf("HeapAlloc: %vM, TotalAlloc: %vM, HeapSys: %vM",
		mem.HeapAlloc/1024/1024,
		mem.TotalAlloc/1024/1024,
		mem.HeapSys/1024/1024)
}

func HeapAlloc() uint64 {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem.HeapAlloc
}

// たぶん「はい」じゃないかな．
func MaybeYes(a string) bool {
	return a == "" || a == "Y" || a == "y" || a == "yes" || a == "YES"
}
