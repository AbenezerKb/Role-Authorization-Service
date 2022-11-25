package permission

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

type permission struct {
	logger           logger.Logger
	permissionModule module.Permission
}

func Init(logger logger.Logger, permissionModule module.Permission) rest.Permission {
	return &permission{
		logger:           logger,
		permissionModule: permissionModule,
	}
}

// CreatePermission is used to register new permissions.
// @Summary      register a new permission.
// @Description  this function registers the service if it does already exist.
// @Description  if the process finishes with out any error it returns true.
// @Description  if the process finishes with any error it returns false.
// @Tags         permissions
// @Accept       json
// @Produce      json
// @param 		 creatnewpermission body dto.CreatePermission true "register permission request body"
// @Success      201  boolean true "successfully register the permission"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized service"
// @Failure      403  {object}  model.ErrorResponse "service is not active"
// @Router       /permissions [post]
// @security 	 BasicAuth
func (p *permission) CreatePermission(ctx *gin.Context) {
	permission := dto.CreatePermission{}
	err := ctx.ShouldBind(&permission)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.logger.Info(ctx, "couldn't bind to dto.Service body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	if err := p.permissionModule.CreatePermission(ctx, permission); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusCreated, nil, nil)
}

// ListPermissions is used to get the list of permissions under the tenant.
// @Summary      returns a list of permission.
// @Description  this function return a list of permissions that are under my domin.
// @Tags         permissions
// @Accept       json
// @Produce      json
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  {object} []dto.Permission
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /permissions [get]
// @security 	 BasicAuth
func (p *permission) ListPermissions(ctx *gin.Context) {
	result, err := p.permissionModule.ListPermissions(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, result, nil)
}

// CreatePermissionDependency is used to create dependency between permissions.
// @Summary      create permission dependency.
// @Description  This function creates a dependency between permission.
// @Description  If a permission inherites a permission and a user is granted the hair permission then the user implicitly granted the inherited permission
// @Tags         permissions
// @Accept       json
// @Produce      json
// @param 		 creatpermissiondependency body []dto.CreatePermissionDependency true "create permission dependency request body"
// @Success      200  boolean true "successfully created dependency"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /permissions/inherit [post]
// @security 	 BasicAuth
func (p *permission) CreatePermissionDependency(ctx *gin.Context) {
	createPermissionInheritance := []dto.CreatePermissionDependency{}
	err := ctx.ShouldBind(&createPermissionInheritance)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.logger.Info(ctx, "couldn't bind to CreatePermissionInherit  body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	if err := p.permissionModule.CreatePermissionDependency(ctx, createPermissionInheritance); err != nil {
		_ = ctx.Error(err)
		return
	}
	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// DeletePermission is used to delete permissions.
// @Summary      deletes permission.
// @Description  This function deletes the permissions if it does already exist.
// @Tags         permissions
// @Accept       json
// @Produce      json
// @param 		 id path string true "permission id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  boolean true "permission successfully deleted"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized service"
// @Failure      403  {object}  model.ErrorResponse "service is not active"
// @Router       /permissions/{id} [delete]
// @security 	 BasicAuth
func (p *permission) DeletePermission(ctx *gin.Context) {

	if err := p.permissionModule.DeletePermission(ctx, ctx.Param("id")); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// GetPermission  returns a permission with the given id
// @Summary      returns a permission with the given id
// @Tags         permissions
// @Accept       json
// @Produce      json
// @param 		 id 	path string true "permission id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200 {object} dto.Permission "successfully returns a permission detail"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Failure      404  {object}  model.ErrorResponse "role not found"
// @Router       /permissions/{id} [get]
// @Security	 BasicAuth
func (p *permission) GetPermission(ctx *gin.Context) {
	permissionId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid permission id")
		p.logger.Info(ctx, "invalid permission id", zap.Error(err), zap.Any("permission id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}

	result, err := p.permissionModule.GetPermission(ctx, permissionId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, result, nil)
}
