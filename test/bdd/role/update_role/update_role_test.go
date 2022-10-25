package updaterole

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

type updateRoleTest struct {
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
	updatedRole           dto.UpdateRole
	per                   []dto.CreatePermission
	result                struct {
		OK   bool     `json:"ok"`
		Data dto.Role `json:"data"`
	}
}

func TestUpdateRole(t *testing.T) {
	c := &updateRoleTest{}
	c.TestInstance = test.Initiate(context.Background(), "../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "update role test", "feature/update_role.feature", c.InitializeScenario)
}

func (u *updateRoleTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
	domain, err := u.apiTest.ReadCellString(domainAndTenant, "domain")
	if err != nil {
		return err
	}
	result, err := u.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      domain,
		ServiceID: u.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	u.domain.ID = result.ID

	tenant, err := u.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	if err = u.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		ServiceID:  u.createdService.ServiceID,
		DomainID:   u.domain.ID,
	}); err != nil {
		return err
	}
	u.tenant = tenant

	return nil
}

func (u *updateRoleTest) iHaveARoleInTenantWithThePermissionsBelow(role, tenant string, permissions *godog.Table) error {
	body, err := u.apiTest.ReadRow(permissions, []src.Type{
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
	u.apiTest.UnmarshalJSON([]byte(body), &u.permission)

	statement, _ := u.permission.Statement.Value()
	var result uuid.UUID
	result, err = u.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
		Name:        u.permission.Name,
		Description: u.permission.Description,
		ServiceID:   u.createdService.ServiceID,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
	})
	if err != nil {
		return err
	}
	u.createdPermissionId = result

	res, err := u.DB.CreateRole(context.Background(), db.CreateRoleParams{
		TenantName: tenant,
		ServiceID:  u.createdService.ServiceID,
		Column4: []uuid.UUID{
			u.createdPermissionId,
		},
		Name: u.createRole.Name,
	})
	if err != nil {
		return err
	}
	u.createdRoleResponseId = res.RoleID
	return nil
}

func (u *updateRoleTest) iHaveServiceWith(service *godog.Table) error {
	body, err := u.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = u.apiTest.UnmarshalJSON([]byte(body), &u.service); err != nil {
		return err
	}
	if u.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	createdService, err := u.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     u.service.Name,
		Password: u.service.Password,
		UserID:   uuid.MustParse(u.service.UserId),
	})
	if err != nil {
		return err
	}
	u.createdService.ServiceID = createdService.ServiceID
	if _, err := u.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", u.createdService.ServiceID); err != nil {
		return err
	}

	if err := u.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", u.service.Name)); err != nil {
		return err
	}
	return nil
}

func (u *updateRoleTest) iSendARequestToUpdateTheRole() error {
	u.apiTest.SetHeader("Authorization", "Basic "+u.BasicAuth(u.createdService.ServiceID.String(), "123456"))
	u.apiTest.SetHeader("x-subject", u.service.UserId)
	u.apiTest.SetHeader("x-action", "*")
	u.apiTest.SetHeader("x-resource", "*")
	u.apiTest.SetHeader("x-tenant", u.tenant)

	u.apiTest.SendRequest()
	return nil
}

func (u *updateRoleTest) iWantToUpdateTheRoleWithTheFollowingPermissions(role string, permissions *godog.Table) error {
	body, err := u.apiTest.ReadRows(permissions, []src.Type{
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
	var pId []uuid.UUID
	u.apiTest.UnmarshalJSON([]byte(body), &u.per)
	for _, p := range u.per {
		statement, _ := p.Statement.Value()
		var result uuid.UUID
		result, err = u.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
			Name:        p.Name,
			Description: p.Description,
			ServiceID:   u.createdService.ServiceID,
			Statement: pgtype.JSON{
				Bytes:  statement,
				Status: pgtype.Present,
			},
		})
		if err != nil {
			return err
		}
		pId = append(pId, result)
	}
	u.updatedRole.RoleID = u.createdRoleResponseId
	u.updatedRole.PermissionsID = pId
	requestBody, err := json.Marshal(u.updatedRole)
	if err != nil {
		return err
	}

	u.apiTest.Body = string(requestBody)
	u.apiTest.SetHeader("Authorization", "Basic "+u.BasicAuth(u.createdService.ServiceID.String(), "123456"))
	u.apiTest.URL = "/v1/roles/" + u.createdRoleResponseId.String()

	return nil
}

func (u *updateRoleTest) theRoleShouldBeUpdated() error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	return nil
}

func (u *updateRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		u.apiTest.Method = http.MethodPut
		u.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^A registered domain and tenant$`, u.aRegisteredDomainAndTenant)
	ctx.Step(`^i have a role "([^"]*)" in tenant "([^"]*)" with the permissions below$`, u.iHaveARoleInTenantWithThePermissionsBelow)
	ctx.Step(`^I have service with$`, u.iHaveServiceWith)
	ctx.Step(`^i send a request to update the role$`, u.iSendARequestToUpdateTheRole)
	ctx.Step(`^I  want to update the role "([^"]*)" with the following permissions:$`, u.iWantToUpdateTheRoleWithTheFollowingPermissions)
	ctx.Step(`^the role should be updated$`, u.theRoleShouldBeUpdated)
}
