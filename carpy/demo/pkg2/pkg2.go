package pkg2

type Pkg2 struct {
	Name  string
	Count int32
	Pkg   map[string]Pkg4
	Pkg2  map[string]map[string]Pkg4
	Pkg3  []map[string]map[string]Pkg4
	Pkg4  []map[string][]map[string]Pkg4
}

type Pkg4 struct {
	Name  string
	Count int32
}
