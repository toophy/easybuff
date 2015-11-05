package help

import (
	"errors"
	"fmt"
	//"glob"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// 存储单个文件信息
type FileInfo struct {
	Path        string // 文件所在的路径
	os.FileInfo        // 文件的 os.FileInfo 信息
}

// 用于处理搜索结果的函数
type OperateFunc func(fi FileInfo) error

// 定义要处理的对象，用于 FileSearch.Include
const (
	IncludeFile = 1 << (iota)              // 需要处理文件
	IncludeDir                             // 需要处理目录
	IncludeAll  = IncludeFile + IncludeDir // 两者都需要处理
)

// 用来执行搜索操作的主结构体
type FileSearch struct {
	Dir              string      // 要搜索的目录
	KeyWord          string      // 要搜索的关键字
	CaseMind         bool        // 是否区分大小写（默认忽略大小写，仅作用于通配符模式）
	SubDir           bool        // 是否搜索子目录（默认处理子目录）
	Include          int         // 是否处理文件或目录（默认只处理文件）
	OpFunc           OperateFunc // 处理搜索结果的函数
	GetAbsPath       bool        // 是否获取绝对路径（默认获取相对路径）
	RecursionAfterOp bool        // 遇到目录时先递归还是先处理（默认先递归）
}

// 创建新的 FileSearch 结构体
func NewFileSearch() *FileSearch {
	fs := new(FileSearch)
	fs.SubDir = true
	fs.Include = IncludeFile
	fs.GetAbsPath = false
	return fs
}

// 搜索文件或目录，并调用处理函数处理找到的文件或目录，使用通配符 * 和 ?
func (fs *FileSearch) Search() error {
	return fs.search(fs.Dir, "", false)
}

// 搜索文件或目录，并调用处理函数处理找到的文件或目录，使用正则表达式
func (fs *FileSearch) RegSearch() error {
	return fs.search(fs.Dir, "", true)
}

// 搜索文件或目录，并调用处理函数处理找到的文件或目录
func (fs *FileSearch) search(basePath, subPath string, useRegexp bool) error {
	// 未指定文件处理函数
	if fs.OpFunc == nil {
		return errors.New("未指定文件处理函数")
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return err
	}

	// 添加结尾的斜杠符号
	if absPath[len(absPath)-1] != os.PathSeparator {
		absPath += string(os.PathSeparator)
	}

	// 开始获取文件列表（不包括子目录）
	fl, err := ioutil.ReadDir(absPath + subPath)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Fprintf(os.Stderr, "权限不够：%s\n", absPath+subPath)
		} else if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "路径不存在：%s\n", absPath+subPath)
		} else {
			fmt.Fprintf(os.Stderr, "读取出错：%s：%v\n", absPath+subPath, err)
		}
	}

	// 添加结尾的斜杠符号
	if subPath != "" && subPath[len(subPath)-1] != os.PathSeparator {
		subPath += string(os.PathSeparator)
	}

	// 保存文件信息的对象，避免在 for 循环内部反复创建变量
	var Fi FileInfo

	// 遍历文件列表并进行处理
	for _, fi := range fl {
		// 获取文件完整路径
		newSubPath := subPath + fi.Name()

		// 填写文件信息
		Fi.FileInfo = fi
		if fs.GetAbsPath {
			Fi.Path = absPath + subPath
		} else {
			Fi.Path = subPath
		}

		// 先递归处理子目录
		if !fs.RecursionAfterOp && fs.SubDir && Fi.IsDir() {
			err = fs.search(absPath, newSubPath, useRegexp)
			if err != nil {
				return err
			}
		}

		// 开始匹配
		var matched bool
		if useRegexp {
			// 使用正则表达式
			matched, err = regexp.MatchString(fs.KeyWord, Fi.Name())
			if err != nil {
				return err
			}
		} /* else {
			// 使用通配符 * 和 ?
			matched = glob.MatchString(fs.KeyWord, Fi.Name(), fs.CaseMind)
		}*/

		// 匹配成功
		if matched {
			// 如果用户要处理目录
			if Fi.IsDir() && (fs.Include&IncludeDir) == IncludeDir ||
				// 或者用户要处理文件
				!Fi.IsDir() && (fs.Include&IncludeFile) == IncludeFile {
				// 调用处理函数
				err = fs.OpFunc(Fi)
				// 如果处理函数返回错误，则停止执行
				if err != nil {
					return err
				}
			}
		}

		// 后递归处理子目录
		if fs.RecursionAfterOp && fs.SubDir && Fi.IsDir() {
			// 判断目录是否存在，防止用户在 OpFunc 操作中重命名或删除目录
			if !DirExists(absPath + newSubPath) {
				continue
			}
			err = fs.search(absPath, newSubPath, useRegexp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 检查目录是否存在
func DirExists(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

// 用来保存文件列表
var fileList []string

// 搜索指定目录，并返回结果列表，使用通配符 * 和 ?
func (fs *FileSearch) SearchToList() ([]string, error) {
	return fs.searchToList(false)
}

// 搜索指定目录，并返回结果列表，使用正则表达式
func (fs *FileSearch) RegSearchToList() ([]string, error) {
	return fs.searchToList(true)
}

// 搜索指定目录，并返回结果列表
func (fs *FileSearch) searchToList(useRegexp bool) ([]string, error) {
	// 创建用来保存结果的列表
	fileList = make([]string, 0)
	// 指定搜索处理函数
	fs.OpFunc = func(fi FileInfo) error {
		// 获取绝对路径
		fn := fi.Path + fi.Name()
		if fi.IsDir() {
			fn += string(os.PathSeparator)
		}
		fileList = append(fileList, fn)
		return nil
	}
	// 开始搜索
	err := fs.search(fs.Dir, "", useRegexp)
	if err != nil {
		return nil, err
	}
	return fileList, nil
}
