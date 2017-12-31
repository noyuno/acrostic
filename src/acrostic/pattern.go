package acrostic

import (
	"errors"
	"hash/crc32"
	"sort"

	"github.com/noyuno/lgo/algo"
	log "github.com/sirupsen/logrus"
)

// Pattern : 重みつけを持たない有向グラフを同格同士シャッフルするための構造体
// グラフ中に少なくとも1つの頂点が終点でなければならない
type Pattern struct {
	// Matrix : 隣接行列
	Matrix [][]bool
	// Length : 大きさ
	Length int

	// InitialOrder : 初期の並び順
	InitialOrder []int
	// Routes : 経路, 逆順
	//Routes [][]int
	// EndEdges : 終点, Matrix[N][len(Matrix)-1](N=0..len(Matrix)-1)に同じ
	//EndEdges []int

	// Orders : 同格をシャッフルしたときの並び替えパターン, 正順
	Orders [][]int

	// SubIndex : サブパターンの場所のリスト
	// matの中に2, 3, 4行目および10, 11行目にサブパターンがあれば，
	// [[2 3 4] [10 11]]となる
	SubIndex [][]int

	SubPatterns [][][]int

	//OrderTable []OrderItem

	Options *Options
}

// NewPattern : constructor
func NewPattern(o *Options, len int, init []int, mat [][]bool, subindex [][]int) *Pattern {
	ret := new(Pattern)
	ret.Options = o
	ret.Length = len
	ret.InitialOrder = init
	ret.Matrix = mat
	ret.SubIndex = subindex
	return ret
}

// Shuffle : シャッフルする
func (p *Pattern) Shuffle() error {
	// check input
	if p.Length != len(p.Matrix) {
		return errors.New("Length and Matrix length do not match")
	}

	// サブパターンを1つの要素に圧縮した隣接行列を作成する
	// まずはroot行列のインデックスを作成する
	rootmat := make([]int, 0)
	subpatflag := make([]bool, 0)
	compressed := make([]int, len(p.SubIndex))
	revcompressed := map[int]int{}
	for c := range compressed {
		compressed[c] = -1
	}
	// rootmatの逆変換をする配列を作成
	revrootmat := make([]int, len(p.Matrix))
	revsubpatflag := make([]bool, len(p.Matrix))
	// orderからspの変換
	//ordertosp := map[int]int{}

	for i := range p.Matrix {
		found := false
		for ri, r := range p.SubIndex {
			for _, c := range r {
				if i == c {
					// 圧縮
					found = true
					if compressed[ri] == -1 {
						// 最初なので，突っ込む
						//log.Debugf("first of subpattern %v at %v", ri, i)
						compressed[ri] = len(rootmat)
						revcompressed[i] = ri
						rootmat = append(rootmat, i)
						subpatflag = append(subpatflag, true)
					} else {
						//log.Debugf("it has already added by subpattern %v at %v", ri, i)
					}
					revsubpatflag[i] = true
					break
				}
			}
			if found {
				break
			}
		}
		if found == false {
			// サブパターンにないので，圧縮しない
			rootmat = append(rootmat, i)
			subpatflag = append(subpatflag, false)
		}
		revrootmat[i] = len(rootmat) - 1
	}
	// debug
	//if p.SubIndex != nil {
	//	log.Debugf("subpat=%v, compressed=%v, rootmat=%v, subpatflag=%v"+
	//		"revrootmat=%v, revsubpatflag=%v",
	//		p.SubIndex, compressed, rootmat, subpatflag, revrootmat, revsubpatflag)
	//}

	// サブパターン出現の最初に突っ込んだが，このままmatを作成すると，
	// サブパターン最初の要素のみmatに入るため，よくないので，
	// サブパターンのorを取る．
	// rootmatをもとにmatを作成
	// kはとばすので，できないのでは
	// 初期化
	mat := make([][]bool, len(rootmat))
	for i := range mat {
		mat[i] = make([]bool, len(rootmat))
	}
	// 格納
	for i := range p.Matrix {
		for k := range p.Matrix[i] {
			if mat[revrootmat[i]][revrootmat[k]] || p.Matrix[i][k] {
				mat[revrootmat[i]][revrootmat[k]] = true
			}
		}
	}

	// XXX
	log.Debug("mat")
	for _, r := range mat {
		o := ""
		for _, c := range r {
			if c {
				o += "1 "
			} else {
				o += "0 "
			}
		}
		log.Debug(o)
	}

	p.SubPatterns = make([][][]int, len(p.SubIndex))
	for i := range p.SubIndex {
		p.SubPatterns[i] = make([][]int, 0)
		for e := range algo.Permutations(p.SubIndex[i]) {
			p.SubPatterns[i] = append(p.SubPatterns[i], e)
		}
	}
	log.Debug("Pattern.SubPatterns")
	for i := range p.SubPatterns {
		log.Debugf("%v: %v", i, p.SubPatterns[i])
	}

	// rootmatを計算する

	endedges := make([]int, 0)
	routes := make([][]int, 0)
	for i := range rootmat {
		foundedge := false
		for j := range mat[i] {
			if mat[i][j] {
				foundedge = true
				break
			}
		}
		if !foundedge {
			endedges = append(endedges, i)
		}
	}
	log.Debugf("endedges: %v", endedges)
	for _, end := range endedges {
		row := make([]int, len(mat))
		for i := range row {
			row[i] = -1
		}
		row[0] = end
		log.Debugf("row: %v", row)
		routes = append(routes, calculateRoutes(row, 0, mat)...)
	}
	if len(routes) == 0 {
		return errors.New("Pattern.Routes length is 0")
	}
	for i := range routes {
		log.Debugf("route[%v]: %v", i, routes[i])
	}

	rootorders := calculatePattern(routes)

	// 一個ずつ復元は，デバッグしにくいので，一気にやる
	log.Debug("restore")
	//log.Debugf("rootmat = %v, rootorders = %v", rootmat, rootorders)
	realorders := make([][]int, len(rootorders))
	for i := range realorders {
		realorders[i] = make([]int, len(rootorders[i]))
		for o := range realorders[i] {
			realorders[i][o] = rootmat[rootorders[i][o]]
		}
		//log.Debugf("%v -> %v", rootorders[i], realorders[i])
	}
	for _, order := range realorders {
		//log.Debugf("order = %v", order)
		out := make([][]int, 0)
		for _, realo := range order {
			//// 復元
			//log.Debugf("order[o]=%v -> realo=%v", order[o], realo)
			if revsubpatflag[realo] {
				// サブパターンをここに展開する
				pattern := p.SubPatterns[revcompressed[realo]]
				//log.Debugf("%v is subpattern = %v", realo, pattern)
				// outとpatternの順序あり組み合わせ
				out = appendPattern(out, pattern)
			} else {
				// サブパターンではない
				// outに配列が存在していればそれらの末尾に入れるが，
				// そうでなければ，新しく追加する
				updated := false
				for u := range out {
					updated = true
					out[u] = append(out[u], realo)
				}
				if updated == false {
					u := make([]int, 1)
					u[0] = realo
					out = append(out, u)
				}
			}
		}
		//log.Debugf("out   = %v", out)
		p.Orders = append(p.Orders, out...)
	}
	for o := range p.Orders {
		algo.Reverse(p.Orders[o])
	}

	return nil
}

// srcにpatを順序に拘束して追加
// 例：src = [[1 0 3], [1 3 0]], pat = [[2 4] [4 2]]ならば戻り値は
// [[1 0 3 2 4] [1 0 3 4 2] [1 3 0 2 4] [1 3 0 4 2]]となる
func appendPattern(src [][]int, pat [][]int) [][]int {
	o := make([][]int, 0, len(pat))
	for s := range src {
		for p := range pat {
			t := make([]int, 0, len(src[s])+len(pat[p]))
			t = append(t, src[s]...)
			t = append(t, pat[p]...)
			o = append(o, t)
		}
	}
	return o
}

func copyArray(src []int, length int) [][]int {
	ret := make([][]int, 1)
	ret[0] = make([]int, length)
	for i := range ret[0] {
		ret[0][i] = -1
	}
	copy(ret[0], src)
	return ret
}

// calculateRoutes : 経路を計算する再帰関数
// end: 終端から探索済みのルートをたどっている固定長配列．
// 大きさはp.Lengthにしなければならない．
// depth: 終端を0としたときの終端からの深さ
// mat: 行列
// length 長さ
// return: パターン配列の配列（パターン配列は逆順になっている）
// 例えば，始めは f([7], 0)として呼ばれる
// f([7], 0)の中では次にf([7 3], 1)およびf([7 5], 1)が呼ばれる
// これがdepthがp.Lengthになるまで繰り返される
// 最後に，これらの配列を集めたものが帰ってくる．
// 当然，逆順になっているため，利用前に戻さなければならないだろう．
func calculateRoutes(end []int, depth int, mat [][]bool) [][]int {
	length := len(end)
	var ret [][]int
	if depth+1 >= length {
		log.Debugf("calculateRoutes: end, %v", end)
		return copyArray(end, length)
	}
	e := end[depth]
	for i := 0; i < length; i++ {
		if mat[i][e] {
			// check whether already exists
			found := false
			for ci, c := range end {
				if c == i {
					found = true
					log.Debugf("calculateRoutes: already exists %v in %v at %v. break.", i, end, ci)
					break
				}
			}
			if found {
				break
			}
			newend := make([]int, length)
			copy(newend, end)
			newend[depth+1] = i
			fret := calculateRoutes(newend, depth+1, mat)
			if fret == nil {
				log.Fatal("returned null")
				return nil
			}
			ret = append(ret, fret...)
		}
	}

	if len(ret) == 0 {
		return copyArray(end, length)
	}
	return ret
}

func calculatePattern(routes [][]int) [][]int {
	ret := calculatePatternBase(routes)
	// 出力チェック
	// 数字がすべて含まれているか
	index := make([]int, 0)
	for i := range routes {
		for k := range routes[i] {
			if routes[i][k] != -1 && !algo.Contains(index, routes[i][k]) {
				index = append(index, routes[i][k])
			}
		}
	}
	newret := make([][]int, 0, len(ret))
	for i := range ret {
		findex := true
		for m := range index {
			if !algo.Contains(ret[i], index[m]) {
				findex = false
				break
			}
		}
		if findex {
			// ret[i]にすべてのindexが含まれている
			newret = append(newret, ret[i])
		}
	}
	return newret
}

// calculatePattern : パターン（順列）を計算する
// depth : 末尾を0とした末尾から始まる深さ(Pattern.EndEdgesの要素)
// routes: ルート
// return : パターン
func calculatePatternBase(routes [][]int) [][]int {
	var ret [][]int

	// mは同じ深さの同じ数字ごとに集められたマップ
	m := map[int][]int{}
	keys := []int{}
	//log.Debugf("routes=%v", routes)
	for i := range routes {
		m[routes[i][0]] = append(m[routes[i][0]], i)
		//log.Debugf("i = %v, depth = %v, routes[i][depth] = %v, m[routes[i][depth]] = %v",
		//	i, depth, routes[i][depth], m[routes[i][depth]])
	}
	for i := range m {
		keys = append(keys, i)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	first := true
	for _, key := range keys {
		if first {
			first = false
		}
		ret = append(ret, []int{key})
		// 再帰の戻り値
		recurseret := map[int][][]int{}
		if len(m[key]) >= 2 {
			// 17 16 8
			// 17 16 8
			// 17 16 9 のように，次の列にまた同じ数字が現れるときは再帰
			found := map[int][]int{}
			for i := range m[key] {
				if len(routes[m[key][i]]) >= 2 {
					if _, ok := found[routes[m[key][i]][1]]; !ok {
						found[routes[m[key][i]][1]] = []int{}
					}
					found[routes[m[key][i]][1]] = append(found[routes[m[key][i]][1]], m[key][i])
				}
			}
			for fk, fv := range found {
				//log.Debugf("fk=%v, fv=%v", fk, fv)
				if len(fv) >= 2 {
					newroutes := make([][]int, len(fv))
					for i := range newroutes {
						//log.Debugf("m[key]=%v, fv[i]=%v", m[key], fv[i])
						newroutes[i] = make([]int, len(routes[fv[i]]))
						for k := 0; k+1 < len(routes[fv[i]]); k++ {
							newroutes[i][k] = routes[fv[i]][k+1]
						}
					}

					r := calculatePatternBase(newroutes)
					//log.Debugf("append to recurseret key: %v, value: %v", fk, r)
					recurseret[fk] = r
				}
			}
		}
		//log.Debugf("recurseret=%v", recurseret)
		newret := [][]int{}
		t := [][]int{}
		for i := range ret {
			found := false
			for k := range ret[i] {
				if algo.Contains(m[key], ret[i][k]) {
					found = true
					break
				}
			}
			if found == false {
				t = append(t, ret[i])
			} else {
				newret = append(t, ret[i])
			}
		}
		ret = newret
		//log.Debugf("t=%v", t)
		for e := range algo.Permutations(m[key]) {
			//log.Debugf("e = %v", e)

			o := make([][]int, len(t))
			for i := range t {
				o[i] = make([]int, len(t[i]))
				copy(o[i], t[i])
			}
			expanded := map[int]bool{}
			// recurseretに入っていれば，それを展開
			for _, ee := range e {
				//log.Debugf("routes[ee][1]=%v, recurseret=%v", routes[ee][1], recurseret)
				ok := false
				var v [][]int
				if len(routes[ee]) > 1 {
					v, ok = recurseret[routes[ee][1]]
				}
				if ok {
					if _, oke := expanded[routes[ee][1]]; !oke {
						// 展開
						for ti := range o {
							//log.Debugf("expand %v=%v into o[%v]=%v", ee, recurseret[routes[ee][1]], ti, o[ti])
							orig := make([]int, len(o[ti]))
							copy(orig, o[ti])
							for i := range v {
								if i == 0 {
									// 追加
									o[ti] = append(o[ti], v[i]...)
								} else {
									//if len(ret)+len(o) < limit {
									//log.Debugf("len(ret)(%v)+len(o)(%v)=%v has not exceeded limit=%v",
									//	len(ret), len(o), len(ret)+len(o), limit)
									// 新規
									tt := []int{}
									tt = append(tt, orig...)
									tt = append(tt, v[i]...)
									//log.Debugf("tt=%v", tt)
									o = append(o, tt)
									//} else {
									//log.Debugf("len(ret)(%v)+len(o)(%v)=%v has exceeded limit=%v",
									//	len(ret), len(o), len(ret)+len(o), limit)
									//	break
									//}
								}
							}
						}
						expanded[routes[ee][1]] = true
					} else {
						// do not anything
						//log.Debugf("ignore %v=%v", ee, recurseret[routes[ee][1]])
					}
				} else {
					// 追加
					//log.Debugf("key=%v, m[key]=%v, t=%v", key, m[key], t)
					// -1になるまでまたは終端まで
					cookie := []int{}
					for i := 0; i < len(routes[ee]); i++ {
						if routes[ee][i] == -1 {
							break
						}
						cookie = append(cookie, routes[ee][i])
					}
					for ti := range o {
						//log.Debugf("append %v=%v into o[%v]=%v", ee, cookie, ti, o[ti])
						if algo.Contains(o[ti], cookie[0]) {
							o[ti] = append(o[ti], cookie[1:]...)
						} else {
							o[ti] = append(o[ti], cookie...)
						}
					}
				}
			}
			//log.Debugf("o=%v", o)
			ret = append(ret, o...)
		}
	}
	//log.Debugf("ret = %v", ret)
	//return ret
	// 重複防止
	newret := make([][]int, 0)
	var intsize uint
	intsize = 32 << (^uint(0) >> 63)
	crc32q := crc32.MakeTable(0xD5828281)
	dup := map[uint32]bool{}
	for i := range ret {
		b := make([]byte, intsize*uint(len(ret[i])))
		var k uint
		for k = 0; k < uint(len(ret[i])); k++ {
			var a uint
			for a = 0; a < intsize; a++ {
				b[k*intsize+a] = byte(ret[i][k] >> (a * 8))
			}
		}
		cs := crc32.Checksum(b, crc32q)
		if _, ok := dup[cs]; !ok {
			newret = append(newret, ret[i])
			dup[cs] = true
		}
	}
	//log.Debugf("ret=%v, newret=%v", len(ret), len(newret))
	return newret
}
