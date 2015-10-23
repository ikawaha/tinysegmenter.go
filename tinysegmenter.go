package tinysegmenter

import (
	"regexp"
	"unicode/utf8"
)

type CharType struct {
	name rune
	re   *regexp.Regexp
}

type Segmenter struct {
	chartypes []CharType
	model     map[string]int
}

const (
	Achar = 'A'
	Ichar = 'I'
	Hchar = 'H'
	Ochar = 'O'
	Uchar = 'U'
	Bchar = 'B'
)

func NewSegmenter() *Segmenter {
	s := new(Segmenter)
	regs := []struct {
		reg      string
		category rune
	}{
		{"[一二三四五六七八九十百千万億兆]", 'M'},
		{"[一-龠々〆ヵヶ]", 'H'},
		{"[ぁ-ん]", 'I'},
		{"[ァ-ヴーｱ-ﾝﾞｰ]", 'K'},
		{"[a-zA-Zａ-ｚＡ-Ｚ]", 'A'},
		{"[0-9０-９]", 'N'},
	}

	for _, v := range regs {
		re, _ := regexp.Compile(v.reg)
		s.chartypes = append(s.chartypes, CharType{name: v.category, re: re})
	}

	return s
}

func (s *Segmenter) gettype(str rune) rune {
	for _, v := range s.chartypes {
		if v.re.MatchString(string(str)) {
			return v.name
		}
	}
	return 'O'
}

func (s *Segmenter) Segment(text string) []string {
	result := make([]string, 0, len(text))
	if text == "" {
		return result
	}

	wordstart := 0
	pos := wordstart

	p1, w1, c1 := Uchar, B3, Ochar
	p2, w2, c2 := Uchar, B2, Ochar
	p3, w3, c3 := Uchar, B1, Ochar

	w5, c5 := E1, Ochar
	w6, c6 := E2, Ochar

	var pos1, pos2, pos3, size int
	w4, pos1 := utf8.DecodeRuneInString(text[pos:]) // rune を現在の位置から取ってくる
	c4 := s.gettype(w4)
	if pos1 < len(text) {
		w5, size = utf8.DecodeRuneInString(text[pos1:])
		c5 = s.gettype(w5)
		pos2 = pos1 + size
		if pos2 < len(text) {
			w6, size = utf8.DecodeRuneInString(text[pos2:])
			c6 = s.gettype(w6)
			pos3 = pos2 + size
		} else {
			w6, c6 = E1, Ochar
		}
	}

	for pos < len(text) {
		w1, w2, w3, w4, w5 = w2, w3, w4, w5, w6
		c1, c2, c3, c4, c5 = c2, c3, c4, c5, c6

		if pos1 == len(text) {
			w6, c6 = E2, Ochar
		} else if pos2 == len(text) {
			w6, c6 = E1, Ochar
		} else {
			pos3 = pos2 + utf8.RuneLen(w5)
			w6, _ = utf8.DecodeRuneInString(text[pos3:])
			c6 = s.gettype(w6)
		}

		score := BIAS
		if p1 == Ochar {
			score += -214 //score += get(UP1, p1, 0)
		}
		if p2 == Bchar {
			score += 69
		} else if p2 == Ochar {
			score += 935 //score += get(UP2, p2, 0)
		}
		if p3 == Bchar {
			score += 189 //score += get(UP3, p3, 0)
		}

		score += BP1[pair{p1, p2}]
		score += BP2[pair{p2, p3}]
		score += UW1[w1]
		score += UW2[w2]
		score += UW3[w3]
		score += UW4[w4]
		score += UW5[w5]
		score += UW6[w6]
		score += BW1[pair{w2, w3}]
		score += BW2[pair{w3, w4}]
		score += BW3[pair{w4, w5}]
		score += TW1[triple{w1, w2, w3}]
		score += TW2[triple{w2, w3, w4}]
		score += TW3[triple{w3, w4, w5}]
		score += TW4[triple{w4, w5, w6}]
		score += UC1[c1]
		score += UC2[c2]
		if c3 == Achar {
			score += -1370
		} else if c3 == Ichar {
			score += 2311 //score += get(UC3, c3, 0)
		}
		score += UC4[c4]
		score += UC5[c5]
		score += UC6[c6]
		score += BC1[pair{c2, c3}]
		score += BC2[pair{c3, c4}]
		score += BC3[pair{c4, c5}]
		score += TC1[triple{c1, c2, c3}]
		score += TC2[triple{c2, c3, c4}]
		score += TC3[triple{c3, c4, c5}]
		score += TC4[triple{c4, c5, c6}]
		//score += TC5[triple{c5,c6,c7}]
		score += UQ1[pair{p1, c1}]
		score += UQ2[pair{p2, c2}]
		score += UQ3[pair{p3, c3}]
		score += BQ1[triple{p2, c2, c3}]
		score += BQ2[triple{p2, c3, c4}]
		score += BQ3[triple{p3, c2, c3}]
		score += BQ4[triple{p3, c3, c4}]
		score += TQ1[quadra{p2, c1, c2, c3}]
		score += TQ2[quadra{p2, c2, c3, c4}]
		score += TQ3[quadra{p3, c1, c2, c3}]
		score += TQ4[quadra{p3, c2, c3, c4}]

		p := Ochar
		if score > 0 {
			result = append(result, text[wordstart:pos1])
			wordstart = pos1
			p = Bchar
		}
		p1, p2, p3 = p2, p3, p
		pos, pos1, pos2 = pos1, pos2, pos3

	}
	if wordstart != len(text) {
		result = append(result, text[wordstart:])
	}
	return result
}
