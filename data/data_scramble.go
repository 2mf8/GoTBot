package data

import (
	"io/ioutil"
	"net/http"
)

type Scramble struct {
	Instruction string
	ShortName   string
	ShowName    string
}

func Tnoodle(scramble string) (s Scramble) {
	switch scramble {
	case "2":
		return Scramble{"2", "222", "2阶"}
	case "3":
		return Scramble{"3", "333", "3阶"}
	case "4":
		return Scramble{"4", "444", "4阶"}
	case "5":
		return Scramble{"5", "555", "5阶"}
	case "6":
		return Scramble{"6", "666", "6阶"}
	case "7":
		return Scramble{"7", "777", "7阶"}
	case "py":
		return Scramble{"py", "pyram", "pyram"}
	case "sk":
		return Scramble{"sk", "skewb", "skewb"}
	case "sq":
		return Scramble{"sq", "sq1", "sq1"}
	case "cl":
		return Scramble{"cl", "clock", "clock"}
	case "mx":
		return Scramble{"mx", "minx", "minx"}
	case "fm":
		return Scramble{"fm", "333fm", "333fm"}
	case "minx":
		return Scramble{"minx", "minx", "minx"}
	default:
		return Scramble{"instruction", "shortName", "showName"}
	}
}

func GetScramble(s string) string {
	resp, err := http.Get("http://localhost:2014/scramble/.txt?=" + s)
	if err != nil {
		return "获取失败"
	}
	v, _ := ioutil.ReadAll(resp.Body)
	return string(v)
}
