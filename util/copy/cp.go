package cp

import (
	"comm/cartime"
	"errors"
	"time"

	"github.com/jinzhu/copier"
)

var StringToTime = copier.Option{IgnoreEmpty: true, Converters: []copier.TypeConverter{
	{
		SrcType: copier.String, DstType: &time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(string)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			if s == "" {
				return nil, nil
			}
			t, err := time.Parse("2006-01-02 15:04:05", s)
			return &t, err
		},
	},
	{
		SrcType: copier.String, DstType: time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(string)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			if s == "" {
				return nil, nil
			}
			t, err := time.Parse("2006-01-02 15:04:05", s)
			return t, err
		},
	},
}}

const carTimeStr = cartime.TimeStr("")

const i64 = int64(0)

var StringTimeToIntTime = copier.Option{IgnoreEmpty: true, Converters: []copier.TypeConverter{
	{
		SrcType: carTimeStr, DstType: i64,
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(cartime.TimeStr)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			if s == "" {
				return nil, nil
			}
			var time int64
			if len(s) == 10 {
				time = cartime.StrToInt(string(s), cartime.DefaultFormatDate)
			} else if len(s) == 19 {
				time = cartime.StrToInt(string(s), cartime.DefaultFormat)
			}
			return time, nil
		},
	},
	{
		SrcType: i64, DstType: carTimeStr,
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(int64)
			if !ok {
				return "", errors.New("src type not matching")
			}
			if s == 0 {
				return "", nil
			}
			var time cartime.TimeStr
			time = cartime.TimeStr(cartime.IntToStr(s, cartime.DefaultFormat))
			return time, nil
		},
	},
}}

var TimeToString = copier.Option{IgnoreEmpty: true, Converters: []copier.TypeConverter{
	{
		SrcType: &time.Time{}, DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			if src != nil {
				s, ok := src.(*time.Time)
				if !ok {
					return nil, errors.New("src type not matching")
				}
				return s.Format("2006-01-02 15:04:05"), nil
			}
			return "", nil
		},
	},
	{
		SrcType: time.Time{}, DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			if src != nil {
				s, ok := src.(time.Time)
				if !ok {
					return nil, errors.New("src type not matching")
				}
				return s.Format("2006-01-02 15:04:05"), nil
			}
			return "", nil
		},
	},
}}

var CopyTimeStringOpt = func(format ...string) copier.Option {
	if len(format) == 0 {
		format = []string{time.RFC3339}
	}
	return copier.Option{IgnoreEmpty: true, Converters: []copier.TypeConverter{{
		SrcType: &time.Time{}, DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			if src != nil {
				s, ok := src.(*time.Time)
				if !ok {
					return nil, errors.New("src type not matching")
				}
				return s.Format(format[0]), nil
			}
			return "", nil
		},
	},
		{
			SrcType: copier.String, DstType: &time.Time{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(string)
				if !ok {
					return nil, errors.New("src type not matching")
				}
				if s == "" {
					return nil, nil
				}
				t, err := time.Parse(format[0], s)
				return &t, err
			},
		},
		{
			SrcType: time.Time{}, DstType: copier.String,
			Fn: func(src interface{}) (interface{}, error) {
				if src != nil {
					s, ok := src.(time.Time)
					if !ok {
						return nil, errors.New("src type not matching")
					}
					return s.Format(format[0]), nil
				}
				return "", nil
			},
		},
		{
			SrcType: copier.String, DstType: time.Time{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(string)
				if !ok {
					return nil, errors.New("src type not matching")
				}
				if s == "" {
					return nil, nil
				}
				t, err := time.Parse(format[0], s)
				return t, err
			},
		},
	}}
}

func CopyWithTime2String(to any, from any) error {
	return copier.CopyWithOption(to, from, TimeToString)
}
