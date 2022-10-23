package permission

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/platform/argon"
	"2f-authorization/test"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type createTestPermission struct {
	test.TestInstance
	apiTest        src.ApiTest
	service        dto.CreateService
	createdService dto.CreateServiceResponse
	domainrequest  dto.Domain
	permission     dto.CreatePermission
	result         struct {
		OK bool `json:"ok"`
	}
}

func TestCreateDomain(t *testing.T) {
	c := &createTestPermission{}
	c.TestInstance = test.Initiate(context.Background(), "../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "create permission test", "feature/create_permission.feature", c.InitializeScenario)
}

func (c *createTestPermission) aRegisteredDomain(domain *godog.Table) error {
	body, err := c.apiTest.ReadRow(domain, nil, false)
	if err != nil {
		return err
	}
	if err = c.apiTest.UnmarshalJSON([]byte(body), &c.domainrequest); err != nil {
		return err
	}
	createddomain, err := c.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      c.domainrequest.Name,
		ServiceID: c.createdService.ServiceID,
	})
	c.domainrequest.ID = createddomain.ID
	if err != nil {
		return err
	}
	return nil
}

func (c *createTestPermission) iCreateAPermmissionInTheDomain(permission *godog.Table) error {
	body, err := c.apiTest.ReadRow(permission, []src.Type{
		{
			Column: "name",
			Kind:   src.Any,
		},
		{
			Column: "description",
			Kind:   src.Any,
		},
		{
			WithName: "statement",
			Columns:  []string{"action", "resource", "effect"},
			Kind:     src.Object,
		},
	},
		true)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(body), &c.permission)

	c.permission.Domain = []uuid.UUID{c.domainrequest.ID}
	requestBody, err := json.Marshal(c.permission)
	if err != nil {
		return err
	}

	c.apiTest.Body = string(requestBody)
	c.apiTest.SetHeader("Authorization", "Basic "+basicAuth(c.createdService.ServiceID.String(), "123456"))

	c.apiTest.SendRequest()
	return nil
}

func (c *createTestPermission) iHaveARegisteredService(service *godog.Table) error {
	body, err := c.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = c.apiTest.UnmarshalJSON([]byte(body), &c.service); err != nil {
		return err
	}
	if c.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	createdService, err := c.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     c.service.Name,
		Password: c.service.Password,
		UserID:   uuid.MustParse(c.service.UserId),
	})
	if err != nil {
		return err
	}
	c.createdService.ServiceID = createdService.ServiceID
	if _, err := c.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", c.createdService.ServiceID); err != nil {
		return err
	}

	if err := c.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", c.service.Name)); err != nil {
		return err
	}
	return nil
}

func (c *createTestPermission) theRequestShouldFailWithError(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (c *createTestPermission) theServiceShouldBeDeleted() error {
	if err := c.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}

	return nil
}

func (c *createTestPermission) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/permissions"
		c.apiTest.Method = http.MethodPost
		c.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = c.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^A registered domain$`, c.aRegisteredDomain)
	ctx.Step(`^I create a permmission in the domain:$`, c.iCreateAPermmissionInTheDomain)
	ctx.Step(`^I have a registered service$`, c.iHaveARegisteredService)
	ctx.Step(`^The request should fail with error "([^"]*)"$`, c.theRequestShouldFailWithError)
	ctx.Step(`^The service should be deleted$`, c.theServiceShouldBeDeleted)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
