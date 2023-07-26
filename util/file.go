package util

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"sync"
)

func CutPathLast(path string, count int) (string, error) {
	for i := 0; i < count; i++ {
		idx := strings.LastIndex(path, `/`)
		if idx == -1 {
			idx = strings.LastIndex(path, `\`)
			if idx == -1 {
				return "", errors.New("count than max")
			}
		}
		path = path[:idx]
	}
	return path, nil
}

func CutPath(path string, count int) (string, error) {
	for i := 0; i < count; i++ {
		idx := strings.Index(path, `/`)
		if idx == -1 {
			idx = strings.Index(path, `\`)
			if idx == -1 {
				return "", errors.New("count than max")
			}
		}
		path = path[idx+1:]
	}
	return path, nil
}

func FixPathSeparator(path string) string {
	if string(os.PathSeparator) == "\\" {
		return strings.ReplaceAll(path, "/", string(os.PathSeparator))
	} else {
		return strings.ReplaceAll(path, "\\", string(os.PathSeparator))
	}
}

func LastName(path string) string {
	idx := strings.LastIndex(path, `/`)
	if idx == -1 {
		idx = strings.LastIndex(path, `\`)
		if idx == -1 {
			return path
		}
	}
	return path[idx+1:]
}

var mutex sync.Mutex

func WriteByteFile(filePath string, data []byte) error {
	mutex.Lock()
	defer mutex.Unlock()
	fileCheck, err := os.Stat(filePath)
	if err == nil {
		if fileCheck.IsDir() {
			return errors.New("path is directory")
		}
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666) //修改模式,变更为写入模式O_WRONLY和清除模式O_TRUNC
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.Write(data)
		if err == nil {
			return err
		}
	}
	err = CreateFilePath(filePath)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

func CreateFilePath(path string) error {
	if !IsFile(path) {
		path = strings.Replace(path, "\\", "/", -1)
		path = path[:strings.LastIndex(path, "/")]
		if !IsExist(path) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("path is file")
	}
	return nil
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	boolean := os.IsExist(err)
	return err == nil || boolean
}

func FileInfo(f string) os.FileInfo {
	fi, e := os.Stat(f)
	if e != nil {
		return nil
	}
	if !fi.IsDir() {
		return fi
	}
	return nil
}

func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

func ProjectPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("not find code file path")
	}
	base, err := CutPathLast(filename, 2)
	if err != nil {
		return "", err
	}
	return base, nil
}

func GetAllPath(path string) ([]string, error) {
	infos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := path + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := GetAllPath(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			paths = append(paths, path)
			continue
		}
	}
	return paths, nil
}

func GetFilePath(filepath string, fileName string) ([]string, error) {
	infos, err := os.ReadDir(filepath)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := filepath + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := GetFilePath(path, fileName)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		if info.Name() == fileName {
			paths = append(paths, path)
		}
	}
	return paths, nil
}
