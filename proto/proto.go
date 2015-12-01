package proto

import (
	"fmt"
	"github.com/toophy/easybuff/help"
	"os"
	"os/exec"
	"sort"
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
	case "string":
		func_name = "WriteString"
	}
	return func_name
}

func TypeToGolang(s string) string {
	type_name := s
	switch s {
	case "int8":
		type_name = "int8"
	case "int16":
		type_name = "int16"
	case "int24":
		type_name = "int32"
	case "int32":
		type_name = "int32"
	case "int40":
		type_name = "int64"
	case "int48":
		type_name = "int64"
	case "int56":
		type_name = "int64"
	case "int64":
		type_name = "int64"
	case "uint8":
		type_name = "uint8"
	case "uint16":
		type_name = "uint16"
	case "uint24":
		type_name = "uint32"
	case "uint32":
		type_name = "uint32"
	case "uint40":
		type_name = "uint64"
	case "uint48":
		type_name = "uint64"
	case "uint56":
		type_name = "uint64"
	case "uint64":
		type_name = "uint64"
	case "string":
		type_name = "string"
	}
	return type_name
}

func TypeToClang(s string) string {
	type_name := s
	switch s {
	case "int8":
		type_name = "char"
	case "int16":
		type_name = "short"
	case "int24":
		type_name = "int"
	case "int32":
		type_name = "int"
	case "int40":
		type_name = "long long"
	case "int48":
		type_name = "long long"
	case "int56":
		type_name = "long long"
	case "int64":
		type_name = "long long"
	case "uint8":
		type_name = "unsigned char"
	case "uint16":
		type_name = "unsigned short"
	case "uint24":
		type_name = "unsigned int"
	case "uint32":
		type_name = "unsigned int"
	case "uint40":
		type_name = "unsigned long long"
	case "uint48":
		type_name = "unsigned long long"
	case "uint56":
		type_name = "unsigned long long"
	case "uint64":
		type_name = "unsigned long long"
	case "string":
		type_name = "std::string"
	}
	return type_name
}

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
	MsgId   string
	Members map[string]*EB_MsgMember
}

// 枚举
type EB_Enum struct {
	EB_Base
	Value string
}

// 引用文件
type EB_Include struct {
	EB_Base
	Value string
}

type EB_ParseTable struct {
	Cells    map[string]interface{}
	CurrCell string
	Comment  []string
	MsgId    string
}

func ParseToNewGolang(d string, fd string, f string) {
	// 结构,枚举唯一
	var table EB_ParseTable
	table.Cells = make(map[string]interface{}, 1000)
	table.Comment = make([]string, 20)

	rows := strings.Split(d, "\n")

	for k, _ := range rows {
		ParseToNewGolangRow(k, rows[k], &table)
	}

	// 写出代码

	// 检查log目录
	if !help.IsExist(fd) {
		os.MkdirAll(fd, os.ModeDir)
	}

	target_file := fd + "/" + f

	if help.IsExist(target_file) {
		os.Remove(target_file)
	}
	file, err := os.OpenFile(target_file, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	// 文件头
	file.WriteString(
		`// easybuff
		// 不要修改本文件, 每次消息有变动, 请手动生成本文件
		// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp)


		package proto

		`)

	// 按照字母顺序
	keys := make([]string, len(table.Cells))
	i := 0
	for key, _ := range table.Cells {
		keys[i] = key
		i++
	}

	sort.Sort(sort.StringSlice(keys))

	// 引用文件

	file.WriteString("import (\n")

	for _, key := range keys {

		switch table.Cells[key].(type) {
		case *EB_Include:
			d := table.Cells[key].(*EB_Include)
			file.WriteString(fmt.Sprintf(".\"%s\"\n", d.Value))
		}
	}
	file.WriteString(")\n\n")

	// 枚举
	for _, key := range keys {

		switch table.Cells[key].(type) {
		case *EB_Enum:
			d := table.Cells[key].(*EB_Enum)
			for k, _ := range d.Comment {
				if len(d.Comment[k]) > 0 {
					file.WriteString(fmt.Sprintf("// %s\n", d.Comment[k]))
				}
			}
			file.WriteString(fmt.Sprintf("const %s = %s \n", d.Name, d.Value))
		}
	}
	file.WriteString("\n\n")

	for _, key := range keys {

		switch table.Cells[key].(type) {
		case *EB_Message:
			d := table.Cells[key].(*EB_Message)

			for k, _ := range d.Comment {
				if len(d.Comment[k]) > 0 {
					file.WriteString(fmt.Sprintf("// %s\n", d.Comment[k]))
				}
			}

			if len(d.MsgId) > 0 {
				file.WriteString(fmt.Sprintf("const %s_Id = %s\n", d.Name, d.MsgId))
			}

			file.WriteString(fmt.Sprintf("type %s struct {\n", d.Name))
			// 循环, 顺序输出
			// 按照字母顺序
			mbkeys := make([]string, len(d.Members))
			for k, _ := range d.Members {
				mbkeys[d.Members[k].Sort] = k
			}

			for _, v := range mbkeys {
				m := d.Members[v]
				if len(m.Range) <= 0 {
					file.WriteString(fmt.Sprintf("%s %s // %s\n", m.Name, TypeToGolang(m.Type), m.Desc))
				} else {
					file.WriteString(fmt.Sprintf("%s []%s // %s\n", m.Name, TypeToGolang(m.Type), m.Desc))
				}
			}

			file.WriteString("}\n")

			// read
			file.WriteString(fmt.Sprintf("func (t *%s) Read(s *Stream) {\n", d.Name))
			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetReadFunc(m.Type)
				if fn == "" {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" t.%s.Read(s)\n", m.Name))
					} else {
						file.WriteString(fmt.Sprintf(`
													  len_%s := int(s.%s())
													  t.%s = make([]%s,len_%s)
													  for i:=0;i<len_%s;i++ {
													  	t.%s[i].Read(s)
													  }
													  
													  `, m.Name, GetReadFunc("uint8"), m.Name, TypeToGolang(m.Type), m.Name, m.Name, m.Name))
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" t.%s = s.%s()\n", m.Name, fn))
					} else {
						file.WriteString(fmt.Sprintf(`
													  len_%s := int(s.%s())
													  t.%s = make([]%s,len_%s)
													  for i:=0;i<len_%s;i++ {
													  	 t.%s = s.%s()
													  }
													  
													  `, m.Name, GetReadFunc("uint8"), m.Name, TypeToGolang(m.Type), m.Name, m.Name, m.Name, fn))
					}
				}
			}

			file.WriteString("}\n\n")

			// write
			file.WriteString(fmt.Sprintf("func (t *%s) Write(s *Stream) {\n", d.Name))
			// write len, msg_id
			// old_pos := s.GetPos()
			// s.Seek(old_pos+help.MsgHeaderSize)
			// ...
			// last_pos := s.GetPos()
			// s.Seek(old_pos)
			// header := len<<16 | uint32(x_Id)
			// s.WriteUint32(header)
			// s.Seek(last_pos)
			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetWriteFunc(m.Type)
				if fn == "" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" t.%s.Write(s)", m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  s.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	t.%s[k].Write(s)
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  s.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	t.%s[i].Write(s)
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, m.Name))
						}
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" s.%s(&t.%s)\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  s.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	s.%s(&t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  s.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	s.%s(&t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" s.%s(t.%s)\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  s.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	s.%s(t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  s.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	s.%s(t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				}
			}

			file.WriteString("}\n\n\n")
		}
	}

	file.Close()

	cmd_data := exec.Command("gofmt", "-w", target_file)
	err = cmd_data.Run()
	if err != nil {
		fmt.Println(target_file + "," + err.Error())
	}
}

func ParseToNewGolangRow(row_id int, d string, table *EB_ParseTable) {

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
	case "message":
		if lens < 2 {
			panic("文件格式错误 : message 行错误 [" + r3 + "]")
		}

		// message name id {
		t := &EB_Message{}
		t.Name = m[1]

		if len(table.CurrCell) > 0 {
			panic("文件格式错误 : message [" + table.CurrCell + "]还没有结束定义.")
		}

		if _, ok := table.Cells[t.Name]; ok {
			panic("文件内容错误 : message 重名 [" + r3 + "]")
		}

		table.CurrCell = t.Name
		table.Cells[table.CurrCell] = t
		t.Comment = make([]string, 10)
		if len(table.Comment) > 0 {
			for k, _ := range table.Comment {
				if len(table.Comment[k]) > 0 {
					t.Comment = append(t.Comment, table.Comment[k])
				}
			}
			table.Comment = table.Comment[0:0]
		}

		t.MsgId = table.MsgId
		table.MsgId = ""

		t.Members = make(map[string]*EB_MsgMember, 10)

	case "enum":
		if lens < 3 {
			panic("文件格式错误 : enum 行错误 [" + r3 + "]")
		}

		t := &EB_Enum{}
		t.Name = m[1]
		t.Value = m[2]

		if _, ok := table.Cells[t.Name]; ok {
			panic("文件内容错误 : enum 重名 [" + r3 + "]")
		}

		table.Cells[t.Name] = t

		t.Comment = make([]string, 10)
		if len(table.Comment) > 0 {
			for k, _ := range table.Comment {
				if len(table.Comment[k]) > 0 {
					t.Comment = append(t.Comment, table.Comment[k])
				}
			}
			table.Comment = table.Comment[0:0]
		}

	case "include":
		if lens < 3 {
			panic("文件格式错误 : include 行错误 [" + r3 + "]")
		}

		t := &EB_Include{}
		t.Name = m[1]
		t.Value = m[2]

		if _, ok := table.Cells[t.Name]; ok {
			panic("文件内容错误 : include 重名 [" + r3 + "]")
		}

		table.Cells[t.Name] = t

		t.Comment = make([]string, 10)
		if len(table.Comment) > 0 {
			for k, _ := range table.Comment {
				if len(table.Comment[k]) > 0 {
					t.Comment = append(t.Comment, table.Comment[k])
				}
			}
			table.Comment = table.Comment[0:0]
		}

	case "}":
		// message 结束符号
		if len(table.CurrCell) == 0 {
			panic("文件格式错误 : 多余的结束符 } .")
		}
		table.CurrCell = ""

	case "--":
		// 注释行, 本行注释, 作用给下一行
		if lens > 1 {
			table.Comment = append(table.Comment, m[1])
		}

	case "id":
		// 消息ID
		if lens > 1 {
			table.MsgId = m[1]
		}

	default:
		if lens < 3 {
			panic("文件格式错误 : member 行错误 [" + r3 + "]")
		}
		if len(table.CurrCell) == 0 {
			panic("文件格式错误 : member 行错误 [" + r3 + "], 没有归属消息")
		}
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

		if _, ok := table.Cells[table.CurrCell].(*EB_Message).Members[mb.Name]; ok {
			panic("文件格式错误 : member 重名 [" + r3 + "] ")
		}

		mb.Sort = len(table.Cells[table.CurrCell].(*EB_Message).Members)

		table.Cells[table.CurrCell].(*EB_Message).Members[mb.Name] = mb
	}
}
