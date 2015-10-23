package tinysegmenter

import (
	"unicode"
	"unicode/utf8"
)

var (
	kannumTable = []rune{
		'一', '二', '三', '四', '五', '六', '七', '八', '九', '十',
		'百', '千', '万', '億', '兆',
	}
	kanjiTable = &unicode.RangeTable{ //[一-龠々〆ヵヶ]
		R16: []unicode.Range16{
			{0x3005, 0x3006, 1},
			{0x30F5, 0x30F6, 1},
			{0x4E00, 0xff9E, 1},
		},
	}
	hiraganaTable = &unicode.RangeTable{ //[ぁ-ん]
		R16: []unicode.Range16{
			{0x3041, 0x3093, 1},
		},
	}
	katakanaTable = &unicode.RangeTable{ //[ァ-ヴーｱ-ﾝﾞｰ]
		R16: []unicode.Range16{
			{0x30A1, 0x30F4, 1},
			{0x30FC, 0x30FC, 1},
			{0xFF70, 0xFF9E, 1},
		},
	}
	alphabetTable = &unicode.RangeTable{ //[a-zA-Zａ-ｚＡ-Ｚ]
		R16: []unicode.Range16{
			{0x41, 0x5A, 1},
			{0x61, 0x7A, 1},
			{0xFF21, 0xFF3A, 1},
			{0xFF41, 0xFF5A, 1},
		},
	}
	numberTable = &unicode.RangeTable{ //[0-9０-９]
		R16: []unicode.Range16{
			{0x30, 0x39, 1},
			{0xFF10, 0xFF19, 1},
		},
	}
)

func gettype(c rune) rune {
	for _, x := range kannumTable {
		if x == c {
			return 'M'
		}
	}
	switch {
	case unicode.In(c, kanjiTable):
		return 'H'
	case unicode.In(c, hiraganaTable):
		return 'I'
	case unicode.In(c, katakanaTable):
		return 'K'
	case unicode.In(c, alphabetTable):
		return 'A'
	case unicode.In(c, numberTable):
		return 'N'
	}
	return 'O'
}

func Segment(input string) []string {
	ret := make([]string, 0, len(input))
	if input == "" {
		return ret
	}

	wordstart := 0
	pos := wordstart

	p1, w1, c1 := 'U', B3, 'O'
	p2, w2, c2 := 'U', B2, 'O'
	p3, w3, c3 := 'U', B1, 'O'

	w5, c5 := E1, 'O'
	w6, c6 := E2, 'O'

	var pos1, pos2, pos3, size int
	w4, pos1 := utf8.DecodeRuneInString(input[pos:])
	c4 := gettype(w4)
	if pos1 < len(input) {
		w5, size = utf8.DecodeRuneInString(input[pos1:])
		c5 = gettype(w5)
		pos2 = pos1 + size
		if pos2 < len(input) {
			w6, size = utf8.DecodeRuneInString(input[pos2:])
			c6 = gettype(w6)
			pos3 = pos2 + size
		} else {
			w6, c6 = E1, 'O'
		}
	}

	for pos < len(input) {
		w1, w2, w3, w4, w5 = w2, w3, w4, w5, w6
		c1, c2, c3, c4, c5 = c2, c3, c4, c5, c6

		if pos1 == len(input) {
			w6, c6 = E2, 'O'
		} else if pos2 == len(input) {
			w6, c6 = E1, 'O'
		} else {
			pos3 = pos2 + utf8.RuneLen(w5)
			w6, _ = utf8.DecodeRuneInString(input[pos3:])
			c6 = gettype(w6)
		}

		score := BIAS
		if p1 == 'O' { //score += UP1[p1]
			score += -214
		}
		if p2 == 'B' { //score += UP2[p2]
			score += 69
		} else if p2 == 'O' {
			score += 935
		}
		if p3 == 'B' { //score += UP3[p3]
			score += 189
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
		if c3 == 'A' { //score += UC3[c3]
			score += -1370
		} else if c3 == 'I' {
			score += 2311
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

		p := 'O'
		if score > 0 {
			ret = append(ret, input[wordstart:pos1])
			wordstart = pos1
			p = 'B'
		}
		p1, p2, p3 = p2, p3, p
		pos, pos1, pos2 = pos1, pos2, pos3

	}
	if wordstart != len(input) {
		ret = append(ret, input[wordstart:])
	}
	return ret
}
