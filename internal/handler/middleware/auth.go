package middleware

import (
	"2f-authorization/internal/constants"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/argon"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthMiddeleware interface {
	BasicAuth() gin.HandlerFunc
	Authorize() gin.HandlerFunc
}
type authMiddeleware struct {
	logger  logger.Logger
	service storage.Service
	opa     opa.Opa
}

func InitAuthMiddleware(logger logger.Logger, service storage.Service, opa opa.Opa) AuthMiddeleware {
	return &authMiddeleware{
		logger:  logger,
		service: service,
		opa:     opa,
	}
}

func (a *authMiddeleware) BasicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Id, secret, ok := ctx.Request.BasicAuth()
		if !ok {
			err := errors.ErrInternalServerError.New("could not get extract service credentials")
			a.logger.Error(ctx, "extract error", zap.Error(err))
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		serviceId, err := uuid.Parse(Id)
		if err != nil {
			err := errors.ErrInvalidUserInput.Wrap(err, "invalid service id")
			a.logger.Error(ctx, "parse error", zap.Error(err), zap.Any("service-id", Id))
			ctx.Error(err)
			ctx.Abort()
			return
		}

		service, err := a.service.GetServiceById(ctx, dto.Service{ID: serviceId})
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		if service.Status != constants.Active {
			Err := errors.ErrAuthError.New("Your service is not active, Please consult the system administrator to activate your service")
			a.logger.Warn(ctx, "service is inactive", zap.String("service-id", service.ID.String()))
			ctx.Error(Err)
			ctx.Abort()
			return
		}

		if ok, _ := argon.ComparePasswordAndHash(secret, service.Password); !ok {
			err = errors.ErrAcessError.New("unauthorized_service")
			a.logger.Warn(ctx, "unauthorized_service", zap.Error(err), zap.String("service-id", service.ID.String()), zap.String("provided-password", secret))
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set("x-service-id", service.ID.String())
		ctx.Next()
	}
}

func (a *authMiddeleware) Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.Request{}
		if err := ctx.ShouldBind(&req); err != nil {
			err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
			a.logger.Info(ctx, "couldn't bind to dto.request body", zap.Error(err))
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		if err := req.Validate(); err != nil {
			err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
			a.logger.Info(ctx, "invalid input", zap.Error(err))
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		req.Service = ctx.GetString("x-service-id")

		ok, err := a.opa.Allow(ctx, req)
		if err != nil {
			err := errors.ErrAcessError.Wrap(err, "unable to perform operation")
			a.logger.Error(ctx, "error while enforcing policy", zap.Error(err), zap.String("service-id", req.Service), zap.String("user-id", req.Subject))
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		if !ok {
			err := errors.ErrAcessError.Wrap(err, "Access denied")
			a.logger.Info(ctx, "access denied", zap.Error(err), zap.String("service-id", req.Service), zap.String("user-id", req.Subject))
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
