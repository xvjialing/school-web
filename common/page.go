package common

type Page struct {
	Total    int64         `json:"total"`
	List     []interface{} `json:"list"`
	PageNum  int64         `json:"pageNum"`
	PageSize int64         `json:"pageSize"`
	//Size int64 `json:"size"`
	//StartRow int `json:"startRow"`
	//EndRow int `json:"endRow"`
	Pages int `json:"pages"`
	//	"prePage":0
	//	"nextPage":0
	//	"isFirstPage":true
	//	"isLastPage":true
	//	"hasPreviousPage":false
	//	"hasNextPage":false
	//	"navigatePages":3
	//	"navigatepageNums":[
	//	1
	//]
	//	"navigateFirstPage":1
	//	"navigateLastPage":1
	//	"firstPage":1
	//	"lastPage":1
}

//func tansform(data []interface{},pageNum ,pageSize int) (page Page) {
//	page.PageNum=pageNum
//	page.PageSize=pageSize
//	size := len(data)
//	page.Total =size
//	if size ==0 {
//		page.Pages=0
//		return page
//	}
//	if pageSize <= 0 || pageNum <= 0{
//		page.Pages=0
//		return page
//	}
//
//	totalPage := int(math.Ceil(float64(size) / float64(pageSize)))
//	page.Pages=totalPage
//
//
//}
