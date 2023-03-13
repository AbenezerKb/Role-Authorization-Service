package getroles

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
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

type getRolesTest struct {
	test.TestInstance
	apiTest               src.ApiTest
	service               dto.CreateService
	createdService        dto.CreateServiceResponse
	domain                dto.Domain
	tenant                string
	permission            dto.CreatePermission
	createdRoleResponseId uuid.UUID
	createdPermissionId   uuid.UUID
	result                struct {
		OK   bool       `json:"ok"`
		Data []dto.Role `json:"data"`
	}
}

func TestGetRoles(t *testing.T) {
	g := &getRolesTest{}
	g.TestInstance = test.Initiate(context.Background(), "../../../../")
	g.apiTest.InitializeServer(g.Server)
	g.apiTest.InitializeTest(t, "get roles test", "feature/get_roles.feature", g.InitializeScenario)
}
func (g *getRolesTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
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

	err = g.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		DomainID:   result.ID,
		ServiceID:  g.createdService.ServiceID,
	})

	if err != nil {
		return err
	}
	g.tenant = tenant

	return nil
}

func (g *getRolesTest) iHaveServiceWith(service *godog.Table) error {
	body, err := g.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = g.apiTest.UnmarshalJSON([]byte(body), &g.service); err != nil {
		return err
	}
	g.service.Password = "123456"
	createdService, err := g.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     g.service.Name,
		Password: g.service.Password,
		UserID:   uuid.MustParse(g.service.UserId),
	})
	if err != nil {
		return err
	}
	g.createdService.ServiceID = createdService.ServiceID
	if _, err := g.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE' where id = $1", g.createdService.ServiceID); err != nil {
		return err
	}

	if err := g.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", g.service.Name)); err != nil {
		return err
	}
	return nil
}

func (g *getRolesTest) iRequestToGetAllRolesUnderMyTenant() error {
	g.apiTest.SetHeader("Authorization", "Basic "+g.BasicAuth(g.createdService.ServiceID.String(), "123456"))
	g.apiTest.SetHeader("x-subject", g.service.UserId)
	g.apiTest.SetHeader("x-action", "*")
	g.apiTest.SetHeader("x-resource", "*")
	g.apiTest.SetHeader("x-tenant", g.tenant)

	g.apiTest.SendRequest()
	return nil
}

func (g *getRolesTest) iShouldGetAllRolesInMyTenant(roles *godog.Table) error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	if err := g.apiTest.UnmarshalJSON(g.apiTest.ResponseBody, &g.result); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(g.result.Data[0].ID, g.createdRoleResponseId); err != nil {
		return err

	}

	return nil
}

func (g *getRolesTest) theRoleIsRegisteredWithTheFollowingPermissionInTheTenant(role string, permission *godog.Table) error {
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
			WithName: "statement",
			Columns:  []string{"action", "resource", "effect"},
			Kind:     src.Object,
		},
	},
		true)
	if err != nil {
		return err
	}
	g.apiTest.UnmarshalJSON([]byte(body), &g.permission)

	statement, _ := g.permission.Statement.Value()
	var result uuid.UUID
	result, err = g.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
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

	createdRoleResponse, err := g.DB.CreateRole(context.Background(), db.CreateRoleParams{
		TenantName: g.tenant,
		ServiceID:  g.createdService.ServiceID,
		Column4:    []uuid.UUID{g.createdPermissionId},
		Name:       role,
	})
	if err != nil {
		return err
	}
	g.createdRoleResponseId = createdRoleResponse.RoleID
	return nil
}

func (g *getRolesTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		g.apiTest.URL = "/v1/roles"
		g.apiTest.Method = http.MethodGet
		g.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^A registered domain and tenant$`, g.aRegisteredDomainAndTenant)
	ctx.Step(`^I have service with$`, g.iHaveServiceWith)
	ctx.Step(`^I request to get all roles under my tenant$`, g.iRequestToGetAllRolesUnderMyTenant)
	ctx.Step(`^I should get all roles in my tenant$`, g.iShouldGetAllRolesInMyTenant)
	ctx.Step(`^The role "([^"]*)" is registered with the following permission in the tenant$`, g.theRoleIsRegisteredWithTheFollowingPermissionInTheTenant)
}
