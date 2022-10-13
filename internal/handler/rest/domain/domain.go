package domain

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

type domain struct {
	logger       logger.Logger
	domainModule module.Domain
}

func Init(log logger.Logger, domainModule module.Domain) rest.Domain {
	return &domain{
		logger:       log,
		domainModule: domainModule,
	}
}

// CreateDomain is used to create new domain within the service.
// @Summary      create new domain.
// @Description  this function create new domain within the service if not exist.
// @Tags         domain
// @Accept       json
// @Produce      json
// @param 		 createdomain body dto.Domain true "create domain request body"
// @Success      201 {object} dto.Domain "successfully create new domain"
// @Failure      400  {object}  model.ErrorResponse "required field error,bad request error"
// @Router       /domains [post]
func (d *domain) CreateDomain(ctx *gin.Context) {

	domain := dto.Domain{}
	if err := ctx.ShouldBind(&domain); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		d.logger.Info(ctx, "couldn't bind to dto.Domain body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	createdDomain, err := d.domainModule.CreateDomain(ctx, domain)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constants.SuccessResponse(ctx, http.StatusCreated, createdDomain, nil)

}
