package role

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
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type createRoleTest struct {
	test.TestInstance
	apiTest             src.ApiTest
	service             dto.CreateService
	createdService      dto.CreateServiceResponse
	domain              dto.Domain
	tenant              string
	permission          dto.CreatePermission
	createdPermissionId uuid.UUID
	role                dto.CreateRole
	result              struct {
		OK   bool     `json:"ok"`
		Data dto.Role `json:"data"`
	}
}

func TestCreateRole(t *testing.T) {
	c := &createRoleTest{}
	c.TestInstance = test.Initiate(context.Background(), "../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "create role test", "feature/create_role.feature", c.InitializeScenario)
}

func (c *createRoleTest) aPermissionsRegisteredOnTheDomain(permission *godog.Table) error {
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
		{
			Column: "domains",
			Kind:   src.Array,
		},
	},
		true)
	if err != nil {
		return err
	}
	c.apiTest.UnmarshalJSON([]byte(body), &c.permission)

	statement, _ := c.permission.Statement.Value()
	var result uuid.UUID
	result, err = c.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
		Name:        c.permission.Name,
		Description: c.permission.Description,
		ServiceID:   c.createdService.ServiceID,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
	})
	if err != nil {
		return err
	}
	c.createdPermissionId = result
	return nil
}

func (c *createRoleTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
	domain, err := c.apiTest.ReadCellString(domainAndTenant, "domain")
	if err != nil {
		return err
	}
	result, err := c.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      domain,
		ServiceID: c.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	c.domain.ID = result.ID

	tenant, err := c.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	if err = c.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		ServiceID:  c.createdService.ServiceID,
		DomainID:   c.domain.ID,
	}); err != nil {
		return err
	}
	c.tenant = tenant

	return nil
}

func (c *createRoleTest) iHaveServiceWith(service *godog.Table) error {
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

func (c *createRoleTest) iRequestToCreateARoleWithThePermissions(permission string, role *godog.Table) error {
	body, err := c.apiTest.ReadRow(role, nil, false)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(body), &c.role)
	c.role.PermissionID = []uuid.UUID{
		c.createdPermissionId,
	}
	c.role.ServiceID = c.createdService.ServiceID
	c.role.TenantName = c.tenant
	data, err := json.Marshal(c.role)
	if err != nil {
		return err
	}
	c.apiTest.Body = string(data)
	c.apiTest.SetHeader("Authorization", "Basic "+c.BasicAuth(c.createdService.ServiceID.String(), "123456"))
	c.apiTest.SetHeader("x-subject", c.service.UserId)
	c.apiTest.SetHeader("x-action", "*")
	c.apiTest.SetHeader("x-resource", "*")
	c.apiTest.SetHeader("x-tenant", "administrator")

	c.apiTest.SendRequest()
	return nil
}

func (c *createRoleTest) myRequestShouldFailWith(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (c *createRoleTest) theRoleShouldSuccessfullyBeCreated() error {
	if err := c.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("ok", "true"); err != nil {
		return err
	}

	return nil
}

func (c *createRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/roles"
		c.apiTest.Method = http.MethodPost
		c.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = c.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})

	ctx.Step(`^a permissions registered on the domain$`, c.aPermissionsRegisteredOnTheDomain)
	ctx.Step(`^A registered domain and tenant$`, c.aRegisteredDomainAndTenant)
	ctx.Step(`^I have service with$`, c.iHaveServiceWith)
	ctx.Step(`^I request to create a role with the "([^"]*)" permissions$`, c.iRequestToCreateARoleWithThePermissions)
	ctx.Step(`^my request should fail with "([^"]*)"$`, c.myRequestShouldFailWith)
	ctx.Step(`^the role should successfully be created$`, c.theRoleShouldSuccessfullyBeCreated)
}
