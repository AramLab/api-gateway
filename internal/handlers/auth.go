package handlers

import (
	"context"
	"github.com/AramLab/api-gateway/internal/dto"
	"github.com/AramLab/api-gateway/internal/validator"
	pb "github.com/AramLab/protos/gen/go/auth"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type AuthClient interface {
	Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)
	Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Close() error
}

type AuthHandler struct {
	grpcClient AuthClient
	logger     *zap.SugaredLogger
}

func NewAuthHandler(grpcClient AuthClient, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{grpcClient: grpcClient, logger: logger}
}

func (a *AuthHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.UserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Вызов gRPC
		grpcReq := &pb.RegisterRequest{
			Username:  req.Username,
			Password:  req.Password,
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		resp, err := a.grpcClient.Register(c.Context(), grpcReq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": st.Message(),
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})
		}

		return c.JSON(fiber.Map{"user_id": resp.UserID})
	}
}

func (a *AuthHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.LoginRequestDTO
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		validate := validator.New()
		if err := validate.StructPartial(req, "Username", "Password"); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Вызов gRPC
		grpcReq := &pb.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		}

		resp, err := a.grpcClient.Login(c.Context(), grpcReq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": st.Message(),
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})
		}

		return c.JSON(fiber.Map{
			"token": resp.Token,
		})
	}
}
