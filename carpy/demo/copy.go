package demo

import (
	"github.com/carlos-yuan/cargen/carpy"
	"github.com/carlos-yuan/cargen/carpy/demo/pkg1"
	"github.com/carlos-yuan/cargen/carpy/demo/pkg2"
)

var cp carpy.Copy

func TestCopy() {
	var c1 = Pkg{Name: "1", Count: 1}
	var c2 = pkg1.Pkg1{Name: "2", Count: 2}
	var c3 = pkg2.Pkg2{Name: "3", Count: 3}
	err := cp.Copy(&c1, &pkg1.Pkg1{Name: "2", Count: 2})
	if err != nil {
		panic(err)
	}
	err = cp.Copy(&c2, &c3)
	if err != nil {
		panic(err)
	}
	err = cp.Copy(&c3, &c1, func(to any, from any) (any, error) {
		err := cp.Copy(&c2, &c3)
		if err != nil {
			panic(err)
		}
		switch to.(type) {
		case int32:
			switch from.(type) {
			case int:
				return int32(from.(int)), nil
			case int64:
				return int32(from.(int64)), nil
			}
		}
		return to, nil
	})
	if err != nil {
		panic(err)
	}
}

type Pkg struct {
	Name  string
	Count int64
}

func (p *Pkg) CopyTest() {
	var c1 = Pkg{Name: "1", Count: 1}
	var c2 = pkg1.Pkg1{Name: "2", Count: 2}
	cp.Copy(&c1, &c2)
}
