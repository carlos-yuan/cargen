package demo

import (
	"errors"
	"github.com/carlos-yuan/cargen/carpy"
)

func init() {
	cp = &mainCp{}
}

type mainCp struct {
}

func (c mainCp) Copy(to any, from any) error {
	switch to := to.(type) {
	case *Cp1:
		switch from := from.(type) {
		case *Cp2:
			CopyCp2ToCp1(to, from)
		case *Cp3:
			CopyCp3ToCp1(to, from)
		case *Cp4:
			CopyCp4ToCp1(to, from)
		case *Cp5:
			CopyCp5ToCp1(to, from)
		case *Cp6:
			CopyCp6ToCp1(to, from)
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}
	case *Cp2:
		switch from := from.(type) {
		case *Cp2:
			CopyCp2ToCp2(to, from)
		case *Cp3:
			CopyCp3ToCp2(to, from)
		case *Cp4:
			CopyCp4ToCp2(to, from)
		case *Cp5:
			CopyCp5ToCp2(to, from)
		case *Cp6:
			CopyCp6ToCp2(to, from)
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}
	case *Cp3:
		switch from := from.(type) {
		case *Cp2:
			CopyCp2ToCp3(to, from)
		case *Cp3:
			CopyCp3ToCp3(to, from)
		case *Cp4:
			CopyCp4ToCp3(to, from)
		case *Cp5:
			CopyCp5ToCp3(to, from)
		case *Cp6:
			CopyCp6ToCp3(to, from)
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}
	case *Cp4:
		switch from := from.(type) {
		case *Cp2:
			CopyCp2ToCp4(to, from)
		case *Cp3:
			CopyCp3ToCp4(to, from)
		case *Cp4:
			CopyCp4ToCp4(to, from)
		case *Cp5:
			CopyCp5ToCp4(to, from)
		case *Cp6:
			CopyCp6ToCp4(to, from)
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}
	case *Cp5:
		switch from := from.(type) {
		case *Cp2:
			CopyCp2ToCp5(to, from)
		case *Cp3:
			CopyCp3ToCp5(to, from)
		case *Cp4:
			CopyCp4ToCp5(to, from)
		case *Cp5:
			CopyCp5ToCp5(to, from)
		case *Cp6:
			CopyCp6ToCp5(to, from)
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}
	case *Cp6:
		switch from := from.(type) {
		case *Cp2:
			CopyCp2ToCp6(to, from)
		case *Cp3:
			CopyCp3ToCp6(to, from)
		case *Cp4:
			CopyCp4ToCp6(to, from)
		case *Cp5:
			CopyCp5ToCp6(to, from)
		case *Cp6:
			CopyCp6ToCp6(to, from)
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}
	default:
		return errors.New("unknown copy to " + carpy.GetTypeName(to))
	}
	return nil
}

func CopyCp6ToCp6(to *Cp6, from *Cp6) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp5ToCp6(to *Cp6, from *Cp5) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp4ToCp6(to *Cp6, from *Cp4) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp3ToCp6(to *Cp6, from *Cp3) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp2ToCp6(to *Cp6, from *Cp2) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp6ToCp5(to *Cp5, from *Cp6) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp5ToCp5(to *Cp5, from *Cp5) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp4ToCp5(to *Cp5, from *Cp4) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp3ToCp5(to *Cp5, from *Cp3) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp2ToCp5(to *Cp5, from *Cp2) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp6ToCp4(to *Cp4, from *Cp6) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp5ToCp4(to *Cp4, from *Cp5) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp4ToCp4(to *Cp4, from *Cp4) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp3ToCp4(to *Cp4, from *Cp3) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp2ToCp4(to *Cp4, from *Cp2) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp6ToCp3(to *Cp3, from *Cp6) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp5ToCp3(to *Cp3, from *Cp5) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
	to.Name7 = from.Name7
	to.Name8 = from.Name8
	to.Name9 = from.Name9
	to.Name10 = from.Name10
	to.Name11 = from.Name11
	to.Name12 = from.Name12
	to.Name13 = from.Name13
	to.Name14 = from.Name14
	to.Name15 = from.Name15
	to.Name16 = from.Name16
	to.Name17 = from.Name17
	to.Name18 = from.Name18
	to.Name19 = from.Name19
	to.Name20 = from.Name20
	to.Name21 = from.Name21
	to.Name22 = from.Name22
	to.Name23 = from.Name23
	to.Name24 = from.Name24
	to.Name25 = from.Name25
	to.Name26 = from.Name26
	to.Name27 = from.Name27
	to.Name28 = from.Name28
	to.Name29 = from.Name29
	to.Name30 = from.Name30
}

func CopyCp4ToCp3(to *Cp3, from *Cp4) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp6ToCp2(to *Cp2, from *Cp6) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp3ToCp3(to *Cp3, from *Cp3) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp5ToCp2(to *Cp2, from *Cp5) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp4ToCp2(to *Cp2, from *Cp4) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp3ToCp2(to *Cp2, from *Cp3) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp2ToCp2(to *Cp2, from *Cp2) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp6ToCp1(to *Cp1, from *Cp6) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp5ToCp1(to *Cp1, from *Cp5) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp4ToCp1(to *Cp1, from *Cp4) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp3ToCp1(to *Cp1, from *Cp3) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp2ToCp3(to *Cp3, from *Cp2) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

func CopyCp2ToCp1(to *Cp1, from *Cp2) {
	to.Name = from.Name
	to.Name2 = from.Name2
	to.Name3 = from.Name3
	to.Name4 = from.Name4
	to.Name5 = from.Name5
	to.Name6 = from.Name6
}

type Cp1 struct {
	Name   string
	Name2  string
	Name3  string
	Name4  string
	Name5  string
	Name6  string
	Name7  string
	Name8  string
	Name9  string
	Name10 string
	Name11 string
	Name12 string
	Name13 string
	Name14 string
	Name15 string
	Name16 string
	Name17 string
	Name18 string
	Name19 string
	Name20 string
	Name21 string
	Name22 string
	Name23 string
	Name24 string
	Name25 string
	Name26 string
	Name27 string
	Name28 string
	Name29 string
	Name30 string
}

type Cp2 struct {
	Name   string
	Name2  string
	Name3  string
	Name4  string
	Name5  string
	Name6  string
	Name7  string
	Name8  string
	Name9  string
	Name10 string
	Name11 string
	Name12 string
	Name13 string
	Name14 string
	Name15 string
	Name16 string
	Name17 string
	Name18 string
	Name19 string
	Name20 string
	Name21 string
	Name22 string
	Name23 string
	Name24 string
	Name25 string
	Name26 string
	Name27 string
	Name28 string
	Name29 string
	Name30 string
}

type Cp3 struct {
	Name   string
	Name2  string
	Name3  string
	Name4  string
	Name5  string
	Name6  string
	Name7  string
	Name8  string
	Name9  string
	Name10 string
	Name11 string
	Name12 string
	Name13 string
	Name14 string
	Name15 string
	Name16 string
	Name17 string
	Name18 string
	Name19 string
	Name20 string
	Name21 string
	Name22 string
	Name23 string
	Name24 string
	Name25 string
	Name26 string
	Name27 string
	Name28 string
	Name29 string
	Name30 string
}

type Cp4 struct {
	Name   string
	Name2  string
	Name3  string
	Name4  string
	Name5  string
	Name6  string
	Name7  string
	Name8  string
	Name9  string
	Name10 string
	Name11 string
	Name12 string
	Name13 string
	Name14 string
	Name15 string
	Name16 string
	Name17 string
	Name18 string
	Name19 string
	Name20 string
	Name21 string
	Name22 string
	Name23 string
	Name24 string
	Name25 string
	Name26 string
	Name27 string
	Name28 string
	Name29 string
	Name30 string
}
type Cp5 struct {
	Name   string
	Name2  string
	Name3  string
	Name4  string
	Name5  string
	Name6  string
	Name7  string
	Name8  string
	Name9  string
	Name10 string
	Name11 string
	Name12 string
	Name13 string
	Name14 string
	Name15 string
	Name16 string
	Name17 string
	Name18 string
	Name19 string
	Name20 string
	Name21 string
	Name22 string
	Name23 string
	Name24 string
	Name25 string
	Name26 string
	Name27 string
	Name28 string
	Name29 string
	Name30 string
}

type Cp6 struct {
	Name   string
	Name2  string
	Name3  string
	Name4  string
	Name5  string
	Name6  string
	Name7  string
	Name8  string
	Name9  string
	Name10 string
	Name11 string
	Name12 string
	Name13 string
	Name14 string
	Name15 string
	Name16 string
	Name17 string
	Name18 string
	Name19 string
	Name20 string
	Name21 string
	Name22 string
	Name23 string
	Name24 string
	Name25 string
	Name26 string
	Name27 string
	Name28 string
	Name29 string
	Name30 string
}
