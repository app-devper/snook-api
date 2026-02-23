package request

type MenuCategory struct {
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}
