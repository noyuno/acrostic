package acrostic

/*
func TestGetSynonyms(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{DisableSorting: false, QuoteEmptyFields: true})
	expected := []WordNetResult{
		WordNetResult{Surface: []rune("アイコ"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("アイス"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("アイスコーヒー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("アイリッシュコーヒー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("インスタントコーヒー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("エスプレッソ"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("カップチーノ"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("カフェ"), Link: WNSynonym, Language: "jpn"},
		WordNetResult{Surface: []rune("カフェイン"), Link: WNMsub, Language: "jpn"},
		WordNetResult{Surface: []rune("カフェインレスコーヒー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("カフェオレ"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("カフェロワイヤル"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("カフェー"), Link: WNSynonym, Language: "jpn"},
		WordNetResult{Surface: []rune("カフエ"), Link: WNSynonym, Language: "jpn"},
		WordNetResult{Surface: []rune("カプチーノ"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("キャフェ"), Link: WNSynonym, Language: "jpn"},
		WordNetResult{Surface: []rune("コーヒー"), Link: WNSynonym, Language: "jpn"},
		WordNetResult{Surface: []rune("コーヒー豆"), Link: WNMsub, Language: "jpn"},
		WordNetResult{Surface: []rune("コールコーヒー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("デカフェ"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("デミタス"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("トルココーヒー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("ドリップコーヒー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("モカ"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("代替コーヒー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("代用コーヒー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("冷コー"), Link: WNHypo, Language: "jpn"},
		WordNetResult{Surface: []rune("水もの"), Link: WNHype, Language: "jpn"},
		WordNetResult{Surface: []rune("水物"), Link: WNHype, Language: "jpn"},
		WordNetResult{Surface: []rune("珈琲"), Link: WNSynonym, Language: "jpn"},
		WordNetResult{Surface: []rune("飲み料"), Link: WNHype, Language: "jpn"},
		WordNetResult{Surface: []rune("飲み物"), Link: WNHype, Language: "jpn"},
		WordNetResult{Surface: []rune("飲料"), Link: WNHype, Language: "jpn"},
		WordNetResult{Surface: []rune("飲物"), Link: WNHype, Language: "jpn"},
	}
	format := "WordNetResult{Surface: []rune(\"%v\"), Link: %v, Language: \"%v\", Part: \"%v\"},\n"

	o := &Options{
		JumanCommand:         "jumanpp",
		KnpCommand:           "knp -tab -anaphora",
		JumanDirectory:       "/usr/local/share/juman",
		JumanPPDirectory:     "/usr/local/share/jumanpp",
		WordNetDatabase:      "../../third-party/wnjpn/wnjpn.db",
		SynonymsJapaneseOnly: true,
		UseSynsetList:        false,
		Interactive:          false,
	}
	i := &Instance{}
	jk, err := NewJumanKnp(o, i)
	i.JumanKnp = jk

	if err != nil {
		t.Errorf("cannot initialize JumanKnp, error: %v", err.Error())
		return
	}

	var wn *WordNet
	wn, err = NewWordNet(o, i)
	if err != nil {
		t.Errorf("constructor error: %v", err.Error())
		return
	}
	if err = wn.DB.Ping(); err != nil {
		t.Errorf("ping failure on sqlite3(%v): %v", wn.Options.WordNetDatabase, err.Error())
	}
	var ret []WordNetResult
	bp := &BasicPhrase{
		Origin:  []rune("コーヒー"),
		Surface: []rune("コーヒーを"),
		ID:      0,
	}
	ret, err = wn.GetSynonymsBase(bp, []WordNetLink{WNSynonym, WNHype})
	if err != nil {
		t.Errorf("error: %v", err.Error())
		return
	}

	match := make([]bool, len(expected))
	for _, r := range ret {
		fmt.Printf(format, string(r.Surface), r.Link.String(), r.Language, r.Part.String())
		flag := false
		for i, e := range expected {
			if runes.Compare(r.Surface, e.Surface) {
				if match[i] {
					t.Errorf("already added")
				} else {
					match[i] = true
					flag = true
					break
				}
			}
		}
		if flag == false {
			t.Errorf("returned unexpected item\n"+format,
				string(r.Surface), r.Link.String(), r.Language, r.Part.String())
		}
	}
	for i := range match {
		if match[i] == false {
			t.Errorf("expected %v, but did not return it", string(expected[i].Surface))
		}
	}
}
*/
