package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shuza/Personalised-ai-agent/internal/domain"
	"github.com/shuza/Personalised-ai-agent/internal/service"
)

type Handler struct {
	users    *service.UsersService
	workflow *service.WorkflowService
}

func NewHandler(users *service.UsersService, workflow *service.WorkflowService) *Handler {
	return &Handler{
		users:    users,
		workflow: workflow,
	}
}

func (h *Handler) Router() http.Handler {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", h.handleHealth)
	router.GET("/v1/users/:userID", h.handleGetUser)
	router.POST("/v1/events", h.handleEvent)

	return router
}

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) handleGetUser(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	user, err := h.users.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) handleEvent(c *gin.Context) {
	var event domain.Event
	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.workflow.HandleEvent(c.Request.Context(), event)

	if err != nil {
		status := http.StatusInternalServerError
		errMsg := err.Error()
		if containsValidationError(errMsg) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusAccepted, result)
}

func containsValidationError(errMsg string) bool {
	return containsAny(errMsg, []string{"required", "not found"})
}

func containsAny(input string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(input, p) {
			return true
		}
	}
	return false
}
