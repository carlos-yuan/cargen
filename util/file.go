package util

import (
	"errors"
	"strings"
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
