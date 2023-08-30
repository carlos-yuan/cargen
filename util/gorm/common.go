package gormutil

type Paging struct {
	Page  int64 `json:"page"`  //页数
	Size  int64 `json:"size"`  //单页大小
	Order int32 `json:"order"` //排序
}

func (p *Paging) Offset() int {
	return int((p.Page - 1) * p.Size)
}
