package gen

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func GormGen(path, dsn, name string, tables []string) {
	g := gen.NewGenerator(gen.Config{
		OutPath:       path + "/orm/" + name + "/query",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		FieldNullable: true,
		ModelPkgPath:  path + "/orm/" + name + "/model",
	})
	g.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
		return ToCamelFirstLowerCase(columnName)
	})
	gormdb, _ := gorm.Open(mysql.Open(dsn))
	g.UseDB(gormdb) // reuse your gorm db

	// Generate basic type-safe API for struct `model.User` following conventions
	//g.ApplyBasic(model.User{})

	// Generate Order Safe API with hand-optimized SQL defined on Querier interface for `model.User` and `model.Company`
	//g.ApplyInterface(func(Querier){}, model.User{})
	var tablesObj []interface{}
	if len(tables) > 0 {
		for _, table := range tables {
			tablesObj = append(tablesObj, g.GenerateModel(table))
		}
	} else {
		tablesObj = g.GenerateAllTable()
	}
	//g.ApplyBasic(g.GenerateAllTable()...)
	g.ApplyBasic(tablesObj...)

	// Generate the code
	g.Execute()
}
