package role

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
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
	"go.uber.org/zap"
)

type role struct {
	logger     logger.Logger
	roleModule module.Role
}

func Init(logger logger.Logger, roleModule module.Role) rest.Role {
	return &role{
		logger:     logger,
		roleModule: roleModule,
	}
}

// CreateRole is used to create new role.
// @Summary      add new role.
// @Description  This function creates new role if the role doesn't exist.
// @Tags         roles
// @Accept       json
// @Produce      json
// @param 		 createrole body dto.CreateRole true "create role request body"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  {object} dto.Role "successfully creates the role"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /roles [post]
// @security 	 BasicAuth
func (r *role) CreateRole(ctx *gin.Context) {
	role := dto.CreateRole{}
	err := ctx.ShouldBind(&role)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "couldn't bind to dto.CreateRole body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	createdRole, err := r.roleModule.CreateRole(ctx, role)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusCreated, createdRole, nil)
}

// AssignRole is used to create new role.
// @Summary      assign role to a user.
// @Description  This function assign new role if the role  dosen't assigned.
// @Tags         roles
// @Accept       json
// @Produce      json
// @param 		 userid path string true "user id"
// @param 		 id path string true "role id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  boolean 	true "successfully assigned role"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /roles/{id}/users/{userid} [post]
// @security 	 BasicAuth
func (r *role) AssignRole(ctx *gin.Context) {
	var err error
	role := dto.TenantUsersRole{}

	role.UserID, err = uuid.Parse(ctx.Param("userid"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("user id", ctx.Param("userid")))
		_ = ctx.Error(err)
		return
	}

	role.RoleID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("role id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}

	err = r.roleModule.AssignRole(ctx, role)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// RevokeRole is used to revoke user role.
// @Summary      revoke user role.
// @Description  This function revoke user's role if it is given.
// @Tags         roles
// @Accept       json
// @Produce      json
// @param 		 userid path string true "user id"
// @param 		 id path string true "role id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  boolean true "successfully assigned role"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /roles/{id}/users/{userid} [patch]
// @security 	 BasicAuth
func (r *role) RevokeRole(ctx *gin.Context) {
	var err error
	role := dto.TenantUsersRole{}

	role.UserID, err = uuid.Parse(ctx.Param("userid"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("user id", ctx.Param("userid")))
		_ = ctx.Error(err)
		return
	}

	role.RoleID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("role id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}

	err = r.roleModule.RevokeRole(ctx, role)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// UpdateRole is used to update the existing role.
// @Summary      update role.
// @Description  This function updates the given role.
// @Tags         roles
// @Accept       json
// @Produce      json
// @param 		 role id path string true "role id"
// @param 		 updaterolepermissionslist body dto.UpdateRole true "update role request body"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  boolean true "successfully updated role"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /roles/{id} [put]
// @security 	 BasicAuth
func (r *role) UpdateRole(ctx *gin.Context) {
	role := dto.UpdateRole{}
	err := ctx.ShouldBind(&role)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "couldn't bind to dto.UpdateRole body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	role.RoleID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("role id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}

	if err := r.roleModule.UpdateRole(ctx, role); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// DeleteRole is used to delete the existing role.
// @Summary      delete role.
// @Description  This function deletes the given role.
// @Tags         roles
// @Accept       json
// @Produce      json
// @param 		 role id path string true "role id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  {object}  dto.Role 			"successfully deleted role"
// @Failure      400  {object}  model.ErrorResponse "invalid input error"
// @Failure      404  {object}  model.ErrorResponse "role not found"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /roles/{id} [delete]
// @security 	 BasicAuth
func (r *role) DeleteRole(ctx *gin.Context) {
	roleId := ctx.Param("id")

	role, err := r.roleModule.DeleteRole(ctx, roleId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constants.SuccessResponse(ctx, http.StatusOK, role, nil)
}

// ListRoles is used to get the list of roles under the tenant.
// @Summary      returns a list of roles.
// @Description  this function return a list of roles.
// @Tags         roles
// @Accept       json
// @Produce      json
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  {object} []dto.Role
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /roles [get]
// @security 	 BasicAuth
func (r *role) ListRoles(ctx *gin.Context) {
	param := db_pgnflt.PgnFltQueryParams{}

	err := ctx.BindQuery(&param)
	if err != nil {
		er := errors.ErrInvalidUserInput.Wrap(err, "invalid user input")
		r.logger.Info(ctx, "unable to bind role query", zap.Error(err))
		_ = ctx.Error(er)
		return
	}

	result, metaData, err := r.roleModule.ListRoles(ctx, param)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, result, metaData)
}

// UpdateRoleStatus updates role status
// @Summary      changes role status
// @Tags         roles
// @Accept       json
// @Produce      json
// @param status body dto.UpdateRoleStatus true "status"
// @param 		 id 	path string true "role id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200 boolean true "successfully updates the role's status"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /roles/{id}/status [patch]
// @Security	 BasicAuth
func (r *role) UpdateRoleStatus(ctx *gin.Context) {
	updateStatusParam := dto.UpdateRoleStatus{}
	err := ctx.ShouldBind(&updateStatusParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "unable to bind role status", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	roleId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid role id")
		r.logger.Info(ctx, "invalid role id", zap.Error(err), zap.Any("role id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}

	if err := r.roleModule.UpdateRoleStatus(ctx, updateStatusParam, roleId); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// GetRole  returns a role with the given id
// @Summary      returns a role with the given id
// @Tags         roles
// @Accept       json
// @Produce      json
// @param 		 id 	path string true "role id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200 {object} dto.Role "successfully returns a role detail"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Failure      404  {object}  model.ErrorResponse "role not found"
// @Router       /roles/{id} [get]
// @Security	 BasicAuth
func (r *role) GetRole(ctx *gin.Context) {
	roleId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid role id")
		r.logger.Info(ctx, "invalid role id", zap.Error(err), zap.Any("role id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}

	result, err := r.roleModule.GetRole(ctx, roleId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, result, nil)
}

// SystemAssignRole is used by the system to assign a role to a user.
// @Summary      assign role to a user.
// @Description  This function assign new role if the role  dosen't assigned.
// @Tags         roles
// @Accept       json
// @Produce      json
// @param 		 userid path string true "user id"
// @param 		 role body dto.RoleTenant true "role"
// @Success      200  boolean 	true 				"successfully assigned role"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /system/users/{userid}/roles [post]
// @security 	 BasicAuth
func (r *role) SystemAssignRole(ctx *gin.Context) {
	roleTenant := dto.RoleTenant{}

	err := ctx.ShouldBind(&roleTenant)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "unable to bind role body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("user id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}

	err = r.roleModule.AssignRole(ctx, dto.TenantUsersRole{
		UserID:     userID,
		RoleTenant: roleTenant,
	})
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}
