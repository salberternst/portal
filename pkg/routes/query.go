package routes

type QueryParams struct {
	Page      int    `form:"page" binding:"required"`
	PageSize  int    `form:"page_size"  binding:"required"`
	SortField string `form:"field"`
	SortOrder string `form:"order"`
}
