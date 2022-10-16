package registeruser

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/platform/argon"
	"2f-authorization/test"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type registerUser struct {
	test.TestInstance
	apiTest        src.ApiTest
	service        dto.CreateService
	createdService dto.CreateServiceResponse
	userId         uuid.UUID
}

func TestCreateTenant(t *testing.T) {
	r := &registerUser{}
	r.TestInstance = test.Initiate(context.Background(), "../../../../")
	r.apiTest.InitializeServer(r.Server)
	r.apiTest.InitializeTest(t, "register user test", "feature/register_user.feature", r.InitializeScenario)
}
func (r *registerUser) iHaveServiceWith(service *godog.Table) error {
	body, err := r.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = r.apiTest.UnmarshalJSON([]byte(body), &r.service); err != nil {
		return err
	}
	if r.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	createdService, err := r.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     r.service.Name,
		Password: r.service.Password,
		UserID:   uuid.MustParse(r.service.UserId),
	})
	if err != nil {
		return err
	}
	r.createdService.ServiceID = createdService.ServiceID
	if _, err := r.DB.Pool.Exec(context.Background(), "UPDATE services set status = true where id = $1", r.createdService.ServiceID); err != nil {
		return err
	}

	if err := r.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", r.service.Name)); err != nil {
		return err
	}
	return nil

}

func (r *registerUser) iHaveTheUserId(user_id *godog.Table) error {

	body, err := r.apiTest.ReadRow(user_id, nil, false)
	if err != nil {
		return err
	}

	r.apiTest.Body = body
	r.apiTest.SetHeader("Authorization", "Basic "+r.BasicAuth(r.createdService.ServiceID.String(), "123456"))
	r.apiTest.SetHeader("x-subject", r.service.UserId)
	r.apiTest.SetHeader("x-action", "*")
	r.apiTest.SetHeader("x-resource", "*")
	r.apiTest.SetHeader("x-tenant", "administrator")

	return nil
}

func (r *registerUser) iSendTheRequestToAddTheUser() error {
	r.apiTest.SendRequest()
	return nil
}

func (r *registerUser) theRequestShouldBeSuccessfull() error {
	if err := r.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("ok", "true"); err != nil {
		return err
	}

	return nil
}

func (r *registerUser) theRequestShouldFailWithError(message string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}
func (r *registerUser) iHaveAUserRegisteredOnMyService(user *godog.Table) error {

	if err := r.iHaveTheUserId(user); err != nil {
		return nil
	}

	return r.iSendTheRequestToAddTheUser()
}

func (r *registerUser) iWantToAddTheSameUserAgain(user *godog.Table) error {
	return r.iHaveTheUserId(user)
}
func (r *registerUser) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.URL = "/v1/users"
		r.apiTest.Method = http.MethodPost
		r.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^I have service with$`, r.iHaveServiceWith)
	ctx.Step(`^I have a user registered on my service:$`, r.iHaveAUserRegisteredOnMyService)
	ctx.Step(`^I have the user id:$`, r.iHaveTheUserId)
	ctx.Step(`^I want to add the same user again:$`, r.iWantToAddTheSameUserAgain)
	ctx.Step(`^i send the request to add the user$`, r.iSendTheRequestToAddTheUser)
	ctx.Step(`^the request should be successfull$`, r.theRequestShouldBeSuccessfull)
	ctx.Step(`^the request should fail with error "([^"]*)"$`, r.theRequestShouldFailWithError)
}
