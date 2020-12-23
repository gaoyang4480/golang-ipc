package ipc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func IsDir(path string) (bool, error) {
	fi, err := GetFileInfo(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func Rename(srcFileFullPath, dstFileFullPath string) error {
	return os.Rename(srcFileFullPath, dstFileFullPath)
}

func DelFile(path string) error {
	return os.Remove(path)
}

func DelDir(path string, ifDelSelf bool) error {
	if ifDelSelf { // 删除自身.
		return os.RemoveAll(path)
	} else { // 不删除自身.
		// TODO:
	}

	return nil
}

func CreateDirIfNotExists(dstFilePath string) error {
	if !IsDirExists(dstFilePath) {
		err := CreateDir(dstFilePath, 666)
		if err != nil {
			errInfo := "创建目录" + dstFilePath + "失败, 原因: " + err.Error()
			return errors.New(errInfo)
		}
	}

	return nil
}

func CreateDir(path string, mode os.FileMode) error {
	err := os.MkdirAll(path, mode)
	return err
}

func CopyFile(srcPath, dstPath string) (w int64, err error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)

	if err != nil {
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func CopyDir(srcPath, dstPath string) error {
	// TODO: ...
	return nil
}

func GetFileInfo(path string) (fi os.FileInfo, err error) {
	fi, err = os.Stat(path)
	return fi, err
}

func GetSubDirlist(path string) (subDirNames []string) {
	err := filepath.Walk(path, func(tmpPath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() && path != tmpPath {
			subDirNames = append(subDirNames, tmpPath)
			return nil
		}
		//println(tmpPath)
		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return nil
	} else {
		return subDirNames
	}
}

// 获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤.
func ListDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

// 获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤.
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil { //忽略错误
			return err
		}
		if fi == nil {
			return nil
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func ReadAll(path string) ([]byte, error) {
	fi, err := os.Open(path)
	if err != nil {
		//panic(err)
		return nil, err
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		//panic(err)
		return nil, err
	}

	return fd, nil
	// fmt.Println(string(fd))
	//return string(fd)
}

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		//err1 := handler(line)
		//if err1 != nil {
		//	return err1
		//}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

	return nil
}

func WriteFile(filePath string, content []byte) error {
	err := ioutil.WriteFile(filePath, content, 0666)
	if err != nil {
		errInfo := "写文件" + filePath + "失败, 原因: " + err.Error()
		return errors.New(errInfo)
	}

	return nil
}
