package demo

import (
	"github.com/carlos-yuan/cargen/carpy"
	"github.com/carlos-yuan/cargen/carpy/demo/pkg1"
	"github.com/carlos-yuan/cargen/carpy/demo/pkg2"
	"github.com/carlos-yuan/cargen/enum"
	"github.com/jinzhu/copier"
)

var cp carpy.Copy

func TestCopy() {
	var count1 int64 = 1
	name1 := "1"
	var c1 = Pkg{Name: &name1, Count: []*int64{&count1}}
	count := 2
	var c2 = pkg1.Pkg1{Name: "2", Count: []int{count}, TestDict: map[string]enum.Dict{"1": {Name: "1"}}, Pkg: map[string]pkg2.Pkg2{"2": {Name: name1}}}
	var c3 = pkg2.Pkg2{Name: "3", Count: 3, Pkg: map[string]pkg2.Pkg4{"3": {Name: name1}}}
	err := copier.Copy(&c3, &c2)
	if err != nil {
		panic(err)
	}
	err = cp.Copy(&c1, &pkg1.Pkg1{Name: "2", Count: []int{count}, TestDict: map[string]enum.Dict{"1": {Name: "1"}}}, func(to any, from any) (any, error) {
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
	err = cp.Copy(&c2, &c3)
	if err != nil {
		panic(err)
	}
	err = cp.Copy(&c3, &c1, func(to any, from any) (any, error) {
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
	Name      *string
	Count     []*int64
	TestDict  *map[string]enum.Dict
	TestDict2 enum.Dict
}
