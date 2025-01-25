package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/internal/interceptor"
	"github.com/welovemedia/ffmate/internal/service"
	"github.com/welovemedia/ffmate/sev"
	"github.com/welovemedia/ffmate/sev/exceptions"
)

type TaskController struct {
	sev.Controller
	sev    *sev.Sev
	Prefix string
}

func (c *TaskController) Setup(s *sev.Sev) {
	c.sev = s
	s.Gin().GET(c.Prefix+c.getEndpoint(), interceptor.PageLimit, c.listTasks)
	s.Gin().POST(c.Prefix+c.getEndpoint(), c.addTask)
	s.Gin().POST(c.Prefix+c.getEndpoint()+"/batch", c.addTasks)
	s.Gin().GET(c.Prefix+c.getEndpoint()+"/:uuid", c.getTask)
	s.Gin().GET(c.Prefix+c.getEndpoint()+"/batch/:uuid", interceptor.PageLimit, c.getTasks)
	s.Gin().DELETE(c.Prefix+c.getEndpoint()+"/:uuid", c.deleteTask)
	s.Gin().PATCH(c.Prefix+c.getEndpoint()+"/:uuid/cancel", c.cancelTask)
}

// @Summary List all tasks
// @Description List all existing tasks
// @Tags tasks
// @Param page query int false "the page of a pagination request (min 0)"
// @Param perPage query int false "the amount of results of a pagination request (min 1; max: 100)"
// @Produce json
// @Success 200 {object} []dto.Task
// @Router /tasks [get]
func (c *TaskController) listTasks(gin *gin.Context) {
	tasks, total, err := service.TaskService().ListTasks(gin.GetInt("page"), gin.GetInt("perPage"))
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	gin.Header("X-Total", fmt.Sprintf("%d", total))

	// Transform each task to its DTO
	var taskDTOs = []dto.Task{}
	for _, task := range *tasks {
		taskDTOs = append(taskDTOs, *task.ToDto())
	}

	gin.JSON(200, taskDTOs)
}

// @Summary Delete a task
// @Description Delete a task by its uuid
// @Tags tasks
// @Param uuid path string true "the tasks uuid"
// @Produce json
// @Success 204
// @Router /tasks/{uuid} [delete]
func (c *TaskController) deleteTask(gin *gin.Context) {
	uuid := gin.Param("uuid")
	err := service.TaskService().DeleteTask(uuid)

	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	gin.AbortWithStatus(204)
}

// @Summary Add a batch of tasks
// @Description	Add a batch of new tasks to the queue
// @Tags tasks
// @Accept json
// @Param request body []dto.NewTask true "new tasks"
// @Produce json
// @Success 200 {object} []dto.Task
// @Router /tasks/batch [post]
func (c *TaskController) addTasks(gin *gin.Context) {
	newTasks := &[]dto.NewTask{}
	c.sev.Validate().BindWithoutValidation(gin, newTasks)

	// bind and validation in a single step throws a nil error, so we separate those tasks
	for _, t := range *newTasks {
		c.sev.Validate().ValidateOnly(gin, &t)
	}

	tasks, err := service.TaskService().NewTasks(newTasks)
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	// Transform each task to its DTO
	var taskDTOs = []dto.Task{}
	for _, task := range *tasks {
		taskDTOs = append(taskDTOs, *task.ToDto())
	}

	gin.JSON(200, taskDTOs)
}

// @Summary Add a new task
// @Description	Add a new tasks to the queue
// @Tags tasks
// @Accept json
// @Param request body dto.NewTask true "new task"
// @Produce json
// @Success 200 {object} dto.Task
// @Router /tasks [post]
func (c *TaskController) addTask(gin *gin.Context) {
	newTask := &dto.NewTask{}
	c.sev.Validate().Bind(gin, newTask)

	task, err := service.TaskService().NewTask(newTask, "", "api")
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	gin.JSON(200, task.ToDto())
}

// @Summary Get single task
// @Description	Get a single task by its uuid
// @Tags tasks
// @Param uuid path string true "the tasks uuid"
// @Produce json
// @Success 200 {object} dto.Task
// @Router /tasks/{uuid} [get]
func (c *TaskController) getTask(gin *gin.Context) {
	uuid := gin.Param("uuid")
	task, err := service.TaskService().GetTaskById(uuid)
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	gin.JSON(200, task.ToDto())
}

// @Summary Get tasks for batch
// @Description	Get tasks by batch uuid
// @Tags tasks
// @Param uuid path string true "the batch uuid"
// @Produce json
// @Success 200 {object} []dto.Task
// @Router /tasks/batch/{uuid} [get]
func (c *TaskController) getTasks(gin *gin.Context) {
	uuid := gin.Param("uuid")
	tasks, total, err := service.TaskService().GetTasksByBatchId(uuid, gin.GetInt("page"), gin.GetInt("perPage"))
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	gin.Header("X-Total", fmt.Sprintf("%d", total))

	// Transform each task to its DTO
	var taskDTOs = []dto.Task{}
	for _, task := range *tasks {
		taskDTOs = append(taskDTOs, *task.ToDto())
	}

	gin.JSON(200, taskDTOs)
}

// @Summary Cancel a task
// @Description Cancel a task by its uuid
// @Tags tasks
// @Param uuid path string true "the tasks uuid"
// @Produce json
// @Success 200 {object} dto.Task
// @Router /tasks/{uuid}/cancel [patch]
func (c *TaskController) cancelTask(gin *gin.Context) {
	uuid := gin.Param("uuid")
	task, err := service.TaskService().CancelTask(uuid)
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	gin.JSON(200, task.ToDto())
}

func (c *TaskController) GetName() string {
	return "task"
}

func (c *TaskController) getEndpoint() string {
	return "/v1/tasks"
}
