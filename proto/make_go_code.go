package proto

import (
	"fmt"
	"github.com/toophy/easybuff/help"
	"os"
	"os/exec"
	"sort"
	"strings"
)

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
	case "float32":
		type_name = "float32"
	case "float64":
		type_name = "float64"
	case "string":
		type_name = "string"
	}
	return type_name
}

func WriteAllEpdGoCode() {
	for _, v := range g_ParseTable.Files {
		writeGoCode(v)
	}
}

func writeInclude(file *os.File, table *EB_FileTable) {

	// 按照字母顺序
	keys := make([]int, len(table.MyIncludeIdx))
	i := 0
	for key, _ := range table.MyIncludeIdx {
		keys[i] = key
		i++
	}

	sort.Sort(sort.IntSlice(keys))

	// 引用文件
	file.WriteString("import (\n")

	for _, key := range keys {
		file.WriteString(fmt.Sprintf(".\"%s\"\n", table.MyIncludeIdx[key].Name))
	}
	file.WriteString(")\n\n")
}

func writeStructs(file *os.File, table *EB_FileTable) {
	file.WriteString("// ------ 普通结构\n")

	keys := make([]int, len(table.MyStructIds))
	i := 0
	for key, _ := range table.MyStructIds {
		keys[i] = key
		i++
	}

	sort.Sort(sort.IntSlice(keys))

	// 消息结构
	for _, key := range keys {
		struct_name := table.MyStructIds[key]

		switch g_ParseTable.Cells[struct_name].(type) {
		case *EB_Message:
			d := g_ParseTable.Cells[struct_name].(*EB_Message)

			for k, _ := range d.Comment {
				if len(d.Comment[k]) > 0 {
					file.WriteString(fmt.Sprintf("// %s\n", d.Comment[k]))
				}
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
			file.WriteString(fmt.Sprintf(`func (t *%s) Read(p *PacketReader) bool {
				defer RecoverRead("%s")

				`, d.Name, d.Name))
			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetReadFunc(m.Type)
				if fn == "" {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" t.%s.Read(p)\n", m.Name))
					} else {
						file.WriteString(fmt.Sprintf(`
													  len_%s := int(p.%s())
													  t.%s = make([]%s,len_%s)
													  for i:=0;i<len_%s;i++ {
													  	t.%s[i].Read(p)
													  }
													  
													  `, m.Name, GetReadFunc("uint8"), m.Name, TypeToGolang(m.Type), m.Name, m.Name, m.Name))
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" t.%s = p.%s()\n", m.Name, fn))
					} else {
						file.WriteString(fmt.Sprintf(`
													  len_%s := int(p.%s())
													  t.%s = make([]%s,len_%s)
													  for i:=0;i<len_%s;i++ {
													  	 t.%s = p.%s()
													  }
													  
													  `, m.Name, GetReadFunc("uint8"), m.Name, TypeToGolang(m.Type), m.Name, m.Name, m.Name, fn))
					}
				}
			}

			file.WriteString("\nreturn true\n}\n\n")

			// write
			file.WriteString(fmt.Sprintf(`func (t *%s) Write(p *PacketWriter) bool {
				defer RecoverWrite("%s")
				`, d.Name, d.Name))

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetWriteFunc(m.Type)
				if fn == "" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" t.%s.Write(p)\n", m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	t.%s[k].Write(p)
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	t.%s[i].Write(p)
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, m.Name))
						}
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" p.%s(&t.%s)\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	p.%s(&t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	p.%s(&t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" p.%s(t.%s)\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	p.%s(t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	p.%s(t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				}
			}

			// write over
			file.WriteString("\n\nreturn true\n}\n\n\n")
		}
	}

	file.WriteString("\n\n")
}

func writeMessages(file *os.File, table *EB_FileTable) {

	file.WriteString("// ------ 消息结构\n")

	// 按照字母顺序
	keys := make([]int, len(table.MyMsgIds))
	i := 0
	for key, _ := range table.MyMsgIds {
		keys[i] = key
		i++
	}

	sort.Sort(sort.IntSlice(keys))

	for _, key := range keys {
		msg_name := table.MyMsgIds[key]

		switch g_ParseTable.Cells[msg_name].(type) {
		case *EB_Message:
			d := g_ParseTable.Cells[msg_name].(*EB_Message)

			for k, _ := range d.Comment {
				if len(d.Comment[k]) > 0 {
					file.WriteString(fmt.Sprintf("// %s\n", d.Comment[k]))
				}
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
			file.WriteString(fmt.Sprintf(`func (t *%s) Read(p *PacketReader) bool {
				defer RecoverRead("%s")

				`, d.Name, d.Name))
			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetReadFunc(m.Type)
				if fn == "" {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" t.%s.Read(p)\n", m.Name))
					} else {
						file.WriteString(fmt.Sprintf(`
													  len_%s := int(p.%s())
													  t.%s = make([]%s,len_%s)
													  for i:=0;i<len_%s;i++ {
													  	t.%s[i].Read(p)
													  }
													  
													  `, m.Name, GetReadFunc("uint8"), m.Name, TypeToGolang(m.Type), m.Name, m.Name, m.Name))
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" t.%s = p.%s()\n", m.Name, fn))
					} else {
						file.WriteString(fmt.Sprintf(`
													  len_%s := int(p.%s())
													  t.%s = make([]%s,len_%s)
													  for i:=0;i<len_%s;i++ {
													  	 t.%s = p.%s()
													  }
													  
													  `, m.Name, GetReadFunc("uint8"), m.Name, TypeToGolang(m.Type), m.Name, m.Name, m.Name, fn))
					}
				}
			}

			file.WriteString("\nreturn true\n}\n\n")

			// write
			if strings.HasPrefix(d.Name, "S2G_") || strings.HasPrefix(d.Name, "G2S_") || strings.HasPrefix(d.Name, "S2C_") {
				file.WriteString(fmt.Sprintf(`func (t *%s) Write(p *PacketWriter, tgid uint64) bool {
				defer RecoverWrite("%s")
				p.SetsubTgid(tgid)
				p.WriteMsgId(%s_Id)
				`, d.Name, d.Name, d.Name))
			} else {
				file.WriteString(fmt.Sprintf(`func (t *%s) Write(p *PacketWriter) bool {
				defer RecoverWrite("%s")
				p.WriteMsgId(%s_Id)
				`, d.Name, d.Name, d.Name))
			}
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
						file.WriteString(fmt.Sprintf(" t.%s.Write(p)\n", m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	t.%s[k].Write(p)
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	t.%s[i].Write(p)
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, m.Name))
						}
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" p.%s(&t.%s)\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	p.%s(&t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	p.%s(&t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" p.%s(t.%s)\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(len(t.%s)))
														  for k,_ := range t.%s {
														  	p.%s(t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s))
														  len_%s := len(t.%s)
														  for i:=0; i<%s && i<len_%s; i++ {
														  	p.%s(t.%s[i])
														  }
														  
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				}
			}

			// write over
			file.WriteString("p.WriteMsgOver()\n\nreturn true\n}\n\n\n")
		}
	}
}

func writeGoCode(table *EB_FileTable) {
	// 检查log目录
	if !help.IsExist(table.FileDir) {
		os.MkdirAll(table.FileDir, os.ModeDir)
	}

	target_file := table.FileDir + "/" + table.FileName

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
	keys := make([]int, len(table.MyMsgIds))
	i := 0
	for key, _ := range table.MyMsgIds {
		keys[i] = key
		i++
	}

	sort.Sort(sort.IntSlice(keys))

	// 引用文件
	writeInclude(file, table)

	// 枚举
	file.WriteString("// ------ 枚举\n")
	file.WriteString("const(\n")
	for _, key := range keys {
		msg_name := table.MyMsgIds[key]

		switch g_ParseTable.Cells[msg_name].(type) {
		case *EB_Enum:
			d := g_ParseTable.Cells[msg_name].(*EB_Enum)
			comment_string := ""
			for k, _ := range d.Comment {
				if len(d.Comment[k]) > 0 {
					if len(comment_string) > 0 {
						comment_string += ","
					}
					comment_string += d.Comment[k]
				}
			}
			if len(comment_string) > 0 {
				file.WriteString(fmt.Sprintf("%s = %s // %s\n", d.Name, d.Value, comment_string))
			} else {
				file.WriteString(fmt.Sprintf("%s = %s\n", d.Name, d.Value))
			}
		}
	}
	file.WriteString(")\n\n")

	// 消息ID
	file.WriteString("// ------ 消息ID\n")
	file.WriteString("const(\n")

	for _, key := range keys {
		msg_name := table.MyMsgIds[key]

		switch g_ParseTable.Cells[msg_name].(type) {
		case *EB_Message:
			d := g_ParseTable.Cells[msg_name].(*EB_Message)

			comment_string := ""
			for k, _ := range d.Comment {
				if len(d.Comment[k]) > 0 {
					if len(comment_string) > 0 {
						comment_string += ", "
					}
					comment_string += d.Comment[k]
				}
			}

			if d.MsgId > 0 {
				if len(comment_string) > 0 {
					file.WriteString(fmt.Sprintf("%s_Id = %d// %s\n", d.Name, d.MsgId, comment_string))
				} else {
					file.WriteString(fmt.Sprintf("%s_Id = %d\n", d.Name, d.MsgId))
				}

			} else {
				panic(fmt.Sprintf("文件格式错误 : message [%s]的ID(%d)非法", d.Name, d.MsgId))
			}
		}
	}

	file.WriteString(")\n\n")

	// 普通结构
	writeStructs(file, table)

	// 消息结构
	writeMessages(file, table)

	file.Close()

	cmd_data := exec.Command("gofmt", "-w", target_file)
	err = cmd_data.Run()
	if err != nil {
		fmt.Println(target_file + "," + err.Error())
	}
}
