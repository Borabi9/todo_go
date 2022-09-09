package util

func GetPageInfo(page, navLen int, total, limit int64) PageInfo {
	totalPage := total / limit
	if total%limit != 0 {
		totalPage += 1
	}

	firstPage := ((page / navLen) * navLen)
	lastPage := firstPage
	if firstPage+navLen < int(totalPage) {
		lastPage += navLen
	} else {
		lastPage += int(totalPage) % navLen
	}

	pageSlice := []Page{}
	for i := firstPage; i < lastPage; i++ {
		pageSlice = append(pageSlice, Page{i + 1, (i + 1) == page})
	}

	if len(pageSlice) == 0 {
		return PageInfo{}
	}

	return PageInfo{
		PageSlice: pageSlice,
		TotalPage: totalPage,
		NavLen:    navLen,
		FirstPage: pageSlice[0].PageNum,
		LastPage:  pageSlice[len(pageSlice)-1].PageNum,
		Previous:  pageSlice[0].PageNum - navLen,
		Next:      pageSlice[len(pageSlice)-1].PageNum + 1,
	}
}

type PageInfo struct {
	PageSlice []Page
	TotalPage int64
	NavLen    int
	FirstPage int
	LastPage  int
	Previous  int
	Next      int
}

type Page struct {
	PageNum    int
	IsSelected bool
}
