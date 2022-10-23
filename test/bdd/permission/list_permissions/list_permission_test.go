package listpermissions

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/platform/argon"
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

type ListPermission struct {
	test.TestInstance
	apiTest             src.ApiTest
	service             dto.CreateService
	createdService      dto.CreateServiceResponse
	domain              dto.Domain
	tenant              string
	permission          dto.CreatePermission
	respPermission      dto.Permission
	createdPermissionId uuid.UUID
	result              struct {
		OK   bool `json:"ok"`
		Data []dto.Permission
	}
}

func TestCreateDomain(t *testing.T) {
	l := &ListPermission{}
	l.TestInstance = test.Initiate(context.Background(), "../../../../")
	l.apiTest.InitializeServer(l.Server)
	l.apiTest.InitializeTest(t, "list permission test", "feature/list_permissions.feature", l.InitializeScenario)
}
func (l *ListPermission) aPermissionsRegisteredOnTheDomain(permission *godog.Table) error {
	body, err := l.apiTest.ReadRow(permission, []src.Type{
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
	l.apiTest.UnmarshalJSON([]byte(body), &l.permission)
	statment, _ := l.permission.Statement.Value()
	var result uuid.UUID
	result, err = l.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
		Name:        l.permission.Name,
		Description: l.permission.Description,
		ServiceID:   l.createdService.ServiceID,
		Statment: pgtype.JSON{
			Bytes:  statment,
			Status: pgtype.Present,
		},
		Column5: []uuid.UUID{
			l.domain.ID,
		},
	})
	if err != nil {
		return err
	}
	l.createdPermissionId = result
	return nil
}

func (l *ListPermission) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
	domain, err := l.apiTest.ReadCellString(domainAndTenant, "domain")
	if err != nil {
		return err
	}
	result, err := l.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      domain,
		ServiceID: l.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	l.domain.ID = result.ID

	tenant, err := l.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	if err = l.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		ServiceID:  l.createdService.ServiceID,
		DomainID:   l.domain.ID,
	}); err != nil {
		return err
	}
	l.tenant = tenant

	return nil
}

func (l *ListPermission) iHaveServiceWith(service *godog.Table) error {
	body, err := l.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = l.apiTest.UnmarshalJSON([]byte(body), &l.service); err != nil {
		return err
	}
	if l.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	createdService, err := l.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     l.service.Name,
		Password: l.service.Password,
		UserID:   uuid.MustParse(l.service.UserId),
	})
	if err != nil {
		return err
	}
	l.createdService.ServiceID = createdService.ServiceID
	if _, err := l.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", l.createdService.ServiceID); err != nil {
		return err
	}

	if err := l.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", l.service.Name)); err != nil {
		return err
	}
	return nil
}

func (l *ListPermission) iRequestToGetAllPermissionsUnderMyTenant() error {
	l.apiTest.SetHeader("Authorization", "Basic "+l.BasicAuth(l.createdService.ServiceID.String(), "123456"))
	l.apiTest.SetHeader("x-subject", l.service.UserId)
	l.apiTest.SetHeader("x-action", "*")
	l.apiTest.SetHeader("x-resource", "*")
	l.apiTest.SetHeader("x-tenant", l.tenant)

	l.apiTest.SendRequest()
	return nil
}

func (l *ListPermission) iShouldGetAllPermissionsInMyTenant(permissions *godog.Table) error {
	if err := l.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	body, err := l.apiTest.ReadRow(permissions, []src.Type{
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
	if err = l.apiTest.UnmarshalJSON([]byte(body), &l.respPermission); err != nil {
		return err
	}

	if err = l.apiTest.UnmarshalJSON(l.apiTest.ResponseBody, &l.result); err != nil {
		return err
	}
	fmt.Println("here", l.result)
	for _, p := range l.result.Data {
		if err := l.apiTest.AssertEqual(p.Name, l.respPermission.Name); err != nil {
			return err
		}
		if err := l.apiTest.AssertEqual(p.Description, l.respPermission.Description); err != nil {
			return err
		}
		if err := l.apiTest.AssertEqual(p.Statement.Action, l.respPermission.Statement.Action); err != nil {
			return err
		}
		if err := l.apiTest.AssertEqual(p.Statement.Effect, l.respPermission.Statement.Effect); err != nil {
			return err
		}
		if err := l.apiTest.AssertEqual(p.Statement.Resource, l.respPermission.Statement.Resource); err != nil {
			return err
		}

	}

	return nil
}

func (l *ListPermission) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		l.apiTest.URL = "/v1/permissions"
		l.apiTest.Method = http.MethodGet
		l.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = l.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^a permissions registered on the domain$`, l.aPermissionsRegisteredOnTheDomain)
	ctx.Step(`^A registered domain and tenant$`, l.aRegisteredDomainAndTenant)
	ctx.Step(`^I have service with$`, l.iHaveServiceWith)
	ctx.Step(`^I request to get all permissions under my tenant$`, l.iRequestToGetAllPermissionsUnderMyTenant)
	ctx.Step(`^I should get all permissions in my tenant$`, l.iShouldGetAllPermissionsInMyTenant)
}
