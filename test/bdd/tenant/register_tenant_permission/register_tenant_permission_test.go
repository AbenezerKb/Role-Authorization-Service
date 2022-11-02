package registertenantpermission

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

type registerTenantPermission struct {
	test.TestInstance
	apiTest        src.ApiTest
	service        dto.CreateService
	createdService dto.CreateServiceResponse
	domain         dto.Domain
	tenant         string
	permission     dto.CreatePermission
	result         struct {
		OK   bool             `json:"ok"`
		Data []dto.Permission `json:"data"`
	}
}

func TestRegisterTenantPermission(t *testing.T) {
	r := &registerTenantPermission{}
	r.TestInstance = test.Initiate(context.Background(), "../../../../")
	r.apiTest.InitializeServer(r.Server)
	r.apiTest.InitializeTest(t, "register tenant permission test", "feature/register_tenant_permission.feature", r.InitializeScenario)
}
func (r *registerTenantPermission) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
	domain, err := r.apiTest.ReadCellString(domainAndTenant, "domain")
	if err != nil {
		return err
	}
	result, err := r.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      domain,
		ServiceID: r.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	r.domain.ID = result.ID

	tenant, err := r.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	err = r.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		DomainID:   result.ID,
		ServiceID:  r.createdService.ServiceID,
	})

	if err != nil {
		return err
	}
	r.tenant = tenant

	return nil
}

func (r *registerTenantPermission) iHaveServiceWith(service *godog.Table) error {
	body, err := r.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = r.apiTest.UnmarshalJSON([]byte(body), &r.service); err != nil {
		return err
	}
	if r.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	createdService, err := r.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     r.service.Name,
		Password: r.service.Password,
		UserID:   uuid.MustParse(r.service.UserId),
	})
	if err != nil {
		return err
	}
	r.createdService.ServiceID = createdService.ServiceID
	if _, err := r.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE' where id = $1", r.createdService.ServiceID); err != nil {
		return err
	}

	if err := r.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", r.service.Name)); err != nil {
		return err
	}
	return nil
}

func (r *registerTenantPermission) iSendRequestToAddThePermission() error {
	r.apiTest.SendRequest()
	return nil
}

func (r *registerTenantPermission) iWantToRegisterTheFollowingPermissionUnderTheTenant(permission *godog.Table) error {
	body, err := r.apiTest.ReadRow(permission, []src.Type{
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
	json.Unmarshal([]byte(body), &r.permission)

	r.permission.Domain = []uuid.UUID{r.domain.ID}
	requestBody, err := json.Marshal(r.permission)
	if err != nil {
		return err
	}

	r.apiTest.Body = string(requestBody)
	r.apiTest.SetHeader("Authorization", "Basic "+r.BasicAuth(r.createdService.ServiceID.String(), "123456"))
	r.apiTest.SetHeader("x-subject", r.service.UserId)
	r.apiTest.SetHeader("x-action", "*")
	r.apiTest.SetHeader("x-resource", "*")
	r.apiTest.SetHeader("x-tenant", r.tenant)

	return nil
}

func (r *registerTenantPermission) thePermissionShouldBeAccessibleThroughTheTenant() error {

	r.apiTest.URL = "/v1/permissions"
	r.apiTest.Method = http.MethodGet
	r.apiTest.SendRequest()

	if err := r.apiTest.UnmarshalJSON(r.apiTest.ResponseBody, &r.result); err != nil {
		return err
	}

	if err := r.apiTest.AssertEqual(r.result.OK, true); err != nil {
		return err
	}

	for _, p := range r.result.Data {
		if err := r.apiTest.AssertEqual(p.Name, r.permission.Name); err != nil {
			return err
		}
		if err := r.apiTest.AssertEqual(p.Description, r.permission.Description); err != nil {
			return err
		}
		if err := r.apiTest.AssertEqual(p.Statement.Action, r.permission.Statement.Action); err != nil {
			return err
		}
		if err := r.apiTest.AssertEqual(p.Statement.Effect, r.permission.Statement.Effect); err != nil {
			return err
		}
		if err := r.apiTest.AssertEqual(p.Statement.Resource, r.permission.Statement.Resource); err != nil {
			return err
		}
	}
	return nil
}

func (r *registerTenantPermission) theRequestShouldSuccessfull() error {
	if err := r.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}

	return nil
}

func (r *registerTenantPermission) theRequestShouldFailWithError(message string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (r *registerTenantPermission) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.URL = "/v1/tenants/permissions"
		r.apiTest.Method = http.MethodPost
		r.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^A registered domain and tenant$`, r.aRegisteredDomainAndTenant)
	ctx.Step(`^I have service with$`, r.iHaveServiceWith)
	ctx.Step(`^The request should fail with error "([^"]*)"$`, r.theRequestShouldFailWithError)

	ctx.Step(`^I send request to add the permission$`, r.iSendRequestToAddThePermission)
	ctx.Step(`^I want to register the following permission under the tenant$`, r.iWantToRegisterTheFollowingPermissionUnderTheTenant)
	ctx.Step(`^the permission should be accessible through the tenant$`, r.thePermissionShouldBeAccessibleThroughTheTenant)
	ctx.Step(`^the request should successfull$`, r.theRequestShouldSuccessfull)
}
