package validator

import (
	"strconv"
	"github.com/syronz/infrastructure/server/utils/debug"
)

const PAGE_SIZE int = 100

func PaginationPageSize(v string) int{
	perPage, err := strconv.Atoi(v)
	if err != nil {
		debug.Log(err)
		perPage = PAGE_SIZE
	}

	if perPage < 0 {
		perPage = PAGE_SIZE
	}

	return perPage
}

func PaginationPageNumber(v string) int{
	page, err := strconv.Atoi(v)
	if err != nil {
		debug.Log(err)
		page = 1
	}

	if page < 1 {
		page = 1
	}

	return page
}
