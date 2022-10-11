package createservice

import (
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/test"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type createServiceTest struct {
	test.TestInstance
	apiTest src.ApiTest
	service struct {
		OK   bool                      `json:"ok"`
		Data dto.CreateServiceResponse `json:"data"`
	}
}

func TestLogin(t *testing.T) {
	c := &createServiceTest{}
	c.TestInstance = test.Initiate(context.Background(), "../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "create service test", "feature/create_service.feature", c.InitializeScenario)
}

func (c *createServiceTest) iAmASystemUser() error {
	return nil
}

func (c *createServiceTest) iSendTheRequest(service *godog.Table) error {
	body, err := c.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	c.apiTest.Body = body

	c.apiTest.SendRequest()
	return nil
}

func (c *createServiceTest) theRequestShouldFailWithErrorMessage(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (c *createServiceTest) theResultShouldBeSuccessfull(arg1 string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}

	err := json.Unmarshal(c.apiTest.ResponseBody, &c.service)
	if err != nil {
		return err
	}
	if err := c.apiTest.AssertEqual(arg1, "true"); err != nil {
		return err
	}

	if err := c.apiTest.AssertColumnExists("data.service_id"); err != nil {
		return err
	}
	return nil
}

func (c *createServiceTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/services"
		c.apiTest.Method = http.MethodPost
		c.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = c.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})

	ctx.Step(`^i am a system user$`, c.iAmASystemUser)
	ctx.Step(`^i send the request:$`, c.iSendTheRequest)
	ctx.Step(`^the request should fail with error message "([^"]*)"$`, c.theRequestShouldFailWithErrorMessage)
	ctx.Step(`^the result should be successfull "([^"]*)"$`, c.theResultShouldBeSuccessfull)
}
