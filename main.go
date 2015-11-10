// main.go
package main

import (
	"flag"
	"fmt"
	"github.com/toophy/easybuff/help"
	"github.com/toophy/easybuff/proto"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	// "time"
)

// easybuff -s 描述文件目录 -d 目标文件目录 -l 语言(go,cpp)
var g_Source = flag.String("s", "", "描述文件目录")
var g_Target = flag.String("t", "", "目标文件目录")
var g_Language = flag.String("l", "go", "语言(go,cpp)")

var g_CountMutex sync.Mutex
var g_Count int

func main() {

	flag.Parse()

	fmt.Println(*g_Source, *g_Target, *g_Language)

	if !help.IsExist(*g_Source) {
		fmt.Printf("找不到描述文件目录[%s]\n", g_Source)
		return
	}

	cd, _ := os.Getwd()
	fmt.Println(cd)

	fs := help.NewFileSearch()
	fs.Dir = *g_Source
	fs.KeyWord = "/*.epd"
	fs.SubDir = false
	fl, _ := fs.RegSearchToList()

	for _, key := range fl {
		g_CountMutex.Lock()
		g_Count++
		g_CountMutex.Unlock()

		go func(k string) {
			file_path := *g_Source + k

			fi, err := os.Open(file_path)
			if err != nil {
				fmt.Println("读文件失败: %s", err.Error())
				return
			}
			fd, err := ioutil.ReadAll(fi)
			fi.Close()

			switch *g_Language {
			case "go":
				proto.ParseToNewGolang(string(fd), *g_Target, k+".go")
			case "cpp":
				proto.ParseToCpplang(string(fd), *g_Target, k+".cpp")
			}
			g_CountMutex.Lock()
			g_Count--
			g_CountMutex.Unlock()
		}(key)
	}

	for {
		g_CountMutex.Lock()
		if g_Count == 0 {
			g_CountMutex.Unlock()
			break
		}
		g_CountMutex.Unlock()
		//time.Sleep(300 * time.Millisecond)
		runtime.Gosched()
	}
}
