package deleterole

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

type deleteRoleTest struct {
	test.TestInstance
	apiTest               src.ApiTest
	service               dto.CreateService
	createdService        dto.CreateServiceResponse
	domain                dto.Domain
	createdUser           dto.RegisterUser
	tenant                string
	assignedRole          dto.TenantUsersRole
	createRole            dto.CreateRole
	permission            dto.CreatePermission
	createdRoleResponseId uuid.UUID
	createdPermissionId   uuid.UUID
	result                struct {
		OK   bool     `json:"ok"`
		Data dto.Role `json:"data"`
	}
}

func TestDeleteRole(t *testing.T) {
	d := &deleteRoleTest{}
	d.TestInstance = test.Initiate(context.Background(), "../../../../")
	d.apiTest.InitializeServer(d.Server)
	d.apiTest.InitializeTest(t, "delete role test", "feature/delete_role.feature", d.InitializeScenario)
}

func (d *deleteRoleTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {

	domain, err := d.apiTest.ReadCellString(domainAndTenant, "domain")
	if err != nil {
		return err
	}
	result, err := d.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      domain,
		ServiceID: d.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	d.domain.ID = result.ID

	tenant, err := d.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	err = d.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		DomainID:   result.ID,
		ServiceID:  d.createdService.ServiceID,
	})

	if err != nil {
		return err
	}
	d.tenant = tenant

	return nil

}

func (d *deleteRoleTest) iHaveARoleInTenantWithThePermissionsBelow(role, tenant string, permissions *godog.Table) error {
	body, err := d.apiTest.ReadRow(permissions, []src.Type{
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
	d.apiTest.UnmarshalJSON([]byte(body), &d.permission)

	statement, _ := d.permission.Statement.Value()
	var result uuid.UUID
	result, err = d.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
		Name:        d.permission.Name,
		Description: d.permission.Description,
		ServiceID:   d.createdService.ServiceID,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
	})
	if err != nil {
		return err
	}
	d.createdPermissionId = result

	res, err := d.DB.CreateRole(context.Background(), db.CreateRoleParams{
		TenantName: tenant,
		ServiceID:  d.createdService.ServiceID,
		Column4: []uuid.UUID{
			d.createdPermissionId,
		},
		Name: d.createRole.Name,
	})
	if err != nil {
		return err
	}
	d.createdRoleResponseId = res.RoleID
	return nil
}

func (d *deleteRoleTest) iHaveServiceWith(service *godog.Table) error {
	body, err := d.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = d.apiTest.UnmarshalJSON([]byte(body), &d.service); err != nil {
		return err
	}
	d.service.Password = "123456"

	createdService, err := d.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     d.service.Name,
		Password: d.service.Password,
		UserID:   uuid.MustParse(d.service.UserId),
	})
	if err != nil {
		return err
	}
	d.createdService.ServiceID = createdService.ServiceID
	if _, err := d.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE' where id = $1", d.createdService.ServiceID); err != nil {
		return err
	}

	if err := d.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", d.service.Name)); err != nil {
		return err
	}
	return nil

}

func (d *deleteRoleTest) iSendARequestToDeletedTheRole() error {
	d.apiTest.SetHeader("Authorization", "Basic "+d.BasicAuth(d.createdService.ServiceID.String(), "123456"))
	d.apiTest.SetHeader("x-subject", d.service.UserId)
	d.apiTest.SetHeader("x-action", "*")
	d.apiTest.SetHeader("x-resource", "*")
	d.apiTest.SetHeader("x-tenant", d.tenant)

	d.apiTest.SendRequest()
	return nil
}

func (d *deleteRoleTest) iWantToDeleteTheRole(role string) error {
	d.apiTest.URL = "/v1/roles/" + d.createdRoleResponseId.String()
	return nil
}

func (d *deleteRoleTest) theRequestShouldFailWithError(message string) error {
	if err := d.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := d.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}

func (d *deleteRoleTest) theRoleDoesNotExistsInTheSystem() error {
	d.apiTest.URL = "/v1/roles/8b68aba9-1a3b-475a-a0f8-93bc5bea5f7b"
	return nil
}

func (d *deleteRoleTest) theRoleShouldBeDelete() error {
	if err := d.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	return nil
}
func (d *deleteRoleTest) theRoleShouldNotBeFoundInTheTenant() error {
	d.apiTest.SetHeader("Authorization", "Basic "+d.BasicAuth(d.createdService.ServiceID.String(), "123456"))
	d.apiTest.SetHeader("x-subject", d.service.UserId)
	d.apiTest.SetHeader("x-action", "*")
	d.apiTest.SetHeader("x-resource", "*")
	d.apiTest.SetHeader("x-tenant", d.tenant)

	d.apiTest.SendRequest()
	if err := d.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return nil
}

func (d *deleteRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		d.apiTest.Method = http.MethodDelete
		d.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = d.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})

	ctx.Step(`^A registered domain and tenant$`, d.aRegisteredDomainAndTenant)
	ctx.Step(`^i have a role "([^"]*)" in tenant "([^"]*)" with the permissions below$`, d.iHaveARoleInTenantWithThePermissionsBelow)
	ctx.Step(`^I have service with$`, d.iHaveServiceWith)
	ctx.Step(`^i send a request to deleted the role$`, d.iSendARequestToDeletedTheRole)
	ctx.Step(`^I  want to delete the role "([^"]*)"$`, d.iWantToDeleteTheRole)
	ctx.Step(`^the request should fail with error "([^"]*)"$`, d.theRequestShouldFailWithError)
	ctx.Step(`^the role does not exists in the system$`, d.theRoleDoesNotExistsInTheSystem)
	ctx.Step(`^the role should be delete$`, d.theRoleShouldBeDelete)
	ctx.Step(`^the role should not be found in the tenant$`, d.theRoleShouldNotBeFoundInTheTenant)
}
