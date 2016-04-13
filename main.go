// main.go
package main

import (
	"flag"
	"fmt"
	"github.com/toophy/easybuff/help"
	"github.com/toophy/easybuff/proto"
	"io/ioutil"
	"os"
	"strings"
)

// easybuff -s 描述文件目录 -d 目标文件目录 -l 语言(go,cpp)
var g_Source = flag.String("s", "proto_desc", "描述文件目录")
var g_Target = flag.String("t", "proto", "目标文件目录")
var g_Language = flag.String("l", "go", "语言(go,cpp)")

// var g_CountMutex sync.Mutex
// var g_Count int

func main() {

	flag.Parse()

	source_dir := *g_Source
	target_dir := *g_Target

	source_dir = strings.TrimSuffix(source_dir, "/")
	target_dir = strings.TrimSuffix(target_dir, "/")

	fmt.Println(source_dir, target_dir, *g_Language)

	if !help.IsExist(source_dir) {
		fmt.Printf("找不到描述文件目录[%s]\n", source_dir)
		return
	}

	cd, _ := os.Getwd()
	fmt.Println(cd)

	fs := help.NewFileSearch()
	fs.Dir = source_dir
	fs.KeyWord = "/*.epd"
	fs.SubDir = false
	fl, _ := fs.RegSearchToList()

	for _, key := range fl {

		file_path := source_dir + "/" + key

		fi, err := os.Open(file_path)
		if err != nil {
			fmt.Println("读文件失败: %s", err.Error())
			return
		}
		fd, err := ioutil.ReadAll(fi)
		fi.Close()

		switch *g_Language {
		case "go":
			proto.ParseToNewGolang(string(fd), target_dir, key+".go")
			proto.WriteAllEpdGoCode()
		case "c#":
			proto.ParseToNewGolang(string(fd), target_dir, key+".cs")
			proto.WriteAllEpdCSharpCode()
			// case "cpp":
			// 	proto.ParseToCpplang(string(fd), target_dir, key+".cpp")
		}
	}

}
