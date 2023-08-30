package gormutil

import (
	"errors"
	"fmt"
	"regexp"
	"runtime/debug"
	"time"

	e "github.com/carlos-yuan/cargen/core/error"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// 非零值过滤
func nzero(f field.Expr, v interface{}) gen.Condition {
	switch a := v.(type) {
	default:
	case string:
		if a == "" {
			return nil
		} else {
			return f.(field.String).Eq(a)
		}
	case int:
		if a == 0 {
			return nil
		} else {
			return f.(field.Int).Eq(a)
		}
	case int32:
		if a == 0 {
			return nil
		} else {
			return f.(field.Int32).Eq(a)
		}
	case int64:
		if a == 0 {
			return nil
		} else {
			return f.(field.Int64).Eq(a)
		}
	case *int:
		if *a == 0 || a == nil {
			return nil
		} else {
			return f.(field.Int).Eq(*a)
		}
	case *int32:
		if *a == 0 || a == nil {
			return nil
		} else {
			return f.(field.Int32).Eq(*a)
		}
	case *int64:
		if *a == 0 || a == nil {
			return nil
		} else {
			return f.(field.Int64).Eq(*a)
		}
	case []string:
		if len(a) == 0 {
			return nil
		} else {
			return f.(field.String).In(a...)
		}
	case *[]string:
		if len(*a) == 0 || a == nil {
			return nil
		} else {
			a := *a
			return f.(field.String).In(a...)
		}
	case []int:
		if len(a) == 0 {
			return nil
		} else {
			return f.(field.Int).In(a...)
		}
	case *[]int:
		if len(*a) == 0 || a == nil {
			return nil
		} else {
			a := *a
			return f.(field.Int).In(a...)
		}
	case []int32:
		if len(a) == 0 || a == nil {
			return nil
		} else {
			return f.(field.Int32).In(a...)
		}
	case *[]int32:
		if len(*a) == 0 || a == nil {
			return nil
		} else {
			a := *a
			return f.(field.Int32).In(a...)
		}
	case []int64:
		if len(a) == 0 {
			return nil
		} else {
			return f.(field.Int64).In(a...)
		}
	case *[]int64:
		if len(*a) == 0 || a == nil {
			return nil
		} else {
			a := *a
			return f.(field.Int64).In(a...)
		}
	}
	return nil
}

// 非零值过滤
func nzeroExpr(c gen.Condition, v interface{}) gen.Condition {
	switch a := v.(type) {
	default:
	case string:
		if a == "" {
			return nil
		} else {
			return c
		}
	case int:
		if a == 0 {
			return nil
		} else {
			return c
		}
	case int32:
		if a == 0 {
			return nil
		} else {
			return c
		}
	case int64:
		if a == 0 {
			return nil
		} else {
			return c
		}
	case *int:
		if *a == 0 || a == nil {
			return nil
		} else {
			return c
		}
	case *int32:
		if *a == 0 || a == nil {
			return nil
		} else {
			return c
		}
	case *int64:
		if *a == 0 || a == nil {
			return nil
		} else {
			return c
		}
	case []string:
		if len(a) == 0 {
			return nil
		} else {
			return c
		}
	case *[]string:
		if len(*a) == 0 || a == nil {
			return nil
		} else {
			return c
		}
	case []int:
		if len(a) == 0 {
			return nil
		} else {
			return c
		}
	case *[]int:
		if len(*a) == 0 || a == nil {
			return nil
		} else {
			return c
		}
	case []int32:
		if len(a) == 0 || a == nil {
			return nil
		} else {
			return c
		}
	case *[]int32:
		if len(*a) == 0 || a == nil {
			return nil
		} else {
			return c
		}
	case []int64:
		if len(a) == 0 {
			return nil
		} else {
			return c
		}
	case *[]int64:
		if len(*a) == 0 || a == nil {
			return nil
		} else {
			return c
		}
	}
	return nil
}

// 非零可变参数过滤
func some(args ...interface{}) []gen.Condition {
	if len(args)%2 != 0 {
		panic("arg length error")
	}
	var conditions = make([]gen.Condition, 0, len(args))
	for i := 0; i < len(args)/2; i++ {
		idx := i * 2
		a := args[idx]
		switch a.(type) {
		case field.Int:
			c := nzero(a.(field.Int), args[idx+1])
			if c != nil {
				conditions = append(conditions, c)
			}
		case field.Int32:
			c := nzero(a.(field.Int32), args[idx+1])
			if c != nil {
				conditions = append(conditions, c)
			}
		case field.Int64:
			c := nzero(a.(field.Int64), args[idx+1])
			if c != nil {
				conditions = append(conditions, c)
			}
		case field.String:
			c := nzero(a.(field.String), args[idx+1])
			if c != nil {
				conditions = append(conditions, c)
			}
		}
	}
	return conditions
}

// 非零值过滤小于
func nzeroLt(f field.Expr, v interface{}) gen.Condition {
	switch a := f.(type) {
	default:
	case field.Int:
		vt, ok := v.(int)
		if ok {
			return a.Lt(vt)
		} else {
			vp, ok := v.(*int)
			if ok {
				return a.Lt(*vp)
			}
		}
	case field.Int32:
		vt, ok := v.(int32)
		if ok {
			return a.Lt(vt)
		} else {
			vp, ok := v.(*int32)
			if ok {
				return a.Lt(*vp)
			}
		}
	case field.Int64:
		vt, ok := v.(int64)
		if ok {
			return a.Lt(vt)
		} else {
			vp, ok := v.(*int64)
			if ok {
				return a.Lt(*vp)
			}
		}
	case field.String:
		vt, ok := v.(string)
		if ok {
			return a.Lt(vt)
		} else {
			vp, ok := v.(*string)
			if ok {
				return a.Lt(*vp)
			}
		}
	case field.Time:
		vt, ok := v.(string)
		if ok {
			if vt != "" {
				t, _ := time.Parse("2006-01-02 15:04:05", vt)
				return a.Lt(t)
			}
		} else {
			vp, ok := v.(*string)
			if ok {
				if vp != nil && *vp != "" {
					t, _ := time.Parse("2006-01-02 15:04:05", *vp)
					return a.Lt(t)
				}
			} else {
				vts, ok := v.(time.Time)
				if ok && !vts.IsZero() {
					return a.Lt(vts)
				} else {
					vpts, ok := v.(*time.Time)
					if ok && vpts != nil && !vpts.IsZero() {
						return a.Lt(*vpts)
					}
				}
			}
		}
	}
	return nil
}

// 非零值过滤小于
func nzeroLte(f field.Expr, v interface{}) gen.Condition {
	switch a := f.(type) {
	default:
	case field.Int:
		vt, ok := v.(int)
		if ok {
			return a.Lte(vt)
		} else {
			vp, ok := v.(*int)
			if ok {
				return a.Lte(*vp)
			}
		}
	case field.Int32:
		vt, ok := v.(int32)
		if ok {
			return a.Lte(vt)
		} else {
			vp, ok := v.(*int32)
			if ok {
				return a.Lte(*vp)
			}
		}
	case field.Int64:
		vt, ok := v.(int64)
		if ok {
			return a.Lte(vt)
		} else {
			vp, ok := v.(*int64)
			if ok {
				return a.Lte(*vp)
			}
		}
	case field.String:
		vt, ok := v.(string)
		if ok {
			return a.Lte(vt)
		} else {
			vp, ok := v.(*string)
			if ok {
				return a.Lte(*vp)
			}
		}
	case field.Time:
		vt, ok := v.(string)
		if ok {
			if vt != "" {
				t, _ := time.Parse("2006-01-02 15:04:05", vt)
				return a.Lte(t)
			}
		} else {
			vp, ok := v.(*string)
			if ok {
				if vp != nil && *vp != "" {
					t, _ := time.Parse("2006-01-02 15:04:05", *vp)
					return a.Lte(t)
				}
			} else {
				vts, ok := v.(time.Time)
				if ok && !vts.IsZero() {
					return a.Lte(vts)
				} else {
					vpts, ok := v.(*time.Time)
					if ok && vpts != nil && !vpts.IsZero() {
						return a.Lte(*vpts)
					}
				}
			}
		}
	}
	return nil
}

// 非零值过滤大于
func nzeroGt(f field.Expr, v interface{}) gen.Condition {
	switch a := f.(type) {
	default:
	case field.Int:
		vt, ok := v.(int)
		if ok {
			return a.Gt(vt)
		} else {
			vp, ok := v.(*int)
			if ok {
				return a.Gt(*vp)
			}
		}
	case field.Int32:
		vt, ok := v.(int32)
		if ok {
			return a.Gt(vt)
		} else {
			vp, ok := v.(*int32)
			if ok {
				return a.Gt(*vp)
			}
		}
	case field.Int64:
		vt, ok := v.(int64)
		if ok {
			return a.Gt(vt)
		} else {
			vp, ok := v.(*int64)
			if ok {
				return a.Gt(*vp)
			}
		}
	case field.String:
		vt, ok := v.(string)
		if ok {
			return a.Gt(vt)
		} else {
			vp, ok := v.(*string)
			if ok {
				return a.Gt(*vp)
			}
		}
	case field.Time:
		vt, ok := v.(string)
		if ok {
			if vt != "" {
				t, _ := time.Parse("2006-01-02 15:04:05", vt)
				return a.Gt(t)
			}
		} else {
			vp, ok := v.(*string)
			if ok {
				if vp != nil && *vp != "" {
					t, _ := time.Parse("2006-01-02 15:04:05", *vp)
					return a.Gt(t)
				}
			} else {
				vts, ok := v.(time.Time)
				if ok && !vts.IsZero() {
					return a.Gt(vts)
				} else {
					vpts, ok := v.(*time.Time)
					if ok && vpts != nil && !vpts.IsZero() {
						return a.Gt(*vpts)
					}
				}
			}
		}
	}
	return nil
}

// 非零值过滤大于
func nzeroGte(f field.Expr, v interface{}) gen.Condition {
	switch a := f.(type) {
	default:
	case field.Int:
		vt, ok := v.(int)
		if ok {
			return a.Gte(vt)
		} else {
			vp, ok := v.(*int)
			if ok {
				return a.Gte(*vp)
			}
		}
	case field.Int32:
		vt, ok := v.(int32)
		if ok {
			return a.Gte(vt)
		} else {
			vp, ok := v.(*int32)
			if ok {
				return a.Gte(*vp)
			}
		}
	case field.Int64:
		vt, ok := v.(int64)
		if ok {
			return a.Gte(vt)
		} else {
			vp, ok := v.(*int64)
			if ok {
				return a.Gte(*vp)
			}
		}
	case field.String:
		vt, ok := v.(string)
		if ok {
			return a.Gte(vt)
		} else {
			vp, ok := v.(*string)
			if ok {
				return a.Gte(*vp)
			}
		}
	case field.Time:
		vt, ok := v.(string)
		if ok {
			if vt != "" {
				t, _ := time.Parse("2006-01-02 15:04:05", vt)
				return a.Gte(t)
			}
		} else {
			vp, ok := v.(*string)
			if ok {
				if vp != nil && *vp != "" {
					t, _ := time.Parse("2006-01-02 15:04:05", *vp)
					return a.Gte(t)
				}
			} else {
				vts, ok := v.(time.Time)
				if ok && !vts.IsZero() {
					return a.Gte(vts)
				} else {
					vpts, ok := v.(*time.Time)
					if ok && vpts != nil && !vpts.IsZero() {
						return a.Gte(*vpts)
					}
				}
			}
		}
	}
	return nil
}

func checkInjection(data string) error {
	str := `(?:')|(?:--)|(/\*(?:.|[\\n\\r])*?\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	if re, err := regexp.Compile("(?i)" + str); err != nil {
		return errors.New("check injection fail")
	} else {
		if re.MatchString(data) {
			return errors.New("check injection fail")
		}
	}
	return nil
}

func RollBackFn(tx *gorm.DB, me e.Err, err *error) {
	r := recover()
	if r != nil {
		me = me.SetRecover(r)
		*err = &me
		s := string(debug.Stack())
		fmt.Printf("err=%v, stack=%s\n", err, s)
	}
	if err != nil && *err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
