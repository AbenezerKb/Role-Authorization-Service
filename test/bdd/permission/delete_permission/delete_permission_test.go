package deletepermission

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

type deletePermissionTest struct {
	test.TestInstance
	apiTest             src.ApiTest
	service             dto.CreateService
	createdService      dto.CreateServiceResponse
	domain              dto.Domain
	tenant              string
	permission          dto.CreatePermission
	createdPermissionId uuid.UUID
	// result              struct {
	// 	OK   bool `json:"ok"`
	// 	Data []dto.Permission
	// }
}

func TestDeletePermissionTest(t *testing.T) {
	d := &deletePermissionTest{}
	d.TestInstance = test.Initiate(context.Background(), "../../../../")
	d.apiTest.InitializeServer(d.Server)
	d.apiTest.InitializeTest(t, "delete permission test", "feature/delete_permission.feature", d.InitializeScenario)
}
func (d *deletePermissionTest) aPermissionRegisteredOnTheTenant(permission *godog.Table) error {
	body, err := d.apiTest.ReadRow(permission, []src.Type{
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
	if err := d.apiTest.UnmarshalJSON([]byte(body), &d.permission); err != nil {
		return err
	}
	statement, _ := d.permission.Statement.Value()
	result, err := d.DB.TenantRegisterPermission(context.Background(), db.TenantRegisterPermissionParams{
		Name:        d.permission.Name,
		Description: d.permission.Description,
		ServiceID:   d.createdService.ServiceID,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
		TenantName: d.tenant,
	})
	if err != nil {
		return err
	}
	d.createdPermissionId = result.ID

	return nil
}

func (d *deletePermissionTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
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

	if err = d.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		ServiceID:  d.createdService.ServiceID,
		DomainID:   d.domain.ID,
	}); err != nil {
		return err
	}
	d.tenant = tenant

	return nil
}

func (d *deletePermissionTest) iHaveARegisteredService(service *godog.Table) error {
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
	if _, err := d.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", d.createdService.ServiceID); err != nil {
		return err
	}

	if err := d.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", d.service.Name)); err != nil {
		return err
	}
	return nil
}
func (d *deletePermissionTest) iSendARequestToDeleteThePermission() error {
	d.apiTest.SetHeader("Authorization", "Basic "+d.BasicAuth(d.createdService.ServiceID.String(), "123456"))
	d.apiTest.SetHeader("x-subject", d.service.UserId)
	d.apiTest.SetHeader("x-action", "*")
	d.apiTest.SetHeader("x-resource", "*")
	d.apiTest.SetHeader("x-tenant", d.tenant)

	d.apiTest.SendRequest()
	return nil
}

func (d *deletePermissionTest) iWantToDeleteThePermission() error {
	d.apiTest.URL = "/v1/permissions/" + d.createdPermissionId.String()
	d.apiTest.Method = http.MethodDelete
	return nil
}

func (d *deletePermissionTest) thePermissionShouldBeDeleted() error {
	count, err := d.DB.CheckIfPermissionExistsInTenant(context.Background(), db.CheckIfPermissionExistsInTenantParams{
		TenantName: d.tenant,
		ServiceID:  d.createdService.ServiceID,
		Name:       d.permission.Name,
	})
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("permission should be deleted")
	}
	return nil
}

func (d *deletePermissionTest) theRequestShouldBeSuccessfull() error {
	if err := d.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	return nil
}

func (d *deletePermissionTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		d.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = d.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^A permission registered on the tenant$`, d.aPermissionRegisteredOnTheTenant)
	ctx.Step(`^A registered domain and tenant$`, d.aRegisteredDomainAndTenant)
	ctx.Step(`^I have a registered service$`, d.iHaveARegisteredService)
	ctx.Step(`^I send a request to delete the permission$`, d.iSendARequestToDeleteThePermission)
	ctx.Step(`^I want to delete the permission$`, d.iWantToDeleteThePermission)
	ctx.Step(`^the permission should be deleted$`, d.thePermissionShouldBeDeleted)
	ctx.Step(`^The request should be successfull$`, d.theRequestShouldBeSuccessfull)
}
