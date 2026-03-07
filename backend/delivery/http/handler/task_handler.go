package handler

import (
	"net/http"
	"strconv"
	"taskmanager/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskUseCase usecase.TaskUseCase
}

func NewTaskHandler(taskUseCase usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{taskUseCase: taskUseCase}
}

type createTaskReq struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline" binding:"required"`
	Checkpoints []string  `json:"checkpoints" binding:"required,len=4"`
	CategoryID  *int      `json:"category_id"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req createTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskUseCase.CreateTask(c.Request.Context(), userID, req.Title, req.Description, req.Deadline, req.Checkpoints, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID := c.GetInt("user_id")

	tasks, err := h.taskUseCase.GetTasks(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

type completeCheckpointReq struct {
	Completed bool `json:"completed"`
}

func (h *TaskHandler) CompleteCheckpoint(c *gin.Context) {
	userID := c.GetInt("user_id")

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	cpID, err := strconv.Atoi(c.Param("cpId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkpoint id"})
		return
	}

	var req completeCheckpointReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.taskUseCase.CompleteCheckpoint(c.Request.Context(), userID, taskID, cpID, req.Completed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "checkpoint updated successfully"})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID := c.GetInt("user_id")

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	if err := h.taskUseCase.DeleteTask(c.Request.Context(), userID, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

type updateTaskReq struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline" binding:"required"`
	Checkpoints []string  `json:"checkpoints" binding:"required,len=4"`
	CategoryID  *int      `json:"category_id"`
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID := c.GetInt("user_id")

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req updateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.taskUseCase.UpdateTask(c.Request.Context(), userID, taskID, req.Title, req.Description, req.Deadline, req.Checkpoints, req.CategoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}
