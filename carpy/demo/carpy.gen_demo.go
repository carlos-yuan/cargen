package demo

import (
	"errors"
	"github.com/carlos-yuan/cargen/carpy"
	pkg1a20b31f2907910b1 "github.com/carlos-yuan/cargen/carpy/demo/pkg1"
	pkg28013fd659cb047fb "github.com/carlos-yuan/cargen/carpy/demo/pkg2"
)

func init() {
	cp = &carpydemo{}
}

type carpydemo struct{}

func (c *carpydemo) Copy(to any, from any, opts ...carpy.CopyOption) error {
	if to == nil || from == nil {
		return nil
	}
	switch to := to.(type) {

	case *Pkg:
		switch from := from.(type) {

		case *pkg1a20b31f2907910b1.Pkg1:
			return CopyPkg1ToPkg5034dfa3363c1a27(to, from, opts...)

		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}

	case *pkg1a20b31f2907910b1.Pkg1:
		switch from := from.(type) {

		case *Pkg:
			return CopyPkgToPkg1f69e3fbc82aa5597(to, from, opts...)

		case *pkg28013fd659cb047fb.Pkg2:
			return CopyPkg2ToPkg12b34611c8d60e091(to, from, opts...)

		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}

	case *pkg28013fd659cb047fb.Pkg2:
		switch from := from.(type) {

		case *Pkg:
			return CopyPkgToPkg2b11976f4e87833b6(to, from, opts...)

		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}

	default:
		return errors.New("unknown copy to " + carpy.GetTypeName(to))

	}
}

func CopyPkg1ToPkg5034dfa3363c1a27(to *Pkg, from *pkg1a20b31f2907910b1.Pkg1, opts ...carpy.CopyOption) (err error) {
	to.Name = from.Name
	to.Count = int64(from.Count)

	return
}

func CopyPkgToPkg1f69e3fbc82aa5597(to *pkg1a20b31f2907910b1.Pkg1, from *Pkg, opts ...carpy.CopyOption) (err error) {
	to.Name = from.Name
	to.Count = int(from.Count)

	return
}

func CopyPkg2ToPkg12b34611c8d60e091(to *pkg1a20b31f2907910b1.Pkg1, from *pkg28013fd659cb047fb.Pkg2, opts ...carpy.CopyOption) (err error) {
	to.Name = from.Name
	to.Count = int(from.Count)

	return
}

func CopyPkgToPkg2b11976f4e87833b6(to *pkg28013fd659cb047fb.Pkg2, from *Pkg, opts ...carpy.CopyOption) (err error) {
	to.Name = from.Name

	return
}
