package getpermissiondetail

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
	"github.com/jackc/pgtype"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type getPermissionDetails struct {
	test.TestInstance
	apiTest             src.ApiTest
	service             dto.CreateService
	createdService      dto.CreateServiceResponse
	domain              dto.Domain
	tenant              string
	permission          dto.CreatePermission
	createdPermissionId uuid.UUID
	result              struct {
		OK   bool           `json:"ok"`
		Data dto.Permission `json:"data"`
	}
	expectedPermission dto.Permission
}

func TestGetPermissionDetails(t *testing.T) {
	g := &getPermissionDetails{}
	g.TestInstance = test.Initiate(context.Background(), "../../../../")
	g.apiTest.InitializeServer(g.Server)
	g.apiTest.InitializeTest(t, "get permission detail test", "feature/get_permission_detail.feature", g.InitializeScenario)
}
func (g *getPermissionDetails) aPermissionRegisteredOnTheDomain(permission *godog.Table) error {
	body, err := g.apiTest.ReadRow(permission, []src.Type{
		{
			Column: "name",
			Kind:   src.Any,
		},
		{
			Column: "description",
			Kind:   src.Any,
		},
		{
			Column: "fields",
			Kind:   src.Array,
		},
		{
			WithName: "statement",
			Columns:  []string{"action", "resource", "effect", "fields"},
			Kind:     src.Object,
		},
	},
		true)
	if err != nil {
		return err
	}
	if err := g.apiTest.UnmarshalJSON([]byte(body), &g.permission); err != nil {
		return err
	}
	statement, _ := g.permission.Statement.Value()
	result, err := g.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
		Name:        g.permission.Name,
		Description: g.permission.Description,
		ServiceID:   g.createdService.ServiceID,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
		Column5: []uuid.UUID{
			g.domain.ID,
		},
	})
	if err != nil {
		return err
	}
	g.createdPermissionId = result

	return nil
}

func (g *getPermissionDetails) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
	domain, err := g.apiTest.ReadCellString(domainAndTenant, "domain")
	if err != nil {
		return err
	}
	result, err := g.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      domain,
		ServiceID: g.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	g.domain.ID = result.ID

	tenant, err := g.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	if err = g.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		ServiceID:  g.createdService.ServiceID,
		DomainID:   g.domain.ID,
	}); err != nil {
		return err
	}
	g.tenant = tenant

	return nil
}

func (g *getPermissionDetails) iHaveServiceWith(service *godog.Table) error {
	body, err := g.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = g.apiTest.UnmarshalJSON([]byte(body), &g.service); err != nil {
		return err
	}
	if g.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	createdService, err := g.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     g.service.Name,
		Password: g.service.Password,
		UserID:   uuid.MustParse(g.service.UserId),
	})
	if err != nil {
		return err
	}
	g.createdService.ServiceID = createdService.ServiceID
	if _, err := g.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", g.createdService.ServiceID); err != nil {
		return err
	}

	if err := g.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", g.service.Name)); err != nil {
		return err
	}
	return nil
}

func (g *getPermissionDetails) iSendTheRequest() error {
	g.apiTest.SetHeader("Authorization", "Basic "+g.BasicAuth(g.createdService.ServiceID.String(), "123456"))
	g.apiTest.SetHeader("x-subject", g.service.UserId)
	g.apiTest.SetHeader("x-action", "*")
	g.apiTest.SetHeader("x-resource", "*")
	g.apiTest.SetHeader("x-tenant", g.tenant)

	g.apiTest.SendRequest()
	return nil
}

func (g *getPermissionDetails) iShouldGetThePermissionDetail(permissions *godog.Table) error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	body, err := g.apiTest.ReadRow(permissions, []src.Type{
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
	if err = g.apiTest.UnmarshalJSON([]byte(body), &g.expectedPermission); err != nil {
		return err
	}

	if err = g.apiTest.UnmarshalJSON(g.apiTest.ResponseBody, &g.result); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(g.result.Data.Name, g.expectedPermission.Name); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(g.result.Data.Description, g.expectedPermission.Description); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(g.result.Data.Statement.Action, g.expectedPermission.Statement.Action); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(g.result.Data.Statement.Effect, g.expectedPermission.Statement.Effect); err != nil {
		return err
	}
	return g.apiTest.AssertEqual(g.result.Data.Statement.Resource, g.expectedPermission.Statement.Resource)

}

func (g *getPermissionDetails) iWantToGetThePermissionDetail() error {
	g.apiTest.Method = http.MethodGet
	g.apiTest.URL = "/v1/permissions/" + g.createdPermissionId.String()
	return nil
}

func (g *getPermissionDetails) theRequestShouldBeSuccessfull() error {
	return g.apiTest.AssertStatusCode(http.StatusOK)
}
func (g *getPermissionDetails) iSendTheRequestToGetThePermissionDetails() error {
	g.apiTest.SetHeader("Authorization", "Basic "+g.BasicAuth(g.createdService.ServiceID.String(), "123456"))
	g.apiTest.SetHeader("x-subject", g.service.UserId)
	g.apiTest.SetHeader("x-action", "*")
	g.apiTest.SetHeader("x-resource", "*")
	g.apiTest.SetHeader("x-tenant", g.tenant)

	g.apiTest.SendRequest()
	return nil
}

func (g *getPermissionDetails) iShouldGetAnErrorWithMessage(message string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return g.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (g *getPermissionDetails) thePermissionDoesNotExists(permission string) error {
	g.apiTest.Method = http.MethodGet
	g.apiTest.URL = "/v1/permissions/" + permission
	return nil
}
func (g *getPermissionDetails) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		g.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^I send the request to get the permission details$`, g.iSendTheRequestToGetThePermissionDetails)
	ctx.Step(`^I should get an error with message "([^"]*)"$`, g.iShouldGetAnErrorWithMessage)
	ctx.Step(`^the permission does not exist "([^"]*)"$`, g.thePermissionDoesNotExists)
	ctx.Step(`^A permission registered on the domain$`, g.aPermissionRegisteredOnTheDomain)
	ctx.Step(`^A registered domain and tenant$`, g.aRegisteredDomainAndTenant)
	ctx.Step(`^I have service with$`, g.iHaveServiceWith)
	ctx.Step(`^I send the request$`, g.iSendTheRequest)
	ctx.Step(`^I should get the permission detail$`, g.iShouldGetThePermissionDetail)
	ctx.Step(`^I want to get the permission detail$`, g.iWantToGetThePermissionDetail)
	ctx.Step(`^The request should be successfull$`, g.theRequestShouldBeSuccessfull)
}
