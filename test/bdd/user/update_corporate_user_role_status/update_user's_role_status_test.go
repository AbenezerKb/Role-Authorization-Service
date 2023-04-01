package updateuserrolestatus

import (
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

type updateUserRoleStatusTest struct {
	test.TestInstance
	apiTest               src.ApiTest
	OK                    bool `json:"ok"`
	service               db.CreateServiceParams
	userRoleStatus        dto.UpdateUserRoleStatus
	createdService        db.CreateServiceRow
	createdUser           dto.RegisterUser
	Tur                   dto.TenantUsersRole
	domain                dto.Domain
	tenant                string
	createRole            dto.CreateRole
	permission            dto.CreatePermission
	createdRoleResponseId uuid.UUID
	createdPermissionId   uuid.UUID
	Name                  dto.Role
}

func TestUpdateCorporateUserRoleStatus(t *testing.T) {
	u := &updateUserRoleStatusTest{}
	u.TestInstance = test.Initiate(context.Background(), "../../../../")
	u.apiTest.InitializeServer(u.Server)
	u.apiTest.InitializeTest(t, "update user's role status test", "feature/update_user_role_status.feature", u.InitializeScenario)
}
func (u *updateUserRoleStatusTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {

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

func (u *updateUserRoleStatusTest) iHaveServiceWith(service *godog.Table) error {
	body, err := u.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = u.apiTest.UnmarshalJSON([]byte(body), &u.service); err != nil {
		return err
	}

	u.service.Password = "123456"
	if u.createdService, err = u.DB.CreateService(context.Background(), u.service); err != nil {
		return err
	}

	if _, err := u.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", u.createdService.ServiceID); err != nil {
		return err
	}

	if err := u.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", u.service.Name)); err != nil {
		return err
	}

	return nil
}

func (u *updateUserRoleStatusTest) iSendTheRequestToUpdateTheStatus() error {
	u.apiTest.SetHeader("Authorization", "Basic "+u.BasicAuth(u.createdService.ServiceID.String(), "123456"))
	u.apiTest.SetHeader("x-subject", u.service.UserID.String())
	u.apiTest.SetHeader("x-action", "*")
	u.apiTest.SetHeader("x-resource", "*")
	u.apiTest.SetHeader("x-tenant", u.tenant)
	u.apiTest.SendRequest()
	return nil
}

func (u *updateUserRoleStatusTest) iWantToUpdateTheUsersRoleStatusTo(status string) error {
	u.apiTest.URL = "/v1/tenants/corporate/" + u.tenant + "/users/" + u.createdUser.UserId.String() + "/roles/" + u.createdRoleResponseId.String() + "/status"
	fmt.Println("the URL: ", u.apiTest.URL)
	u.apiTest.Method = http.MethodPatch
	u.userRoleStatus.Status = status
	body, err := json.Marshal(&u.userRoleStatus)
	if err != nil {
		return err
	}
	u.apiTest.Body = string(body)
	return nil
}

func (u *updateUserRoleStatusTest) theRoleStatusShouldUpdateTo(status string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	return nil
}

func (u *updateUserRoleStatusTest) theUserHasTheFollowingRoleInTheFollowingTenant(userRoleTenant *godog.Table) error {
	body, err := u.apiTest.ReadRow(userRoleTenant, nil, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body), &u.createdUser)
	if err != nil {
		return err
	}
	err = u.DB.RegisterUser(context.Background(), db.RegisterUserParams{
		UserID:    u.createdUser.UserId,
		ServiceID: u.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	tenant, err := u.apiTest.ReadCellString(userRoleTenant, "tenant")
	if err != nil {
		return err
	}
	u.tenant = tenant

	if err := u.DB.AssignRole(context.Background(), db.AssignRoleParams{
		UserID:     u.createdUser.UserId,
		TenantName: u.tenant,
		ID:         u.createdRoleResponseId,
		ServiceID:  u.createdService.ServiceID,
	}); err != nil {
		return err
	}

	return nil
}

func (u *updateUserRoleStatusTest) aRoleInTenantWithTheFollowingPermissions(role, tenant string, permissions *godog.Table) error {
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
	},
		true)
	if err != nil {
		return err
	}
	u.apiTest.UnmarshalJSON([]byte(body), &u.permission)

	statement, _ := u.permission.Statement.Value()

	u.createdPermissionId, err = u.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
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

	_, err = u.DB.CreateRole(context.Background(), db.CreateRoleParams{
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
	adminRole, err := u.DB.GetRoleByNameAndTenantName(context.Background(), db.GetRoleByNameAndTenantNameParams{
		TenantName: u.tenant,
		Name:       "admin",
	})
	if err != nil {
		return err
	}
	u.createdRoleResponseId = adminRole
	return nil
}
func (u *updateUserRoleStatusTest) thenIShouldGetAFieldErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (u *updateUserRoleStatusTest) theUserHasACTIVEAdminRoleInTheFollowingTenant(userRoleTenant *godog.Table) error {
	body, err := u.apiTest.ReadRow(userRoleTenant, nil, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body), &u.createdUser)
	if err != nil {
		return err
	}
	err = u.DB.RegisterUser(context.Background(), db.RegisterUserParams{
		UserID:    u.createdUser.UserId,
		ServiceID: u.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	tenant, err := u.apiTest.ReadCellString(userRoleTenant, "tenant")
	if err != nil {
		return err
	}
	u.tenant = tenant
	role, err := u.DB.GetRoleByNameAndTenantName(context.Background(), db.GetRoleByNameAndTenantNameParams{
		Name:       "admin",
		TenantName: tenant,
	})
	if err != nil {
		return err
	}
	if err := u.DB.AssignRole(context.Background(), db.AssignRoleParams{
		UserID:     u.createdUser.UserId,
		TenantName: u.tenant,
		ID:         role,
		ServiceID:  u.createdService.ServiceID,
	}); err != nil {
		return err
	}
	u.createdRoleResponseId = role
	return nil
}

func (u *updateUserRoleStatusTest) theRoleStatusShouldFailToUpdateWith(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusInternalServerError); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}

func (u *updateUserRoleStatusTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		u.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^Then I should get a field error with message "([^"]*)"$`, u.thenIShouldGetAFieldErrorWithMessage)
	ctx.Step(`^a role "([^"]*)" in tenant "([^"]*)" with the following permissions$`, u.aRoleInTenantWithTheFollowingPermissions)
	ctx.Step(`^A registered domain and tenant$`, u.aRegisteredDomainAndTenant)
	ctx.Step(`^I have service with$`, u.iHaveServiceWith)
	ctx.Step(`^I send the request to update the status$`, u.iSendTheRequestToUpdateTheStatus)
	ctx.Step(`^I want to update the user\'s role status to "([^"]*)"$`, u.iWantToUpdateTheUsersRoleStatusTo)
	ctx.Step(`^the role status should fail to update with "([^"]*)"$`, u.theRoleStatusShouldFailToUpdateWith)
	ctx.Step(`^the role status should update to "([^"]*)"$`, u.theRoleStatusShouldUpdateTo)
	ctx.Step(`^the user has ACTIVE admin role in the following tenant$`, u.theUserHasACTIVEAdminRoleInTheFollowingTenant)
	ctx.Step(`^the user has admin role in the following tenant$`, u.theUserHasTheFollowingRoleInTheFollowingTenant)
	ctx.Step(`^the user has the following role in the following tenant$`, u.theUserHasTheFollowingRoleInTheFollowingTenant)
}
