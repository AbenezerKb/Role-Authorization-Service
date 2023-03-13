package updatepermissionstatus

import (
	"2f-authorization/internal/constants/dbinstance"
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
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

type updatePermissionStatusTest struct {
	test.TestInstance
	apiTest             src.ApiTest
	service             dto.CreateService
	createdService      dto.CreateServiceResponse
	domain              dto.Domain
	tenant              string
	CreatePermission    dto.RegisterTenantPermission
	permission          dto.CreatePermission
	createdPermissionId uuid.UUID
	permissionStatus    dto.UpdatePermissionStatus
}

func TestUpdatePermissionStatus(t *testing.T) {
	u := &updatePermissionStatusTest{}
	u.TestInstance = test.Initiate(context.Background(), "../../../../")
	u.apiTest.InitializeServer(u.Server)
	u.apiTest.InitializeTest(t, "update permission status test", "feature/update_permission_status.feature", u.InitializeScenario)
}
func (u *updatePermissionStatusTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {

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

func (u *updatePermissionStatusTest) iHaveServiceWith(service *godog.Table) error {
	body, err := u.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = u.apiTest.UnmarshalJSON([]byte(body), &u.service); err != nil {
		return err
	}
	u.service.Password = "123456"

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

func (u *updatePermissionStatusTest) iHaveTheFollowingPermissionInTenant(tenant string, permissions *godog.Table) error {
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
	u.apiTest.UnmarshalJSON([]byte(body), &u.permission)

	statement, _ := u.permission.Statement.Value()

	result, err := u.DB.TenantRegisterPermission(context.Background(), db.TenantRegisterPermissionParams{
		Name:        u.permission.Name,
		Description: u.permission.Description,
		ServiceID:   u.createdService.ServiceID,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
		TenantName: u.tenant,
	})
	if err != nil {
		return err
	}
	u.createdPermissionId = result.ID

	return nil
}

func (u *updatePermissionStatusTest) iSendTheRequestToUpdateTheStatus() error {
	u.apiTest.SetHeader("Authorization", "Basic "+u.BasicAuth(u.createdService.ServiceID.String(), "123456"))
	u.apiTest.SetHeader("x-subject", u.service.UserId)
	u.apiTest.SetHeader("x-action", "*")
	u.apiTest.SetHeader("x-resource", "*")
	u.apiTest.SetHeader("x-tenant", u.tenant)

	u.apiTest.SendRequest()
	return nil
}

func (u *updatePermissionStatusTest) iWantToUpdateThePermissionsStatusTo(status string) error {
	u.permissionStatus.Status = status
	u.apiTest.URL = "/v1/permissions/" + u.createdPermissionId.String() + "/status"
	body, err := json.Marshal(&u.permissionStatus)
	if err != nil {
		return err
	}
	u.apiTest.Body = string(body)
	return nil
}

func (u *updatePermissionStatusTest) thePermissionStatusShouldUpdateTo(status string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	permission, err := u.DB.GetPermissionDetails(context.Background(), dbinstance.GetPermissionDetailsParams{
		TenantName: u.tenant,
		ID:         u.createdPermissionId,
		ServiceID:  u.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	if err := u.apiTest.AssertEqual(permission.Status, status); err != nil {
		return err
	}

	return nil
}
func (u *updatePermissionStatusTest) thePermissionIsNotOnTheSystem(id string) error {
	u.apiTest.URL = "/v1/permissions/" + id + "/status"
	return nil
}

func (u *updatePermissionStatusTest) thenIShouldGetAnErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}
func (u *updatePermissionStatusTest) thenIShouldGetAFieldErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}
func (u *updatePermissionStatusTest) InitializeScenario(ctx *godog.ScenarioContext) {
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
	ctx.Step(`^the permission is not on the system "([^"]*)"$`, u.thePermissionIsNotOnTheSystem)
	ctx.Step(`^Then I should get an error with message "([^"]*)"$`, u.thenIShouldGetAnErrorWithMessage)
	ctx.Step(`^A registered domain and tenant$`, u.aRegisteredDomainAndTenant)
	ctx.Step(`^I have service with$`, u.iHaveServiceWith)
	ctx.Step(`^i have the following permission in tenant "([^"]*)"$`, u.iHaveTheFollowingPermissionInTenant)
	ctx.Step(`^I send the request to update the status$`, u.iSendTheRequestToUpdateTheStatus)
	ctx.Step(`^I want to update the permission\'s status to "([^"]*)"$`, u.iWantToUpdateThePermissionsStatusTo)
	ctx.Step(`^the permission status should update to "([^"]*)"$`, u.thePermissionStatusShouldUpdateTo)
}
