package test

import (
	"github.com/carlos-yuan/cargen/cmd/gen"
	"testing"
)

func TestGenDatabase(t *testing.T) {
	gen.Config{Gen: gen.GenDB, DbName: "enterprise",
		DbDsn:  `admin:admin123@tcp(192.168.0.88:3306)/ai_hazardous_chemical_enterprise?charset=utf8&parseTime=True&loc=Local&timeout=1000ms`,
		Tables: `application_file,application_outline,hce_enterprise_info,hce_operation_permit,storage_area,storage_equipment,storage_facility,hazardous_chemicals_catalog,organize,hce_doc,region,dict_data,dict_type,hce_user,hce_enterprise_register,chemicals_catalog_detail,waybill,waybill_goods`,
		Path:   "D:\\carlos\\hc_enterprise_server",
	}.Build()

}
