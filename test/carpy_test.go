package test

import (
	"github.com/carlos-yuan/cargen/carpy"
	"github.com/carlos-yuan/cargen/carpy/demo"
	"testing"
)

func TestCarpy(t *testing.T) {
	carpy.Gen("D:\\carlos\\cargen")
	demo.TestCopy()
}
