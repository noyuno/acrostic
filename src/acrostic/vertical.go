// acrostic : 縦読み化プログラム
package acrostic

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/noyuno/lgo/color"
	flag "github.com/ogier/pflag"
	log "github.com/sirupsen/logrus"
)

// Options : プログラムオプション
type Options struct {
	// TextFileName : テキストファイル名
	TextFileName string
	// KeywordFileName : キーワードファイル名
	KeywordFileName string
	// OutFileName : 出力ファイル名
	OutFileName string
	// StdoutResult : 結果を標準出力するかどうか
	//StdoutResult bool
	// Mode : 解析ツール（未実装）
	Mode string

	JumanCommand string

	KnpCommand string

	KakasiCommand string

	MeCabCommand string

	// KnpOnly : KNPの実行だけして終了する
	KnpOnly bool
	// CaseAnalysis : 格解析をする
	CaseAnalysis bool
	// Synonyms : WordNetを使って類語で言い換える
	Synonyms bool
	// WordNetDatabase : WordNetデータベースのファイル名
	WordNetDatabase string
	// SynonymsJapaneseOnly : 類語検索結果は日本語だけとする
	SynonymsJapaneseOnly bool
	// Width : 行の幅
	Width int
	// MaxWidth : 行の最大幅(-1でWidthと同じにする)
	MaxWidth int
	// Height : 最大行(未指定であれば(文字数/Width*2))
	Height int

	// JumanDirectory : 活用形を取得するためのファイル
	JumanDirectory string
	//JumanPPDirectory string

	// 類義語の概念を示すリスト
	// 書式：BasicPhrasePos:Synset,Synset,...;BasicPhrasePos:Synset,Synset,...
	// 例：4:02284544-v;5:01110517-v,01171183-v
	SynsetListText string

	// SynsetListTextが含まれているファイル
	SynsetListFile string

	// 類義語の概念を示すリストを使用するかどうか
	UseSynsetList bool

	// 類義語の概念を示すリスト
	SynsetList map[int]map[string]WordNetAnswer

	// インタラクティブ（対話的）に実行する
	Interactive bool

	// 基本句の類義語Aの文字数がその基本句の他の類義語の文字数と同じで，
	// すでに処理されているときは，Aの探索を省略する
	SkipSameLength bool

	// 処理前にユーザによる確認を行う
	Confirm bool

	// WARNING出力を無効にする
	Quiet bool
	// INFO出力を有効にする
	Verbose bool
	// DEBUG出力を有効にする
	Verbosely bool

	// 類義語画面でかなも出力する
	// (このオプションをtrueにしなくても，UseKanaがtrueであればかなに基づいて探索する)
	PrintKana bool

	UseKana bool

	// かなモード: mecab, juman, kakasiのいづれかひとつ
	// ただし，jumanを指定するときはフォールバック先を選択できる
	// example:
	// mecab
	// juman,mecab
	// juman,mecab,kakasi
	KanaMode string

	KanaModeList  map[string]bool
	KanaModeOrder []string

	// WordNetで検索するリンクを指定
	WordNetLinkString string
	WordNetLink       []WordNetLink

	// 言い換えデータベースのファイル名(内容はcsv)
	ParaphraseDatabase string

	// 並列処理で計算する
	Parallel bool

	// pproofをつかう
	UsePproof bool

	// 拡張構造を有効にする(未実装)
	ExtensionStructure bool

	// BasicPhrase以下の構造体もコピーする(メモリ対策)
	EnableDeepCopy bool

	// パターンごとに出力する
	OutputEachPattern bool

	// ArrangeMatrixで結果を書き出して一掃するかどうか（メモリ使用量対策）
	WipeOut bool

	// ArrangeMatrixで結果を書き出して一掃するタイミング
	WipeOutLength int

	// SynonymsVerb : 動詞の類義語を使うかどうか
	SynonymsVerb bool

	// UsePolite : 丁寧語を使うかどうか
	UsePolite bool

	// MatchLength : キーワードの文字数と出力文の行数を一致させる
	MatchLength bool

	// PatternSize : 文パターンの最大サイズ
	PatternSize int

	// SwapSentences : 文を入れ替えるかどうか（レシートなどの箇条書きに有効）
	SwapSentences bool

	// GCHeapSize : GCするヒープサイズ(ただし，WipeOutではこれに関わらずかならずGCする)
	GCHeapSize uint64

	Progress bool

	//ProgressDepth int

	Language string

	One bool

	ExitCode bool

	WordPatternLimit int

	OnlyKeywords bool

	AllWordLength bool

	UseKanji bool
}

// Instance : 共通インスタンス
type Instance struct {
	// Juman KNP
	JumanKnp *JumanKnp
	// WordNet : 類語検索
	WordNet *WordNet

	Kakasi *Kakasi

	MeCab *MeCab

	Kana *Kana

	Paraphrase *Paraphrase

	Variables *Variables
}

// Acrostic : 構造の根
type Acrostic struct {
	Keywords   [][]rune
	Text       [][]rune
	Options    *Options
	Instance   *Instance
	Paragraphs []Paragraph
	Found      bool
}

// NewOptions : constructor
func NewOptions() (*Options, error) {
	o := new(Options)

	o.i18n()
	T, _ := i18n.Tfunc(o.Language)

	flag.StringVarP(&o.KeywordFileName, "keyword", "k", "", T("f-keyword"))
	flag.StringVarP(&o.TextFileName, "text", "t", "", T("f-text"))
	flag.StringVarP(&o.OutFileName, "out", "o", "", T("f-out"))
	flag.IntVarP(&o.Width, "width", "w", 10, T("f-width"))
	flag.IntVarP(&o.MaxWidth, "max", "m", -1, T("f-max"))
	flag.IntVarP(&o.Height, "height", "h", -1, T("f-height"))

	flag.StringVar(&o.JumanCommand, "juman-command", "jumanpp", T("f-juman-command"))
	flag.StringVar(&o.KnpCommand, "knp-command", "knp -tab -anaphora", T("f-knp-command"))
	flag.StringVar(&o.KakasiCommand, "kakasi-command", "kakasi -JH -iutf-8 -outf-8", T("f-kakasi-command"))
	flag.StringVar(&o.MeCabCommand, "mecab-command", "mecab -d /usr/local/lib/mecab/dic/mecab-ipadic-neologd -O yomi", T("f-mecab-command"))
	flag.BoolVar(&o.KnpOnly, "knp-only", false, T("f-knp-only"))
	flag.BoolVar(&o.CaseAnalysis, "case-analysis", true, T("f-case-analysis"))
	flag.BoolVar(&o.Synonyms, "synonyms", true, T("f-synonyms"))
	flag.StringVar(&o.WordNetDatabase, "wordnetdb", "third-party/wnjpn/wnjpn.db", T("f-wordnetdb"))
	flag.StringVar(&o.JumanDirectory, "juman-directory", "/usr/local/share/juman", T("f-juman-directory"))
	flag.BoolVar(&o.SynonymsJapaneseOnly, "synonyms-jpn-only", true, T("f-synonyms-jpn-only"))
	flag.StringVarP(&o.SynsetListFile, "synset", "s", "", T("f-synset"))
	flag.StringVar(&o.SynsetListText, "synset-string", "", T("f-synset-string"))
	flag.BoolVarP(&o.Interactive, "interactive", "i", false, T("f-interactive"))
	flag.BoolVar(&o.SkipSameLength, "skip-same-length", true, T("f-skip-same-length"))
	flag.BoolVarP(&o.Verbose, "verbose", "v", false, T("f-verbose"))
	flag.BoolVarP(&o.Verbosely, "verbosely", "V", false, T("f-verbosely"))
	flag.BoolVarP(&o.Quiet, "quiet", "q", false, T("f-quiet"))
	flag.BoolVar(&o.PrintKana, "print-kana", false, T("f-print-kana"))
	flag.StringVar(&o.KanaMode, "kana-mode", "mecab", T("f-kana-mode"))
	flag.StringVar(&o.WordNetLinkString, "wordnet-link", "synonyms,hype", T("f-wordnet-link"))
	flag.StringVar(&o.ParaphraseDatabase, "paraphrase", "data/paraphrase.csv", T("f-paraphrase"))
	flag.BoolVarP(&o.Parallel, "parallel", "j", true, T("f-parallel"))
	flag.BoolVar(&o.Confirm, "confirm", true, T("f-confirm"))
	flag.BoolVar(&o.UsePproof, "use-pproof", true, T("f-pproof"))
	flag.BoolVar(&o.ExtensionStructure, "extension-structure", false, T("f-extension-structure"))
	flag.BoolVar(&o.EnableDeepCopy, "deep-copy", false, T("f-deep-copy"))
	flag.BoolVar(&o.OutputEachPattern, "output-each", true, T("f-output-each"))
	flag.IntVar(&o.WipeOutLength, "wipeout-length", 1000000, T("f-wipeout-length"))
	flag.BoolVar(&o.WipeOut, "wipeout", true, T("f-wipeout"))
	flag.BoolVar(&o.SynonymsVerb, "synonyms-verb", false, T("f-synonyms-verb"))
	flag.BoolVar(&o.UsePolite, "polite", true, T("f-polite"))
	flag.BoolVarP(&o.MatchLength, "match-length", "l", true, T("f-match-length"))
	flag.IntVar(&o.PatternSize, "pattern-size", 1000000, T("f-pattern-size"))
	flag.BoolVarP(&o.SwapSentences, "swap", "a", false, T("f-swap"))
	flag.Uint64Var(&o.GCHeapSize, "gc", 10*1024*1024, T("f-gc"))
	flag.BoolVar(&o.Progress, "progress", false, T("f-progress"))
	//flag.IntVar(&o.ProgressDepth, "progress-depth", 20, "depth of progress")
	flag.StringVar(&o.Language, "language", "", T("f-language"))
	flag.BoolVar(&o.One, "one", false, T("f-one"))
	flag.BoolVar(&o.ExitCode, "code", false, T("f-code"))
	flag.IntVar(&o.WordPatternLimit, "word-pattern", 100, T("f-word-pattern"))
	flag.BoolVar(&o.OnlyKeywords, "only-keywords", true, T("f-only-keywords"))
	flag.BoolVar(&o.AllWordLength, "all-word-length", true, T("f-all-word-length"))
	flag.BoolVar(&o.UseKanji, "kanji", false, T("f-kanji"))
	flag.BoolVar(&o.UseKana, "kana", true, T("f-kana"))
	flag.Parse()
	if o.KeywordFileName == "" {
		return nil, errors.New("require keyword (-k)")
	}
	if o.TextFileName == "" {
		return nil, errors.New("require text (-t)")
	}
	// default log level is warning level
	log.SetLevel(log.WarnLevel)
	if o.Quiet {
		log.SetLevel(log.ErrorLevel)
	}
	if o.Verbose {
		log.SetLevel(log.InfoLevel)
	}
	if o.Verbosely {
		log.SetLevel(log.DebugLevel)
	}
	// pproof
	if o.UsePproof {
		go func() {
			log.Debug(http.ListenAndServe("localhost:8080", nil))
			log.Debug("pproof endpoint: http://localhost:6060/debug/pprof/heap?debug=1")
		}()
	}

	err := o.parseSynset()
	if err != nil {
		return nil, err
	}
	err = o.parseKanaMode()
	if err != nil {
		return nil, err
	}
	err = o.parseWordNetLink()
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *Options) i18n() {
	lang := o.Language
	if lang == "" {
		lang = os.Getenv("LANGUAGE")
		if lang == "" {
			lang = os.Getenv("LC_ALL")
			if lang == "" {
				lang = os.Getenv("LANG")
				if lang == "" {
					lang = "en_US"
				}
			}
		}
	}
	lang = strings.Split(lang, ".")[0]
	lang = strings.Replace(lang, "_", "-", -1)
	lang = strings.ToLower(lang)
	o.Language = lang
	SystemLanguage = lang
	i18n.MustLoadTranslationFile("locales/" + lang + ".yaml")
}

func (o *Options) parseSynset() error {
	if o.SynsetListFile != "" {
		fp, err := os.Open(o.SynsetListFile)
		if err != nil {
			return errors.New("failed to open synset file: " + o.SynsetListFile)
		}
		defer fp.Close()
		v := make([]byte, 0)
		v, err = ioutil.ReadAll(fp)
		o.SynsetListText = strings.Replace(string(v), "\n", "", -1)
		if err != nil {
			return err
		}
	}
	if o.SynsetListText == "" {
		return nil
	}
	o.SynsetList = map[int]map[string]WordNetAnswer{}
	for _, word := range strings.Split(o.SynsetListText, ";") {
		//log.Debug(word)
		if word == "" {
			continue
		}
		syns := strings.Split(word, ":")
		if len(syns) != 2 {
			return errors.New("there must be exactly one colon in one BasicPhrase")
		}
		pos, err := strconv.Atoi(syns[0])
		if err != nil {
			return errors.New("cannot convert BasicPhrase position number to int")
		}
		o.SynsetList[pos] = map[string]WordNetAnswer{}
		for _, syn := range strings.Split(syns[1], ",") {
			lr := strings.Split(syn, "=")
			if len(lr) != 2 {
				return errors.New("left=right format error")
			}
			o.SynsetList[pos][lr[0]] = NewWordNetAnswer(lr[1])
			o.UseSynsetList = true
		}
	}
	return nil
}

func (o *Options) parseKanaMode() error {
	o.KanaModeList = map[string]bool{}
	o.KanaModeOrder = []string{}
	if o.KanaMode == "" {
		return nil
	}
	modes := strings.Split(o.KanaMode, ",")
	if len(modes) >= 2 && modes[0] != "juman" {
		return errors.New("kana-mode: only fallback is valid for juman")
	}
	for _, m := range modes {
		if (m == "juman" || m == "mecab" || m == "kakasi") == false {
			return errors.New("kana-mode: only juman, mecab or kakasi")
		}
		o.KanaModeList[m] = true
		o.KanaModeOrder = append(o.KanaModeOrder, m)
	}
	return nil
}

func (o *Options) parseWordNetLink() error {
	o.WordNetLink = []WordNetLink{}
	if o.WordNetLinkString == "" {
		return nil
	}
	for _, l := range strings.Split(o.WordNetLinkString, ",") {
		o.WordNetLink = append(o.WordNetLink, NewWordNetLink(l))
	}
	return nil
}

func NewInstance(o *Options) (*Instance, error) {
	log.Debug("initializing instances")
	var err error
	ret := new(Instance)
	ret.JumanKnp, err = NewJumanKnp(o, ret)
	if err != nil {
		return nil, err
	}
	ret.WordNet, err = NewWordNet(o, ret)
	if err != nil {
		return nil, err
	}
	ret.Kana = NewKana(o, ret)
	if o.KanaModeList["kakasi"] {
		ret.Kakasi, err = NewKakasi(o)
		if err != nil {
			return nil, err
		}
	}
	if o.KanaModeList["mecab"] {
		ret.MeCab, err = NewMeCab(o)
		if err != nil {
			return nil, err
		}
	}
	if o.ParaphraseDatabase != "" {
		ret.Paraphrase, err = NewParaphrase(o)
		if err != nil {
			return nil, err
		}
	}
	ret.Variables = &Variables{}
	return ret, nil
}

// NewVertical : constructor
func NewVertical(o *Options) (*Acrostic, error) {
	var err error
	v := new(Acrostic)
	v.Options = o
	if o.MaxWidth == -1 {
		o.MaxWidth = o.Width
	}
	v.Instance, err = NewInstance(v.Options)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// ファイルから読み取って，配列に1つの要素として格納する
func readFileA1(filename string, slice *[][]rune) error {
	var f, err = os.Open(filename)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	out := ""
	for scanner.Scan() {
		t := scanner.Text()
		if t != "" {
			out += t + "\n"
		}
	}
	*slice = append(*slice, []rune(out+"\n"))
	return scanner.Err()
}

// ファイルから読み取って，配列に改行で区切られた要素を格納する
func readFileA2(filename string, slice *[][]rune) error {
	var f, err = os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if t != "" {
			*slice = append(*slice, []rune(t))
		}
	}
	return scanner.Err()
}

// ReadKeyword : ファイルからキーワードを読み取る
func (v *Acrostic) ReadKeyword() error {
	return readFileA2(v.Options.KeywordFileName, &v.Keywords)
}

// ReadText : ファイルからテキストを読み取る
func (v *Acrostic) ReadText() error {
	err := readFileA1(v.Options.TextFileName, &v.Text)
	if err != nil {
		return err
	}
	if v.Options.Height == -1 {
		// auto
		lflen := 0
		for _, c := range v.Text[len(v.Text)-1] {
			if string(c) == "\n" {
				lflen++
			}
		}
		v.Options.Height = (len(v.Text[len(v.Text)-1]) * lflen) / v.Options.Width * 2
	}
	return nil
}

// Analyze : 解析する
func (v *Acrostic) Analyze() error {
	T, _ := i18n.Tfunc(v.Options.Language)
	for _, t := range v.Text {
		//log.Debugf("Acrostic.Analyze: %v", string(t))
		p := NewParagraph(v.Options, v.Instance, t, v.Keywords)
		err := p.Analyze()
		if err != nil {
			return err
		}
		if v.Options.KnpOnly == false {
			p.PrintAnalyzeResult()
			keywords := make([][]rune, 0)
			for k := range v.Keywords {
				if p.CheckContainsKeyword(v.Keywords[k], k) {
					keywords = append(keywords, v.Keywords[k])
				}
			}
			if len(keywords) == 0 {
				log.Fatalf(T("All keywords not found. Abort."))
			}
			v.Keywords = keywords
			if p.Options.Confirm {
				TypeToContinue()
			}
		}
		v.Paragraphs = append(v.Paragraphs, *p)
	}
	return nil
}

func (v *Acrostic) Generate() error {
	T, _ := i18n.Tfunc(v.Options.Language)
	start := time.Now().UTC()

	for k := range v.Keywords {
		for w := v.Options.Width; w <= v.Options.MaxWidth; w++ {
			fmt.Printf("%v%2v: %v%v (%v: %v)\n",
				color.FGreen, k, string(v.Keywords[k]), color.Reset, T("width"), w)
			for p := range v.Paragraphs {
				if v.Paragraphs[p].FoundBasicPhrase {
					r, err := v.Paragraphs[p].Generate(v.Keywords[k], k, w)
					if err != nil {
						return err
					}
					if r {
						v.Found = true
						if v.Options.One {
							break
						}
					}
				} else {
					log.Fatal("invalid operation to call Generate")
				}
			}
			if v.Options.One && v.Found {
				break
			}
		}
		if v.Options.One && v.Found {
			break
		}
	}
	if v.Options.Verbose || v.Options.Verbosely {
		elapsed := time.Since(start)
		fmt.Printf(T("arrange process time")+": %v\n", elapsed)
	}
	return nil
}
