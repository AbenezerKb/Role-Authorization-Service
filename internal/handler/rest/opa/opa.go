package opa

import (
	"2f-authorization/internal/constants"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/handler/rest"
	"2f-authorization/internal/module"
	"2f-authorization/platform/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type opa struct {
	logger    logger.Logger
	opamodule module.Opa
}

func Init(logger logger.Logger, opamodule module.Opa) rest.Opa {
	return &opa{
		logger:    logger,
		opamodule: opamodule,
	}
}

// Authorize is used to check whether the user is authorized or not.
// @Summary      authorize user.
// @Description  This function check whether the user is authorized or not to perform the action on the resource within the given tenant and service.
// @Tags         authorize
// @Accept       json
// @Produce      json
// @param 		 authorize user body model.Request true "authorization request body"
// @Success      200 {object} boolean "successfully authorize the user"
// @Failure      400  {object}  model.ErrorResponse "required field error,bad request error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /authorize [post]
// @security 	 BasicAuth
func (o *opa) Authorize(ctx *gin.Context) {
	req := model.Request{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "couldn't bind to dto.CreateRole body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	var ok bool
	ok, err = o.opamodule.Authorize(ctx, req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, ok, nil)
}
