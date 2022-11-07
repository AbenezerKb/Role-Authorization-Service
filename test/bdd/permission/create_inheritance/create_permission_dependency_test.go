package createinheritance

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

type createPermissionDependencyTest struct {
	test.TestInstance
	apiTest        src.ApiTest
	service        dto.CreateService
	createdService dto.CreateServiceResponse
	domain         dto.Domain
	tenant         string
	permission     []dto.CreatePermission
	inheritance    []dto.CreatePermissionDependency
	result         struct {
		OK bool `json:"ok"`
	}
}

func TestCreatePermissionDependency(t *testing.T) {
	c := &createPermissionDependencyTest{}
	c.TestInstance = test.Initiate(context.Background(), "../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "create permission dependency test", "feature/create_permission_inheritance.feature", c.InitializeScenario)
}

func (c *createPermissionDependencyTest) aPermissionsRegisteredOnTheDomain(permissions *godog.Table) error {
	body, err := c.apiTest.ReadRows(permissions, []src.Type{
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

	c.apiTest.UnmarshalJSON([]byte(body), &c.permission)
	for _, p := range c.permission {
		statement, _ := p.Statement.Value()
		_, err := c.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
			Name:        p.Name,
			Description: p.Description,
			ServiceID:   c.createdService.ServiceID,
			Statement: pgtype.JSON{
				Bytes:  statement,
				Status: pgtype.Present,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *createPermissionDependencyTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
	domain, err := c.apiTest.ReadCellString(domainAndTenant, "domain")
	if err != nil {
		return err
	}
	result, err := c.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      domain,
		ServiceID: c.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	c.domain.ID = result.ID

	tenant, err := c.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	err = c.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		DomainID:   result.ID,
		ServiceID:  c.createdService.ServiceID,
	})

	if err != nil {
		return err
	}
	c.tenant = tenant

	return nil
}

func (c *createPermissionDependencyTest) iHaveServiceWith(service *godog.Table) error {
	body, err := c.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = c.apiTest.UnmarshalJSON([]byte(body), &c.service); err != nil {
		return err
	}
	if c.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	createdService, err := c.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     c.service.Name,
		Password: c.service.Password,
		UserID:   uuid.MustParse(c.service.UserId),
	})
	if err != nil {
		return err
	}
	c.createdService.ServiceID = createdService.ServiceID
	if _, err := c.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE' where id = $1", c.createdService.ServiceID); err != nil {
		return err
	}

	if err := c.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", c.service.Name)); err != nil {
		return err
	}
	return nil
}

func (c *createPermissionDependencyTest) iSendRequestToCreateTheInheritance() error {
	c.apiTest.SetHeader("Authorization", "Basic "+c.BasicAuth(c.createdService.ServiceID.String(), "123456"))
	c.apiTest.SendRequest()
	return nil
}

func (c *createPermissionDependencyTest) iWantToHaveARelationBetweenPermissionAsAParentAndPermissionAsAChild(parent, child string) error {
	c.inheritance = append(c.inheritance, dto.CreatePermissionDependency{
		PermissionName: parent,
		InheritedPermissions: []string{
			child,
		},
	})

	body, err := json.Marshal(c.inheritance)
	if err != nil {
		return err
	}
	c.apiTest.Body = string(body)
	return nil
}

func (c *createPermissionDependencyTest) theRequestShouldBeSuccessfull() error {
	if err := c.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	if err := c.apiTest.UnmarshalJSON(c.apiTest.ResponseBody, &c.result); err != nil {
		return err
	}

	if err := c.apiTest.AssertEqual(c.result.OK, true); err != nil {
		return err
	}
	return nil
}
func (c *createPermissionDependencyTest) iShouldGetErrorWithMessage(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return c.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (c *createPermissionDependencyTest) iWantToHaveARelationBetweenPermissions(permission *godog.Table) error {
	body, err := c.apiTest.ReadRows(permission, []src.Type{
		{
			Column: "permission",
			Kind:   src.String,
		},
		{
			Column: "inherited_permissions",
			Kind:   src.Array,
		},
	},
		true)
	if err != nil {
		return err
	}
	c.apiTest.Body = body
	return nil
}
func (c *createPermissionDependencyTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/permissions/inherit"
		c.apiTest.Method = http.MethodPost
		c.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = c.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^a permissions registered on the domain$`, c.aPermissionsRegisteredOnTheDomain)
	ctx.Step(`^A registered domain and tenant$`, c.aRegisteredDomainAndTenant)
	ctx.Step(`^I have service with$`, c.iHaveServiceWith)
	ctx.Step(`^I send request to create the inheritance$`, c.iSendRequestToCreateTheInheritance)
	ctx.Step(`^I want to have a relation between "([^"]*)" permission as a parent and "([^"]*)" permission as a child$`, c.iWantToHaveARelationBetweenPermissionAsAParentAndPermissionAsAChild)
	ctx.Step(`^the request should be successfull$`, c.theRequestShouldBeSuccessfull)
	ctx.Step(`^i should get error with message "([^"]*)"$`, c.iShouldGetErrorWithMessage)
	ctx.Step(`^I want to have a relation between permissions$`, c.iWantToHaveARelationBetweenPermissions)

}
