package gormutil

import (
	"reflect"
	"time"

	"github.com/carlos-yuan/cargen/util/cartime"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

// RegisterCartimeCallbacks cartime注入 替换掉原来的时间戳
func RegisterCartimeCallbacks(db *gorm.DB) error {
	db.NowFunc = func() time.Time { //注入cartime
		return time.Unix(cartime.NowToInt(), 0)
	}
	err := db.Callback().Query().Before("gorm:query").Register("find_not_delete", QueryCallback)
	if err != nil {
		return err
	}
	return nil
}

// QueryCallback 自定义查询时间标记
func QueryCallback(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil && !db.Statement.SkipHooks {
		field := db.Statement.Schema.LookUpField("DeletedAt")
		if field != nil && field.FieldType.Kind() == reflect.Int64 {
			if conds := db.Statement.BuildCondition("`" + db.Statement.Table + "`" + ".`deleted_at`=0"); len(conds) > 0 {
				db.Statement.AddClause(clause.Where{Exprs: conds})
			}
		} else {
			callbacks.Query(db)
		}
	}
}
