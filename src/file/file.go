package util

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

/**
文件追加
*/
func AppendToFile(fPath string, content string) {

	dPath := path.Dir(fPath)
	existed, err := CheckExisted(dPath)
	Check(err)
	// 检查目录是否存在，不存在创建之
	if !existed {
		MakeDir(dPath)
	}

	f, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	defer f.Close()
	Check(err)

	// fix end by '\n'
	if len(content) <= 0 {
		return
	}
	if strings.LastIndex(content, "\n") < len(content)-1 {
		content += "\n"
	}

	_, err = f.WriteString(content)
	Check(err)
}

/**
读取文件内容
*/
func ReadFile(fPath string) []byte {
	existed, err := CheckExisted(fPath)
	Check(err)
	if !existed {
		return nil
	}

	bytes, err := ioutil.ReadFile(fPath)
	Check(err)
	return bytes
}

/**
删除指定文件
*/
func DeleteFile(fPath string) {
	existed, err := CheckExisted(fPath)
	Check(err)
	if existed {
		err = os.Remove(fPath)
		Check(err)
	}
}

/**
替换指定行内容
*/
func ReplaceFileLine(fPath, oldContent, newContent string) {
	data, err := ioutil.ReadFile(fPath)
	Check(err)

	content := string(data)
	lines := strings.Split(content, "\n")
	var ln int
	for k, v := range lines {
		if v == oldContent {
			ln = k
		}
	}

	var newLines []string
	newLines = append(newLines, lines[:ln]...)
	newLines = append(newLines, newContent)
	newLines = append(newLines, lines[ln+1:]...)

	ioutil.WriteFile(fPath, []byte(strings.Join(newLines, "\n")), os.ModePerm)
}

/**
检查文件或者目录是否存在
*/
func CheckExisted(anyPath string) (bool, error) {
	_, err := os.Stat(anyPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

/**
创建文件夹
*/
func MakeDir(dPath string) {
	existed, err := CheckExisted(dPath)
	Check(err)
	if !existed {
		err = os.Mkdir(dPath, os.ModePerm)
		Check(err)
	}
}

/**
创建文件
*/
func CreateFile(fPath string) string {
	existed, err := CheckExisted(path.Dir(fPath))
	Check(err)
	if !existed {
		MakeDir(path.Dir(fPath))
	}
	f, err := os.Create(fPath)
	defer f.Close()
	Check(err)
	return f.Name()
}
