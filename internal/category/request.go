package category

type createCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type updateCategoryRequest struct {
	Name        string `json:"name,omitempty" binding:"omitempty,min=1"`
	Description string `json:"description,omitempty" binding:"omitempty,min=1"`
}
