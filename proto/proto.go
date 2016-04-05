package proto

import (
	"fmt"
	"github.com/toophy/easybuff/help"
	"strconv"
	"strings"
)

func GetReadFunc(s string) string {
	func_name := ""
	switch s {
	case "int8":
		func_name = "ReadInt8"
	case "int16":
		func_name = "ReadInt16"
	case "int24":
		func_name = "ReadInt24"
	case "int32":
		func_name = "ReadInt32"
	case "int40":
		func_name = "ReadInt40"
	case "int48":
		func_name = "ReadInt48"
	case "int56":
		func_name = "ReadInt56"
	case "int64":
		func_name = "ReadInt64"
	case "uint8":
		func_name = "ReadUint8"
	case "uint16":
		func_name = "ReadUint16"
	case "uint24":
		func_name = "ReadUint24"
	case "uint32":
		func_name = "ReadUint32"
	case "uint40":
		func_name = "ReadUint40"
	case "uint48":
		func_name = "ReadUint48"
	case "uint56":
		func_name = "ReadUint56"
	case "uint64":
		func_name = "ReadUint64"
	case "float32":
		func_name = "ReadFloat32"
	case "float64":
		func_name = "ReadFloat64"
	case "string":
		func_name = "ReadString"
	}
	return func_name
}

func GetWriteFunc(s string) string {
	func_name := ""
	switch s {
	case "int8":
		func_name = "WriteInt8"
	case "int16":
		func_name = "WriteInt16"
	case "int24":
		func_name = "WriteInt24"
	case "int32":
		func_name = "WriteInt32"
	case "int40":
		func_name = "WriteInt40"
	case "int48":
		func_name = "WriteInt48"
	case "int56":
		func_name = "WriteInt56"
	case "int64":
		func_name = "WriteInt64"
	case "uint8":
		func_name = "WriteUint8"
	case "uint16":
		func_name = "WriteUint16"
	case "uint24":
		func_name = "WriteUint24"
	case "uint32":
		func_name = "WriteUint32"
	case "uint40":
		func_name = "WriteUint40"
	case "uint48":
		func_name = "WriteUint48"
	case "uint56":
		func_name = "WriteUint56"
	case "uint64":
		func_name = "WriteUint64"
	case "float32":
		func_name = "WriteFloat32"
	case "float64":
		func_name = "WriteFloat64"
	case "string":
		func_name = "WriteString"
	}
	return func_name
}

const (
	NullMessage = "!!!****空消息占位使用****!!!"
)

// mesasge  ActorBase     1            角色基础信息 {
//   Name   string        名称
//   Age    int8          年龄
//   Hp     int24         血量
//   Mp     int32         法力
//   Maxhp  int24         最大血量
//   Maxmp  int32         最大法力
//   Exp    int48         经验值
//   MaxExp int56         最多经验值
// }

// enum BagMaxGrid        100          背包格子数量上限
// enum EquipMaxCount     20           当前装备上限
// enum SkillMaxCount     100          技能上限

// 消息成员
type EB_MsgMember struct {
	Sort  int
	Name  string
	Type  string
	Range string
	Desc  string
}

// 类型
type EB_Base struct {
	Name    string   // 名称, 文件夹内, 所有文件中, 不得重名
	File    string   // 所在文件名, 便于后期输出
	Comment []string // 本单元注释
}

// 消息定义
type EB_Message struct {
	EB_Base
	MsgId   int
	Members map[string]*EB_MsgMember
}

// 枚举
type EB_Enum struct {
	EB_Base
	Value string
}

type EB_ParseTable struct {
	Cells    map[string]interface{}
	MsgIds   map[int]string
	CurrCell string
	Comment  []string
	MsgId    int // 按照文件start_id为起始规则递增, 但是整体受MsgIds控制, 不能有任何重复, 否则生成协议失败
	Files    map[string]*EB_FileTable
}

type EB_Include struct {
	Name    string
	Comment []string
	Index   int
}

type EB_FileTable struct {
	FileDir      string
	FileName     string
	MyStructIds  map[int]string
	MyMsgIds     map[int]string
	MyInclude    map[string]*EB_Include
	MyIncludeIdx map[int]*EB_Include
	StructIndex  int // 结构序号
	IncludeIndex int // 引用文件序号
	MsgId        int // 按照文件start_id为起始规则递增, 但是整体受MsgIds控制, 不能有任何重复, 否则生成协议失败
}

var g_ParseTable EB_ParseTable

func init() {
	g_ParseTable.Cells = make(map[string]interface{}, 1000)
	g_ParseTable.MsgIds = make(map[int]string, 1000)
	g_ParseTable.Comment = make([]string, 20)
	g_ParseTable.Files = make(map[string]*EB_FileTable, 100)
}

// MsgId
// 每个epd文件, MsgId按照文件中顺序递增, 空消息[message {}]可以占位,
// 有结构名的消息不认为是空消息, epd文件中"start_id"关键字后面的数值
// 是本文件起始消息Id, 不得随意修改, 修改后, 要求所有管理系统重新编译

func ParseToNewGolang(d string, fd string, f string) {
	// 结构,枚举唯一

	rows := strings.Split(d, "\n")

	file_table := new(EB_FileTable)
	file_table.FileDir = fd
	file_table.FileName = f
	file_table.MyMsgIds = make(map[int]string, 100)
	file_table.MyInclude = make(map[string]*EB_Include, 100)
	file_table.MyIncludeIdx = make(map[int]*EB_Include, 100)
	file_table.MyStructIds = make(map[int]string, 100)

	g_ParseTable.Files[fd+"/"+f] = file_table

	for k, _ := range rows {
		ParseToNewGolangRow(k, rows[k], file_table)
	}
}

func ParseToNewGolangRow(row_id int, d string, table *EB_FileTable) {

	// 捕捉异常
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				println("ParseToNewGolangRow:" + r.(error).Error())
			case string:
				println("ParseToNewGolangRow:" + help.Utf82Gbk(r.(string)))
			}
		}
	}()

	// mesasge -- 规则解释, } 结束符
	// enum    -- 规则解释, } 结束符

	r1 := strings.Replace(d, "\t", " ", -1)
	r2 := strings.Replace(r1, "\r\n", " ", -1)
	r3 := strings.Replace(r2, "\n", " ", -1)

	m := strings.Fields(r3)
	lens := len(m)

	if lens < 1 {
		return
	}

	// message
	switch m[0] {
	case "struct":
		// 结构, 没有msg_id的结构, 其他都一样
		if lens < 2 {
			panic("文件格式错误 : struct 行错误 [" + r3 + "]")
		}

		// struct name {
		t := &EB_Message{}
		t.Name = m[1]

		if len(g_ParseTable.CurrCell) > 0 {
			panic("文件格式错误 : message [" + g_ParseTable.CurrCell + "]还没有结束定义.")
		}

		if _, ok := g_ParseTable.Cells[t.Name]; ok {
			panic("文件内容错误 : message 重名 [" + r3 + "]")
		}

		g_ParseTable.CurrCell = t.Name
		g_ParseTable.Cells[g_ParseTable.CurrCell] = t
		t.Comment = make([]string, 10)
		if len(g_ParseTable.Comment) > 0 {
			for k, _ := range g_ParseTable.Comment {
				if len(g_ParseTable.Comment[k]) > 0 {
					t.Comment = append(t.Comment, g_ParseTable.Comment[k])
				}
			}
			g_ParseTable.Comment = g_ParseTable.Comment[0:0]
		}

		t.MsgId = 0

		t.Members = make(map[string]*EB_MsgMember, 10)

		table.MyStructIds[table.StructIndex] = t.Name
		table.StructIndex++

	case "message":
		if lens < 2 {
			panic("文件格式错误 : message 行错误 [" + r3 + "]")
		}

		// message name id {
		t := &EB_Message{}
		t.Name = m[1]

		if m[1] == "{" {
			t.Name = NullMessage
		}

		if len(g_ParseTable.CurrCell) > 0 {
			panic("文件格式错误 : message [" + g_ParseTable.CurrCell + "]还没有结束定义.")
		}

		if _, ok := g_ParseTable.Cells[t.Name]; ok {
			panic("文件内容错误 : message 重名 [" + r3 + "]")
		}

		if v, ok := g_ParseTable.MsgIds[table.MsgId]; ok {
			panic("文件内容错误 : message Id重复 [" + r3 + "] => " + v)
		}

		g_ParseTable.CurrCell = t.Name
		if t.Name == NullMessage {
			g_ParseTable.Comment = g_ParseTable.Comment[0:0]
		} else {
			g_ParseTable.Cells[g_ParseTable.CurrCell] = t
			t.Comment = make([]string, 10)
			if len(g_ParseTable.Comment) > 0 {
				for k, _ := range g_ParseTable.Comment {
					if len(g_ParseTable.Comment[k]) > 0 {
						t.Comment = append(t.Comment, g_ParseTable.Comment[k])
					}
				}
				g_ParseTable.Comment = g_ParseTable.Comment[0:0]
			}

			t.MsgId = table.MsgId

			t.Members = make(map[string]*EB_MsgMember, 10)

			g_ParseTable.MsgIds[t.MsgId] = t.Name
			table.MyMsgIds[t.MsgId] = t.Name
		}
		table.MsgId++

	case "enum":
		if len(g_ParseTable.CurrCell) > 0 {
			panic(fmt.Sprintf("文件格式错误 : enum 行错误 [%s], 不能在消息内部[%s]", r3, g_ParseTable.CurrCell))
		}

		if lens < 3 {
			panic("文件格式错误 : enum 行错误 [" + r3 + "]")
		}

		t := &EB_Enum{}
		t.Name = m[1]
		t.Value = m[2]

		if _, ok := g_ParseTable.Cells[t.Name]; ok {
			panic("文件内容错误 : enum 重名 [" + r3 + "]")
		}

		g_ParseTable.Cells[t.Name] = t

		t.Comment = make([]string, 10)
		if len(g_ParseTable.Comment) > 0 {
			for k, _ := range g_ParseTable.Comment {
				if len(g_ParseTable.Comment[k]) > 0 {
					t.Comment = append(t.Comment, g_ParseTable.Comment[k])
				}
			}
			g_ParseTable.Comment = g_ParseTable.Comment[0:0]
		}

	case "include":
		if len(g_ParseTable.CurrCell) > 0 {
			panic(fmt.Sprintf("文件格式错误 : include 行错误 [%s], 不能在消息内部[%s]", r3, g_ParseTable.CurrCell))
		}

		if lens < 2 {
			panic("文件格式错误 : include 行错误 [" + r3 + "]")
		}

		t := new(EB_Include)
		t.Index = table.IncludeIndex
		t.Name = m[1]

		if _, ok := table.MyInclude[t.Name]; ok {
			panic("文件内容错误 : include 重名 [" + r3 + "]")
		}

		t.Comment = make([]string, 10)
		if len(g_ParseTable.Comment) > 0 {
			for k, _ := range g_ParseTable.Comment {
				if len(g_ParseTable.Comment[k]) > 0 {
					t.Comment = append(t.Comment, g_ParseTable.Comment[k])
				}
			}
			g_ParseTable.Comment = g_ParseTable.Comment[0:0]
		}
		table.MyInclude[t.Name] = t
		table.MyIncludeIdx[t.Index] = t
		table.IncludeIndex++

	case "}":
		// message 结束符号
		if len(g_ParseTable.CurrCell) == 0 {
			panic("文件格式错误 : 多余的结束符 } .")
		}
		g_ParseTable.CurrCell = ""

	case "--":
		if len(g_ParseTable.CurrCell) > 0 {
			panic(fmt.Sprintf("文件格式错误 : --注释 行错误 [%s], 不能在消息内部[%s]", r3, g_ParseTable.CurrCell))
		}
		// 注释行, 本行注释, 作用给下一行
		if lens > 1 {
			g_ParseTable.Comment = append(g_ParseTable.Comment, m[1])
		}

	case "start_id":
		if len(g_ParseTable.CurrCell) > 0 {
			panic(fmt.Sprintf("文件格式错误 : start_id 行错误 [%s], 不能在消息内部[%s]", r3, g_ParseTable.CurrCell))
		}
		// 消息ID
		if lens > 1 {
			if table.MsgId != 0 {
				panic("文件格式错误 : 不能有多个start_id [" + r3 + "]")
			} else {
				start_id, err_start_id := strconv.Atoi(m[1])
				if err_start_id != nil {
					panic("文件格式错误 : start_id无效 [" + r3 + "]")
				} else {
					if start_id > 0 && start_id < 65536 {
						table.MsgId = start_id
					} else {
						panic("文件格式错误 : start_id无效 [" + r3 + "]")
					}
				}
			}
		}
		g_ParseTable.Comment = g_ParseTable.Comment[0:0]

	default:
		if lens < 3 {
			panic("文件格式错误 : member 行错误 [" + r3 + "]")
		}
		if len(g_ParseTable.CurrCell) == 0 {
			panic("文件格式错误 : member 行错误 [" + r3 + "], 没有归属消息")
		}
		if g_ParseTable.CurrCell != NullMessage {
			mb := &EB_MsgMember{}
			mb.Name = m[0]
			mb.Type = m[1]
			mb.Desc = m[2]

			// Range
			// [1]Type
			if strings.Contains(mb.Type, "]") {
				mn := strings.Split(mb.Type, "]")
				if len(mn) > 1 {
					mb.Type = mn[1]
					if len(mn[0]) > 1 {
						mb.Range = mn[0][1:]
					} else {
						mb.Range = "--ArrayLen"
					}
				}
			}

			if _, ok := g_ParseTable.Cells[g_ParseTable.CurrCell].(*EB_Message).Members[mb.Name]; ok {
				panic("文件格式错误 : member 重名 [" + r3 + "] ")
			}

			mb.Sort = len(g_ParseTable.Cells[g_ParseTable.CurrCell].(*EB_Message).Members)

			g_ParseTable.Cells[g_ParseTable.CurrCell].(*EB_Message).Members[mb.Name] = mb
		}
	}
}
