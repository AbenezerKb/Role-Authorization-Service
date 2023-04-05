package role

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

type UserIdRoleName struct {
	RoleName string    `json:"role_name"`
	UserId   uuid.UUID `json:"user_id"`
}
type assignRoleTest struct {
	test.TestInstance
	apiTest               src.ApiTest
	service               dto.CreateService
	createdService        dto.CreateServiceResponse
	domain                dto.Domain
	createdUser           dto.RegisterUser
	tenant                string
	createRole            dto.CreateRole
	permission            dto.CreatePermission
	createdRoleResponseId uuid.UUID
	createdPermissionId   uuid.UUID
	result                struct {
		OK   bool     `json:"ok"`
		Data dto.Role `json:"data"`
	}
}

func TestAssignRole(t *testing.T) {
	c := &assignRoleTest{}
	c.TestInstance = test.Initiate(context.Background(), "../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "assign role test", "feature/assign_role.feature", c.InitializeScenario)
}

func (r *assignRoleTest) aPermissionsRegisteredOnTheDomain(permission *godog.Table) error {
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
		{
			Column: "domains",
			Kind:   src.Array,
		},
	},
		true)
	if err != nil {
		return err
	}
	r.apiTest.UnmarshalJSON([]byte(body), &r.permission)

	statement, _ := r.permission.Statement.Value()
	var result uuid.UUID
	result, err = r.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
		Name:        r.permission.Name,
		Description: r.permission.Description,
		ServiceID:   r.createdService.ServiceID,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
	})
	if err != nil {
		return err
	}
	r.createdPermissionId = result
	return nil
}

func (r *assignRoleTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
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
func (r *assignRoleTest) iHaveServiceWith(service *godog.Table) error {
	body, err := r.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = r.apiTest.UnmarshalJSON([]byte(body), &r.service); err != nil {
		return err
	}
	r.service.Password = "123456"

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
func (r *assignRoleTest) iHaveUser(user_id *godog.Table) error {

	body, err := r.apiTest.ReadRow(user_id, nil, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body), &r.createdUser)
	if err != nil {
		return err
	}

	r.createdUser.ServiceID = r.createdService.ServiceID

	err = r.DB.RegisterUser(context.Background(), db.RegisterUserParams{
		UserID:    r.createdUser.UserId,
		ServiceID: r.createdUser.ServiceID,
	})
	if err != nil {
		return err
	}
	return nil
}
func (r *assignRoleTest) iHaveRole(role *godog.Table) error {
	body, err := r.apiTest.ReadRow(role, nil, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body), &r.createRole)
	if err != nil {
		return err
	}

	r.createRole.ServiceID = r.createdService.ServiceID
	r.createRole.PermissionID = []uuid.UUID{r.createdPermissionId}

	createdRoleResponse, err := r.DB.CreateRole(context.Background(), db.CreateRoleParams{
		TenantName: r.tenant,
		ServiceID:  r.createdService.ServiceID,
		Column4:    r.createRole.PermissionID,
		Name:       r.createRole.Name,
	})
	if err != nil {
		return err
	}
	r.createdRoleResponseId = createdRoleResponse.RoleID
	return nil
}

func (r *assignRoleTest) iRequestToAssignRoleToUser(user_id_role_name *godog.Table) error {
	user_id, err := r.apiTest.ReadCellString(user_id_role_name, "user_id")
	if err != nil {
		return err
	}

	r.apiTest.URL = fmt.Sprintf("/v1/roles/%s/users/%s", r.createdRoleResponseId, user_id)
	r.apiTest.SetHeader("Authorization", "Basic "+r.BasicAuth(r.createdService.ServiceID.String(), "123456"))
	r.apiTest.SetHeader("x-subject", r.service.UserId)
	r.apiTest.SetHeader("x-action", "*")
	r.apiTest.SetHeader("x-resource", "*")
	r.apiTest.SetHeader("x-tenant", "administrator")

	r.apiTest.SendRequest()
	return nil
}

func (r *assignRoleTest) iRequestToAssignRoleToUserWhileFieldsAreMissing(userIdRoleName *godog.Table) error {

	useridrolename := UserIdRoleName{}
	body, err := r.apiTest.ReadRow(userIdRoleName, nil, false)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(body), &useridrolename)

	requestRoleId := uuid.Nil
	requestUserId := uuid.Nil
	if useridrolename.RoleName != "" {
		requestRoleId = r.createdRoleResponseId
	}
	if useridrolename.UserId != uuid.MustParse("00000000-0000-0000-0000-000000000000") {
		requestUserId = r.createdUser.UserId
	}
	r.apiTest.URL = fmt.Sprintf("/v1/roles/%s/users/%s", requestRoleId, requestUserId)
	r.apiTest.SetHeader("Authorization", "Basic "+r.BasicAuth(r.createdService.ServiceID.String(), "123456"))
	r.apiTest.SetHeader("x-subject", r.service.UserId)
	r.apiTest.SetHeader("x-action", "*")
	r.apiTest.SetHeader("x-resource", "*")
	r.apiTest.SetHeader("x-tenant", "administrator")

	r.apiTest.SendRequest()
	return nil
}

func (r *assignRoleTest) myRequestShouldFailWith(message string) error {

	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (r *assignRoleTest) theRoleShouldSuccessfullyBeAssigned() error {
	if err := r.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("ok", "true"); err != nil {
		return err
	}
	return nil
}

func (r *assignRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.Method = http.MethodPost
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^a permissions registered on the domain$`, r.aPermissionsRegisteredOnTheDomain)
	ctx.Step(`^A registered domain and tenant$`, r.aRegisteredDomainAndTenant)
	ctx.Step(`^i have role$`, r.iHaveRole)
	ctx.Step(`^I have service with$`, r.iHaveServiceWith)
	ctx.Step(`^I have user$`, r.iHaveUser)
	ctx.Step(`^I request to  assign role to user$`, r.iRequestToAssignRoleToUser)
	ctx.Step(`^I request to assign  role to user while fields are missing$`, r.iRequestToAssignRoleToUserWhileFieldsAreMissing)
	ctx.Step(`^my request should fail with "([^"]*)"$`, r.myRequestShouldFailWith)
	ctx.Step(`^the role should successfully be  assigned$`, r.theRoleShouldSuccessfullyBeAssigned)
}
