package tenant

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/test"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type createTenant struct {
	test.TestInstance
	apiTest      src.ApiTest
	servicemodel db.CreateServiceParams
	serviceId    uuid.UUID
	domain       dto.Domain
	tenant       dto.CreateTenent
}

func TestCreateTenant(t *testing.T) {
	c := &createTenant{}
	c.TestInstance = test.Initiate(context.Background(), "../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "create domain test", "feature/create_tenant.feature", c.InitializeScenario)
}

func (c *createTenant) iHaveServiceWith(service *godog.Table) error {

	body, err := c.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}

	err = c.apiTest.UnmarshalJSON([]byte(body), &c.servicemodel)
	if err != nil {
		return err
	}
	c.servicemodel.Password = "password"
	result, err := c.DB.CreateService(context.Background(), c.servicemodel)
	if err != nil {
		return err
	}

	_, err = c.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", result.ServiceID)
	if err != nil {
		return err
	}

	c.serviceId = result.ServiceID

	return nil
}
func (c *createTenant) iSendTheRequest() error {

	c.apiTest.SetHeader("Authorization", "Basic "+basicAuth(c.serviceId.String(), "password"))
	c.apiTest.SendRequest()
	return nil
}

func (c *createTenant) theResultShouldBeSuccessfull(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("ok", message); err != nil {
		return err
	}

	return nil
}

func (c *createTenant) theResultShouldBeEmptyError(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}

	return nil
}

func (c *createTenant) aDomain(domain *godog.Table) error {
	body, err := c.apiTest.ReadRow(domain, nil, false)
	if err != nil {
		return err
	}
	if err = c.apiTest.UnmarshalJSON([]byte(body), &c.domain); err != nil {
		return err
	}
	result, err := c.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      c.domain.Name,
		ServiceID: c.serviceId,
	})
	if err != nil {
		return err
	}
	c.domain.ID = result.ID

	return nil
}
func (c *createTenant) iWantToCreateATenantWithData(tenant *godog.Table) error {

	body, err := c.apiTest.ReadRow(tenant, nil, false)
	if err != nil {
		return err
	}

	if err = c.apiTest.UnmarshalJSON([]byte(body), &c.tenant); err != nil {
		return err
	}
	c.tenant.DomainID = c.domain.ID

	data, err := json.Marshal(c.tenant)
	if err != nil {
		return err
	}

	c.apiTest.Body = string(data)
	return nil
}
func (c *createTenant) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/tenants"
		c.apiTest.Method = http.MethodPost
		c.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = c.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})

	ctx.Step(`^i send the request$`, c.iSendTheRequest)
	ctx.Step(`^a domain$`, c.aDomain)
	ctx.Step(`^i want to create a tenant with data:$`, c.iWantToCreateATenantWithData)
	ctx.Step(`^I have service with$`, c.iHaveServiceWith)
	ctx.Step(`^the result should be successfull "([^"]*)"$`, c.theResultShouldBeSuccessfull)
	ctx.Step(`^the result should be empty error "([^"]*)"$`, c.theResultShouldBeEmptyError)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
