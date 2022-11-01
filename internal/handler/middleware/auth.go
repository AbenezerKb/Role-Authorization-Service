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
	"strings"

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
			err := errors.ErrInternalServerError.Wrap(nil, "could not extract service credentials")
			a.logger.Error(ctx, "extract error", zap.Error(err))
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		serviceId, err := uuid.Parse(Id)
		if err != nil {
			err := errors.ErrInvalidUserInput.Wrap(err, "invalid service id")
			a.logger.Error(ctx, "parse error", zap.Error(err), zap.Any("service-id", Id))
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		service, err := a.service.GetServiceById(ctx, dto.Service{ID: serviceId})
		if err != nil {
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		switch service.Status {
		case constants.InActive:
			Err := errors.ErrAuthError.New("Your service is not active, Please consult the system administrator to activate your service")
			a.logger.Warn(ctx, "service status is not active", zap.String("service-id", service.ID.String()))
			_ = ctx.Error(Err)
			ctx.Abort()
			return
		case constants.Pending:
			Err := errors.ErrAuthError.New("Your service is on pending state, Please consult the system administrator to activate your service")
			a.logger.Warn(ctx, "service status is pending", zap.String("service-id", service.ID.String()))
			_ = ctx.Error(Err)
			ctx.Abort()
			return
		}

		ok, err = argon.ComparePasswordAndHash(secret, service.Password)
		if err != nil {
			err = errors.ErrAcessError.Wrap(err, "unauthorized_service")
			a.logger.Warn(ctx, "unauthorized_service", zap.Error(err), zap.String("service-id", service.ID.String()), zap.String("provided-password", secret))
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		if !ok {
			err = errors.ErrAcessError.Wrap(nil, "unauthorized_service")
			a.logger.Warn(ctx, "unauthorized_service", zap.Error(err), zap.String("service-id", service.ID.String()), zap.String("provided-password", secret))
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set("x-service-id", service.ID.String())
		ctx.Next()
	}
}

func (a *authMiddeleware) Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.Request{
			Tenant:   ctx.GetHeader("x-tenant"),
			Subject:  ctx.GetHeader("x-subject"),
			Action:   ctx.GetHeader("x-action"),
			Resource: ctx.GetHeader("x-resource"),
			Service:  ctx.GetString("x-service-id"),
			Fields:   strings.Split(ctx.GetHeader("x-fields"), ","),
		}

		if len(req.Fields) == 0 {
			req.Fields = []string{"*"}
		}

		if err := req.Validate(); err != nil {
			err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
			a.logger.Info(ctx, "invalid input", zap.Error(err))
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

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
		ctx.Set("x-tenant", ctx.GetHeader("x-tenant"))
		ctx.Next()
	}
}
