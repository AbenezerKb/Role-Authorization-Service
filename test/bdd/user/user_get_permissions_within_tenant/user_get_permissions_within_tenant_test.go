package usergetpermissionswithintenant

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/platform/argon"
	"2f-authorization/test"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type getUserPermissionsWithinTenantTest struct {
	test.TestInstance
	apiTest               src.ApiTest
	service               dto.CreateService
	createdService        dto.CreateServiceResponse
	domain                dto.Domain
	tenant                string
	createRole            dto.CreateRole
	permission            dto.CreatePermission
	createdRoleResponseId uuid.UUID
	createdPermissionId   uuid.UUID
	user                  dto.RegisterUser
	result                struct {
		OK   bool             `json:"ok"`
		Data []dto.Permission `json:"data"`
	}
}

func TestGetUserPermissionsWithinTenant(t *testing.T) {
	g := &getUserPermissionsWithinTenantTest{}
	g.TestInstance = test.Initiate(context.Background(), "../../../../")
	g.apiTest.InitializeServer(g.Server)
	g.apiTest.InitializeTest(t, "user get permission within a tenant test", "feature/user_get_permissions_within_tenant.feature", g.InitializeScenario)
}

func (g *getUserPermissionsWithinTenantTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
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

func (g *getUserPermissionsWithinTenantTest) iAmGrantedAnRole(role string) error {
	if err := g.DB.AssignRole(context.Background(), db.AssignRoleParams{
		ID:     g.createdRoleResponseId,
		UserID:     g.user.UserId,
		TenantName: g.tenant,
		ServiceID:  g.createdService.ServiceID,
	}); err != nil {
		return err
	}

	if err := g.Opa.Refresh(context.Background(), fmt.Sprintf("Assigned role - [%v] to user - [%v]", g.createdRoleResponseId, g.user.UserId)); err != nil {
		return err
	}
	return nil
}

func (g *getUserPermissionsWithinTenantTest) iHaveARoleInTenantWithTheFollowingPermissions(role, tenant string, permissions *godog.Table) error {
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
			Columns:  []string{"action", "resource", "effect", "fields"},
			Kind:     src.Object,
		},
	},
		true)
	if err != nil {
		return err
	}
	stbody, err := g.apiTest.ReadRow(permissions, []src.Type{
		{
			Column: "action",
			Kind:   src.Any,
		},
		{
			Column: "resource",
			Kind:   src.Any,
		},
		{
			Column: "effect",
			Kind:   src.Any,
		},
		{
			Column: "fields",
			Kind:   src.Array,
		},
	},
		true)
	if err != nil {
		return err
	}
	g.apiTest.UnmarshalJSON([]byte(body), &g.permission)
	st := dto.Statement{}
	g.apiTest.UnmarshalJSON([]byte(stbody), &st)

	statement, _ := st.Value()
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

	res, err := g.DB.CreateRole(context.Background(), db.CreateRoleParams{
		TenantName: tenant,
		ServiceID:  g.createdService.ServiceID,
		Column4: []uuid.UUID{
			g.createdPermissionId,
		},
		Name: g.createRole.Name,
	})
	if err != nil {
		return err
	}
	g.createdRoleResponseId = res.RoleID
	return nil
}

func (g *getUserPermissionsWithinTenantTest) iHaveServiceWith(service *godog.Table) error {
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

func (g *getUserPermissionsWithinTenantTest) iRequestToGetMyPermissions() error {
	g.apiTest.SetHeader("Authorization", "Basic "+g.BasicAuth(g.createdService.ServiceID.String(), "123456"))
	g.apiTest.SendRequest()
	return nil
}

func (g *getUserPermissionsWithinTenantTest) iWantToGetMyPermissions() error {
	g.apiTest.URL = "/v1/users/" + g.user.UserId.String() + "/tenants/" + g.tenant + "/permissions"
	return nil
}

func (g *getUserPermissionsWithinTenantTest) theRequestShouldBeSuccessfull() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	if err := json.Unmarshal(g.apiTest.ResponseBody, &g.result); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(g.result.OK, true); err != nil {
		return err
	}
	for _, p := range g.result.Data {
		if err := g.apiTest.AssertEqual(p.Name, g.permission.Name); err != nil {
			return err
		}
		if err := g.apiTest.AssertEqual(p.Description, g.permission.Description); err != nil {
			return err
		}
		if err := g.apiTest.AssertEqual(p.Statement.Action, g.permission.Statement.Action); err != nil {
			return err
		}
		if err := g.apiTest.AssertEqual(p.Statement.Resource, g.permission.Statement.Resource); err != nil {
			return err
		}
		if err := g.apiTest.AssertEqual(p.Statement.Effect, g.permission.Statement.Effect); err != nil {
			return err
		}
	}
	return nil
}

func (g *getUserPermissionsWithinTenantTest) iAmRegisteredOnTheSystem(user *godog.Table) error {
	body, err := g.apiTest.ReadRow(user, nil, false)
	if err != nil {
		return err
	}
	if err = g.apiTest.UnmarshalJSON([]byte(body), &g.user); err != nil {
		return err
	}

	if err := g.DB.RegisterUser(context.Background(), db.RegisterUserParams{
		UserID:    g.user.UserId,
		ServiceID: g.createdService.ServiceID,
	}); err != nil {
		return err
	}

	return nil
}

func (g *getUserPermissionsWithinTenantTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		g.apiTest.Method = http.MethodGet
		g.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^A registered domain and tenant$`, g.aRegisteredDomainAndTenant)
	ctx.Step(`^I am granted an "([^"]*)" role$`, g.iAmGrantedAnRole)
	ctx.Step(`^I am registered on the system$`, g.iAmRegisteredOnTheSystem)
	ctx.Step(`^I have a role "([^"]*)" in tenant "([^"]*)" with the following permissions$`, g.iHaveARoleInTenantWithTheFollowingPermissions)
	ctx.Step(`^I have service with$`, g.iHaveServiceWith)
	ctx.Step(`^I request to get my permissions$`, g.iRequestToGetMyPermissions)
	ctx.Step(`^I want to get my permissions$`, g.iWantToGetMyPermissions)
	ctx.Step(`^The Request should be successfull$`, g.theRequestShouldBeSuccessfull)
}
