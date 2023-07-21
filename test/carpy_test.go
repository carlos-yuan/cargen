package test

import (
	"github.com/carlos-yuan/cargen/carpy"
	"testing"
)

var cp carpy.Copy

func TestCarpy(t *testing.T) {
	carpy.Gen("D:\\carlos\\cargen")
}
