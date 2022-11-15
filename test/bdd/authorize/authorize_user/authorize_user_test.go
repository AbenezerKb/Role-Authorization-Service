package authorize

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

type authorizeUserTest struct {
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
		OK   bool `json:"ok"`
		Data bool `json:"data"`
	}
}

func TestAuthorizeUser(t *testing.T) {
	a := &authorizeUserTest{}
	a.TestInstance = test.Initiate(context.Background(), "../../../../")
	a.apiTest.InitializeServer(a.Server)
	a.apiTest.InitializeTest(t, "authorize user test", "feature/authorize.feature", a.InitializeScenario)
}

func (a *authorizeUserTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
	domain, err := a.apiTest.ReadCellString(domainAndTenant, "domain")
	if err != nil {
		return err
	}
	result, err := a.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      domain,
		ServiceID: a.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	a.domain.ID = result.ID

	tenant, err := a.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	if err = a.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		ServiceID:  a.createdService.ServiceID,
		DomainID:   a.domain.ID,
	}); err != nil {
		return err
	}
	a.tenant = tenant

	return nil
}

func (a *authorizeUserTest) iHaveARoleInTenantWithThePermissionsBelow(role, tenant string, permissions *godog.Table) error {
	body, err := a.apiTest.ReadRow(permissions, []src.Type{
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
	stbody, err := a.apiTest.ReadRow(permissions, []src.Type{
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
	a.apiTest.UnmarshalJSON([]byte(body), &a.permission)
	st := dto.Statement{}
	a.apiTest.UnmarshalJSON([]byte(stbody), &st)

	statement, _ := st.Value()
	var result uuid.UUID
	result, err = a.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
		Name:        a.permission.Name,
		Description: a.permission.Description,
		ServiceID:   a.createdService.ServiceID,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
		Column5: []uuid.UUID{
			a.domain.ID,
		},
	})
	if err != nil {
		return err
	}
	a.createdPermissionId = result

	res, err := a.DB.CreateRole(context.Background(), db.CreateRoleParams{
		TenantName: tenant,
		ServiceID:  a.createdService.ServiceID,
		Column4: []uuid.UUID{
			a.createdPermissionId,
		},
		Name: a.createRole.Name,
	})
	if err != nil {
		return err
	}
	a.createdRoleResponseId = res.RoleID
	return nil
}

func (a *authorizeUserTest) iHaveServiceWith(service *godog.Table) error {
	body, err := a.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = a.apiTest.UnmarshalJSON([]byte(body), &a.service); err != nil {
		return err
	}
	if a.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	createdService, err := a.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     a.service.Name,
		Password: a.service.Password,
		UserID:   uuid.MustParse(a.service.UserId),
	})
	if err != nil {
		return err
	}
	a.createdService.ServiceID = createdService.ServiceID
	if _, err := a.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", a.createdService.ServiceID); err != nil {
		return err
	}

	if err := a.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", a.service.Name)); err != nil {
		return err
	}
	return nil
}

func (a *authorizeUserTest) iSendARequestToAuthorizeTheUser() error {
	a.apiTest.SetHeader("Authorization", "Basic "+a.BasicAuth(a.createdService.ServiceID.String(), "123456"))
	a.apiTest.SendRequest()
	return nil
}

func (a *authorizeUserTest) iWantToAuthorizeTheUserToPerformTheBelowActionOnTheResource(req *godog.Table) error {
	body, err := a.apiTest.ReadRow(req, nil, false)
	if err != nil {
		return err
	}
	a.apiTest.Body = body
	return nil
}

func (a *authorizeUserTest) theUserIsGrantedWithTheRole(role string) error {
	if err := a.DB.AssignRole(context.Background(), db.AssignRoleParams{
		RoleID:     a.createdRoleResponseId,
		UserID:     a.user.UserId,
		TenantName: a.tenant,
		ServiceID: a.createdService.ServiceID,
	}); err != nil {
		return err
	}

	if err := a.Opa.Refresh(context.Background(), fmt.Sprintf("Assigned role - [%v] to user - [%v]", a.createdRoleResponseId, a.user.UserId)); err != nil {
		return err
	}
	return nil
}

func (a *authorizeUserTest) theUserShouldBeAllowed() error {
	if err := a.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	if err := json.Unmarshal(a.apiTest.ResponseBody, &a.result); err != nil {
		return err
	}

	if err := a.apiTest.AssertEqual(a.result.Data, true); err != nil {
		return err
	}

	return nil
}

func (a *authorizeUserTest) thereIsAUserRegisteredOnTheSystem(user *godog.Table) error {
	body, err := a.apiTest.ReadRow(user, nil, false)
	if err != nil {
		return err
	}
	if err = a.apiTest.UnmarshalJSON([]byte(body), &a.user); err != nil {
		return err
	}

	if err := a.DB.RegisterUser(context.Background(), db.RegisterUserParams{
		UserID:    a.user.UserId,
		ServiceID: a.createdService.ServiceID,
	}); err != nil {
		return err
	}

	return nil
}

func (a *authorizeUserTest) theUserIsNotGrantedAnyRole() error {
	return nil
}

func (a *authorizeUserTest) theUserShouldBeDenied() error {
	if err := a.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	if err := json.Unmarshal(a.apiTest.ResponseBody, &a.result); err != nil {
		return err
	}

	if err := a.apiTest.AssertEqual(a.result.Data, false); err != nil {
		return err
	}

	return nil
}

func (a *authorizeUserTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		a.apiTest.URL = "/v1/authorize"
		a.apiTest.Method = http.MethodPost
		a.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = a.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^A registered domain and tenant$`, a.aRegisteredDomainAndTenant)
	ctx.Step(`^i have a role "([^"]*)" in tenant "([^"]*)" with the permissions below$`, a.iHaveARoleInTenantWithThePermissionsBelow)
	ctx.Step(`^I have service with$`, a.iHaveServiceWith)
	ctx.Step(`^the user is not granted any role$`, a.theUserIsNotGrantedAnyRole)
	ctx.Step(`^the user should be denied$`, a.theUserShouldBeDenied)
	ctx.Step(`^i send a request to authorize the user$`, a.iSendARequestToAuthorizeTheUser)
	ctx.Step(`^I  want to authorize the user to perform the below action on the resource:$`, a.iWantToAuthorizeTheUserToPerformTheBelowActionOnTheResource)
	ctx.Step(`^The user is granted with the "([^"]*)" role$`, a.theUserIsGrantedWithTheRole)
	ctx.Step(`^the user should be allowed$`, a.theUserShouldBeAllowed)
	ctx.Step(`^There is a user registered on the system:$`, a.thereIsAUserRegisteredOnTheSystem)
}
