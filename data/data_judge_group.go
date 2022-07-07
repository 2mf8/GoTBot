package data

import(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/gomodule/redigo/redis"
)

type JudgeGroup struct {
	Groups []int64
}

type JudgeGroupSync struct {
	IsTrue        bool
	JudgeGroupSync *JudgeGroup
}

func JudgeGroupId(groupId int64, group JudgeGroup) int64 {
	for _, j_group_id := range group.Groups {
		if j_group_id == groupId {
			return j_group_id
		}
	}
	return 0
}

func JudgeGroupIndex(groupId int64, group JudgeGroup) int {
	for i, j_group_id := range group.Groups {
		if j_group_id == groupId {
			return i
		}
	}
	return -1
}

func GetJudgeGroup() (group JudgeGroupSync, err error) {
	JudgeGroup := JudgeGroup{}
	group = JudgeGroupSync{
		IsTrue: true,
		JudgeGroupSync: &JudgeGroup,
	}
	var vb []byte
	var bw_set []byte

	bw := "JudgeGroup"
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", bw)
	c.Flush()
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		fmt.Println("[查询] 首次查询-守卫", bw)
		jgroup, err := JudgeGroupRead()
		group.JudgeGroupSync = &jgroup
		if err != nil {
			group = JudgeGroupSync{
				IsTrue: false,
				JudgeGroupSync: &JudgeGroup,
			}
			group.JudgeGroupSync.JudgeGroupCreate()
		}
		bw_set, _ = json.Marshal(&group)
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return group, err
	}
	err = json.Unmarshal(vb, &group)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	fmt.Println("[Redis] groupey(", bw, ") Value(", group.IsTrue, *group.JudgeGroupSync, ")")  //测试用
	return
}

func (group *JudgeGroup) JudgeGroupCreate() error {
	output, err := json.MarshalIndent(&group, "", "\t")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return err
	}
	err = ioutil.WriteFile("JudgeGroup.json", output, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
		return err
	}
	return nil
}

func JudgeGroupRead() (group JudgeGroup, err error) {
	jsonFile, err := os.Open("JudgeGroup.json")
	if err != nil {
		fmt.Println("Error reading JSON File:", err)
		return
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	json.Unmarshal(jsonData, &group)
	//fmt.Println(group)
	return
}

func (group *JudgeGroupSync) JudgeGroupUpdate(ug ...int64) error {
	for _, v := range ug {
		if JudgeGroupId(v, *group.JudgeGroupSync) == 0 && v != 0 {
			group.JudgeGroupSync.Groups = append(group.JudgeGroupSync.Groups, v)
		}
	}
	bw := "JudgeGroup"
	var bw_set []byte
	JudgeGroupSync := JudgeGroupSync{
		IsTrue: true,
		JudgeGroupSync: group.JudgeGroupSync,
	}
	bw_set, _ = json.Marshal(&JudgeGroupSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Sprintf("%#v", v)
	err = group.JudgeGroupSync.JudgeGroupCreate()
	return err
}

func (group *JudgeGroupSync) JudgeGroupDelete(dgroup ...int64) {
	for _, v := range dgroup {
		if v == 0 {
			continue
		}
		i := JudgeGroupIndex(v, *group.JudgeGroupSync)
		if i != -1 {
			if group.JudgeGroupSync.Groups[i+1:] != nil {
				group.JudgeGroupSync.Groups = append(group.JudgeGroupSync.Groups[:i], group.JudgeGroupSync.Groups[i+1:]...)
				i--
			}
		}
		bw := "JudgeGroup"
		var bw_set []byte
		JudgeGroupSync := JudgeGroupSync{
			IsTrue: true,
			JudgeGroupSync: group.JudgeGroupSync,
		}
		bw_set, _ = json.Marshal(&JudgeGroupSync)
		c := Pool.Get()
		defer c.Close()
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, err := c.Receive()
		if err != nil {
			fmt.Println("[错误] Receive出错")
		}
		fmt.Sprintf("%#v", v)
		group.JudgeGroupSync.JudgeGroupCreate()
	}
}
