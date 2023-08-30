package gormutil

import (
	"reflect"

	"github.com/carlos-yuan/cargen/util/cartime"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

func RegisterCallbacks(db *gorm.DB) error {
	err := db.Callback().Create().Replace("gorm:before_create", CreateCallback)
	if err != nil {
		return err
	}
	err = db.Callback().Update().Replace("gorm:before_update", UpdateCallback)
	if err != nil {
		return err
	}
	err = db.Callback().Delete().Replace("gorm:before_delete", DeleteCallback)
	if err != nil {
		return err
	}
	err = db.Callback().Query().Before("gorm:query").Register("find_not_delete", QueryCallback)
	if err != nil {
		return err
	}
	return nil
}

// CreateCallback 自定义创建时间标记
func CreateCallback(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil && !db.Statement.SkipHooks {
		field := db.Statement.Schema.LookUpField("CreatedAt")
		if field != nil && field.FieldType.Kind() == reflect.Int64 {
			err := field.Set(db.Statement.Context, reflect.Indirect(db.Statement.ReflectValue), cartime.NowToInt())
			if err != nil {
				db.AddError(err)
			}
		} else {
			callbacks.BeforeCreate(db)
		}
	}
}

// UpdateCallback 自定义更新时间标记
func UpdateCallback(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil && !db.Statement.SkipHooks {
		field := db.Statement.Schema.LookUpField("UpdatedAt")
		if field != nil && field.FieldType.Kind() == reflect.Int64 {
			field.AutoUpdateTime = 0
			err := field.Set(db.Statement.Context, reflect.Indirect(db.Statement.ReflectValue), cartime.NowToInt())
			if err != nil {
				db.AddError(err)
			}
			field := db.Statement.Schema.LookUpField("DeletedAt")
			if field != nil && field.FieldType.Kind() == reflect.Int64 {
				if conds := db.Statement.BuildCondition("`" + db.Statement.Table + "`" + ".`deleted_at`=0"); len(conds) > 0 {
					db.Statement.AddClause(clause.Where{Exprs: conds})
				}
			}
		} else {
			callbacks.BeforeUpdate(db)
		}
	}
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

// DeleteCallback 自定义删除时间标记
func DeleteCallback(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil && !db.Statement.SkipHooks {
		deleteField := db.Statement.Schema.LookUpField("DeletedAt")
		if !db.Statement.Unscoped && deleteField != nil && deleteField.FieldType.Kind() == reflect.Int64 {
			db.Statement.SQL.Grow(100)
			//Soft Delete
			if db.Statement.SQL.String() == "" {
				db.Statement.AddClause(
					clause.Set{{
						Column: clause.Column{Name: deleteField.DBName},
						Value:  cartime.NowToInt(),
					}},
				)
				if conds := db.Statement.BuildCondition("`" + db.Statement.Table + "`" + ".`deleted_at`=0"); len(conds) > 0 {
					db.Statement.AddClause(clause.Where{Exprs: conds})
				}
				db.Statement.AddClauseIfNotExists(clause.Update{})
				db.Statement.BuildClauses = []string{"UPDATE", "SET", "WHERE"}
			}
		} else {
			callbacks.BeforeDelete(db)
		}
	}
}
