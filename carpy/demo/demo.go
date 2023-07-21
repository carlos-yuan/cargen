package demo

import (
	"github.com/carlos-yuan/cargen/carpy"
	"testing"
)

var cp carpy.Copy

func TestCopy(t *testing.T) {
	var c1 Cp1
	var c2 Cp2
	var c2s = &Cp2{}
	cp.Copy(&c1, &c2)
	cp.Copy(&c1, c2)
	cp.Copy(&c1, &Cp2{})
	cp.Copy(&c1, Cp2{})
	cp.Copy(&c1, *c2s)

}
