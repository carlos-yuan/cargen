package ctl

type Paging struct {
	Page  int64 `form:"page"`
	Limit int64 `form:"limit"`
	Order int32 `form:"order"`
}
