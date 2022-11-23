package usergetpermissionsunderthedomain

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/platform/argon"
	"2f-authorization/test"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type getUserPermissionsWithinDomainTest struct {
	test.TestInstance
	apiTest               src.ApiTest
	service               dto.CreateService
	createdService        dto.CreateServiceResponse
	domain                dto.Domain
	permission            dto.CreatePermission
	createdPermissions    []dto.CreatePermission
	createdRoleResponseId []struct {
		roleId uuid.UUID
		tenant string
	}
	user   dto.RegisterUser
	result struct {
		OK   bool                    `json:"ok"`
		Data []dto.DomainPermissions `json:"data"`
	}
}

func TestGetUserPermissionsWithinDomain(t *testing.T) {
	g := &getUserPermissionsWithinDomainTest{}
	g.TestInstance = test.Initiate(context.Background(), "../../../../")
	g.apiTest.InitializeServer(g.Server)
	g.apiTest.InitializeTest(t, "user get permission within a domain test", "feature/user_get_permissions_under_the_domain.feature", g.InitializeScenario)
}
func (g *getUserPermissionsWithinDomainTest) aRegisteredDomainAndTenants(domainAndTenant *godog.Table) error {
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

	tenants, err := g.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	for _, t := range strings.Split(tenants, ",") {
		if err = g.DB.CreateTenent(context.Background(), db.CreateTenentParams{
			TenantName: t,
			ServiceID:  g.createdService.ServiceID,
			DomainID:   g.domain.ID,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (g *getUserPermissionsWithinDomainTest) aRegisteredUserOnTheSystem(user *godog.Table) error {
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

func (g *getUserPermissionsWithinDomainTest) aRoleInTenantWithTheFollowingPermissions(role, tenant string, permissions *godog.Table) error {
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
			Column: "domains",
			Kind:   src.Array,
			Ignore: true,
		},
		{
			WithName: "statement",
			Columns:  []string{"action", "resource", "effect", "fields"},
			Kind:     src.Object,
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
	if err := g.apiTest.UnmarshalJSON([]byte(body), &g.permission); err != nil {
		return err
	}

	statement, err := g.permission.Statement.Value()
	if err != nil {
		return err
	}

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
	g.createdPermissions = append(g.createdPermissions, g.permission)

	res, err := g.DB.CreateRole(context.Background(), db.CreateRoleParams{
		TenantName: tenant,
		ServiceID:  g.createdService.ServiceID,
		Column4: []uuid.UUID{
			result,
		},
		Name: role,
	})
	if err != nil {
		return err
	}
	g.createdRoleResponseId = append(g.createdRoleResponseId, struct {
		roleId uuid.UUID
		tenant string
	}{
		roleId: res.RoleID,
		tenant: tenant,
	})
	return nil
}

func (g *getUserPermissionsWithinDomainTest) theUserIsGrantedTheFollowingRoleInTheRespectiveTenant() error {
	for _, r := range g.createdRoleResponseId {
		if err := g.DB.AssignRole(context.Background(), db.AssignRoleParams{
			RoleID:     r.roleId,
			UserID:     g.user.UserId,
			TenantName: r.tenant,
			ServiceID:  g.createdService.ServiceID,
		}); err != nil {
			return err
		}
		if err := g.Opa.Refresh(context.Background(), fmt.Sprintf("Assigned role - [%v] to user - [%v]", g.createdRoleResponseId, g.user.UserId)); err != nil {
			return err
		}
	}
	return nil
}

func (g *getUserPermissionsWithinDomainTest) iHaveServiceWith(service *godog.Table) error {
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
func (g *getUserPermissionsWithinDomainTest) iRequestToGetMyPermissions() error {
	g.apiTest.SetHeader("Authorization", "Basic "+g.BasicAuth(g.createdService.ServiceID.String(), "123456"))
	g.apiTest.SendRequest()
	return nil
}

func (g *getUserPermissionsWithinDomainTest) iWantToGetMyPermissions() error {
	g.apiTest.URL = "/v1/users/" + g.user.UserId.String() + "/domains/" + g.domain.ID.String() + "/permissions"
	return nil
}

func (g *getUserPermissionsWithinDomainTest) theRequestShouldBeSuccessfull() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	if err := json.Unmarshal(g.apiTest.ResponseBody, &g.result); err != nil {
		return err
	}
	for _, t := range g.result.Data {
		found := false
		for _, p := range t.Permissions {
			for _, cp := range g.createdPermissions {

				if p.Name == cp.Name && p.Description == cp.Description && p.Statement.Action == cp.Statement.Action && p.Statement.Resource == cp.Statement.Resource && p.Statement.Effect == cp.Statement.Effect {
					found = true
					continue
				}

			}
			if !found {
				return fmt.Errorf("expected permission: %v", p)
			}
		}
	}
	return nil
}

func (g *getUserPermissionsWithinDomainTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		g.apiTest.Method = http.MethodGet
		g.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^A registered domain and tenants$`, g.aRegisteredDomainAndTenants)
	ctx.Step(`^A registered user on the system$`, g.aRegisteredUserOnTheSystem)
	ctx.Step(`^A role "([^"]*)" in tenant "([^"]*)" with the following permissions$`, g.aRoleInTenantWithTheFollowingPermissions)
	ctx.Step(`^The user is granted the following role in the respective tenant$`, g.theUserIsGrantedTheFollowingRoleInTheRespectiveTenant)
	ctx.Step(`^I have service with$`, g.iHaveServiceWith)
	ctx.Step(`^I request to get my permissions$`, g.iRequestToGetMyPermissions)
	ctx.Step(`^I want to get my permissions$`, g.iWantToGetMyPermissions)
	ctx.Step(`^The Request should be successfull$`, g.theRequestShouldBeSuccessfull)
}
