package gen

import (
	"fmt"
	"os"

	"github.com/carlos-yuan/cargen/util/convert"
	"github.com/carlos-yuan/cargen/util/fileUtil"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func GormGen(path, dsn, name string, tables []string) {
	outPath := fileUtil.FixPathSeparator(path + "/orm/" + name + "/query")
	modelPkgPath := fileUtil.FixPathSeparator(path + "/orm/" + name + "/model")
	g := gen.NewGenerator(gen.Config{
		OutPath:       outPath,
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		FieldNullable: true,
		ModelPkgPath:  modelPkgPath,
	})
	g.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
		return convert.ToCamelFirstLowerCase(columnName)
	})
	gormdb, _ := gorm.Open(mysql.Open(dsn))
	g.UseDB(gormdb) // reuse your gorm db

	// Generate basic type-safe API for struct `model.User` following conventions
	//g.ApplyBasic(model.User{})

	// Generate Order Safe API with hand-optimized SQL defined on Querier interface for `model.User` and `model.Company`
	//g.ApplyInterface(func(Querier){}, model.User{})
	var tablesObj []interface{}
	var tablesInfo []TableInfo
	if len(tables) > 0 {
		for _, table := range tables {
			t := g.GenerateModel(table)
			tablesInfo = append(tablesInfo, TableInfo{
				TableName:       t.TableName,
				FileName:        t.FileName,
				QueryStructName: t.QueryStructName,
				ModelStructName: t.ModelStructName,
				S:               t.S,
				Generated:       t.Generated,
			})
			tablesObj = append(tablesObj, t)
		}
	} else {
		tablesObj = g.GenerateAllTable()
	}
	//g.ApplyBasic(g.GenerateAllTable()...)
	g.ApplyBasic(tablesObj...)

	// Generate the code
	g.Execute()
	generateDaoFile(outPath, tablesInfo)
}

type TableInfo struct {
	Generated       bool   // whether to generate db model
	FileName        string // generated file name
	S               string // the first letter(lower case)of simple Name (receiver)
	QueryStructName string // internal query struct name
	ModelStructName string // origin/model struct name
	TableName       string // table name in db server
}

func generateDaoFile(outPath string, tablesInf []TableInfo) {
	for _, table := range tablesInf {
		filePath := outPath + string(os.PathSeparator) + table.FileName + ".go"
		if !fileUtil.IsExist(filePath) {
			code := fmt.Sprintf(databaseDaoTemplate,
				table.ModelStructName,
				table.ModelStructName,
				table.ModelStructName,
				table.ModelStructName,
				table.ModelStructName,
				table.ModelStructName,
				table.ModelStructName,
				table.ModelStructName,
				table.ModelStructName,
				table.ModelStructName,
				table.ModelStructName,
			)
			err := fileUtil.WriteByteFile(filePath, []byte(code))
			if err != nil {
				panic(err)
			}
		}
	}
}

const databaseDaoTemplate = `
package query

import (
	"context"
	"gorm.io/gorm"
)

type %sDao struct {
	I%sDo
}

func (q *Query) Get%sDao(ctx context.Context) %sDao {
	return %sDao{q.%s.WithContext(ctx)}
}

func (q *Query) Get%sDaoWithTx(ctx context.Context, tx *gorm.DB) %sDao {
	dao := %sDao{q.%s.WithContext(ctx)}
	dao.I%sDo.ReplaceDB(tx)
	return dao
}
`
