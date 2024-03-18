package database

import "strings"

var eventMap = map[string]string{
	"2":     "222",
	"3":     "333",
	"4":     "444",
	"5":     "555",
	"6":     "666",
	"7":     "777",
	"sk":    "skewb",
	"py":    "pyram",
	"sq":    "sq1",
	"cl":    "clock",
	"mx":    "minx",
	"fm":    "333fm",
	"222":   "222",
	"333":   "333",
	"444":   "444",
	"555":   "555",
	"666":   "666",
	"777":   "777",
	"skewb": "skewb",
	"pyram": "pyram",
	"sq1":   "sq1",
	"clock": "clock",
	"minx":  "minx",
	"333fm": "333fm",
}

var regionMap = map[string]string{
	"wr":  "wr",  // world record
	"nr":  "nr",  // record
	"asr": "asr", // asia record
	"afr": "afr", // africa record
	"er":  "er",  // europe record
	"nar": "nar", // north america record
	"sar": "sar", // south america record
	"ocr": "ocr", // oceania record
}

var typeMap = map[string]string{
	"sin":     "sin", // single result
	"avg":     "avg", // average result
	"single":  "sin",
	"average": "avg",
	"最佳":      "sin",
	"平均":      "avg",
}

var genderMap = map[string]string{
	"m":      "m",
	"f":      "f",
	"all":    "all",
	"male":   "m",
	"female": "f",
	"男":      "m",
	"女":      "f",
}

func ToGetEvent(s string) string {
	tgc, ok := eventMap[s]
	if !ok {
		tgc = ""
	}
	return tgc
}

func ToGetRegion(s string) string {
	tgc, ok := regionMap[s]
	if !ok {
		tgc = ""
	}
	return tgc
}

func ToGetType(s string) string {
	tgc, ok := typeMap[s]
	if !ok {
		tgc = ""
	}
	return tgc
}

func ToGetGender(s string) string {
	tgc, ok := genderMap[s]
	if !ok {
		tgc = ""
	}
	return tgc
}

func ToGetRankInfo(s string) (event, region, rtype, gender string) {
	e := ""
	r := ""
	t := ""
	g := ""
	strs := strings.Split(s, " ")
	for _, v := range strs {
		if e == "" {
			e = ToGetEvent(strings.TrimSpace(v))
		}
		if r == "" {
			r = ToGetRegion(strings.TrimSpace(v))
		}
		if t == "" {
			t = ToGetType(strings.TrimSpace(v))
		}
		if g == "" {
			g = ToGetGender(strings.TrimSpace(v))
		}
	}
	if e == "" {
		e = "333"
	}
	if r == "" {
		r = "nr"
	}
	if t == "" {
		t = "sin"
	}
	if g == "" {
		g = "all"
	}
	return e, r, t, g
}
