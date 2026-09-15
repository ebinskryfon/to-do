package types

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required" example:"Buy groceries"`
	Description string `json:"description" example:"Milk, eggs, and bread"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}
