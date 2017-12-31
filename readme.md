# acrostic

A software to make it possible to read Japanese horizontal writing text vertically, an acrostic.

日本語の横書きのテキストを縦にも読めるようにするソフトウェア．

## Requirements

### Golang

[The Go Programming Language](https://golang.org/)

- Golang v1.9.2
    - [gb](https://getgb.io/)
    - github.com/kr/pty
    - github.com/mattn/go-sqlite3
    - github.com/nicksnyder/go-i18n/i18n
    - github.com/noyuno/lgo
    - github.com/ogier/pflag
    - github.com/sirupsen/logrus

~~~
export GOPATH=$HOME/go
mkdir -p $GOPATH
go get github.com/constabulary/gb/...
gb vendor restore
~~~

### Juman & KNP

Both JUMAN and JUMAN++ are required.
This repository does not contain any third-party source codes/databases.
Please get them yourself.

JUMANとJUMAN++の両方が必要です．
このリポジトリにはサードパーティ製のソースコードおよびデータベースが含まれていません．
それらはご自身で取得してください．

#### JUMAN 7.01

[JUMAN - KUROHASHI-KAWAHARA LAB](http://nlp.ist.i.kyoto-u.ac.jp/index.php?JUMAN)

    ./configure
    make
    sudo make install

It is used for inflection.

#### JUMAN++ 1.02

[JUMAN++ - KUROHASHI-KAWAHARA LAB](http://nlp.ist.i.kyoto-u.ac.jp/index.php?JUMAN++)

    ./configure
    make
    sudo make install

It is used for KNP.

#### KNP 4.18

[KNP - KUROHASHI-KAWAHARA LAB](http://nlp.ist.i.kyoto-u.ac.jp/?KNP)

    ./configure
    make
    sudo make install

It is used for dependency analysis and case analysis.

### WordNet 1.1 Sqlite3 database

[日本語 Wordnet](http://compling.hss.ntu.edu.sg/wnja/)

    mkdir -p third-party/wnjpn
    curl -o /tmp/wnjpn.db.gz http://compling.hss.ntu.edu.sg/wnja/data/1.1/wnjpn.db.gz
    zcat /tmp/wnjpn.db.gz > third-party/wnjpn/wnjpn.db

It is used to get paraphrases.

### MeCab

#### Software

    ./configure --with-charset=utf8
    make
    sudo make install

[MeCab: Yet Another Part-of-Speech and Morphological Analyzer](http://taku910.github.io/mecab/)

It is used to get Kana.

#### IPA dictionary

[MeCab: Yet Another Part-of-Speech and Morphological Analyzer](http://taku910.github.io/mecab/)

    ./configure --with-charset=utf-8
    make
    sudo make install

#### NEologd dictionary

[neologd/mecab-ipadic-neologd: Neologism dictionary based on the language resources on the Web for mecab-ipadic](https://github.com/neologd/mecab-ipadic-neologd)

    git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git
    cd mecab-ipadic-neologd
    ./bin/install-mecab-ipadic-neologd -n

## Build

Ubuntu: `apt install golang`
Arch Linux: `pacman -S go`

    gb build all -f
    ./bin/main -t data/4.sentence -k data/3.keyword

## Run sample

    ./scripts/sample1.sh

~~~
　パックのきみ
つ性などをたか
めて、べいはん
の味や品質を長
持ちさせ、日本
産米の輸出拡大
につなげる。
~~~

## Common usage

~~~
  --confirm
        処理前にユーザによる確認を行う (default true)
  -h, --height int
        最大行(未指定であれば(文字数/Width*2)) (default -1)
  -i, --interactive
        インタラクティブ（対話的）に実行する
  --kana
        かなを使用する (default true)
  --kanji
        漢字を使う
  -k, --keyword string
        キーワードファイル名
  -l, --match-length
        キーワードの文字数と出力文の行数を一致させる (default true)
  -m, --max int
        行の最大幅(-1でWidthと同じにする) (default -1)
  --one
        一つ見つけたら終了する
  --only-keywords
        キーワード限定 (default true)
  -o, --out string
        出力ファイル名
  --paraphrase string
        言い換えデータベースのファイル名(内容はcsv) (default "data/paraphrase.csv")
  --polite
        丁寧語を使うかどうか (default true)
  --print-kana
        類義語画面でかなも表示する
  --progress
        進捗表示
  --skip-same-length
        基本句の類義語Aの文字数がその基本句の他の類義語の文字数と同じで，すでに処理されている ときは，Aの探索を省略する (default true)
  -a, --swap
        文を入れ替えるかどうか（レシートなどの箇条書きに有用）
  --synonyms-verb
        動詞の類義語を使うかどうか
  -s, --synset string
        SynsetListTextが含まれているファイル
  -t, --text string
        テキストファイル名
  -v, --verbose
        INFO出力を有効にする
  -V, --verbosely
        DEBUG出力を有効にする
  -w, --width int
        行の幅 (default 10)
  --word-pattern int
        類義語パターンの最大サイズ (default 100)
  --wordnet-link string
        WordNetで検索するリンクを指定 (default "synonyms,hype")
~~~

Type `./bin/main --help` to show options' descriptions.

## Edit

    go get golang.org/x/tools/cmd/stringer
    stringer -type WordNetLink src/acrostic/wordnetlink.go
    ctags -R
    go get github.com/gosexy/gettext/go-xgettext

### NeoVim

#### Plugins

- [fatih/vim-go: Go development plugin for Vim](https://github.com/fatih/vim-go)
    - `:GoInstallBinaries` to install required binaries.

## Debug

require direnv

    go get github.com/derekparker/delve/cmd/dlv
    dlv debug main -- -k data/3.keyword -t data/4.sentence

## Analyze

require direnv

    go get github.com/uber/go-torch
    git clone https://github.com/brendangregg/FlameGraph.git ~/go/src/github.com/uber/go-torch/FlameGraph
    go-torch --width=5000 -f ~/output/torch.svg

## Test

require direnv

    go test -v ./...

## Help

require direnv

    godoc -http :8080

## direnv

- Ubuntu: `apt install direnv`
- Arch Linux: `pacman -S direnv`
- Zsh: `eval "$(direnv hook zsh)"`

    direnv allow

## How it works

First, generate paraphrases of phrases of input text.
Second, generate sentence patterns from dependency analysis and case analysis.
Last, verify generated sentence pattern can be read vertically.

まず，入力テキストの単語の言い換えを生成する．
次に，係り受け解析および格解析から文パターンを生成する．
最後に，生成された文パターンが縦読み可能かどうかを検証する．

## Verification

「みかん」をキーワードとしたとき，入力243文中縦読み可能なテキストを77文（約31.7%）作成できた．
ただし，入力は読売新聞社のインターネット上のニュース記事のうち，70文字以内に収まる文で，
出力テキストの幅は4から30文字，類義語は1基本句あたり5つまでとした．

これら検証用データは，このリポジトリに含まれていない．

