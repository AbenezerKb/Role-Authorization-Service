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
// @param 		 roleid path string true "role id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200  {object} dto.Role "successfully assigned role"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /roles/{roleid}/users/{userid} [post]
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

	role.RoleID, err = uuid.Parse(ctx.Param("roleid"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("role id", ctx.Param("roleid")))
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
