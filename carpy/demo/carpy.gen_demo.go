package demo

import (
	"errors"
	"github.com/carlos-yuan/cargen/carpy"
	"github.com/carlos-yuan/cargen/carpy/demo/pkg1"
	"github.com/carlos-yuan/cargen/carpy/demo/pkg2"
)

func init() {
	cp = &mainCp{}
}

type mainCp struct{}

func (c *mainCp) Copy(to any, from any, opts ...carpy.CopyOption) error {
	if to == nil || from == nil {
		return nil
	}
	switch to := to.(type) {
	case *Pkg:
		switch from := from.(type) {
		case *pkg1.Pkg1:
			return CopyCargenCarpyDemoPkg1Pkg1ToCargenCarpyDemoPkg(to, from, opts...)
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}
	case *pkg1.Pkg1:
		switch from := from.(type) {
		case *pkg2.Pkg2:
			return CopyCargenCarpyDemoPkg2Pkg2ToCargenCarpyDemoPkg1Pkg1(to, from, opts...)
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}
	case *pkg2.Pkg2:
		switch from := from.(type) {
		case *Pkg:
			return CopyCargenCarpyDemoPkgToCargenCarpyDemoPkg2Pkg2(to, from, opts...)
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}
	default:
		return errors.New("unknown copy to " + carpy.GetTypeName(to))
	}
}

func CopyCargenCarpyDemoPkg2Pkg2ToCargenCarpyDemoPkg1Pkg1(to *pkg1.Pkg1, from *pkg2.Pkg2, opts ...carpy.CopyOption) (err error) {
	to.Name = from.Name
	to.Count = int(from.Count)
	return
}

func CopyCargenCarpyDemoPkgToCargenCarpyDemoPkg2Pkg2(to *pkg2.Pkg2, from *Pkg, opts ...carpy.CopyOption) (err error) {
	to.Name = from.Name
	for _, opt := range opts {
		dst, err := opt(to.Count, from.Count)
		if err != nil {
			return err
		}
		if count, ok := dst.(int32); ok {
			to.Count = count
		}
	}
	return
}

func CopyCargenCarpyDemoPkg1Pkg1ToCargenCarpyDemoPkg(to *Pkg, from *pkg1.Pkg1, opts ...carpy.CopyOption) (err error) {
	to.Name = from.Name
	to.Count = int64(from.Count)
	return
}
