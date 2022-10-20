package domain

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/platform/argon"
	"2f-authorization/test"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type createTestDomain struct {
	test.TestInstance
	apiTest      src.ApiTest
	servicemodel db.CreateServiceParams
	serviceId    uuid.UUID
}

func TestCreateDomain(t *testing.T) {
	c := &createTestDomain{}
	c.TestInstance = test.Initiate(context.Background(), "../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "create domain test", "feature/create_domain.feature", c.InitializeScenario)
}

func (c *createTestDomain) iAmASystemUser() error {
	return nil
}

func (c *createTestDomain) iHaveServiceWith(service *godog.Table) error {

	body, err := c.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}

	c.apiTest.UnmarshalJSON([]byte(body), &c.servicemodel)

	if c.servicemodel.Password, err = argon.CreateHash("password", argon.DefaultParams); err != nil {
		return err
	}

	result, err := c.DB.CreateService(context.Background(), c.servicemodel)
	if err != nil {
		return err
	}

	_, err = c.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", result.ServiceID)
	if err != nil {
		return err
	}
	if err := c.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", c.servicemodel.Name)); err != nil {
		return err
	}
	c.serviceId = result.ServiceID

	return nil
}
func (c *createTestDomain) iSendTheRequest(domain *godog.Table) error {

	body, err := c.apiTest.ReadRow(domain, nil, false)
	if err != nil {
		return err
	}

	c.apiTest.Body = body
	c.apiTest.SetHeader("Authorization", "Basic "+basicAuth(c.serviceId.String(), "password"))
	c.apiTest.SetHeader("x-subject", c.servicemodel.UserID.String())
	c.apiTest.SetHeader("x-action", "*")
	c.apiTest.SetHeader("x-resource", "*")
	c.apiTest.SetHeader("x-tenant", "administrator")

	c.apiTest.SendRequest()
	return nil
}

func (c *createTestDomain) theResultShouldBeSuccessfull(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("ok", message); err != nil {
		return err
	}

	return nil
}

func (c *createTestDomain) theResultShouldBeEmptyError(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}

	return nil
}

func (c *createTestDomain) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/domains"
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
	ctx.Step(`^I have service with$`, c.iHaveServiceWith)
	ctx.Step(`^the result should be successfull "([^"]*)"$`, c.theResultShouldBeSuccessfull)
	ctx.Step(`^the result should be empty error "([^"]*)"$`, c.theResultShouldBeEmptyError)
	ctx.Step(`^the result should be successfull "([^"]*)"$`, c.theResultShouldBeSuccessfull)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
