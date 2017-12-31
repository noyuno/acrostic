#!/bin/bash -e

clean() {
    pushd $1
        make clean
    popd
}

clean slide
clean thesis
clean abstract
clean third-party/juman-7.01
clean third-party/jumanpp-1.02
clean third-party/kakasi-2.3.6
clean third-party/knp-4.18
clean third-party/mecab-0.996
clean third-party/mecab-ipadic-2.7.0-20070801

rm -rvf bin pkg output

