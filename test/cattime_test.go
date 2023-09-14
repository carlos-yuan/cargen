package test

import (
	"github.com/carlos-yuan/cargen/util/cartime"
	"testing"
)

func TestCattime(t *testing.T) {
	println(cartime.IntToStr(2023091412410490808, "20060102150405.000Z07:00"))
	println(cartime.StrToInt("2016-01-12", cartime.DefaultFormatDate))

}
