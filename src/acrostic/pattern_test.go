package acrostic

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func IntToBoolMat(matint [][]int) [][]bool {
	mat := make([][]bool, len(matint))
	for i := range matint {
		mat[i] = make([]bool, len(matint[i]))
		for k := range matint[i] {
			if matint[i][k] == 1 {
				mat[i][k] = true
			} else {
				mat[i][k] = false
			}
		}
	}
	return mat
}

func AssertMat(t *testing.T, expected [][]int, actual [][]int) bool {
	if len(expected) != len(actual) {
		t.Errorf("length mismatch, want %v, but returned %v", len(expected), len(actual))
		return false
	}
	for i := range actual {
		//fmt.Printf("%v: %v\n", i, actual[i])
		if len(expected[i]) != len(actual[i]) {
			t.Errorf("length mismatch")
			return false
		}
		for k := range actual[i] {
			if expected[i][k] != actual[i][k] {
				t.Errorf("want [%v %v] = %v, but returned %v", i, k, expected[i][k], actual[i][k])
				return false
			}
		}
	}
	fmt.Println("AssertMat ok")
	return true
}

func TCR1(t *testing.T) {
	row := []int{3, -1, -1, -1}
	matint := [][]int{
		[]int{0, 0, 0, 1},
		[]int{0, 0, 0, 1},
		[]int{0, 0, 0, 1},
		[]int{0, 0, 0, 0}}

	routesexpected := [][]int{
		[]int{3, 0, -1, -1},
		[]int{3, 1, -1, -1},
		[]int{3, 2, -1, -1}}

	mat := IntToBoolMat(matint)
	routes := calculateRoutes(row, 0, mat)
	AssertMat(t, routesexpected, routes)
}

// 　航空自衛隊の救難ヘリコプターが静岡県浜松市沖に墜落し、
// 乗員４人が行方不明となっている事故で、空自は３日、現場周辺の海底から
// 乗員１人の遺体を引き揚げたと発表した。
func TCR2(t *testing.T) {
	//matint := [][]int{
	//	[]int{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//}
	//submats := []int{4, 6}

	//matmainint := [][]int{
	//	//    0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17
	//	[]int{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //  0
	//	[]int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //  1
	//	[]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //  2
	//	[]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //  3
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //  4
	//	[]int{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //  5
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //  6
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //  7
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0}, //  8
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}, //  9
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}, // 10
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, // 11
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}, // 12
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0}, // 13
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, // 14
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}, // 15
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, // 16
	//	[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 17
	//}

	//rows := [][]int{
	//	[]int{4, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	//	[]int{17, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	//}

	//routesexpected := [][]int{
	//	//    0  1  2  3   4   5   6   7   8   9  10  11  12  13  14  15  16  17
	//	[]int{4, 2, 1, 0, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	//	[]int{4, 3, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	//	[]int{17, 16, 8, 6, 5, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	//}

}

func TestCalculateRoute(t *testing.T) {
	TCR1(t)
}

func TCP1(t *testing.T) {
	fmt.Println("TCP1")
	routes := [][]int{
		[]int{2, 0, -1},
		[]int{2, 1, -1},
	}
	expected := [][]int{
		[]int{2, 0, 1},
		[]int{2, 1, 0},
	}

	ret := calculatePattern(routes)

	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("unmatch slice, expected %v, but returned %v", expected, ret)
	}

}

func TCP2(t *testing.T) {
	fmt.Println("TCP2")
	//assert := assert.New(t)
	routes := [][]int{
		[]int{2, 1, 0},
	}
	expected := [][]int{
		[]int{2, 1, 0},
	}

	ret := calculatePattern(routes)
	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("unmatch slice, expected %v, but returned %v", expected, ret)
	}
}

func TCP3(t *testing.T) {
	fmt.Println("TCP3")
	routes := [][]int{
		//    0  1  2  3   4   5   6   7   8   9  10  11  12  13  14  15  16  17
		[]int{4, 2, 1, 0, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		[]int{4, 3, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		[]int{17, 16, 8, 6, 5, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		[]int{17, 16, 8, 7, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		[]int{17, 16, 9, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		[]int{17, 16, 12, 11, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		[]int{17, 16, 15, 14, 13, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		[]int{17, 10, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	}

	ret := calculatePattern(routes)
	f, err := os.Open(os.Getenv("HOME") + "/labo/acrostic/test/news0-pattern")
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer f.Close()
	expected := make([][]int, 0, len(ret))
	s := bufio.NewScanner(f)
	for s.Scan() {
		row := make([]int, 0)
		sp := strings.Split(s.Text(), " ")
		for i := range sp {
			if sp[i] != "" {
				tt, err := strconv.Atoi(sp[i])
				if err != nil {
					t.Error(err.Error())
					return
				}
				row = append(row, tt)
			}
		}
		expected = append(expected, row)
	}
	AssertMat(t, expected, ret)
	//f, _ := os.Create(os.Getenv("HOME") + "/labo/acrostic/test/news0-pattern")
	//defer f.Close()
	//w := bufio.NewWriter(f)
	//for i := range ret {
	//	fmt.Printf("%v\n", ret[i])
	//	o := ""
	//	for k := range ret[i] {
	//		o += strconv.Itoa(ret[i][k]) + " "
	//	}
	//	w.WriteString(o + "\n")
	//}
	//w.Flush()
}

func TestCalculatePattern(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		DisableSorting:   false,
		QuoteEmptyFields: true,
		ForceColors:      true,
		FullTimestamp:    false,
	})
	TCP1(t)
	TCP2(t)
	TCP3(t)
}
