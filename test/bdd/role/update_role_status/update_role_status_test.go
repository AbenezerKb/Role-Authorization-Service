package updaterolestatus

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

type updateRoleStatusTest struct {
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
	roleStatus            dto.UpdateRoleStatus
}

func TestUpdateRoleStatus(t *testing.T) {
	u := &updateRoleStatusTest{}
	u.TestInstance = test.Initiate(context.Background(), "../../../../")
	u.apiTest.InitializeServer(u.Server)
	u.apiTest.InitializeTest(t, "update role status test", "feature/update_role_status.feature", u.InitializeScenario)
}
func (u *updateRoleStatusTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {

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

	err = u.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		DomainID:   result.ID,
		ServiceID:  u.createdService.ServiceID,
	})

	if err != nil {
		return err
	}
	u.tenant = tenant

	return nil
}

func (u *updateRoleStatusTest) iHaveARoleInTenantWithThePermissionsBelow(role, tenant string, permissions *godog.Table) error {
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

func (u *updateRoleStatusTest) iHaveServiceWith(service *godog.Table) error {
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
	if _, err := u.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE' where id = $1", u.createdService.ServiceID); err != nil {
		return err
	}

	if err := u.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", u.service.Name)); err != nil {
		return err
	}
	return nil
}

func (u *updateRoleStatusTest) iSendTheRequestToUpdateTheStatus() error {
	u.apiTest.SetHeader("Authorization", "Basic "+u.BasicAuth(u.createdService.ServiceID.String(), "123456"))
	u.apiTest.SetHeader("x-subject", u.service.UserId)
	u.apiTest.SetHeader("x-action", "*")
	u.apiTest.SetHeader("x-resource", "*")
	u.apiTest.SetHeader("x-tenant", u.tenant)

	u.apiTest.SendRequest()
	return nil
}

func (u *updateRoleStatusTest) iWantToUpdateTheRolesStatusTo(status string) error {
	u.roleStatus.Status = status
	u.apiTest.URL = "/v1/roles/" + u.createdRoleResponseId.String() + "/status"
	body, err := json.Marshal(&u.roleStatus)
	if err != nil {
		return err
	}
	u.apiTest.Body = string(body)
	return nil
}

func (u *updateRoleStatusTest) theRoleStatusShouldUpdateTo(status string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	return nil
}
func (u *updateRoleStatusTest) theRoleIsNotOnTheSystem(role string) error {
	u.apiTest.URL = "/v1/roles/" + role + "/status"
	return nil
}

func (u *updateRoleStatusTest) thenIShouldGetAnErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}
func (u *updateRoleStatusTest) thenIShouldGetAFieldErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}
func (u *updateRoleStatusTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		u.apiTest.Method = http.MethodPatch
		u.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^Then I should get a field error with message "([^"]*)"$`, u.thenIShouldGetAFieldErrorWithMessage)
	ctx.Step(`^the role is not on the system "([^"]*)"$`, u.theRoleIsNotOnTheSystem)
	ctx.Step(`^Then I should get an error with message "([^"]*)"$`, u.thenIShouldGetAnErrorWithMessage)
	ctx.Step(`^A registered domain and tenant$`, u.aRegisteredDomainAndTenant)
	ctx.Step(`^i have a role "([^"]*)" in tenant "([^"]*)" with the permissions below$`, u.iHaveARoleInTenantWithThePermissionsBelow)
	ctx.Step(`^I have service with$`, u.iHaveServiceWith)
	ctx.Step(`^I send the request to update the status$`, u.iSendTheRequestToUpdateTheStatus)
	ctx.Step(`^I want to update the role\'s status to "([^"]*)"$`, u.iWantToUpdateTheRolesStatusTo)
	ctx.Step(`^the role status should update to "([^"]*)"$`, u.theRoleStatusShouldUpdateTo)
}
