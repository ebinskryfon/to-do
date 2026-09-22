package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"todo/internal/delivery/http/response"
	domainerrors "todo/internal/domain/errors"
	"todo/internal/types"
	"todo/internal/usecase/interfaces"
)

// TodoHandler handles HTTP requests for Todo resources.
type TodoHandler struct {
	createUC interfaces.CreateTodoUsecase
	getUC    interfaces.GetTodoUsecase
	updateUC interfaces.UpdateTodoUsecase
	listUC interfaces.ListTodoUsecase
	log      zerolog.Logger
}

// NewTodoHandler creates a new instance of TodoHandler.
func NewTodoHandler(createUC interfaces.CreateTodoUsecase, getUC interfaces.GetTodoUsecase, updateUC interfaces.UpdateTodoUsecase,listUC interfaces.ListTodoUsecase, log zerolog.Logger) *TodoHandler {
	return &TodoHandler{
		createUC: createUC,
		getUC:    getUC,
		updateUC: updateUC,
		listUC:	  listUC,
		log:      log,
	}
}

// Create handles creating a new Todo item.
// @Summary Create a new todo
// @Description Create a new todo item with title and optional description
// @Tags todos
// @Accept json
// @Produce json
// @Param request body CreateTodoRequest true "Todo creation request payload"
// @Success 201 {object} response.Envelope{data=entity.Todo}
// @Failure 400 {object} response.Envelope{error=response.ErrorBody}
// @Failure 500 {object} response.Envelope{error=response.ErrorBody}
// @Router /api/v1/todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	var req types.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	todo, err := h.createUC.Execute(c.Request.Context(), req.Title, req.Description)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Created(c, todo)
}

func (h *TodoHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerrors.ErrInvalidInput):
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
	case errors.Is(err, domainerrors.ErrTodoNotFound):
		response.Error(c, http.StatusNotFound, "TODO_NOT_FOUND", err.Error())
	default:
		h.log.Error().Err(err).Msg("internal server error")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

// GetByID handles retrieving a specific Todo item by ID.
// @Summary Get a specific todo
// @Description Get a specific todo item by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} response.Envelope{data=entity.Todo}
// @Failure 400 {object} response.Envelope{error=response.ErrorBody}
// @Failure 404 {object} response.Envelope{error=response.ErrorBody}
// @Failure 500 {object} response.Envelope{error=response.ErrorBody}
// @Router /api/v1/todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid todo id")
		return
	}
	todo, err := h.getUC.Execute(c.Request.Context(), parsedID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Success(c, todo)
}

func (h *TodoHandler) Update(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid todo id")
		return
	}

	var req types.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	updateTodo, err := h.updateUC.Execute(c.Request.Context(), parsedID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, updateTodo)

}

func (h *TodoHandler) List(c *gin.Context) {
	var err error
	page := 1
	pageParam := c.Query("page")

	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid page value")
			return
		}
	}

	if page < 1 {
			response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid page number")
			return
		}

	pageSize := 10
	pageSizeParam := c.Query("page_size")
	
	if pageSizeParam != "" {
		pageSize, err = strconv.Atoi(pageSizeParam)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid page size value")
			return
		}
	}

	if pageSize < 1 || pageSize > 100 {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "page size should be in between 1 and 100")
		return
	}

	var cmp bool
	var completed *bool
	completedParam := c.Query("completed")

	if completedParam != "" {
		cmp, err = strconv.ParseBool(completedParam)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "invalid completed value")
			return
		}
		completed = &cmp
	}

	todos, total, err := h.listUC.Execute(c.Request.Context(), page, pageSize, completed)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Paginated(c, todos, total)
}