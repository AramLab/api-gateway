package handlers

import (
	"github.com/AramLab/api-gateway/internal/dto"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TaskHandler struct {
	log     *zap.SugaredLogger
	client  *resty.Client
	baseURL string
}

func NewTaskHandler(logger *zap.SugaredLogger, baseURL string) *TaskHandler {
	return &TaskHandler{
		log:     logger,
		client:  resty.New(),
		baseURL: baseURL,
	}
}

func (h *TaskHandler) GetTasks() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var tasks []dto.TaskResponse

		resp, err := h.client.R().
			SetResult(&tasks).
			Get(h.baseURL + "/tasks")
		if err != nil {
			h.log.Error("Error fetching tasks", zap.Error(err))
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to fetch tasks",
			})
		}

		if resp.StatusCode() >= 400 {
			h.log.Errorw("Unexpected status from tasks API", "status", resp.StatusCode())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "tasks service returned error",
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(dto.Response{
			Status: "success",
			Data:   tasks,
		})
	}
}

func (h *TaskHandler) CreateTask() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req dto.TaskRequest
		var resp dto.TaskResponse

		if err := ctx.BodyParser(&req); err != nil {
			h.log.Error("Error parsing request body", zap.Error(err))
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		respResty, err := h.client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(req).
			SetResult(&resp).
			Post(h.baseURL + "/tasks")

		if err != nil {
			h.log.Error("Error creating task", zap.Error(err))
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create task",
			})
		}

		if respResty.StatusCode() >= 400 {
			h.log.Errorw("Error response from tasks API", "status", respResty.StatusCode())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "tasks service returned error",
			})
		}

		return ctx.Status(fiber.StatusCreated).JSON(dto.Response{
			Status: "success",
			Data:   resp,
		})
	}
}
