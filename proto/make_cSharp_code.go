package proto

import (
	"fmt"
	"github.com/toophy/easybuff/help"
	"os"
	"regexp"
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

// 线程日志 : 致命[F]级别日志
func WriteCSharpString(file *os.File, level int, f string, v ...interface{}) {
	switch level {
	case 0:
		file.WriteString(fmt.Sprintf(f, v...) + "\r\n")
	case 1:
		file.WriteString("    " + fmt.Sprintf(f, v...) + "\r\n")
	case 2:
		file.WriteString("        " + fmt.Sprintf(f, v...) + "\r\n")
	case 3:
		file.WriteString("            " + fmt.Sprintf(f, v...) + "\r\n")
	default:
		file.WriteString(fmt.Sprintf(f, v...) + "\r\n")
	}
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
	// 	file.WriteString(fmt.Sprintf("#include .\"%s\"", table.MyIncludeIdx[key].Name))
	// }
	WriteCSharpString(file, 0, "")
}

func writeMessagesCSharp(file *os.File, table *EB_FileTable, msgEnum string) {

	WriteCSharpString(file, 1, "// ------ 消息结构")

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
					WriteCSharpString(file, 1, "// %s", d.Comment[k])
				}
			}

			WriteCSharpString(file, 1, "class %s", d.Name)
			WriteCSharpString(file, 1, "{")

			// if d.MsgId > 0 {
			// 	WriteCSharpString(file, 2, "public const ushort %s_Id = %d;", d.Name, d.MsgId)
			// } else {
			// 	panic(fmt.Sprintf("文件格式错误 : message [%s]的ID(%d)非法", d.Name, d.MsgId))
			// }

			// 循环, 顺序输出
			// 按照字母顺序
			mbkeys := make([]string, len(d.Members))
			for k, _ := range d.Members {
				mbkeys[d.Members[k].Sort] = k
			}

			for _, v := range mbkeys {
				m := d.Members[v]
				if len(m.Range) <= 0 {
					WriteCSharpString(file, 2, "public %s %s; // %s", TypeToCSharp(m.Type), m.Name, m.Desc)
				} else {
					WriteCSharpString(file, 2, "public %s[] %s; // %s", TypeToCSharp(m.Type), m.Name, m.Desc)
				}
			}

			WriteCSharpString(file, 0, "")

			// 构造
			WriteCSharpString(file, 2, "public %s()", d.Name)
			WriteCSharpString(file, 2, "{")

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetReadFunc(m.Type)
				if fn == "" {
					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "%s = new %s();", m.Name, m.Type)
					} else {
						WriteCSharpString(file, 3, "%s = new %s[%s];", m.Name, m.Type, m.Range)
					}
				}
			}

			WriteCSharpString(file, 2, "}")

			// read
			WriteCSharpString(file, 2, "public void Read(ref PacketReader p)")
			WriteCSharpString(file, 2, "{")

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetReadFunc(m.Type)
				if fn == "" {
					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "%s.Read(ref p);", m.Name)
					} else {
						WriteCSharpString(file, 3, "int len_%s = int(p.%s());", m.Name, GetReadFunc("uint8"))
						WriteCSharpString(file, 3, "for (int i=0;i<len_%s;++i)", m.Name)
						WriteCSharpString(file, 3, "{")
						WriteCSharpString(file, 4, "%s[i].Read(ref p);", m.Name)
						WriteCSharpString(file, 3, "}")
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "p.%s(ref %s);", fn, m.Name)
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							WriteCSharpString(file, 3, "p.%s(uint8(%s.Lenght));", GetReadFunc("uint8"), m.Name)
							WriteCSharpString(file, 3, "for (int i=;i<%s.Lenght;++i)", m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "p.%s(ref %s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")
						} else {
							// 枚举
							WriteCSharpString(file, 3, "int len_%s = int(p.%s());", m.Name, GetReadFunc("uint8"))
							WriteCSharpString(file, 3, "for (int i=0; i<%s && i<len_%s; ++i)", m.Range, m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "%s(ref %s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")
						}
					}
				} else {
					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "%s = p.%s();", m.Name, fn)
					} else {
						WriteCSharpString(file, 3, "int len_%s = int(p.%s());", m.Name, GetReadFunc("uint8"))
						WriteCSharpString(file, 3, "for (int i=0;i<len_%s;++i)", m.Name)
						WriteCSharpString(file, 3, "{")
						WriteCSharpString(file, 3, "%s = p.%s();", m.Name, fn)
						WriteCSharpString(file, 3, "}")
					}
				}
			}

			WriteCSharpString(file, 2, "}")

			// write
			if strings.HasPrefix(d.Name, "S2G_") || strings.HasPrefix(d.Name, "G2S_") || strings.HasPrefix(d.Name, "S2C_") {
				WriteCSharpString(file, 2, "public void Write(ref PacketWriter p, long tgid)")
				WriteCSharpString(file, 2, "{")
				WriteCSharpString(file, 3, "p.WriteMsgId((ushort)%s.%s_Id);", msgEnum, d.Name)
			} else {
				WriteCSharpString(file, 2, "public void Write(ref PacketWriter p)")
				WriteCSharpString(file, 2, "{")
				WriteCSharpString(file, 3, "p.WriteMsgId((ushort)%s.%s_Id);", msgEnum, d.Name)
			}

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetWriteFunc(m.Type)
				if fn == "" {

					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "%s.Write(ref p);", m.Name)
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							WriteCSharpString(file, 3, "p.%s(uint8(%s.Lenght));", GetWriteFunc("uint8"), m.Name)
							WriteCSharpString(file, 3, "for (int i=;i<%s.Lenght;++i)", m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "%s[k].Write(ref p);", m.Name)
							WriteCSharpString(file, 3, "}")
						} else {
							// 枚举
							WriteCSharpString(file, 3, "p.%s(uint8(%s));", GetWriteFunc("uint8"), m.Range)
							WriteCSharpString(file, 3, "int len_%s = %s.Lenght;", m.Name, m.Name)
							WriteCSharpString(file, 3, "for (int i=0; i<%s && i<len_%s; ++i)", m.Range, m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "%s[i].Write(ref p);", m.Name)
							WriteCSharpString(file, 3, "}")
						}
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "p.%s(ref %s);", fn, m.Name)
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							WriteCSharpString(file, 3, "p.%s(uint8(%s.Lenght));", GetWriteFunc("uint8"), m.Name)
							WriteCSharpString(file, 3, "for (int i=;i<%s.Lenght;++i) {", m.Name)
							WriteCSharpString(file, 4, "p.%s(ref %s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")
						} else {
							// 枚举
							WriteCSharpString(file, 3, "p.%s(uint8(%s));", GetWriteFunc("uint8"), m.Range)
							WriteCSharpString(file, 3, "int len_%s = %s.Lenght;", m.Name, m.Name)
							WriteCSharpString(file, 3, "for (int i=0; i<%s && i<len_%s; ++i) {", m.Range, m.Name)
							WriteCSharpString(file, 4, "%s(ref %s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")
						}
					}
				} else {
					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "p.%s(%s);", fn, m.Name)
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							WriteCSharpString(file, 3, "p.%s(uint8(%s.Lenght));", GetWriteFunc("uint8"), m.Name)
							WriteCSharpString(file, 3, "for (int i=0;i<%s.Lenght;++i) {", m.Name)
							WriteCSharpString(file, 4, "p.%s(%s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")

						} else {
							// 枚举
							WriteCSharpString(file, 3, "p.%s(uint8(%s));", GetWriteFunc("uint8"), m.Range)
							WriteCSharpString(file, 3, "int len_%s = %s.Lenght;", m.Name, m.Name)
							WriteCSharpString(file, 3, "for (int i=0; i<%s && i<len_%s; ++i) {", m.Range, m.Name)
							WriteCSharpString(file, 4, "p.%s(%s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")
						}
					}
				}
			}

			// write over
			WriteCSharpString(file, 3, "p.WriteMsgOver();")
			WriteCSharpString(file, 2, "}")
			WriteCSharpString(file, 1, "}")
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

	// 文件名标志
	reg := regexp.MustCompile(`[^/.]+`)
	file_flag := reg.FindString(table.FileName)

	// 文件头
	WriteCSharpString(file, 0, "// easybuff")
	WriteCSharpString(file, 0, "// 不要修改本文件, 每次消息有变动, 请手动生成本文件")
	WriteCSharpString(file, 0, "// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp,c#)")
	WriteCSharpString(file, 0, "")
	WriteCSharpString(file, 0, "using NetMsg;")
	WriteCSharpString(file, 0, "using System;")

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
	WriteCSharpString(file, 0, "namespace NetMsg")
	WriteCSharpString(file, 0, "{")

	// 枚举
	WriteCSharpString(file, 1, "// ------ 枚举")
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
	// 			file.WriteString(fmt.Sprintf("static int %s = %s; // %s", d.Name, d.Value, comment_string))
	// 		} else {
	// 			file.WriteString(fmt.Sprintf("static int %s = %s;", d.Name, d.Value))
	// 		}
	// 	}
	// }

	// 消息ID
	WriteCSharpString(file, 1, "// ------ 消息ID")
	WriteCSharpString(file, 1, fmt.Sprintf("enum MsgId_%s", file_flag))
	WriteCSharpString(file, 1, fmt.Sprintf("{"))

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
					WriteCSharpString(file, 2, fmt.Sprintf("%s_Id = %d, // %s", d.Name, d.MsgId, comment_string))
				} else {
					WriteCSharpString(file, 2, fmt.Sprintf("%s_Id = %d,", d.Name, d.MsgId))
				}

			} else {
				panic(fmt.Sprintf("文件格式错误 : message [%s]的ID(%d)非法", d.Name, d.MsgId))
			}
		}
	}

	WriteCSharpString(file, 1, fmt.Sprintf("}"))

	// 普通结构
	writeStructsCSharp(file, table)

	// 消息结构
	writeMessagesCSharp(file, table, fmt.Sprintf("MsgId_%s", file_flag))

	WriteCSharpString(file, 0, "}")

	file.Close()

	// cmd_data := exec.Command("gofmt", "-w", target_file)
	// err = cmd_data.Run()
	// if err != nil {
	// 	fmt.Println(target_file + "," + err.Error())
	// }
}

func writeStructsCSharp(file *os.File, table *EB_FileTable) {
	WriteCSharpString(file, 1, "// ------ 普通结构")

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
					WriteCSharpString(file, 1, "// %s", d.Comment[k])
				}
			}

			WriteCSharpString(file, 1, "class %s", d.Name)
			WriteCSharpString(file, 1, "{")

			// 循环, 顺序输出
			// 按照字母顺序
			mbkeys := make([]string, len(d.Members))
			for k, _ := range d.Members {
				mbkeys[d.Members[k].Sort] = k
			}

			for _, v := range mbkeys {
				m := d.Members[v]
				if len(m.Range) <= 0 {
					WriteCSharpString(file, 2, "public %s %s; // %s", TypeToCSharp(m.Type), m.Name, m.Desc)
				} else {
					WriteCSharpString(file, 2, "public %s[] %s; // %s", TypeToCSharp(m.Type), m.Name, m.Desc)
				}
			}

			// 构造
			WriteCSharpString(file, 2, "public %s()", d.Name)
			WriteCSharpString(file, 2, "{")

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetReadFunc(m.Type)
				if fn == "" {
					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "%s = new %s();", m.Name, m.Type)
					} else {
						WriteCSharpString(file, 3, "%s = new %s[%s];", m.Name, m.Type, m.Range)
					}
				}
			}

			WriteCSharpString(file, 2, "}")

			// read
			WriteCSharpString(file, 2, "public void Read(ref PacketReader p)")
			WriteCSharpString(file, 2, "{")

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetReadFunc(m.Type)
				if fn == "" {
					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "%s.Read(ref p);", m.Name)
					} else {
						WriteCSharpString(file, 3, "int len_%s = int(p.%s());", m.Name, GetReadFunc("uint8"))
						WriteCSharpString(file, 3, "for (int i=0;i<len_%s;++i)", m.Name)
						WriteCSharpString(file, 3, "{")
						WriteCSharpString(file, 4, "%s[i].Read(ref p);", m.Name)
						WriteCSharpString(file, 3, "}")
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "p.%s(ref %s);", fn, m.Name)
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							WriteCSharpString(file, 3, "p.%s(uint8(%s.Lenght));", GetReadFunc("uint8"), m.Name)
							WriteCSharpString(file, 3, "for (int i=;i<%s.Lenght;++i)", m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "p.%s(ref %s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")
						} else {
							// 枚举
							WriteCSharpString(file, 3, "int len_%s = int(p.%s());", m.Name, GetReadFunc("uint8"))
							WriteCSharpString(file, 3, "for (int i=0; i<%s && i<len_%s; ++i)", m.Range, m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "%s(ref %s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")
						}
					}
				} else {
					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "%s = p.%s();", m.Name, fn)
					} else {
						WriteCSharpString(file, 3, "int len_%s = int(p.%s());", m.Name, GetReadFunc("uint8"))
						WriteCSharpString(file, 3, "for (int i=0;i<len_%s;++i)", m.Name)
						WriteCSharpString(file, 3, "{")
						WriteCSharpString(file, 4, "%s = p.%s();", m.Name, fn)
						WriteCSharpString(file, 3, "}")
					}
				}
			}

			WriteCSharpString(file, 2, "}")

			// write
			WriteCSharpString(file, 2, "public void Write(ref PacketWriter p)")
			WriteCSharpString(file, 2, "{")

			for _, v := range mbkeys {
				m := d.Members[v]
				fn := GetWriteFunc(m.Type)
				if fn == "" {

					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "%s.Write(ref p);", m.Name)
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							WriteCSharpString(file, 3, "p.%s(byte(%s.Lenght));", GetWriteFunc("uint8"), m.Name)
							WriteCSharpString(file, 3, "for (int i=0;i<%s.Lenght;++i)", m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "%s[k].Write(ref p);", m.Name)
							WriteCSharpString(file, 3, "}")

						} else {
							// 枚举
							WriteCSharpString(file, 3, "p.%s(byte(%s));", GetWriteFunc("uint8"), m.Range)
							WriteCSharpString(file, 3, "int len_%s = %s.Lenght;", m.Name, m.Name)
							WriteCSharpString(file, 3, "for (int i=0; i<%s && i<len_%s; i++)", m.Range, m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "%s[i].Write(ref p);", m.Name)
							WriteCSharpString(file, 3, "}")
						}
					}
				} else if m.Type == "string" {

					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "p.%s(ref %s);", fn, m.Name)
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							WriteCSharpString(file, 3, "p.%s(uint8(%s.Lenght));", GetWriteFunc("uint8"), m.Name)
							WriteCSharpString(file, 3, "for (int i=0;i<%s.Lenght;++i)", m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "p.%s(ref %s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")

						} else {
							// 枚举
							WriteCSharpString(file, 3, "p.%s(uint8(%s));", GetWriteFunc("uint8"), m.Range)
							WriteCSharpString(file, 3, "int len_%s = %s.Lenght;", m.Name, m.Name)
							WriteCSharpString(file, 3, "for (int i=0; i<%s && i<len_%s; i++)", m.Range, m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "p.%s(ref %s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")
						}
					}
				} else {
					if len(m.Range) <= 0 {
						WriteCSharpString(file, 3, "p.%s(%s);", fn, m.Name)
					} else {
						if m.Range == "--ArrayLen" {
							// 自动范围
							WriteCSharpString(file, 3, "p.%s(uint8(%s.Lenght));", GetWriteFunc("uint8"), m.Name)
							WriteCSharpString(file, 3, "for(int i=0;i<%s.Lenght;++i)", m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 4, "p.%s(%s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")

						} else {
							// 枚举
							WriteCSharpString(file, 3, "p.%s(uint8(%s));", GetWriteFunc("uint8"), m.Range)
							WriteCSharpString(file, 3, "int len_%s = %s.Lenght;", m.Name, m.Name)
							WriteCSharpString(file, 3, "for (int i=0; i<%s && i<len_%s; i++)", m.Range, m.Name)
							WriteCSharpString(file, 3, "{")
							WriteCSharpString(file, 3, "p.%s(%s[i]);", fn, m.Name)
							WriteCSharpString(file, 3, "}")
						}
					}
				}
			}

			// write over
			WriteCSharpString(file, 2, "}")

			// class over
			WriteCSharpString(file, 1, "}")
		}
	}

	WriteCSharpString(file, 0, "")
}
