package user

import (
	"2f-authorization/internal/constants"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/handler/rest"
	"2f-authorization/internal/module"
	"2f-authorization/platform/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type user struct {
	logger     logger.Logger
	userModule module.User
}

func Init(logger logger.Logger, userModule module.User) rest.User {
	return &user{
		logger:     logger,
		userModule: userModule,
	}
}

// RegisterUser is used to register user.
// @Summary      add new user to the system.
// @Description  This function registers new user if the user doesn't exist.
// @Description  If the process finishes with out any error it returns true.
// @Tags         users
// @Accept       json
// @Produce      json
// @param 		 registeruser body dto.RegisterUser true "Register user request body"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  boolean true "successfully register the user"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /users [post]
// @security 	 BasicAuth
func (u *user) RegisterUser(ctx *gin.Context) {
	user := dto.RegisterUser{}
	if err := ctx.ShouldBind(&user); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.logger.Info(ctx, "couldn't bind to dto.User body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	if err := u.userModule.RegisterUser(ctx, user); err != nil {
		_ = ctx.Error(err)
		return
	}
	constants.SuccessResponse(ctx, http.StatusCreated, nil, nil)
}

// UpdateUserStatus updates user status
// @Summary      changes user status
// @Tags         users
// @Accept       json
// @Produce      json
// @param status body dto.UpdateUserStatus true "status"
// @Success      200 boolean true "successfully updates the user status"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /users/status [patch]
// @Security	 BasicAuth
func (u *user) UpdateUserStatus(ctx *gin.Context) {
	updateStatusParam := dto.UpdateUserStatus{}
	err := ctx.ShouldBind(&updateStatusParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.logger.Info(ctx, "unable to bind user status", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	if err := u.userModule.UpdateUserStatus(ctx, updateStatusParam); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// GetPermissionWithInTenant returns user's permissions within the specified tenant
// @Summary      returns user's permissions
// @Tags         users
// @Accept       json
// @Produce      json
// @param 		 id 	path string true "user id"
// @param 		 tenant-id 	path string true "tenant id"
// @Success      200  {object} []dto.Permission "return permissions list"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Router       /users/{id}/tenants/{tenant-id}/permissions [get]
// @Security	 BasicAuth
func (u *user) GetPermissionWithInTenant(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid role id")
		u.logger.Info(ctx, "invalid role id", zap.Error(err), zap.Any("role id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}
	tenant := ctx.Param("tenant-id")

	permission, err := u.userModule.GetPermissionWithInTenant(ctx, tenant, userId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, permission, nil)
}

// UpdateUserRoleStatus updates user's role status
// @Summary      changes user's role status
// @Tags         users
// @Accept       json
// @Produce      json
// @param status body dto.UpdateUserRoleStatus true "status"
// @param 		 role-id 	path string true "role id"
// @param 		 id 	path string true "user id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200 boolean true "successfully updates the user's role status"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /users/{id}/roles/{role-id}/status [patch]
// @Security	 BasicAuth
func (u *user) UpdateUserRoleStatus(ctx *gin.Context) {
	updateStatusParam := dto.UpdateUserRoleStatus{}
	err := ctx.ShouldBind(&updateStatusParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.logger.Info(ctx, "unable to bind role status", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	roleId, err := uuid.Parse(ctx.Param("role-id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid role id")
		u.logger.Info(ctx, "invalid role id", zap.Error(err), zap.Any("role id", ctx.Param("role-id")))
		_ = ctx.Error(err)
		return
	}
	userId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid user id")
		u.logger.Info(ctx, "invalid role id", zap.Error(err), zap.Any("user id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}

	if err := u.userModule.UpdateUserRoleStatus(ctx, updateStatusParam, roleId, userId); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

func (u *user) GetPermissionWithInDomain(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid role id")
		u.logger.Info(ctx, "invalid role id", zap.Error(err), zap.Any("role id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}
	domain := ctx.Param("domain-id")

	permission, err := u.userModule.GetPermissionWithInDomain(ctx, domain, userId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, permission, nil)

}
