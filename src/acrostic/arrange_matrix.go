package acrostic

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"

	"github.com/noyuno/lgo/algo"
	"github.com/noyuno/lgo/runes"
	log "github.com/sirupsen/logrus"
)

type ArrangeMatrix struct {
	// 文パターンの番号
	Number int
	// 検索したい文字列
	Keyword []rune
	// キーワードの番号
	KeywordNumber int
	// 検索したい文字列の位置（文字）
	KeywordIndex int
	// 基本句列
	BasicPhrases []BasicPhrase
	// これから検索される基本句の位置
	BasicPhraseIndex int
	Surface          []rune
	// キーワードの探索が終了したかどうか
	FinishedSearch bool
	// オプション
	Options *Options
	// 行列
	Matrix [][]rune

	NewLine bool
	// 行列の位置
	MatrixIndex []int
	// 行列のとりうる最大位置
	MatrixIndexMax []int
	// 縦読み行列
	KeywordEnd []int
	// 縦読み列が設定されているかどうか
	IsTargeting bool
	// うまくいったやつ
	MatrixResult []ArrangeMatrixResult
	// CPUの数
	NumCPU int
	// スキップした基本句の文字数
	//SkipLength map[int]bool
	// SkipLengthのMutex
	//SkipLengthMutex sync.RWMutex

	// Writer
	Writer *ArrangeWriter
	// WriterのMutex
	WriterMutex sync.RWMutex

	// 並列処理を無効化
	DisableParallel bool

	// wipe outした数
	WipedLength int

	//SearchKeyword bool

	// 探索したテキストの文字数
	TextIndex int

	// 進めているパターンのスタック
	PatternStack []int
	// Bパターンに進んだかどうかを記録
	BranchStack []int

	Progress *ArrangeProgress

	ProgressID int

	Width int

	// 親
	Parent *ArrangeMatrix
}

type ArrangeMatrixResult struct {
	Matrix      [][]rune
	MatrixIndex []int
	KeywordEnd  []int
	//Surface     []rune
	PatternStack []int
	BranchStack  []int
}

func NewArrangeMatrix(o *Options,
	keyword []rune,
	keywordNumber int,
	keywordindex int,
	s []BasicPhrase,
	bpindex int,
	textindex int,
	number int,
	mat [][]rune,
	matpos []int,
	writer *ArrangeWriter,
	progress *ArrangeProgress,
	progressid int,
	stack []int,
	bstack []int,
	width int,
) (*ArrangeMatrix, error) {
	ret := new(ArrangeMatrix)
	ret.Options = o
	if keywordindex >= len(keyword) {
		return nil, fmt.Errorf("keywordindex %v >= len(keyword) %v", keywordindex, len(keyword))
	}
	ret.Keyword = keyword
	ret.KeywordIndex = keywordindex
	if bpindex >= len(s) {
		return nil, fmt.Errorf("bpindex %v >= len(s) %v", bpindex, len(s))
	}
	ret.BasicPhrases = s
	ret.Surface = make([]rune, 0, 100)
	for i := range ret.BasicPhrases {
		if ret.BasicPhrases[i].NewLine && i != 0 {
			ret.Surface = append(ret.Surface, []rune("\\n")...)
		}
		ret.Surface = append(ret.Surface, ret.BasicPhrases[i].Surface...)
	}
	ret.BasicPhraseIndex = bpindex
	ret.TextIndex = textindex
	if len(matpos) != 2 {
		return nil, fmt.Errorf("len(matpos) != 2")
	}
	ret.Number = number
	ret.Matrix = mat
	ret.MatrixIndex = matpos
	ret.MatrixIndexMax = []int{ret.Options.Height, width}
	ret.MatrixResult = make([]ArrangeMatrixResult, 0)
	ret.KeywordEnd = make([]int, 2)
	ret.NumCPU = runtime.NumCPU()
	//ret.SkipLength = map[int]bool{}
	//ret.SkipLengthMutex = sync.RWMutex{}
	ret.Writer = writer
	ret.WriterMutex = sync.RWMutex{}
	ret.Progress = progress
	ret.ProgressID = progressid
	//ret.SearchKeyword = true
	ret.PatternStack = stack
	ret.BranchStack = bstack
	ret.Width = width
	return ret, nil
}

func (m *ArrangeMatrix) expectedLine() int {
	maxremain := 0
	lflen := 0
	for i := m.BasicPhraseIndex; i < len(m.BasicPhrases); i++ {
		maxremain += m.BasicPhrases[i].PatternMaxLength
		if m.BasicPhrases[i].NewLine {
			lflen++
		}
	}
	expectedline := int(math.Ceil((float64(maxremain))/float64(m.Width) + float64(lflen)))
	return expectedline
}

// SearchContext : 検索
// indent: インデント
// k: キーワード
// pi: パターンの番号
// p: パターン
// return int: 刈った数(成功した数ではない．)
// return error: エラー
func (m *ArrangeMatrix) SearchContext(
	indent string,
	k []rune,
	pi int,
	p []rune,
	progressid int,
	oldstack []int,
	foundnum int,
) (int, error) {
	// oldstack よりも小さければ終了
	//if ArrangeForward(m.Progress.Get(), append(m.PatternStack, pi)) == false {
	//	log.Fatalf("old: %v, new: %v", oldstack, m.PatternStack)
	//}
	//if foundnum == 0 {
	//	if m.Progress.Forward(append(m.PatternStack, pi)) == false {
	//		log.Fatalf(indent+"bpi: %v, pi: %v, old: %v(%v), new: %v",
	//			m.BasicPhraseIndex, pi, m.Progress.LastStack, m.Progress.LastValue,
	//			append(m.PatternStack, pi))
	//		//log.Fatalf("old: %v, new: %v", m.Progress.LastStack, append(m.PatternStack, pi))
	//	}
	//}
	// すでに処理していれば終了
	//if v, ok := progress.IsProcessed(append(m.PatternStack, pi), m.KeywordIndex); ok {
	//	if v == m.KeywordIndex {
	//		log.Fatalf("already processed: %v, %v, prev %v current %v",
	//			m.PatternStack, pi, v, m.KeywordIndex)
	//	} else {
	//		//log.Warnf("already processed: %v, %v, prev %v current %v",
	//		//	m.PatternStack, pi, v, m.KeywordIndex)
	//	}
	//	//return 0, errors.New("already processed")
	//	//return 0, fmt.Errorf("already processed: %v, %v, prev %v current %v",
	//	//	m.PatternStack, pi, v, m.KeywordIndex)
	//	return 0, nil
	//}
	//progress.MarkProcessed(append(m.PatternStack, pi), m.KeywordIndex)

	//log.Debugf("BasicPhraseIndex: %v", m.BasicPhraseIndex)
	//if m.FinishedSearch {
	//	log.Debugf("KeywordEnd: %v", m.KeywordEnd)
	//}

	// もし，pにkが入っていてもこれを飛ばすルートが存在する．
	// 例: キーワード = 「たこ」
	// ・・た・・た
	// ・・・・・こ
	// ・・・こ・

	// 追記：上は「縦読み」ではない
	// 追記：MatchLength = falseであれば上を許容する

	//log.Debugf(indent+"KeywordEnd: %v", m.KeywordEnd)

	//under := m.under()
	//if under != 0 {
	//	log.Debugf("under: %v", under)
	//}

	// 枝刈り：残りキーワードの文字数が残り行数よりも大きければ終了
	// 切り上げ
	expline := m.expectedLine()
	if len(m.Keyword)-m.KeywordIndex-1 > expline {
		//log.Debugf(indent+"pruning k=%v, t=%v",
		//	len(m.Keyword)-m.KeywordIndex-1, expline)
		return 0, nil
	}

	// 枝刈り：MatchLengthで1行目にキーワードがなければ終了
	// MatrixIndexがKeywordEndの列を超えていれば終了
	if m.Options.MatchLength {
		if m.MatrixIndex[0] > 0 && m.IsTargeting == false {
			//log.Debugf(indent + "pruning MatchLength: not appearing keyword at line 0")
			return 0, nil
		}
		if m.MatrixIndex[0] > m.KeywordEnd[0] && m.MatrixIndex[1] >= m.KeywordEnd[1] {
			//log.Debugf(indent + "pruning MatchLength: MatrixIndex exceeded keyword column")
			return 0, nil
		}
	}

	//log.Debugf(indent+"bpindex: %v, surface: %v, pattern: %v:%v, index: %v, kindex: %v, stack: %v",
	//	m.BasicPhraseIndex, string(m.BasicPhrases[m.BasicPhraseIndex].Surface), pi, string(p),
	//	m.MatrixIndex, m.KeywordIndex, m.PatternStack)
	// 検索
	r := -1
	if m.FinishedSearch == false {
		r = runes.Index(p, k, 0)
	}
	foundn := 0

	//log.Debugf(indent+"p=%v, k=%v, r=%v", string(p), string(k), r)
	//PrintMatrix(m.Matrix, m.MatrixIndex)
	for r != -1 {
		//log.Debugf(indent+"found %v in pattern[%v] = %v at %v",
		//	string(k), pi, string(p), r)
		// 文字列中に検索文字が見つかったとしても，同じ列かつ次の行でなければならない
		y := m.MatrixIndex[0]
		//x := m.KeywordEnd[0]
		x := m.MatrixIndex[1] + r
		if m.MatrixIndex[1]+r != 0 && m.BasicPhrases[m.BasicPhraseIndex].NewLine {
			y++
			x = r
			//log.Debugf("newline")
		} else if m.MatrixIndex[1]+r >= m.Width {
			y++
			x = (m.MatrixIndex[1] + r) % m.Width
			//log.Debugf("orikaesiline")
		}
		// 行数不一致ならはじく
		if m.Options.MatchLength {
			if y != m.KeywordIndex {
				// 行数不一致
				//log.Debugf(indent+"begin line mismatch, y: %v, m.KeywordIndex: %v",
				//	y, m.KeywordIndex)
				//return false, nil
				break
			}
		}
		// キーワードが連続していなければはじく
		if m.IsTargeting {
			if m.KeywordEnd[0]+1 != y {
				break
			}
		}
		// 列数不一致ならはじく
		if m.IsTargeting {
			if m.KeywordEnd[1] != x {
				// 列不一致
				//log.Debugf(indent+"column mismatch, KeywordEnd[0]=%v, x=%v",
				//	m.KeywordEnd[0], x)
				break
			}
		}
		//log.Debugf(indent+"y=%v, x=%v, IsTargeting=%v", y, x, m.IsTargeting)
		// copy matrix and append found phrase to this matrix
		//mat := CopyMatrix(m.Matrix)
		mat := make([][]rune, 1)
		mat[0] = make([]rune, m.Width)
		matpos := []int{0, 0}
		matstart := []int{0, 0}
		newline := false
		//surface := make([]rune, len(p))
		copy(matpos, m.MatrixIndex)
		copy(matstart, m.MatrixIndex)
		//PrintMatrix(mat, matpos)
		//copy(surface, p)
		if matpos[1] != 0 && m.BasicPhrases[m.BasicPhraseIndex].NewLine {
			// 改行
			if m.MatrixIndexMax[0] > matpos[0]+1 {
				newline = true
				//log.Debug("new line")
				matpos[0]++
				matpos[1] = 0
				//mat = append(mat, make([]rune, m.Width))
			} else {
				//log.Fatalf("array filled up (lf)")
			}
		}
		//log.Debugf(indent+"A: append '%v' (%v char) to mat", string(p), len(p))
		keywordend := []int{y, x}
		keywordindex := m.KeywordIndex
		matline := matpos[0]
		//log.Debugf("keywordend: %v, keywordindex: %v", keywordend, keywordindex)
		for i := range p {
			//log.Debugf("B [%v %v] %v '%v'", matpos[0], matpos[1], i, string(p[i]))
			if matpos[1] == x && matpos[0] > y {
				// 折り返し
				//log.Debugf("折り返し")
				if len(m.Keyword) > keywordindex+1 {
					keywordend[0]++
					keywordindex++
					//log.Debugf("B: keywordend: %v, keywordindex: %v", keywordend, keywordindex)
					if m.Keyword[keywordindex] != p[i] {
						//PrintMatrix(mat, matpos)
						//log.Debugf("折り返したがキーワードに一致しない: k: %v, m: %v",
						//	string(m.Keyword[keywordindex]), string(p[i]))
						return 0, nil
					}
				} else {
					if m.Options.MatchLength {
						//PrintMatrix(mat, matpos)
						log.Fatal("範囲外")
					}
				}
			}
			mat[matpos[0]-matline][matpos[1]] = p[i]
			flag := 0
			matpos, flag = algo.SliceAdderR(matpos, m.MatrixIndexMax, len(matpos))
			if flag == 3 {
				log.Fatal("array filled up")
			} else if flag == 1 && len(p) > i+1 {
				mat = append(mat, make([]rune, m.Width))
			}
		}
		//PrintMatrix(mat, matpos)
		//TypeToContinue()
		//log.Debugf("matpos=%v", matpos)
		//PrintMatrix(mat, matpos)
		finished := false
		//log.Debugf("KeywordIndex: %v", m.KeywordIndex)
		if keywordindex+1 >= len(m.Keyword) && y+1 >= len(m.Keyword) {
			//log.Debugf(indent + "finished search keyword")
			finished = true
		}

		//log.Debugf(indent+"keywordend: %v", keywordend)

		// BUG: keywordend=[0 0]なのにkeywordindex=1なのはおかしくて，0であるべきだ．
		if keywordend[0] < keywordindex {
			log.Fatalf("Oops! keywordend=%v, but keywordindex=%v", keywordend, keywordindex)
		}
		if m.BasicPhraseIndex+1 >= len(m.BasicPhrases) {
			// BasicPhrase探索終了
			if /*m.FinishedSearch || */ finished {
				accept := false
				// 行数不一致なら不一致
				if m.Options.MatchLength {
					if matpos[1] >= keywordend[1] {
						// 文末端より左側
						if matpos[0] == keywordend[0] {
							accept = true
						} else {
							//log.Debugf(indent+"left  side matpos[0]: %v, ke[0]: %v",
							//	matpos[0], keywordend[0])
						}
					} else {
						// 文末端より右側
						if matpos[0] == keywordend[0]+1 {
							accept = true
						} else {
							//log.Debugf(indent+"right side matpos[0]: %v, ke[0]+1: %v",
							//	matpos[0], keywordend[0]+1)
						}
					}
				} else {
					accept = true
				}
				// 受理
				//PrintMatrix(mat, matpos)
				//log.Debugf("accept: %v", accept)
				//TypeToContinue()
				if accept {
					//log.Debugf(indent + "A: accepted")
					//PrintMatrix(mat, matpos)
					//surface := []rune("")
					//for i := range m.BasicPhrases {
					//	surface = append(surface, m.BasicPhrases[i].Surface...)
					//}
					m.MatrixResult = append(m.MatrixResult, m.makeResult(
						mat, matpos, newline,
						append(m.PatternStack, pi),
						append(m.BranchStack, 0)))
					//m.MatrixResult = append(m.MatrixResult, ArrangeMatrixResult{
					//	Matrix:      mat,
					//	MatrixIndex: matpos,
					//	KeywordEnd:  keywordend,
					//	//Surface:     surface,
					//})
					if m.Options.One {
						return 0, nil
					}
				}
			} else {
				//log.Debugf(indent + "not found")
			}
		} else {
			if finished {
				// キーワードの探索は終了したが，BPの探索は終了していない
				//log.Debugf(indent+"new instance (finished) KeywordIndex=%v, keywordend=%v",
				//	m.KeywordIndex, keywordend)
				newstack := make([]int, len(m.PatternStack)+1)
				copy(newstack, m.PatternStack)
				newstack[len(newstack)-1] = pi
				bstack := make([]int, len(m.BranchStack)+1)
				copy(bstack, m.BranchStack)
				bstack[len(bstack)-1] = 0
				am, err := NewArrangeMatrix(
					m.Options, m.Keyword, m.KeywordNumber,
					keywordindex, m.BasicPhrases,
					m.BasicPhraseIndex+1,
					m.TextIndex+len(p),
					m.Number, mat, matpos, m.Writer,
					m.Progress, progressid, newstack, bstack, m.Width)
				if err != nil {
					return 0, err
				}
				//am.FinishedSearch = finished
				am.IsTargeting = false
				//am.KeywordEnd = []int{matpos[0], keywordend[1]}
				//log.Debugf(indent+"KeywordIndex=%v, finished=%v, keywordend=%v, keyword=%v",
				//	m.KeywordIndex, finished, keywordend, string(m.Keyword))
				am.KeywordEnd = keywordend
				am.DisableParallel = m.DisableParallel
				//am.SearchKeyword = false
				am.FinishedSearch = true
				am.Parent = m
				am.NewLine = newline
				err = am.Search(m.PatternStack, foundn)
				if err != nil {
					return 0, err
				}
				m.MatrixResult = append(m.MatrixResult, am.MatrixResult...)
				m.WipedLength += am.WipedLength
			} else {
				// キーワードの探索およびBPの探索が終わっていない
				//log.Debugf(indent+"A: new instance, KeywordIndex=%v, keywordend=%v, stack=%v",
				//	m.KeywordIndex, keywordend, m.PatternStack)
				newstack := make([]int, len(m.PatternStack)+1)
				copy(newstack, m.PatternStack)
				newstack[len(newstack)-1] = pi
				bstack := make([]int, len(m.BranchStack)+1)
				copy(bstack, m.BranchStack)
				bstack[len(bstack)-1] = 0
				am, err := NewArrangeMatrix(
					m.Options, m.Keyword, m.KeywordNumber,
					keywordindex+1, m.BasicPhrases,
					m.BasicPhraseIndex+1,
					m.TextIndex+len(p),
					m.Number, mat, matpos, m.Writer,
					m.Progress, progressid, newstack, bstack, m.Width)
				if err != nil {
					return 0, err
				}
				//am.FinishedSearch = finished
				am.IsTargeting = true
				am.KeywordEnd = keywordend //[]int{matpos[0], keywordend[1]}
				am.DisableParallel = m.DisableParallel
				am.Parent = m
				am.NewLine = newline
				//log.Debugf(indent+"m=%v, am=%v", m.BasicPhraseIndex, am.BasicPhraseIndex)
				err = am.Search(m.PatternStack, foundn)
				if err != nil {
					return 0, err
				}
				//log.Debugf(indent+"m=%v, am=%v", m.BasicPhraseIndex, am.BasicPhraseIndex)
				m.MatrixResult = append(m.MatrixResult, am.MatrixResult...)
				m.WipedLength += am.WipedLength
				//log.Debugf(indent + "Search end")
			}
			if m.Options.One && len(m.MatrixResult) > 0 {
				return 0, nil
			}
		}
		foundn++
		r = runes.Index(p, k, r+1)
	}

	//if m.Options.SkipSameLength {
	//	m.SkipLengthMutex.RLock()
	//	_, ok := m.SkipLength[len(p)]
	//	m.SkipLengthMutex.RUnlock()
	//	if ok {
	//		//log.Debugf(indent+"not found, the same length(%v) as %v has been processed",
	//		//	len(p), string(p))
	//		//TypeToContinue()
	//		return under, nil
	//	}
	//}
	// skip this phrase
	//log.Debugf(indent+"skip %v length %v at %v", string(b.Surface),
	//	len(p), m.MatrixIndex)
	//mat := CopyMatrix(m.Matrix)
	mat := make([][]rune, 1)
	mat[0] = make([]rune, m.Width)
	matpos := []int{0, 0}
	matstart := []int{0, 0}
	newline := false
	//surface := make([]rune, len(p))
	copy(matpos, m.MatrixIndex)
	copy(matstart, m.MatrixIndex)
	//copy(surface, p)
	if matpos[1] != 0 && m.BasicPhrases[m.BasicPhraseIndex].NewLine {
		// 改行
		if m.MatrixIndexMax[0] > matpos[0]+1 {
			//log.Debug(indent + "new line")
			newline = true
			matpos[0]++
			matpos[1] = 0
			//mat = append(mat, make([]rune, m.Width))
		} else {
			log.Fatalf(indent+"array filled up (lf), MatrixIndexMax: %v, matpos[0]+1: %v",
				m.MatrixIndexMax, matpos[0]+1)
		}
	}
	matline := matpos[0]
	//log.Debugf(indent+"B: append '%v' (%v char) to mat", string(p), len(p))
	for i := range p {
		//log.Debugf("B [%v %v] %v '%v'", matpos[0], matpos[1], i, string(p[i]))
		mat[matpos[0]-matline][matpos[1]] = p[i]
		flag := 0
		matpos, flag = algo.SliceAdderR(matpos, m.MatrixIndexMax, len(matpos))
		if flag == 3 {
			log.Fatal(indent + "array filled up")
		} else if flag == 1 && len(p) > i+1 {
			mat = append(mat, make([]rune, m.Width))
		}
	}
	//PrintMatrix(mat, matpos)
	//if len(k) == m.KeywordIndex+1 {
	//	m.FinishedSearch = true
	//}
	if m.BasicPhraseIndex+1 >= len(m.BasicPhrases) {
		//log.Debugf("end, FinishedSearch=%v", m.FinishedSearch)
		//if m.FinishedSearch {
		accept := false
		// 行数不一致なら不一致
		if m.Options.MatchLength {
			if matpos[1] >= m.KeywordEnd[1] {
				// 文末端より左側
				if matpos[0] == m.KeywordEnd[0] {
					accept = true
				} else {
					//log.Debugf(indent+"left  side matpos[0]: %v, m.KeywordEnd[0]: %v",
					//	matpos[0], m.KeywordEnd[0])
				}
			} else {
				// 文末端より右側
				if matpos[0] == m.KeywordEnd[0]+1 {
					accept = true
				} else {
					//log.Debugf(indent+"right side  matpos[0]: %v, m.KeywordEnd[0]+1: %v",
					//	matpos[0], m.KeywordEnd[0]+1)
				}
			}
		} else {
			accept = true
		}
		if accept && m.FinishedSearch {
			// 受理
			//log.Debugf(indent+"B: accepted, keywordend: %v, FinishedSearch: %v",
			//	[]int{m.KeywordEnd[0], m.KeywordEnd[1]}, m.FinishedSearch)
			//surface := []rune("")
			//for i := range m.BasicPhrases {
			//	surface = append(surface, m.BasicPhrases[i].Surface...)
			//}
			m.MatrixResult = append(m.MatrixResult,
				m.makeResult(mat, matpos, newline,
					append(m.PatternStack, pi),
					append(m.BranchStack, 1)))
			//m.MatrixResult = append(m.MatrixResult, ArrangeMatrixResult{
			//	Matrix:      mat,
			//	MatrixIndex: matpos,
			//	KeywordEnd:  []int{m.KeywordEnd[0], m.KeywordEnd[1]},
			//	//Surface:     surface,
			//})
			if m.Options.One {
				return 0, nil
			}
		}
		//} else {
		//	//log.Debugf(indent + "not found")
		//}
	} else {
		//log.Debugf(indent+"A: call BasicPhraseIndex: %v", m.BasicPhraseIndex)
		//if m.BasicPhraseIndex == 21 {
		//	log.Debugf("21: pi: %v, stack: %v, bstack: %v, kindex: %v, tindex: %v",
		//		pi, m.PatternStack, m.BranchStack, m.KeywordIndex, m.TextIndex)
		//}
		newstack := make([]int, len(m.PatternStack)+1)
		copy(newstack, m.PatternStack)
		newstack[len(newstack)-1] = pi
		bstack := make([]int, len(m.BranchStack)+1)
		copy(bstack, m.BranchStack)
		bstack[len(bstack)-1] = 1
		am, err := NewArrangeMatrix(
			m.Options, m.Keyword, m.KeywordNumber,
			m.KeywordIndex,
			m.BasicPhrases, m.BasicPhraseIndex+1, m.TextIndex+len(p),
			m.Number, mat, matpos, m.Writer,
			m.Progress, progressid, newstack, bstack, m.Width)
		if err != nil {
			return 0, err
		}
		am.FinishedSearch = m.FinishedSearch
		am.KeywordEnd = m.KeywordEnd
		am.IsTargeting = m.IsTargeting
		am.DisableParallel = m.DisableParallel
		am.TextIndex = m.TextIndex
		am.Parent = m
		am.NewLine = newline
		//log.Debugf(indent+"m=%v, am=%v", m.BasicPhraseIndex, am.BasicPhraseIndex)
		am.Search(m.PatternStack, 0)
		//log.Debugf(indent+"m=%v, am=%v", m.BasicPhraseIndex, am.BasicPhraseIndex)
		m.MatrixResult = append(m.MatrixResult, am.MatrixResult...)
		m.WipedLength += am.WipedLength
	}
	//m.SkipLengthMutex.Lock()
	//m.SkipLength[len(p)] = true
	//m.SkipLengthMutex.Unlock()

	// m.MatrixResult がいっぱいだったらwipe outする
	m.WriterMutex.Lock()
	if m.Options.WipeOut && m.Options.WipeOutLength <= len(m.MatrixResult) {
		length := len(m.MatrixResult)
		err := m.WipeOut()
		if err != nil {
			m.WriterMutex.Unlock()
			return 0, err
		}
		m.WipedLength += length
	}
	m.WriterMutex.Unlock()
	return 0, nil
}

func (m *ArrangeMatrix) makeResult(
	in [][]rune, matpos []int, newline bool,
	stack []int, bstack []int) ArrangeMatrixResult {
	//log.Debugf("makeResult")
	parents := make([]*ArrangeMatrix, 0)
	parent := m
	for parent != nil {
		parents = append(parents, parent)
		parent = parent.Parent
	}

	mat := make([][]rune, 0)
	lastpos := 0
	for p := len(parents) - 1; p >= 0; p-- {
		//log.Debugf("parents %v", p)
		for n := range parents[p].Matrix {
			if n == 0 && lastpos != 0 && parents[p].NewLine == false {
				//log.Debugf("continue at %v", lastpos)
				for c := lastpos; c < len(parents[p].Matrix[n]); c++ {
					mat[len(mat)-1][c] = parents[p].Matrix[n][c]
				}
			} else {
				//log.Debugf("new line")
				mat = append(mat, make([]rune, m.Width))
				copy(mat[len(mat)-1], parents[p].Matrix[n])
			}
			//log.Debugf("%v: %v", len(mat)-1, string(mat[len(mat)-1]))
		}
		lastpos = parents[p].MatrixIndex[1]
	}
	for n := range in {
		if n == 0 && lastpos != 0 && newline == false {
			//log.Debugf("in continue at %v", lastpos)
			for c := lastpos; c < len(in[n]); c++ {
				mat[len(mat)-1][c] = in[n][c]
			}
		} else {
			//log.Debugf("in new line")
			mat = append(mat, make([]rune, m.Width))
			copy(mat[len(mat)-1], in[n])
		}
	}

	//PrintMatrix(mat, matpos)

	return ArrangeMatrixResult{
		Matrix:       mat,
		MatrixIndex:  matpos,
		KeywordEnd:   []int{m.KeywordEnd[0], m.KeywordEnd[1]},
		PatternStack: stack,
		BranchStack:  bstack,
	}
}

func (m *ArrangeMatrix) Skip() error {

	return nil
}

func (m *ArrangeMatrix) CheckAfterKeyword() bool {
	// 枝刈り：残りの場所に残りキーワードが順番通りに来なければ終了
	found := true
	// keywordとBasicPhrase iの関連
	kwarray := make([][]int, len(m.Keyword)-m.KeywordIndex)
	for ki := m.KeywordIndex; ki < len(m.Keyword); ki++ {
		kwarray[ki-m.KeywordIndex] = make([]int, 0)
		f := false
		for i := m.BasicPhraseIndex; i < len(m.BasicPhrases); i++ {
			//log.Debugf("i=%v, m.KeywordNumber=%v, ki=%v, Pos[m.KeywordNumber]=%v",
			//	i, m.KeywordNumber, ki, m.BasicPhrases[i].PatternKeywordPos[m.KeywordNumber])
			if m.KeywordNumber < len(m.BasicPhrases[i].PatternKeywordPos) {
				if len(m.BasicPhrases[i].PatternKeywordPos[m.KeywordNumber][ki]) != 0 {
					kwarray[ki-m.KeywordIndex] = append(kwarray[ki-m.KeywordIndex], i)
					//log.Debugf("true: %v", string(m.Keyword[ki]))
					f = true
					// DO NOT BREAK
				}
			} else {
				// FIXME
				f = true
			}
		}
		if f == false {
			found = false
			break
		}
	}
	//log.Debugf("found: %v", found)
	if found == false {
		return false
	}

	//log.Debugf("found, kwarray: %v", kwarray)

	// 順番通りに到達できるか？
	var each func(int) [][]int
	each = func(k int) [][]int {
		r := make([][]int, 0)
		if k == 0 {
			for i := range kwarray[k] {
				r = append(r, []int{kwarray[k][i]})
			}
		} else {
			for i := range kwarray[k] {
				for _, e := range each(k - 1) {
					r = append(r, append(e, kwarray[k][i]))
				}
			}
		}
		return r
	}
	direction := each(len(kwarray) - 1)
	//log.Debugf("direction: %v", direction)

	//direction := make([][]int, 0)
	for i := range direction {
		bad := false
		//log.Debugf("direction[%v]: %v", i, direction[i])
		for k := 0; k+1 < len(direction[i]); k++ {
			if direction[i][k] > direction[i][k+1] {
				bad = true
				break
			}
		}
		if bad == false {
			// 一つでも良い
			//log.Debugf("ok")
			return true
		}
	}
	//log.Debugf("found only bad direction")
	return false
	//if found == false {
	//	//log.Debugf("bp %v: not found keyword '%v' under..", m.BasicPhraseIndex, string(m.Keyword))
	//	return 0, nil
	//} else {
	//	//log.Debugf("bp %v: found keyword '%v' under..", m.BasicPhraseIndex, string(m.Keyword))
	//}
}

func (m *ArrangeMatrix) Search(oldstack []int, foundnum int) error {
	var err error
	//if len(m.BasicPhrases)-m.BasicPhraseIndex > 10 {
	//	log.Debugf("Search: %v", m.PatternStack)
	//}
	m.Progress.Set(m.ProgressID, m.PatternStack)
	if m.FinishedSearch == false && m.CheckAfterKeyword() == false {
		//log.Debugf("not found after this")
		return nil
	}
	if m.Options.Parallel &&
		m.DisableParallel == false &&
		len(m.BasicPhrases[m.BasicPhraseIndex].Pattern) >= m.NumCPU {
		m.DisableParallel = true
		log.Debugf("parallel at %v %v items", m.BasicPhraseIndex, len(m.BasicPhrases[m.BasicPhraseIndex].Pattern))
		_, err = m.SearchParallel(oldstack, foundnum)
	} else {
		_, err = m.SearchNormal(oldstack, foundnum)
	}
	return err
}

type ArrangeMatrixPattern struct {
	Index int
	Text  []rune
}

func (m *ArrangeMatrix) getPattern() chan ArrangeMatrixPattern {
	p := make(chan ArrangeMatrixPattern)
	go func(p chan ArrangeMatrixPattern) {
		defer close(p)
		b := m.BasicPhrases[m.BasicPhraseIndex]
		if m.IsTargeting {
			//log.Debugf("getPattern: IsTargeting: true")
			// 次のキーワード列までの差を求めて，それがPatternをまたぎキーワードに一致
			// する可能性があるのであれば，そのPatternを優先的に採用する
			// （これ以外のPatternは除外せず後回しにする）
			sent := map[int]bool{}
			le := map[int]bool{}
			// 差
			diff := 0
			if b.NewLine {
				diff = m.KeywordEnd[1]
			} else {
				diff = m.KeywordEnd[1] - m.MatrixIndex[1]
				if diff < 0 {
					// 折り返し
					diff += m.Width
				}
			}
			if v, ok := b.PatternKeywordPos[m.KeywordNumber][m.KeywordIndex][diff]; ok {
				for i := range v {
					//log.Debugf("send matched at %v, %v: %v (%v)",
					//	diff, v[i], string(b.Pattern[v[i]]), len(b.Pattern[v[i]]))
					p <- ArrangeMatrixPattern{
						Index: v[i],
						Text:  b.Pattern[v[i]],
					}
					sent[v[i]] = true
					le[len(b.Pattern[v[i]])] = true
				}
			} else {
				//log.Debugf("cannot find Pattern match at %v", diff)
			}
			// 残りの送っていない長さのPatternを送る
			// 同じ残り物はいらない
			for i := range b.Pattern {
				if _, ok := sent[i]; !ok {
					if _, ok := le[len(b.Pattern[i])]; !ok {
						//log.Debugf("send remain %v: %v (%v)",
						//	i, string(b.Pattern[i]), len(b.Pattern[i]))
						p <- ArrangeMatrixPattern{
							Index: i,
							Text:  b.Pattern[i],
						}
						le[len(b.Pattern[i])] = true
					}
				}
			}
		} else {
			for ppi, pp := range b.Pattern {
				p <- ArrangeMatrixPattern{
					Index: ppi,
					Text:  pp,
				}
			}
		}
	}(p)
	return p
}

func (m *ArrangeMatrix) SearchNormal(oldstack []int, foundnum int) (int, error) {
	var err error
	indent := Indent(m.BasicPhraseIndex)
	k := m.Keyword[m.KeywordIndex : m.KeywordIndex+1]
	//b := m.BasicPhrases[m.BasicPhraseIndex]
	//log.Debugf(indent+"search %v into %v", string(k), string(b.Surface))

	for p := range m.getPattern() {
		//log.Debugf("BasicPhrases[%v].Pattern[%v]%v", m.BasicPhraseIndex, pi, string(p))
		_, err = m.SearchContext(indent, k, p.Index, p.Text, m.ProgressID, oldstack, foundnum)
		if err != nil {
			log.Warn(err.Error())
			break
		}
		if m.Options.One && len(m.MatrixResult) > 0 {
			break
		}
		//m.Progress.Set(pi, m.PatternStack, m.BasicPhraseIndex)
	}
	ret := 0

	return ret, nil
}

func (m *ArrangeMatrix) SearchParallel(oldstack []int, foundnum int) (int, error) {
	indent := Indent(m.BasicPhraseIndex)
	k := m.Keyword[m.KeywordIndex : m.KeywordIndex+1]
	b := m.BasicPhrases[m.BasicPhraseIndex]
	//log.Debugf(indent+"search parallel %v into %v", string(k), string(b.Surface))
	var wg sync.WaitGroup
	semaphore := make(chan int, m.NumCPU)
	for pi, p := range b.Pattern {
		wg.Add(1)
		go func(pi int, p []rune) {
			defer wg.Done()
			semaphore <- 1
			progressid := m.Progress.Add("parallel" + strconv.Itoa(pi))
			_, err := m.SearchContext(indent, k, pi, p, progressid, oldstack, foundnum)
			m.Progress.Remove(progressid)
			if err != nil {
				log.Warn(err.Error())
			}
			//progress.Set(pi, m.PatternStack, m.BasicPhraseIndex)
			<-semaphore
		}(pi, p)
		if m.Options.One && len(m.MatrixResult) > 0 {
			break
		}
	}
	wg.Wait()
	ret := 0
	return ret, nil
}

func CopyMatrix(in [][]rune) [][]rune {
	ret := make([][]rune, len(in), len(in)+1)
	for c := range in {
		end := false
		ret[c] = make([]rune, len(in[c]))
		for i := range in[c] {
			if string(in[c][i]) == "" {
				end = true
				break
			} else {
				ret[c][i] = in[c][i]
				//log.Debugf("CopyMatrix: [%v][%v] in: %v, out: %v",
				//	c, i, string(m.Matrix[c][i].Surface), string(ret[c][i].Surface))
			}
		}
		if end {
			break
		}
	}
	return ret
}

func PrintMatrix(mat [][]rune, index []int) {
	for a := range mat {
		end := false
		o := ""
		for b := range mat[a] {
			if (a == index[0] && b == index[1]) || (a > index[0]) {
				end = true
				break
			}
			o += string(mat[a][b])
		}
		fmt.Printf("%2v: %v\n", a, o)
		if end == true {
			break
		}
	}
	if len(mat) == 0 {
		fmt.Println("empty mat")
	}
}

func (m *ArrangeMatrix) WipeOut() error {
	//log.Debugf("WipeOut: started, GC: %v, %v items", MemoryInfo(), len(m.MatrixResult))
	var err error
	m.MatrixResult, err = func() ([]ArrangeMatrixResult, error) {
		_, err := m.Writer.OutputPattern(
			m.Keyword, m.Surface, m.Number, m.MatrixResult, false,
			m.WipedLength, m.Width)
		if err != nil {
			return nil, err
		}
		return m.MatrixResult[:0], nil
	}()
	if err != nil {
		return err
	}
	runtime.GC()
	//log.Debugf("WipeOut: end,     GC: %v", MemoryInfo())
	return nil
}
