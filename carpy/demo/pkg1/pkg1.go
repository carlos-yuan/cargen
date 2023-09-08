package pkg1

import (
	"github.com/carlos-yuan/cargen/carpy/demo/pkg2"
	"github.com/carlos-yuan/cargen/enum"
)

type Pkg1 struct {
	Name     string
	Count    []int
	TestDict map[string]enum.Dict
	Pkg      map[string]pkg2.Pkg2
}
