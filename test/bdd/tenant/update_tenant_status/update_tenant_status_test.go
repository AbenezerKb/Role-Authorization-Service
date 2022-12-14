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
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type updateTenantStatusTest struct {
	test.TestInstance
	apiTest        src.ApiTest
	service        dto.CreateService
	createdService dto.CreateServiceResponse
	domain         dto.Domain
	tenant         string
	tenantStatus   dto.UpdateTenantStatus
}

func TestUpdateTenantStatus(t *testing.T) {
	u := &updateTenantStatusTest{}
	u.TestInstance = test.Initiate(context.Background(), "../../../../")
	u.apiTest.InitializeServer(u.Server)
	u.apiTest.InitializeTest(t, "update tenant status test", "feature/update_tenant_status.feature", u.InitializeScenario)
}
func (u *updateTenantStatusTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {

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

func (u *updateTenantStatusTest) iHaveServiceWith(service *godog.Table) error {
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

func (u *updateTenantStatusTest) iSendTheRequestToUpdateTheStatus() error {
	u.apiTest.SetHeader("Authorization", "Basic "+u.BasicAuth(u.createdService.ServiceID.String(), "123456"))
	u.apiTest.SendRequest()
	return nil
}

func (u *updateTenantStatusTest) iWantToUpdateTheTenantsStatusTo(status string) error {
	u.tenantStatus.Status = status
	u.apiTest.URL = "/v1/tenants/" + u.tenant + "/status"
	u.tenant = ""
	body, err := json.Marshal(&u.tenantStatus)
	if err != nil {
		return err
	}
	u.apiTest.Body = string(body)
	return nil
}

func (u *updateTenantStatusTest) theTenantStatusShouldUpdateTo(status string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	return nil
}
func (u *updateTenantStatusTest) theTenantIsNotOnTheSystem(tenant string) error {
	u.tenant = tenant
	return nil
}

func (u *updateTenantStatusTest) thenIShouldGetAnErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}
func (u *updateTenantStatusTest) thenIShouldGetAFieldErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}
func (u *updateTenantStatusTest) InitializeScenario(ctx *godog.ScenarioContext) {
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
	ctx.Step(`^the tenant is not on the system "([^"]*)"$`, u.theTenantIsNotOnTheSystem)
	ctx.Step(`^Then I should get an error with message "([^"]*)"$`, u.thenIShouldGetAnErrorWithMessage)
	ctx.Step(`^A registered domain and tenant$`, u.aRegisteredDomainAndTenant)
	ctx.Step(`^I have service with$`, u.iHaveServiceWith)
	ctx.Step(`^I send the request to update the status$`, u.iSendTheRequestToUpdateTheStatus)
	ctx.Step(`^I want to update the tenant\'s status to "([^"]*)"$`, u.iWantToUpdateTheTenantsStatusTo)
	ctx.Step(`^the tenant status should update to "([^"]*)"$`, u.theTenantStatusShouldUpdateTo)
}
