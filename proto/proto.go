package proto

import (
	"fmt"
	"github.com/toophy/easybuff/help"
	"os"
	"os/exec"
	"sort"
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
	Name    string // 名称, 文件夹内, 所有文件中, 不得重名
	Desc    string // 描述
	Stop    bool   // 解释结束, 不允许再添加成员
	File    string // 所在文件名, 便于后期输出
	Comment string // 本单元注释
}

// 消息定义
type EB_Message struct {
	EB_Base
	Id      string
	Members map[string]*EB_MsgMember
}

// 枚举
type EB_Enum struct {
	EB_Base
	Value string
}

type EB_ParseTable struct {
	Cells    map[string]interface{}
	CurrCell string
}

func ParseToNewGolang(d string, fd string, f string) {
	// 结构,枚举唯一
	var table EB_ParseTable
	table.Cells = make(map[string]interface{}, 1000)

	rows := strings.Split(d, "\n")

	for k, _ := range rows {
		ParseToNewGolangRow(k, rows[k], &table)
	}

	// 写出代码

	// 检查log目录
	if !help.IsExist(fd) {
		os.MkdirAll(fd, os.ModeDir)
	}

	if !help.IsExist(fd + f) {
		os.Create(fd + f)
	}
	file, err := os.OpenFile(fd+f, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	// 文件头
	file.WriteString(
		`// easybuff
		// 不要修改本文件, 每次消息有变动, 请手动生成本文件
		// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp)


		package proto

		import (
			. "github.com/toophy/login/help"
		)

		`)

	// 按照字母顺序
	keys := make([]string, len(table.Cells))
	i := 0
	for key, _ := range table.Cells {
		keys[i] = key
		i++
	}

	sort.Sort(sort.StringSlice(keys))

	for _, key := range keys {

		switch table.Cells[key].(type) {
		case *EB_Enum:
			d := table.Cells[key].(*EB_Enum)
			file.WriteString(fmt.Sprintf("const %s = %s // %s\n", d.Name, d.Value, d.Desc))
		}
	}

	file.WriteString("\n\n")

	for _, key := range keys {

		switch table.Cells[key].(type) {
		case *EB_Message:
			d := table.Cells[key].(*EB_Message)
			file.WriteString(fmt.Sprintf("// %s\n", d.Desc))
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
					file.WriteString(fmt.Sprintf(" t.%s = s.%s()\n", m.Name, fn))
				}
			}

			file.WriteString("}\n")

			// write
			file.WriteString(fmt.Sprintf("func (t *%s) Write(s *Stream) {\n", d.Name))
			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetWriteFunc(m.Type)
				if fn == "" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" t.%s.Write(s)", m.Name))
					} else {
						rg, err := strconv.ParseInt(m.Range, 10, 32)
						if err != nil {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  s.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	t.%s[i].Write(s)
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, m.Name))
						} else if rg <= 0 {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  s.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	t.%s[k].Write(s)
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, m.Name))
						} else {
							// 限定范围(数值)
							file.WriteString(fmt.Sprintf(`
								  						  len_%s := len(t.%s)
														  s.%s(uint8(%s))
														  for i:=0; i<%s && i<len_%s; i++ {
														  	t.%s[i].Write(s)
														  }
														  
														  `, m.Name, m.Name, GetWriteFunc("uint8"), m.Range, m.Name, m.Name))
						}
					}
				} else if m.Type == "string" {
					file.WriteString(fmt.Sprintf(" s.%s(&t.%s)\n", fn, m.Name))
				} else {
					file.WriteString(fmt.Sprintf(" s.%s(t.%s)\n", fn, m.Name))
				}
			}

			file.WriteString("}\n\n")
		}
	}

	file.Close()

	cmd_data := exec.Command("gofmt", "-w", fd+f)
	err = cmd_data.Run()
	if err != nil {
		fmt.Println(fd + f + "," + err.Error())
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

	// mesasge  ActorBase     1            角色基础信息 {

	// message
	switch m[0] {
	case "message":
		if lens < 4 {
			panic("文件格式错误 : message 行错误 [" + r3 + "]")
		}

		// 正常是 5 个元素
		// 如果是 4 个元素, id 就存在
		// message name desc {
		// 必须存在
		t := &EB_Message{}
		t.Name = m[1]
		t.Stop = false

		if lens < 5 {
			t.Desc = m[2]
		} else {
			if len(table.CurrCell) > 0 {
				panic("文件格式错误 : message [" + table.CurrCell + "]还没有结束定义.")
			}

			t.Id = m[2]
			t.Desc = m[3]
		}

		if _, ok := table.Cells[t.Name]; ok {
			panic("文件内容错误 : message 重名 [" + r3 + "]")
		}

		table.CurrCell = t.Name
		table.Cells[table.CurrCell] = t

		t.Members = make(map[string]*EB_MsgMember, 10)

	case "enum":
		if lens < 4 {
			panic("文件格式错误 : enum 行错误 [" + r3 + "]")
		}

		t := &EB_Enum{}
		t.Stop = true
		t.Name = m[1]
		t.Value = m[2]
		t.Desc = m[3]

		if _, ok := table.Cells[t.Name]; ok {
			panic("文件内容错误 : enum 重名 [" + r3 + "]")
		}

		table.Cells[t.Name] = t

	case "}":
		// message 结束符号
		if len(table.CurrCell) == 0 {
			panic("文件格式错误 : 多余的结束符 } .")
		}
		table.Cells[table.CurrCell].(*EB_Message).Stop = true
		table.CurrCell = ""

	case "--":
		// 注释行, 本行注释, 作用给下一行

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
		if lens == 3 {
			mb.Desc = m[2]
		} else {
			mb.Range = m[2]
			mb.Desc = m[3]
		}

		if _, ok := table.Cells[table.CurrCell].(*EB_Message).Members[mb.Name]; ok {
			panic("文件格式错误 : member 重名 [" + r3 + "] ")
		}

		mb.Sort = len(table.Cells[table.CurrCell].(*EB_Message).Members)

		table.Cells[table.CurrCell].(*EB_Message).Members[mb.Name] = mb
	}
}

func ParseToGolang(d string, fd string, f string) {

	r1 := strings.Replace(d, "\t", " ", -1)
	r2 := strings.Replace(r1, "\r\n", " ", -1)
	r3 := strings.Replace(r2, "\n", " ", -1)
	r4 := strings.Replace(r3, "{", " ", -1)
	r5 := strings.Replace(r4, "}", " ", -1)

	m := strings.Fields(r5)

	lens := len(m)

	data_struct := make([]string, 10)
	data_read := make([]string, 10)
	data_write := make([]string, 10)

	data_count := 0

	flag := "key"
	member := ""

	for i := 0; i < lens; i++ {
		switch m[i] {
		case "message":
			switch flag {
			case "key":
			default:
				if flag == "member" {
					data_count++
				} else {
					fmt.Printf("flag=%s,解析失败[%s]\n", flag, data_struct[data_count])
					data_struct[data_count] = ""
					data_read[data_count] = ""
					data_write[data_count] = ""
				}
			}
			flag = "name"

		default:
			switch flag {

			case "name":
				data_struct[data_count] += "type " + m[i] + " struct {\n"
				data_read[data_count] += "func (this *" + m[i] + ") Read(s *Stream) {\n"
				data_write[data_count] += "func (this *" + m[i] + ") Write(s *Stream) {\n"
				flag = "msgId"

			case "msgId":
				data_struct[data_count] += "type " + m[i] + " struct {\n"
				flag = "desc"

			case "desc":
				data_struct[data_count] += "// " + m[i] + "\n"
				flag = "member"

			case "member":
				member = m[i]
				flag = "type"

			case "type":
				if len(member) > 0 {
					type_name := TypeToGolang(m[i])
					if type_name == "" {
						data_struct[data_count] += "\t" + member + "\t" + m[i] + "// " + m[i] + "\n"
					} else {
						data_struct[data_count] += "\t" + member + "\t" + type_name + "// " + m[i] + "\n"
					}

					func_read := GetReadFunc(m[i])
					if func_read == "" {
						data_read[data_count] += "\t" + "this." + member + ".Read(s)\n"
					} else {
						data_read[data_count] += "\t" + "s." + GetReadFunc(m[i]) + "(this." + member + ")\n"
					}

					func_write := GetWriteFunc(m[i])
					if func_write == "" {
						data_write[data_count] += "\t" + "this." + member + ".Write(s)\n"
					} else {
						data_write[data_count] += "\t" + "s." + GetWriteFunc(m[i]) + "(this." + member + ")\n"
					}
				} else {
					fmt.Printf("flag=%s,解析失败[%s]\n", flag, data_struct[data_count])
				}

				flag = "member"
			}
		}
	}

	for i := 0; i <= data_count; i++ {
		data_struct[i] += "}\n\n"
		data_read[i] += "}\n\n"
		data_write[i] += "}\n\n"
	}

	// 检查log目录
	if !help.IsExist(fd) {
		os.MkdirAll(fd, os.ModeDir)
	}

	if !help.IsExist(fd + f) {
		os.Create(fd + f)
	}
	file, err := os.OpenFile(fd+f, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	// 文件头
	file.WriteString(
		`// easybuff
// 不要修改本文件, 每次消息有变动, 请手动生成本文件
// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp)


package proto

import (
	. "github.com/toophy/login/help"
)

`)

	for i := 0; i <= data_count; i++ {
		file.WriteString(data_struct[i])
		file.WriteString(data_read[i])
		file.WriteString(data_write[i])
	}

	file.Close()

	cmd_data := exec.Command("gofmt", "-w", fd+f)
	err = cmd_data.Run()
	if err != nil {
		fmt.Println(fd + f + "," + err.Error())
	}
}

func ParseToCpplang(d string, fd string, f string) {

	r1 := strings.Replace(d, "\t", " ", -1)
	r2 := strings.Replace(r1, "\r\n", " ", -1)
	r3 := strings.Replace(r2, "\n", " ", -1)
	r4 := strings.Replace(r3, "{", " ", -1)
	r5 := strings.Replace(r4, "}", " ", -1)

	m := strings.Fields(r5)

	lens := len(m)

	data_struct := make([]string, 10)
	data_read := make([]string, 10)
	data_write := make([]string, 10)

	data_count := 0

	flag := "key"
	member := ""

	for i := 0; i < lens; i++ {
		switch m[i] {
		case "-message-":
			switch flag {
			case "key":
			default:
				if flag == "member" {
					data_count++
				} else {
					fmt.Printf("flag=%s,解析失败[%s]\n", flag, data_struct[data_count])
					data_struct[data_count] = ""
					data_read[data_count] = ""
					data_write[data_count] = ""
				}
			}
			flag = "desc"

		default:
			switch flag {
			case "desc":
				data_struct[data_count] += "// " + m[i] + "\n"
				flag = "name"

			case "name":
				data_struct[data_count] += "class " + m[i] + " {\n"
				data_read[data_count] += "void " + m[i] + "::Read(s *Stream) {\n"
				data_write[data_count] += "void " + m[i] + "::Write(s *Stream) {\n"
				flag = "member"

			case "member":
				member = m[i]
				flag = "type"

			case "type":
				if len(member) > 0 {
					type_name := TypeToClang(m[i])
					if type_name == "" {
						data_struct[data_count] += "\t" + m[i] + "\t" + member + "; // " + m[i] + "\n"
					} else {
						data_struct[data_count] += "\t" + type_name + "\t" + member + "; // " + m[i] + "\n"
					}

					func_read := GetReadFunc(m[i])
					if func_read == "" {
						data_read[data_count] += "\t" + "this." + member + ".Read(s);\n"
					} else {
						data_read[data_count] += "\t" + "s." + GetReadFunc(m[i]) + "(this." + member + ");\n"
					}

					func_write := GetWriteFunc(m[i])
					if func_write == "" {
						data_write[data_count] += "\t" + "this." + member + ".Write(s);\n"
					} else {
						data_write[data_count] += "\t" + "s." + GetWriteFunc(m[i]) + "(this." + member + ");\n"
					}
				} else {
					fmt.Printf("flag=%s,解析失败[%s]\n", flag, data_struct[data_count])
				}

				flag = "member"
			}
		}
	}

	for i := 0; i <= data_count; i++ {
		data_struct[i] += "};\n\n"
		data_read[i] += "}\n\n"
		data_write[i] += "}\n\n"
	}
	// 检查log目录
	if !help.IsExist(fd) {
		os.MkdirAll(fd, os.ModeDir)
	}

	if !help.IsExist(fd + f) {
		os.Create(fd + f)
	}
	file, err := os.OpenFile(fd+f, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i <= data_count; i++ {
		file.WriteString(data_struct[i])
		file.WriteString(data_read[i])
		file.WriteString(data_write[i])
	}

	file.Close()
}
