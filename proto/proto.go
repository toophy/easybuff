package proto

import (
	"fmt"
	"github.com/toophy/easybuff/help"
	"os"
	"os/exec"
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
	type_name := ""
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
	type_name := ""
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
