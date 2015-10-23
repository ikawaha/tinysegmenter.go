package tinysegmenter

import (
	"regexp"
	"strings"
)

type CharType struct {
	name string
	re   *regexp.Regexp
}

type Segmenter struct {
	chartypes []CharType
	model     map[string]int
}

func NewSegmenter() *Segmenter {
	s := new(Segmenter)
	regs := [][]string{
		{"[一二三四五六七八九十百千万億兆]", "M"},
		{"[一-龠々〆ヵヶ]", "H"},
		{"[ぁ-ん]", "I"},
		{"[ァ-ヴーｱ-ﾝﾞｰ]", "K"},
		{"[a-zA-Zａ-ｚＡ-Ｚ]", "A"},
		{"[0-9０-９]", "N"},
	}

	for _, v := range regs {
		re, _ := regexp.Compile(v[0])
		s.chartypes = append(s.chartypes, CharType{v[1], re})
	}

	s.model = map[string]int{"BC1:HH": 6, "BC1:II": 2461, "BC1:KH": 406, "BC1:OH": -1377, "BC2:AA": -3266, "BC2:AI": 2744, "BC2:AN": -877, "BC2:HH": -4069, "BC2:HM": -1710, "BC2:HN": 4012, "BC2:HO": 3761, "BC2:IA": 1327, "BC2:IH": -1183, "BC2:II": -1331, "BC2:IK": 1721, "BC2:IO": 5492, "BC2:KI": 3831, "BC2:KK": -8740, "BC2:MH": -3131, "BC2:MK": 3334, "BC2:OO": -2919, "BC3:HH": 996, "BC3:HI": 626, "BC3:HK": -720, "BC3:HN": -1306, "BC3:HO": -835, "BC3:IH": -300, "BC3:KK": 2762, "BC3:MK": 1079, "BC3:MM": 4034, "BC3:OA": -1651, "BC3:OH": 266, "BIAS": -331, "BP1:BB": 295, "BP1:OB": 304, "BP1:OO": -124, "BP1:UB": 352, "BP2:BO": 60, "BP2:OO": -1761, "BQ1:BHH": 1150, "BQ1:BHM": 1521, "BQ1:BII": -1157, "BQ1:BIM": 886, "BQ1:BMH": 1208, "BQ1:BNH": 449, "BQ1:BOH": -90, "BQ1:BOO": -2596, "BQ1:OHI": 451, "BQ1:OIH": -295, "BQ1:OKA": 1851, "BQ1:OKH": -1019, "BQ1:OKK": 904, "BQ1:OOO": 2965, "BQ2:BHH": 118, "BQ2:BHI": -1158, "BQ2:BHM": 466, "BQ2:BIH": -918, "BQ2:BKK": -1719, "BQ2:BKO": 864, "BQ2:OHH": -1138, "BQ2:OHM": -180, "BQ2:OIH": 153, "BQ2:UHI": -1145, "BQ3:BHH": -791, "BQ3:BHI": 2664, "BQ3:BII": -298, "BQ3:BKI": 419, "BQ3:BMH": 937, "BQ3:BMM": 8335, "BQ3:BNN": 998, "BQ3:BOH": 775, "BQ3:OHH": 2174, "BQ3:OHM": 439, "BQ3:OII": 280, "BQ3:OKH": 1798, "BQ3:OKI": -792, "BQ3:OKO": -2241, "BQ3:OMH": -2401, "BQ3:OOO": 11699, "BQ4:BHH": -3894, "BQ4:BIH": 3761, "BQ4:BII": -4653, "BQ4:BIK": 1348, "BQ4:BKK": -1805, "BQ4:BMI": -3384, "BQ4:BOO": -12395, "BQ4:OAH": 926, "BQ4:OHH": 266, "BQ4:OHK": -2035, "BQ4:ONN": -972, "BW1:,と": 660, "BW1:,同": 727, "BW1:B1あ": 1404, "BW1:B1同": 542, "BW1:、と": 660, "BW1:、同": 727, "BW1:」と": 1682, "BW1:あっ": 1505, "BW1:いう": 1743, "BW1:いっ": -2054, "BW1:いる": 672, "BW1:うし": -4816, "BW1:うん": 665, "BW1:から": 3472, "BW1:がら": 600, "BW1:こう": -789, "BW1:こと": 2083, "BW1:こん": -1261, "BW1:さら": -4142, "BW1:さん": 4573, "BW1:した": 2641, "BW1:して": 1104, "BW1:すで": -3398, "BW1:そこ": 1977, "BW1:それ": -870, "BW1:たち": 1122, "BW1:ため": 601, "BW1:った": 3463, "BW1:つい": -801, "BW1:てい": 805, "BW1:てき": 1249, "BW1:でき": 1127, "BW1:です": 3445, "BW1:では": 844, "BW1:とい": -4914, "BW1:とみ": 1922, "BW1:どこ": 3887, "BW1:ない": 5713, "BW1:なっ": 3015, "BW1:など": 7379, "BW1:なん": -1112, "BW1:にし": 2468, "BW1:には": 1498, "BW1:にも": 1671, "BW1:に対": -911, "BW1:の一": -500, "BW1:の中": 741, "BW1:ませ": 2448, "BW1:まで": 1711, "BW1:まま": 2600, "BW1:まる": -2154, "BW1:やむ": -1946, "BW1:よっ": -2564, "BW1:れた": 2369, "BW1:れで": -912, "BW1:をし": 1860, "BW1:を見": 731, "BW1:亡く": -1885, "BW1:京都": 2558, "BW1:取り": -2783, "BW1:大き": -2603, "BW1:大阪": 1497, "BW1:平方": -2313, "BW1:引き": -1335, "BW1:日本": -194, "BW1:本当": -2422, "BW1:毎日": -2112, "BW1:目指": -723, "BW1:Ｂ１あ": 1404, "BW1:Ｂ１同": 542, "BW1:｣と": 1682, "BW2:..": -11821, "BW2:11": -668, "BW2:――": -5729, "BW2:−−": -13174, "BW2:いう": -1608, "BW2:うか": 2490, "BW2:かし": -1349, "BW2:かも": -601, "BW2:から": -7193, "BW2:かれ": 4612, "BW2:がい": 853, "BW2:がら": -3197, "BW2:きた": 1941, "BW2:くな": -1596, "BW2:こと": -8391, "BW2:この": -4192, "BW2:させ": 4533, "BW2:され": 13168, "BW2:さん": -3976, "BW2:しい": -1818, "BW2:しか": -544, "BW2:した": 5078, "BW2:して": 972, "BW2:しな": 939, "BW2:その": -3743, "BW2:たい": -1252, "BW2:たた": -661, "BW2:ただ": -3856, "BW2:たち": -785, "BW2:たと": 1224, "BW2:たは": -938, "BW2:った": 4589, "BW2:って": 1647, "BW2:っと": -2093, "BW2:てい": 6144, "BW2:てき": 3640, "BW2:てく": 2551, "BW2:ては": -3109, "BW2:ても": -3064, "BW2:でい": 2666, "BW2:でき": -1527, "BW2:でし": -3827, "BW2:です": -4760, "BW2:でも": -4202, "BW2:とい": 1890, "BW2:とこ": -1745, "BW2:とと": -2278, "BW2:との": 720, "BW2:とみ": 5168, "BW2:とも": -3940, "BW2:ない": -2487, "BW2:なが": -1312, "BW2:など": -6508, "BW2:なの": 2614, "BW2:なん": 3099, "BW2:にお": -1614, "BW2:にし": 2748, "BW2:にな": 2454, "BW2:によ": -7235, "BW2:に対": -14942, "BW2:に従": -4687, "BW2:に関": -11387, "BW2:のか": 2093, "BW2:ので": -7058, "BW2:のに": -6040, "BW2:のの": -6124, "BW2:はい": 1073, "BW2:はが": -1032, "BW2:はず": -2531, "BW2:ばれ": 1813, "BW2:まし": -1315, "BW2:まで": -6620, "BW2:まれ": 5409, "BW2:めて": -3152, "BW2:もい": 2230, "BW2:もの": -10712, "BW2:らか": -943, "BW2:らし": -1610, "BW2:らに": -1896, "BW2:りし": 651, "BW2:りま": 1620, "BW2:れた": 4270, "BW2:れて": 849, "BW2:れば": 4114, "BW2:ろう": 6067, "BW2:われ": 7901, "BW2:を通": -11876, "BW2:んだ": 728, "BW2:んな": -4114, "BW2:一人": 602, "BW2:一方": -1374, "BW2:一日": 970, "BW2:一部": -1050, "BW2:上が": -4478, "BW2:会社": -1115, "BW2:出て": 2163, "BW2:分の": -7757, "BW2:同党": 970, "BW2:同日": -912, "BW2:大阪": -2470, "BW2:委員": -1249, "BW2:少な": -1049, "BW2:年度": -8668, "BW2:年間": -1625, "BW2:府県": -2362, "BW2:手権": -1981, "BW2:新聞": -4065, "BW2:日新": -721, "BW2:日本": -7067, "BW2:日米": 3372, "BW2:曜日": -600, "BW2:朝鮮": -2354, "BW2:本人": -2696, "BW2:東京": -1542, "BW2:然と": -1383, "BW2:社会": -1275, "BW2:立て": -989, "BW2:第に": -1611, "BW2:米国": -4267, "BW2:１１": -668, "BW3:あた": -2193, "BW3:あり": 719, "BW3:ある": 3846, "BW3:い.": -1184, "BW3:い。": -1184, "BW3:いい": 5308, "BW3:いえ": 2079, "BW3:いく": 3029, "BW3:いた": 2056, "BW3:いっ": 1883, "BW3:いる": 5600, "BW3:いわ": 1527, "BW3:うち": 1117, "BW3:うと": 4798, "BW3:えと": 1454, "BW3:か.": 2857, "BW3:か。": 2857, "BW3:かけ": -742, "BW3:かっ": -4097, "BW3:かに": -668, "BW3:から": 6520, "BW3:かり": -2669, "BW3:が,": 1816, "BW3:が、": 1816, "BW3:がき": -4854, "BW3:がけ": -1126, "BW3:がっ": -912, "BW3:がら": -4976, "BW3:がり": -2063, "BW3:きた": 1645, "BW3:けど": 1374, "BW3:こと": 7397, "BW3:この": 1542, "BW3:ころ": -2756, "BW3:さい": -713, "BW3:さを": 976, "BW3:し,": 1557, "BW3:し、": 1557, "BW3:しい": -3713, "BW3:した": 3562, "BW3:して": 1449, "BW3:しな": 2608, "BW3:しま": 1200, "BW3:す.": -1309, "BW3:す。": -1309, "BW3:する": 6521, "BW3:ず,": 3426, "BW3:ず、": 3426, "BW3:ずに": 841, "BW3:そう": 428, "BW3:た.": 8875, "BW3:た。": 8875, "BW3:たい": -593, "BW3:たの": 812, "BW3:たり": -1182, "BW3:たる": -852, "BW3:だ.": 4098, "BW3:だ。": 4098, "BW3:だっ": 1004, "BW3:った": -4747, "BW3:って": 300, "BW3:てい": 6240, "BW3:てお": 855, "BW3:ても": 302, "BW3:です": 1437, "BW3:でに": -1481, "BW3:では": 2295, "BW3:とう": -1386, "BW3:とし": 2266, "BW3:との": 541, "BW3:とも": -3542, "BW3:どう": 4664, "BW3:ない": 1796, "BW3:なく": -902, "BW3:など": 2135, "BW3:に,": -1020, "BW3:に、": -1020, "BW3:にし": 1771, "BW3:にな": 1906, "BW3:には": 2644, "BW3:の,": -723, "BW3:の、": -723, "BW3:の子": -999, "BW3:は,": 1337, "BW3:は、": 1337, "BW3:べき": 2181, "BW3:まし": 1113, "BW3:ます": 6943, "BW3:まっ": -1548, "BW3:まで": 6154, "BW3:まれ": -792, "BW3:らし": 1479, "BW3:られ": 6820, "BW3:るる": 3818, "BW3:れ,": 854, "BW3:れ、": 854, "BW3:れた": 1850, "BW3:れて": 1375, "BW3:れば": -3245, "BW3:れる": 1091, "BW3:われ": -604, "BW3:んだ": 606, "BW3:んで": 798, "BW3:カ月": 990, "BW3:会議": 860, "BW3:入り": 1232, "BW3:大会": 2217, "BW3:始め": 1681, "BW3:市": 965, "BW3:新聞": -5054, "BW3:日,": 974, "BW3:日、": 974, "BW3:社会": 2024, "BW3:ｶ月": 990, "TC1:AAA": 1093, "TC1:HHH": 1029, "TC1:HHM": 580, "TC1:HII": 998, "TC1:HOH": -389, "TC1:HOM": -330, "TC1:IHI": 1169, "TC1:IOH": -141, "TC1:IOI": -1014, "TC1:IOM": 467, "TC1:MMH": 187, "TC1:OOI": -1831, "TC2:HHO": 2088, "TC2:HII": -1022, "TC2:HMM": -1153, "TC2:IHI": -1964, "TC2:KKH": 703, "TC2:OII": -2648, "TC3:AAA": -293, "TC3:HHH": 346, "TC3:HHI": -340, "TC3:HII": -1087, "TC3:HIK": 731, "TC3:HOH": -1485, "TC3:IHH": 128, "TC3:IHI": -3040, "TC3:IHO": -1934, "TC3:IIH": -824, "TC3:IIM": -1034, "TC3:IOI": -541, "TC3:KHH": -1215, "TC3:KKA": 491, "TC3:KKH": -1216, "TC3:KOK": -1008, "TC3:MHH": -2693, "TC3:MHM": -456, "TC3:MHO": 123, "TC3:MMH": -470, "TC3:NNH": -1688, "TC3:NNO": 662, "TC3:OHO": -3392, "TC4:HHH": -202, "TC4:HHI": 1344, "TC4:HHK": 365, "TC4:HHM": -121, "TC4:HHN": 182, "TC4:HHO": 669, "TC4:HIH": 804, "TC4:HII": 679, "TC4:HOH": 446, "TC4:IHH": 695, "TC4:IHO": -2323, "TC4:IIH": 321, "TC4:III": 1497, "TC4:IIO": 656, "TC4:IOO": 54, "TC4:KAK": 4845, "TC4:KKA": 3386, "TC4:KKK": 3065, "TC4:MHH": -404, "TC4:MHI": 201, "TC4:MMH": -240, "TC4:MMM": 661, "TC4:MOM": 841, "TQ1:BHHH": -226, "TQ1:BHHI": 316, "TQ1:BHIH": -131, "TQ1:BIHH": 60, "TQ1:BIII": 1595, "TQ1:BNHH": -743, "TQ1:BOHH": 225, "TQ1:BOOO": -907, "TQ1:OAKK": 482, "TQ1:OHHH": 281, "TQ1:OHIH": 249, "TQ1:OIHI": 200, "TQ1:OIIH": -67, "TQ2:BIHH": -1400, "TQ2:BIII": -1032, "TQ2:BKAK": -542, "TQ2:BOOO": -5590, "TQ3:BHHH": 478, "TQ3:BHHM": -1072, "TQ3:BHIH": 222, "TQ3:BHII": -503, "TQ3:BIIH": -115, "TQ3:BIII": -104, "TQ3:BMHI": -862, "TQ3:BMHM": -463, "TQ3:BOMH": 620, "TQ3:OHHH": 346, "TQ3:OHHI": 1729, "TQ3:OHII": 997, "TQ3:OHMH": 481, "TQ3:OIHH": 623, "TQ3:OIIH": 1344, "TQ3:OKAK": 2792, "TQ3:OKHH": 587, "TQ3:OKKA": 679, "TQ3:OOHH": 110, "TQ3:OOII": -684, "TQ4:BHHH": -720, "TQ4:BHHM": -3603, "TQ4:BHII": -965, "TQ4:BIIH": -606, "TQ4:BIII": -2180, "TQ4:OAAA": -2762, "TQ4:OAKK": 180, "TQ4:OHHH": -293, "TQ4:OHHI": 2446, "TQ4:OHHO": 480, "TQ4:OHIH": -1572, "TQ4:OIHH": 1935, "TQ4:OIHI": -492, "TQ4:OIIH": 626, "TQ4:OIII": -4006, "TQ4:OKAK": -8155, "TW1:につい": -4680, "TW1:東京都": 2026, "TW2:ある程": -2048, "TW2:いった": -1255, "TW2:ころが": -2433, "TW2:しょう": 3873, "TW2:その後": -4429, "TW2:だって": -1048, "TW2:ていた": 1833, "TW2:として": -4656, "TW2:ともに": -4516, "TW2:もので": 1882, "TW2:一気に": -791, "TW2:初めて": -1511, "TW2:同時に": -8096, "TW2:大きな": -1254, "TW2:対して": -2720, "TW2:社会党": -3215, "TW3:いただ": -1733, "TW3:してい": 1314, "TW3:として": -4313, "TW3:につい": -5482, "TW3:にとっ": -5988, "TW3:に当た": -6246, "TW3:ので,": -726, "TW3:ので、": -726, "TW3:のもの": -599, "TW3:れから": -3751, "TW3:十二月": -2286, "TW4:いう.": 8576, "TW4:いう。": 8576, "TW4:からな": -2347, "TW4:してい": 2958, "TW4:たが,": 1516, "TW4:たが、": 1516, "TW4:ている": 1538, "TW4:という": 1349, "TW4:ました": 5543, "TW4:ません": 1097, "TW4:ようと": -4257, "TW4:よると": 5865, "UC1:A": 484, "UC1:K": 93, "UC1:M": 645, "UC1:O": -504, "UC2:A": 819, "UC2:H": 1059, "UC2:I": 409, "UC2:M": 3987, "UC2:N": 5775, "UC2:O": 646, "UC3:A": -1369, "UC3:I": 2311, "UC4:A": -2642, "UC4:H": 1809, "UC4:I": -1031, "UC4:K": -3449, "UC4:M": 3565, "UC4:N": 3876, "UC4:O": 6646, "UC5:H": 313, "UC5:I": -1237, "UC5:K": -798, "UC5:M": 539, "UC5:O": -830, "UC6:H": -505, "UC6:I": -252, "UC6:K": 87, "UC6:M": 247, "UC6:O": -386, "UP1:O": -213, "UP2:B": 69, "UP2:O": 935, "UP3:B": 189, "UQ1:BH": 21, "UQ1:BI": -11, "UQ1:BK": -98, "UQ1:BN": 142, "UQ1:BO": -55, "UQ1:OH": -94, "UQ1:OI": 477, "UQ1:OK": 410, "UQ1:OO": -2421, "UQ2:BH": 216, "UQ2:BI": 113, "UQ2:OK": 1759, "UQ3:BA": -478, "UQ3:BH": 42, "UQ3:BI": 1913, "UQ3:BK": -7197, "UQ3:BM": 3160, "UQ3:BN": 6427, "UQ3:BO": 14761, "UQ3:OI": -826, "UQ3:ON": -3211, "UW1:,": 156, "UW1:、": 156, "UW1:「": -462, "UW1:あ": -940, "UW1:う": -126, "UW1:が": -552, "UW1:き": 121, "UW1:こ": 505, "UW1:で": -200, "UW1:と": -546, "UW1:ど": -122, "UW1:に": -788, "UW1:の": -184, "UW1:は": -846, "UW1:も": -465, "UW1:や": -469, "UW1:よ": 182, "UW1:ら": -291, "UW1:り": 208, "UW1:れ": 169, "UW1:を": -445, "UW1:ん": -136, "UW1:・": -134, "UW1:主": -401, "UW1:京": -267, "UW1:区": -911, "UW1:午": 871, "UW1:国": -459, "UW1:大": 561, "UW1:委": 729, "UW1:市": -410, "UW1:日": -140, "UW1:理": 361, "UW1:生": -407, "UW1:県": -385, "UW1:都": -717, "UW1:｢": -462, "UW1:･": -134, "UW2:,": -828, "UW2:、": -828, "UW2:〇": 892, "UW2:「": -644, "UW2:」": 3145, "UW2:あ": -537, "UW2:い": 505, "UW2:う": 134, "UW2:お": -501, "UW2:か": 1454, "UW2:が": -855, "UW2:く": -411, "UW2:こ": 1141, "UW2:さ": 878, "UW2:ざ": 540, "UW2:し": 1529, "UW2:す": -674, "UW2:せ": 300, "UW2:そ": -1010, "UW2:た": 188, "UW2:だ": 1837, "UW2:つ": -948, "UW2:て": -290, "UW2:で": -267, "UW2:と": -980, "UW2:ど": 1273, "UW2:な": 1063, "UW2:に": -1763, "UW2:の": 130, "UW2:は": -408, "UW2:ひ": -1272, "UW2:べ": 1261, "UW2:ま": 600, "UW2:も": -1262, "UW2:や": -401, "UW2:よ": 1639, "UW2:り": -578, "UW2:る": -693, "UW2:れ": 571, "UW2:を": -2515, "UW2:ん": 2095, "UW2:ア": -586, "UW2:カ": 306, "UW2:キ": 568, "UW2:ッ": 831, "UW2:三": -757, "UW2:不": -2149, "UW2:世": -301, "UW2:中": -967, "UW2:主": -860, "UW2:事": 492, "UW2:人": -122, "UW2:会": 978, "UW2:保": 362, "UW2:入": 548, "UW2:初": -3024, "UW2:副": -1565, "UW2:北": -3413, "UW2:区": -421, "UW2:大": -1768, "UW2:天": -864, "UW2:太": -482, "UW2:子": -1518, "UW2:学": 760, "UW2:実": 1023, "UW2:小": -2008, "UW2:市": -812, "UW2:年": -1059, "UW2:強": 1067, "UW2:手": -1518, "UW2:揺": -1032, "UW2:政": 1522, "UW2:文": -1354, "UW2:新": -1681, "UW2:日": -1814, "UW2:明": -1461, "UW2:最": -629, "UW2:朝": -1842, "UW2:本": -1649, "UW2:東": -930, "UW2:果": -664, "UW2:次": -2377, "UW2:民": -179, "UW2:気": -1739, "UW2:理": 752, "UW2:発": 529, "UW2:目": -1583, "UW2:相": -241, "UW2:県": -1164, "UW2:立": -762, "UW2:第": 810, "UW2:米": 509, "UW2:自": -1352, "UW2:行": 838, "UW2:西": -743, "UW2:見": -3873, "UW2:調": 1010, "UW2:議": 1198, "UW2:込": 3041, "UW2:開": 1758, "UW2:間": -1256, "UW2:｢": -644, "UW2:｣": 3145, "UW2:ｯ": 831, "UW2:ｱ": -586, "UW2:ｶ": 306, "UW2:ｷ": 568, "UW3:,": 4889, "UW3:1": -799, "UW3:−": -1722, "UW3:、": 4889, "UW3:々": -2310, "UW3:〇": 5827, "UW3:」": 2670, "UW3:〓": -3572, "UW3:あ": -2695, "UW3:い": 1006, "UW3:う": 2342, "UW3:え": 1983, "UW3:お": -4863, "UW3:か": -1162, "UW3:が": 3271, "UW3:く": 1004, "UW3:け": 388, "UW3:げ": 401, "UW3:こ": -3551, "UW3:ご": -3115, "UW3:さ": -1057, "UW3:し": -394, "UW3:す": 584, "UW3:せ": 3685, "UW3:そ": -5227, "UW3:た": 842, "UW3:ち": -520, "UW3:っ": -1443, "UW3:つ": -1080, "UW3:て": 6167, "UW3:で": 2318, "UW3:と": 1691, "UW3:ど": -898, "UW3:な": -2787, "UW3:に": 2745, "UW3:の": 4056, "UW3:は": 4555, "UW3:ひ": -2170, "UW3:ふ": -1797, "UW3:へ": 1199, "UW3:ほ": -5515, "UW3:ま": -4383, "UW3:み": -119, "UW3:め": 1205, "UW3:も": 2323, "UW3:や": -787, "UW3:よ": -201, "UW3:ら": 727, "UW3:り": 649, "UW3:る": 5905, "UW3:れ": 2773, "UW3:わ": -1206, "UW3:を": 6620, "UW3:ん": -517, "UW3:ア": 551, "UW3:グ": 1319, "UW3:ス": 874, "UW3:ッ": -1349, "UW3:ト": 521, "UW3:ム": 1109, "UW3:ル": 1591, "UW3:ロ": 2201, "UW3:ン": 278, "UW3:・": -3793, "UW3:一": -1618, "UW3:下": -1758, "UW3:世": -2086, "UW3:両": 3815, "UW3:中": 653, "UW3:主": -757, "UW3:予": -1192, "UW3:二": 974, "UW3:人": 2742, "UW3:今": 792, "UW3:他": 1889, "UW3:以": -1367, "UW3:低": 811, "UW3:何": 4265, "UW3:作": -360, "UW3:保": -2438, "UW3:元": 4858, "UW3:党": 3593, "UW3:全": 1574, "UW3:公": -3029, "UW3:六": 755, "UW3:共": -1879, "UW3:円": 5807, "UW3:再": 3095, "UW3:分": 457, "UW3:初": 2475, "UW3:別": 1129, "UW3:前": 2286, "UW3:副": 4437, "UW3:力": 365, "UW3:動": -948, "UW3:務": -1871, "UW3:化": 1327, "UW3:北": -1037, "UW3:区": 4646, "UW3:千": -2308, "UW3:午": -782, "UW3:協": -1005, "UW3:口": 483, "UW3:右": 1233, "UW3:各": 3588, "UW3:合": -240, "UW3:同": 3906, "UW3:和": -836, "UW3:員": 4513, "UW3:国": 642, "UW3:型": 1389, "UW3:場": 1219, "UW3:外": -240, "UW3:妻": 2016, "UW3:学": -1355, "UW3:安": -422, "UW3:実": -1007, "UW3:家": 1078, "UW3:小": -512, "UW3:少": -3101, "UW3:州": 1155, "UW3:市": 3197, "UW3:平": -1803, "UW3:年": 2416, "UW3:広": -1029, "UW3:府": 1605, "UW3:度": 1452, "UW3:建": -2351, "UW3:当": -3884, "UW3:得": 1905, "UW3:思": -1290, "UW3:性": 1822, "UW3:戸": -487, "UW3:指": -3972, "UW3:政": -2012, "UW3:教": -1478, "UW3:数": 3222, "UW3:文": -1488, "UW3:新": 1764, "UW3:日": 2099, "UW3:旧": 5792, "UW3:昨": -660, "UW3:時": -1247, "UW3:曜": -950, "UW3:最": -936, "UW3:月": 4125, "UW3:期": 360, "UW3:李": 3094, "UW3:村": 364, "UW3:東": -804, "UW3:核": 5156, "UW3:森": 2438, "UW3:業": 484, "UW3:氏": 2613, "UW3:民": -1693, "UW3:決": -1072, "UW3:法": 1868, "UW3:海": -494, "UW3:無": 979, "UW3:物": 461, "UW3:特": -3849, "UW3:生": -272, "UW3:用": 914, "UW3:町": 1215, "UW3:的": 7313, "UW3:直": -1834, "UW3:省": 792, "UW3:県": 6293, "UW3:知": -1527, "UW3:私": 4231, "UW3:税": 401, "UW3:立": -959, "UW3:第": 1201, "UW3:米": 7767, "UW3:系": 3066, "UW3:約": 3663, "UW3:級": 1384, "UW3:統": -4228, "UW3:総": 1163, "UW3:線": 1255, "UW3:者": 6457, "UW3:能": 725, "UW3:自": -2868, "UW3:英": 785, "UW3:見": 1044, "UW3:調": -561, "UW3:財": -732, "UW3:費": 1777, "UW3:車": 1835, "UW3:軍": 1375, "UW3:込": -1503, "UW3:通": -1135, "UW3:選": -680, "UW3:郎": 1026, "UW3:郡": 4404, "UW3:部": 1200, "UW3:金": 2163, "UW3:長": 421, "UW3:開": -1431, "UW3:間": 1302, "UW3:関": -1281, "UW3:雨": 2009, "UW3:電": -1044, "UW3:非": 2066, "UW3:駅": 1620, "UW3:１": -799, "UW3:｣": 2670, "UW3:･": -3793, "UW3:ｯ": -1349, "UW3:ｱ": 551, "UW3:ｸﾞ": 1319, "UW3:ｽ": 874, "UW3:ﾄ": 521, "UW3:ﾑ": 1109, "UW3:ﾙ": 1591, "UW3:ﾛ": 2201, "UW3:ﾝ": 278, "UW4:,": 3930, "UW4:.": 3508, "UW4:―": -4840, "UW4:、": 3930, "UW4:。": 3508, "UW4:〇": 4999, "UW4:「": 1895, "UW4:」": 3798, "UW4:〓": -5155, "UW4:あ": 4752, "UW4:い": -3434, "UW4:う": -639, "UW4:え": -2513, "UW4:お": 2405, "UW4:か": 530, "UW4:が": 6006, "UW4:き": -4481, "UW4:ぎ": -3820, "UW4:く": -3787, "UW4:け": -4375, "UW4:げ": -4733, "UW4:こ": 2255, "UW4:ご": 1979, "UW4:さ": 2864, "UW4:し": -842, "UW4:じ": -2505, "UW4:す": -730, "UW4:ず": 1251, "UW4:せ": 181, "UW4:そ": 4091, "UW4:た": 5034, "UW4:だ": 5408, "UW4:ち": -3653, "UW4:っ": -5881, "UW4:つ": -1658, "UW4:て": 3994, "UW4:で": 7410, "UW4:と": 4547, "UW4:な": 5433, "UW4:に": 6499, "UW4:ぬ": 1853, "UW4:ね": 1413, "UW4:の": 7396, "UW4:は": 8578, "UW4:ば": 1940, "UW4:ひ": 4249, "UW4:び": -4133, "UW4:ふ": 1345, "UW4:へ": 6665, "UW4:べ": -743, "UW4:ほ": 1464, "UW4:ま": 1051, "UW4:み": -2081, "UW4:む": -881, "UW4:め": -5045, "UW4:も": 4169, "UW4:ゃ": -2665, "UW4:や": 2795, "UW4:ょ": -1543, "UW4:よ": 3351, "UW4:ら": -2921, "UW4:り": -9725, "UW4:る": -14895, "UW4:れ": -2612, "UW4:ろ": -4569, "UW4:わ": -1782, "UW4:を": 13150, "UW4:ん": -2351, "UW4:カ": 2145, "UW4:コ": 1789, "UW4:セ": 1287, "UW4:ッ": -723, "UW4:ト": -402, "UW4:メ": -1634, "UW4:ラ": -880, "UW4:リ": -540, "UW4:ル": -855, "UW4:ン": -3636, "UW4:・": -4370, "UW4:ー": -11869, "UW4:一": -2068, "UW4:中": 2210, "UW4:予": 782, "UW4:事": -189, "UW4:井": -1767, "UW4:人": 1036, "UW4:以": 544, "UW4:会": 950, "UW4:体": -1285, "UW4:作": 530, "UW4:側": 4292, "UW4:先": 601, "UW4:党": -2005, "UW4:共": -1211, "UW4:内": 584, "UW4:円": 788, "UW4:初": 1347, "UW4:前": 1623, "UW4:副": 3879, "UW4:力": -301, "UW4:動": -739, "UW4:務": -2714, "UW4:化": 776, "UW4:区": 4517, "UW4:協": 1013, "UW4:参": 1555, "UW4:合": -1833, "UW4:和": -680, "UW4:員": -909, "UW4:器": -850, "UW4:回": 1500, "UW4:国": -618, "UW4:園": -1199, "UW4:地": 866, "UW4:場": -1409, "UW4:塁": -2093, "UW4:士": -1412, "UW4:多": 1067, "UW4:大": 571, "UW4:子": -4801, "UW4:学": -1396, "UW4:定": -1056, "UW4:寺": -808, "UW4:小": 1910, "UW4:屋": -1327, "UW4:山": -1499, "UW4:島": -2055, "UW4:川": -2666, "UW4:市": 2771, "UW4:年": 374, "UW4:庁": -4555, "UW4:後": 456, "UW4:性": 553, "UW4:感": 916, "UW4:所": -1565, "UW4:支": 856, "UW4:改": 787, "UW4:政": 2182, "UW4:教": 704, "UW4:文": 522, "UW4:方": -855, "UW4:日": 1798, "UW4:時": 1829, "UW4:最": 845, "UW4:月": -9065, "UW4:木": -484, "UW4:来": -441, "UW4:校": -359, "UW4:業": -1042, "UW4:氏": 5388, "UW4:民": -2715, "UW4:気": -909, "UW4:沢": -938, "UW4:済": -542, "UW4:物": -734, "UW4:率": 672, "UW4:球": -1266, "UW4:生": -1285, "UW4:産": -1100, "UW4:田": -2899, "UW4:町": 1826, "UW4:的": 2586, "UW4:目": 922, "UW4:省": -3484, "UW4:県": 2997, "UW4:空": -866, "UW4:立": -2111, "UW4:第": 788, "UW4:米": 2937, "UW4:系": 786, "UW4:約": 2171, "UW4:経": 1146, "UW4:統": -1168, "UW4:総": 940, "UW4:線": -993, "UW4:署": 749, "UW4:者": 2145, "UW4:能": -729, "UW4:般": -851, "UW4:行": -791, "UW4:規": 792, "UW4:警": -1183, "UW4:議": -243, "UW4:谷": -999, "UW4:賞": 730, "UW4:車": -1480, "UW4:軍": 1158, "UW4:輪": -1432, "UW4:込": -3369, "UW4:近": 929, "UW4:道": -1290, "UW4:選": 2596, "UW4:郎": -4865, "UW4:都": 1192, "UW4:野": -1099, "UW4:銀": -2212, "UW4:長": 357, "UW4:間": -2343, "UW4:院": -2296, "UW4:際": -2603, "UW4:電": -877, "UW4:領": -1658, "UW4:題": -791, "UW4:館": -1983, "UW4:首": 1749, "UW4:高": 2120, "UW4:｢": 1895, "UW4:｣": 3798, "UW4:･": -4370, "UW4:ｯ": -723, "UW4:ｰ": -11869, "UW4:ｶ": 2145, "UW4:ｺ": 1789, "UW4:ｾ": 1287, "UW4:ﾄ": -402, "UW4:ﾒ": -1634, "UW4:ﾗ": -880, "UW4:ﾘ": -540, "UW4:ﾙ": -855, "UW4:ﾝ": -3636, "UW5:,": 465, "UW5:.": -298, "UW5:1": -513, "UW5:E2": -32767, "UW5:]": -2761, "UW5:、": 465, "UW5:。": -298, "UW5:「": 363, "UW5:あ": 1655, "UW5:い": 331, "UW5:う": -502, "UW5:え": 1199, "UW5:お": 527, "UW5:か": 647, "UW5:が": -420, "UW5:き": 1624, "UW5:ぎ": 1971, "UW5:く": 312, "UW5:げ": -982, "UW5:さ": -1536, "UW5:し": -1370, "UW5:す": -851, "UW5:だ": -1185, "UW5:ち": 1093, "UW5:っ": 52, "UW5:つ": 921, "UW5:て": -17, "UW5:で": -849, "UW5:と": -126, "UW5:ど": 1682, "UW5:な": -786, "UW5:に": -1223, "UW5:の": -634, "UW5:は": -577, "UW5:べ": 1001, "UW5:み": 502, "UW5:め": 865, "UW5:ゃ": 3350, "UW5:ょ": 854, "UW5:り": -207, "UW5:る": 429, "UW5:れ": 504, "UW5:わ": 419, "UW5:を": -1263, "UW5:ん": 327, "UW5:イ": 241, "UW5:ル": 451, "UW5:ン": -342, "UW5:中": -870, "UW5:京": 722, "UW5:会": -1152, "UW5:党": -653, "UW5:務": 3519, "UW5:区": -900, "UW5:告": 848, "UW5:員": 2104, "UW5:大": -1295, "UW5:学": -547, "UW5:定": 1785, "UW5:嵐": -1303, "UW5:市": -2990, "UW5:席": 921, "UW5:年": 1763, "UW5:思": 872, "UW5:所": -813, "UW5:挙": 1618, "UW5:新": -1681, "UW5:日": 218, "UW5:月": -4352, "UW5:査": 932, "UW5:格": 1356, "UW5:機": -1507, "UW5:氏": -1346, "UW5:田": 240, "UW5:町": -3911, "UW5:的": -3148, "UW5:相": 1319, "UW5:省": -1051, "UW5:県": -4002, "UW5:研": -996, "UW5:社": -277, "UW5:空": -812, "UW5:統": 1955, "UW5:者": -2232, "UW5:表": 663, "UW5:語": -1072, "UW5:議": 1219, "UW5:選": -1017, "UW5:郎": -367, "UW5:長": 786, "UW5:間": 1191, "UW5:題": 2368, "UW5:館": -688, "UW5:１": -513, "UW5:Ｅ２": -32767, "UW5:｢": 363, "UW5:ｲ": 241, "UW5:ﾙ": 451, "UW5:ﾝ": -342, "UW6:,": 227, "UW6:.": 808, "UW6:1": -269, "UW6:E1": 306, "UW6:、": 227, "UW6:。": 808, "UW6:あ": -306, "UW6:う": 189, "UW6:か": 241, "UW6:が": -72, "UW6:く": -120, "UW6:こ": -199, "UW6:じ": 1782, "UW6:す": 383, "UW6:た": -427, "UW6:っ": 573, "UW6:て": -1013, "UW6:で": 101, "UW6:と": -104, "UW6:な": -252, "UW6:に": -148, "UW6:の": -416, "UW6:は": -235, "UW6:も": -205, "UW6:り": 187, "UW6:る": -134, "UW6:を": 195, "UW6:ル": -672, "UW6:ン": -495, "UW6:一": -276, "UW6:中": 201, "UW6:件": -799, "UW6:会": 624, "UW6:前": 302, "UW6:区": 1792, "UW6:員": -1211, "UW6:委": 798, "UW6:学": -959, "UW6:市": 887, "UW6:広": -694, "UW6:後": 535, "UW6:業": -696, "UW6:相": 753, "UW6:社": -506, "UW6:福": 974, "UW6:空": -821, "UW6:者": 1811, "UW6:連": 463, "UW6:郎": 1082, "UW6:１": -269, "UW6:Ｅ１": 306, "UW6:ﾙ": -672, "UW6:ﾝ": -495}

	return s
}

func (s *Segmenter) gettype(str string) string {
	for _, v := range s.chartypes {
		if v.re.MatchString(str) {
			return v.name
		}
	}
	return "O"
}

func (s *Segmenter) Segment(str string) []string {
	o := strings.Split(str, "")
	m := s.model

	ctype := []string{"O", "O", "O"}
	for _, v := range o {
		ctype = append(ctype, s.gettype(v))
	}
	ctype = append(ctype, "O", "O", "O")

	seg := []string{"B3", "B2", "B1"}
	seg = append(seg, o...)
	seg = append(seg, "E1", "E2", "E3")

	var result []string
	word := []string{seg[3]}
	p1, p2, p3 := "U", "U", "U"

	for i := 4; i < len(seg)-3; i++ {
		score := m["BIAS"]
		w1 := seg[i-3]
		w2 := seg[i-2]
		w3 := seg[i-1]
		w4 := seg[i]
		w5 := seg[i+1]
		w6 := seg[i+2]
		c1 := ctype[i-3]
		c2 := ctype[i-2]
		c3 := ctype[i-1]
		c4 := ctype[i]
		c5 := ctype[i+1]
		c6 := ctype[i+2]
		score += m["UP1:"+p1]
		score += m["UP2:"+p2]
		score += m["UP3:"+p3]
		score += m["BP1:"+p1+p2]
		score += m["BP2:"+p2+p3]
		score += m["UW1:"+w1]
		score += m["UW2:"+w2]
		score += m["UW3:"+w3]
		score += m["UW4:"+w4]
		score += m["UW5:"+w5]
		score += m["UW6:"+w6]
		score += m["BW1:"+w2+w3]
		score += m["BW2:"+w3+w4]
		score += m["BW3:"+w4+w5]
		score += m["TW1:"+w1+w2+w3]
		score += m["TW2:"+w2+w3+w4]
		score += m["TW3:"+w3+w4+w5]
		score += m["TW4:"+w4+w5+w6]
		score += m["UC1:"+c1]
		score += m["UC2:"+c2]
		score += m["UC3:"+c3]
		score += m["UC4:"+c4]
		score += m["UC5:"+c5]
		score += m["UC6:"+c6]
		score += m["BC1:"+c2+c3]
		score += m["BC2:"+c3+c4]
		score += m["BC3:"+c4+c5]
		score += m["TC1:"+c1+c2+c3]
		score += m["TC2:"+c2+c3+c4]
		score += m["TC3:"+c3+c4+c5]
		score += m["TC4:"+c4+c5+c6]
		score += m["UQ1:"+p1+c1]
		score += m["UQ2:"+p2+c2]
		score += m["UQ3:"+p3+c3]
		score += m["BQ1:"+p2+c2+c3]
		score += m["BQ2:"+p2+c3+c4]
		score += m["BQ3:"+p3+c2+c3]
		score += m["BQ4:"+p3+c3+c4]
		score += m["TQ1:"+p2+c1+c2+c3]
		score += m["TQ2:"+p2+c2+c3+c4]
		score += m["TQ3:"+p3+c1+c2+c3]
		score += m["TQ4:"+p3+c2+c3+c4]
		p := "O"
		if score > 0 {
			result = append(result, strings.Join(word, ""))
			word = word[0:0]
			p = "S"
		}
		p1, p2, p3 = p2, p3, p
		word = append(word, seg[i])
	}
	result = append(result, strings.Join(word, ""))

	return result
}
