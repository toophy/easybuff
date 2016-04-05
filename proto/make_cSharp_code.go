package proto

import (
	"fmt"
	"github.com/toophy/easybuff/help"
	"os"
	"sort"
	"strings"
)

func TypeToCSharp(s string) string {
	type_name := s
	switch s {
	case "int8":
		type_name = "sbyte"
	case "int16":
		type_name = "short"
	case "int24":
		type_name = "int"
	case "int32":
		type_name = "int"
	case "int40":
		type_name = "long"
	case "int48":
		type_name = "long"
	case "int56":
		type_name = "long"
	case "int64":
		type_name = "long"
	case "uint8":
		type_name = "byte"
	case "uint16":
		type_name = "ushort"
	case "uint24":
		type_name = "uint"
	case "uint32":
		type_name = "uint"
	case "uint40":
		type_name = "ulong"
	case "uint48":
		type_name = "ulong"
	case "uint56":
		type_name = "ulong"
	case "uint64":
		type_name = "ulong"
	case "float32":
		type_name = "float"
	case "float64":
		type_name = "double"
	case "string":
		type_name = "string"
	}
	return type_name
}

func WriteAllEpdCSharpCode() {
	for _, v := range g_ParseTable.Files {
		writeCSharpCode(v)
	}
}

func writeIncludeCSharp(file *os.File, table *EB_FileTable) {

	// 按照字母顺序
	keys := make([]int, len(table.MyIncludeIdx))
	i := 0
	for key, _ := range table.MyIncludeIdx {
		keys[i] = key
		i++
	}

	sort.Sort(sort.IntSlice(keys))

	// 引用文件
	// for _, key := range keys {
	// 	file.WriteString(fmt.Sprintf("#include .\"%s\"\n", table.MyIncludeIdx[key].Name))
	// }
	file.WriteString("\n\n")
}

func writeStructsCSharp(file *os.File, table *EB_FileTable) {
	file.WriteString("\n\n// ------ 普通结构\n")

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
					file.WriteString(fmt.Sprintf("	// %s\n", d.Comment[k]))
				}
			}

			file.WriteString(fmt.Sprintf("	class %s\n	{\n", d.Name))
			// 循环, 顺序输出
			// 按照字母顺序
			mbkeys := make([]string, len(d.Members))
			for k, _ := range d.Members {
				mbkeys[d.Members[k].Sort] = k
			}

			for _, v := range mbkeys {
				m := d.Members[v]
				if len(m.Range) <= 0 {
					file.WriteString(fmt.Sprintf("		public %s %s; // %s\n", TypeToCSharp(m.Type), m.Name, m.Desc))
				} else {
					file.WriteString(fmt.Sprintf("		public %s[] %s; // %s\n", TypeToCSharp(m.Type), m.Name, m.Desc))
				}
			}

			file.WriteString("\n")

			// read
			file.WriteString("		public void Read(ref PacketReader p)\n		{\n")

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetReadFunc(m.Type)
				if fn == "" {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf("			%s.Read(ref p);\n", m.Name))
					} else {
						file.WriteString(fmt.Sprintf("			int len_%s = int(p.%s());\n"+
							"			for (int i=0;i<len_%s;++i)\n"+
							"			{\n"+
							"				%s[i].Read(ref p);\n"+
							"			}\n", m.Name, GetReadFunc("uint8"), m.Name, m.Name))
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf("			p.%s(ref %s);\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf("			p.%s(uint8(%s.Lenght));\n"+
								"			for (int i=;i<%s.Lenght;++i)\n"+
								"			{\n"+
								"				p.%s(ref %s[i]);\n"+
								"			}\n", GetReadFunc("uint8"), m.Name, m.Name, fn, m.Name))
						} else {
							// 枚举
							file.WriteString(fmt.Sprintf("			int len_%s = int(p.%s());\n"+
								"			for (int i=0; i<%s && i<len_%s; ++i)\n"+
								"			{\n"+
								"				%s(ref %s[i]);\n"+
								"			}\n", m.Name, GetReadFunc("uint8"), m.Range, m.Name, fn, m.Name))
						}
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf("			%s = p.%s();\n", m.Name, fn))
					} else {
						file.WriteString(fmt.Sprintf("			int len_%s = int(p.%s());\n"+
							"			for (int i=0;i<len_%s;++i)\n"+
							"			{\n"+
							"				%s = p.%s();\n"+
							"			}\n", m.Name, GetReadFunc("uint8"), m.Name, m.Name, fn))
					}
				}
			}

			file.WriteString("		}\n\n")

			// write
			file.WriteString("		public void Write(ref PacketWriter p)\n		{\n")

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetWriteFunc(m.Type)
				if fn == "" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf("			%s.Write(ref p);\n", m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf("			p.%s(byte(%s.Lenght));\n"+
								"			for (int i=0;i<%s.Lenght;++i)\n"+
								"			{\n"+
								"				%s[k].Write(ref p);\n"+
								"			}\n", GetWriteFunc("uint8"), m.Name, m.Name, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf("			p.%s(byte(%s));\n"+
								"			int len_%s = %s.Lenght;\n"+
								"			for (int i=0; i<%s && i<len_%s; i++)\n"+
								"			{\n"+
								"				%s[i].Write(ref p);\n"+
								"			}\n", GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, m.Name))
						}
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf("			p.%s(ref %s);\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf("			p.%s(uint8(%s.Lenght));\n"+
								"			for (int i=0;i<%s.Lenght;++i)\n"+
								"			{\n"+
								"				p.%s(ref %s[i]);\n"+
								"			}\n", GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf("			p.%s(uint8(%s));\n"+
								"			int len_%s = %s.Lenght;\n"+
								"			for (int i=0; i<%s && i<len_%s; i++)\n"+
								"			{\n"+
								"				p.%s(ref %s[i]);\n"+
								"			}\n", GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf("			p.%s(%s);\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf("			p.%s(uint8(%s.Lenght));\n"+
								"			for(int i=0;i<%s.Lenght;++i)\n"+
								"			{\n"+
								"				p.%s(%s[i]);\n"+
								"			}\n", GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf("			p.%s(uint8(%s));\n"+
								"			int len_%s = %s.Lenght;\n"+
								"			for (int i=0; i<%s && i<len_%s; i++)\n"+
								"			{\n"+
								"				p.%s(%s[i]);\n"+
								"			}\n", GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				}
			}

			// write over
			file.WriteString("		}\n")

			// class over
			file.WriteString("	}\n")
		}
	}

	file.WriteString("\n\n")
}

func writeMessagesCSharp(file *os.File, table *EB_FileTable) {

	file.WriteString("\n\n// ------ 消息结构\n")

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

			file.WriteString(fmt.Sprintf("class %s \n{\n", d.Name))

			if d.MsgId > 0 {
				file.WriteString(fmt.Sprintf("public const ushort %s_Id = %d;\n\n", d.Name, d.MsgId))
			} else {
				panic(fmt.Sprintf("文件格式错误 : message [%s]的ID(%d)非法", d.Name, d.MsgId))
			}

			// 循环, 顺序输出
			// 按照字母顺序
			mbkeys := make([]string, len(d.Members))
			for k, _ := range d.Members {
				mbkeys[d.Members[k].Sort] = k
			}

			for _, v := range mbkeys {
				m := d.Members[v]
				if len(m.Range) <= 0 {
					file.WriteString(fmt.Sprintf("public %s %s; // %s\n", TypeToCSharp(m.Type), m.Name, m.Desc))
				} else {
					file.WriteString(fmt.Sprintf("public %s[] %s; // %s\n", TypeToCSharp(m.Type), m.Name, m.Desc))
				}
			}

			file.WriteString("\n")
			// read
			file.WriteString("public void Read(ref PacketReader p)\n{\n")

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetReadFunc(m.Type)
				if fn == "" {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" %s.Read(ref p);\n", m.Name))
					} else {
						file.WriteString(fmt.Sprintf(`
													  int len_%s = int(p.%s());
													  for (int i=0;i<len_%s;++i) {
													  	%s[i].Read(ref p);
													  }
													  `, m.Name, GetReadFunc("uint8"), m.Name, m.Name))
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" p.%s(ref %s);\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s.Lenght));
														  for (int i=;i<%s.Lenght;++i) {
														  	p.%s(ref %s[i]);
														  }
														  `, GetReadFunc("uint8"), m.Name, m.Name, fn, m.Name))
						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  int len_%s = int(p.%s());
														  for (int i=0; i<%s && i<len_%s; ++i) {
														  	%s(ref %s[i]);
														  }
														  `, m.Name, GetReadFunc("uint8"), m.Range, m.Name, fn, m.Name))
						}
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" %s = p.%s();\n", m.Name, fn))
					} else {
						file.WriteString(fmt.Sprintf(`
													  int len_%s = int(p.%s());
													  for (int i=0;i<len_%s;++i) {
													  	 %s = p.%s();
													  }
													  `, m.Name, GetReadFunc("uint8"), m.Name, m.Name, fn))
					}
				}
			}

			file.WriteString("}\n\n")

			// write
			if strings.HasPrefix(d.Name, "S2G_") || strings.HasPrefix(d.Name, "G2S_") || strings.HasPrefix(d.Name, "S2C_") {
				file.WriteString(fmt.Sprintf(`public void Write(ref PacketWriter p, long tgid)
				{
					p.SetsubTgid(tgid);
					p.WriteMsgId(%s_Id);
				`, d.Name))
			} else {
				file.WriteString(fmt.Sprintf(`public void Write(ref PacketWriter p)
				{
					p.WriteMsgId(%s_Id);
				`, d.Name))
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
						file.WriteString(fmt.Sprintf(" %s.Write(ref p);\n", m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s.Lenght));
														  for (int i=;i<%s.Lenght;++i) {
														  	%s[k].Write(ref p);
														  }
														  `, GetWriteFunc("uint8"), m.Name, m.Name, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s));
														  int len_%s = %s.Lenght;
														  for (int i=0; i<%s && i<len_%s; ++i) {
														  	%s[i].Write(ref p);
														  }
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, m.Name))
						}
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" p.%s(ref %s);\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s.Lenght));
														  for (int i=;i<%s.Lenght;++i) {
														  	p.%s(ref %s[i]);
														  }
														  `, GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s));
														  int len_%s = %s.Lenght;
														  for (int i=0; i<%s && i<len_%s; ++i) {
														  	%s(ref %s[i]);
														  }
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				} else {
					if len(m.Range) <= 0 {
						file.WriteString(fmt.Sprintf(" p.%s(%s);\n", fn, m.Name))
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s.Lenght));
														  for (int i=0;i<%s.Lenght;++i) {
														  	p.%s(%s[i]);
														  }
														  `, GetWriteFunc("uint8"), m.Name, m.Name, fn, m.Name))

						} else {
							// 枚举
							file.WriteString(fmt.Sprintf(`
														  p.%s(uint8(%s));
														  int len_%s = %s.Lenght;
														  for (int i=0; i<%s && i<len_%s; ++i) {
														  	p.%s(%s[i]);
														  }
														  `, GetWriteFunc("uint8"), m.Range, m.Name, m.Name, m.Range, m.Name, fn, m.Name))
						}
					}
				}
			}

			// write over
			file.WriteString("p.WriteMsgOver();\n}\n")

			file.WriteString("}\n")
		}
	}
}

func writeCSharpCode(table *EB_FileTable) {
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
	file.WriteString("// easybuff\n" +
		"// 不要修改本文件, 每次消息有变动, 请手动生成本文件\n" +
		"// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp,c#)\n" +
		"using NetMsg;\n" +
		"using System;\n")

	// 按照字母顺序
	keys := make([]int, len(table.MyMsgIds))
	i := 0
	for key, _ := range table.MyMsgIds {
		keys[i] = key
		i++
	}

	sort.Sort(sort.IntSlice(keys))

	// 引用文件
	writeIncludeCSharp(file, table)

	// 文件头
	file.WriteString("namespace NetMsg\n{\n")

	// 枚举
	file.WriteString("// ------ 枚举\n")
	// for _, key := range keys {
	// 	msg_name := table.MyMsgIds[key]

	// 	switch g_ParseTable.Cells[msg_name].(type) {
	// 	case *EB_Enum:
	// 		d := g_ParseTable.Cells[msg_name].(*EB_Enum)
	// 		comment_string := ""
	// 		for k, _ := range d.Comment {
	// 			if len(d.Comment[k]) > 0 {
	// 				if len(comment_string) > 0 {
	// 					comment_string += ","
	// 				}
	// 				comment_string += d.Comment[k]
	// 			}
	// 		}
	// 		if len(comment_string) > 0 {
	// 			file.WriteString(fmt.Sprintf("static int %s = %s; // %s\n", d.Name, d.Value, comment_string))
	// 		} else {
	// 			file.WriteString(fmt.Sprintf("static int %s = %s;\n", d.Name, d.Value))
	// 		}
	// 	}
	// }

	// 普通结构
	writeStructsCSharp(file, table)

	// 消息结构
	writeMessagesCSharp(file, table)

	file.WriteString("}\n\n")

	file.Close()

	// cmd_data := exec.Command("gofmt", "-w", target_file)
	// err = cmd_data.Run()
	// if err != nil {
	// 	fmt.Println(target_file + "," + err.Error())
	// }
}
