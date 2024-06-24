package server

import komputer "github.com/wittano/komputer/api/proto"

func paginationOrDefault(p *komputer.Pagination) *komputer.Pagination {
	if p == nil {
		return &komputer.Pagination{
			Size: 10,
		}
	}

	return p
}
