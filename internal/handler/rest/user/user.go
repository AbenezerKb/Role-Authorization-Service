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
// @Success      200  boolean true "successfully register the permission"
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
