package database

import (
	"fmt"
	"io"
	"net/http"
)

type Scramble struct {
	Instruction string
	ShortName   string
	ShowName    string
}

func Tnoodle(scramble string) (s Scramble) {
	switch scramble {
	case "2", "222":
		s = Scramble{"2", "222", "2阶"}
	case "3", "333":
		s = Scramble{"3", "333", "3阶"}
	case "4", "444":
		s = Scramble{"4", "444", "4阶"}
	case "5", "555":
		s = Scramble{"5", "555", "5阶"}
	case "6", "666":
		s = Scramble{"6", "666", "6阶"}
	case "7", "777":
		s = Scramble{"7", "777", "7阶"}
	case "py", "pyram":
		s = Scramble{"py", "pyram", "pyram"}
	case "sk", "skewb":
		s = Scramble{"sk", "skewb", "skewb"}
	case "sq", "sq1":
		s = Scramble{"sq", "sq1", "sq1"}
	case "cl", "clock":
		s = Scramble{"cl", "clock", "clock"}
	case "mx", "minx":
		s = Scramble{"mx", "minx", "minx"}
	case "fm", "333fm":
		s = Scramble{"fm", "333fm", "333fm"}
	default:
		s = Scramble{"instruction", "shortName", "showName"}
	}
	return
}

func GetScramble(s string) string {
	url := fmt.Sprintf("%s/scramble/.txt?=%s", AllConfig.ScrambleServer, s)
	resp, err := http.Get(url)
	if err != nil {
		return "获取失败"
	}
	v, _ := io.ReadAll(resp.Body)
	return string(v)
}
