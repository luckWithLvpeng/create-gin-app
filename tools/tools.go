package tools

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//create md5 string
func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

//password hash function
func Pwdhash(str string) string {
	return Strtomd5(str)
}

func StringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}

	return jsons
}

func MyIsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func MyRemoveFile(fileName string) error {
	if MyIsExist(fileName) {
		return os.Remove(fileName)
	} else {
		return errors.New("file is not exist")
	}
}

func MyReadFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return string(""), err
	}
	defer file.Close()
	detailByte, err := ioutil.ReadAll(file)
	return string(detailByte), err
}

func MyWriteFile(fileName string, content string) (int, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return file.WriteString(content)
}

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetParentDirectory() (string, error) {
	var dirctory string = ""
	var err error = nil
	dirctory, err = GetCurrentDirectory()
	if err != nil {
		return dirctory, err
	}
	return Substr(dirctory, 0, strings.LastIndex(dirctory, string(filepath.Separator))), err
}

func GetCurrentDirectory() (string, error) {
	var dirctory string = ""
	var err error = nil
	dirctory, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return dirctory, err
	}
	return dirctory, err
}
