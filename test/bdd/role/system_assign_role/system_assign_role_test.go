package systemassignrole

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

type userTenantRole struct {
	User       string `json:"user_id"`
	Role       string `json:"role_name"`
	TenantName string `json:"tenant_name"`
}

type systemAssignRoleTest struct {
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
	userTenantRole dto.TenantUsersRole
}

func TestSystemAssignRole(t *testing.T) {
	s := &systemAssignRoleTest{}
	s.TestInstance = test.Initiate(context.Background(), "../../../../")
	s.apiTest.InitializeServer(s.Server)
	s.apiTest.InitializeTest(t, "system assign role test", "feature/system_assign_role.feature", s.InitializeScenario)
}

func (s *systemAssignRoleTest) aPermissionsRegisteredOnTheDomain(permission *godog.Table) error {
	body, err := s.apiTest.ReadRow(permission, []src.Type{
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
	s.apiTest.UnmarshalJSON([]byte(body), &s.permission)

	statement, _ := s.permission.Statement.Value()
	var result uuid.UUID
	result, err = s.DB.CreateOrGetPermission(context.Background(), db.CreateOrGetPermissionParams{
		Name:        s.permission.Name,
		Description: s.permission.Description,
		ServiceID:   s.createdService.ServiceID,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
	})
	if err != nil {
		return err
	}
	s.createdPermissionId = result
	return nil
}

func (s *systemAssignRoleTest) aRegisteredDomainAndTenant(domainAndTenant *godog.Table) error {
	domain, err := s.apiTest.ReadCellString(domainAndTenant, "domain")
	if err != nil {
		return err
	}
	result, err := s.DB.CreateDomain(context.Background(), db.CreateDomainParams{
		Name:      domain,
		ServiceID: s.createdService.ServiceID,
	})
	if err != nil {
		return err
	}
	s.domain.ID = result.ID

	tenant, err := s.apiTest.ReadCellString(domainAndTenant, "tenant_name")
	if err != nil {
		return err
	}

	err = s.DB.CreateTenent(context.Background(), db.CreateTenentParams{
		TenantName: tenant,
		DomainID:   result.ID,
		ServiceID:  s.createdService.ServiceID,
	})

	if err != nil {
		return err
	}
	s.tenant = tenant

	return nil
}
func (s *systemAssignRoleTest) iHaveServiceWith(service *godog.Table) error {
	body, err := s.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = s.apiTest.UnmarshalJSON([]byte(body), &s.service); err != nil {
		return err
	}
	if s.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	createdService, err := s.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     s.service.Name,
		Password: s.service.Password,
		UserID:   uuid.MustParse(s.service.UserId),
	})
	if err != nil {
		return err
	}
	s.createdService.ServiceID = createdService.ServiceID
	if _, err := s.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE' where id = $1", s.createdService.ServiceID); err != nil {
		return err
	}

	if err := s.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", s.service.Name)); err != nil {
		return err
	}
	return nil

}
func (s *systemAssignRoleTest) iHaveUser(user_id *godog.Table) error {

	body, err := s.apiTest.ReadRow(user_id, nil, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body), &s.createdUser)
	if err != nil {
		return err
	}

	s.createdUser.ServiceID = s.createdService.ServiceID

	err = s.DB.RegisterUser(context.Background(), db.RegisterUserParams{
		UserID:    s.createdUser.UserId,
		ServiceID: s.createdUser.ServiceID,
	})
	if err != nil {
		return err
	}
	return nil
}
func (s *systemAssignRoleTest) iHaveRole(role *godog.Table) error {
	body, err := s.apiTest.ReadRow(role, nil, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body), &s.createRole)
	if err != nil {
		return err
	}

	s.createRole.ServiceID = s.createdService.ServiceID
	s.createRole.PermissionID = []uuid.UUID{s.createdPermissionId}

	createdRoleResponse, err := s.DB.CreateRole(context.Background(), db.CreateRoleParams{
		TenantName: s.tenant,
		ServiceID:  s.createdService.ServiceID,
		Column4:    s.createRole.PermissionID,
		Name:       s.createRole.Name,
	})
	if err != nil {
		return err
	}
	s.createdRoleResponseId = createdRoleResponse.RoleID
	return nil
}

func (s *systemAssignRoleTest) iRequestToAssignRoleToUser(userRole *godog.Table) error {
	body, err := s.apiTest.ReadRow(userRole, nil, false)
	if err != nil {
		return err
	}
	userTenantRole := userTenantRole{}

	if err = s.apiTest.UnmarshalJSON([]byte(body), &userTenantRole); err != nil {
		return err
	}

	s.userTenantRole.TenantName = userTenantRole.TenantName
	s.userTenantRole.RoleID = s.createdRoleResponseId
	reqbody, err := json.Marshal(s.userTenantRole)
	if err != nil {
		return err
	}

	s.apiTest.Body = string(reqbody)
	s.apiTest.URL = "/v1/system/users/" + userTenantRole.User + "/roles"
	s.apiTest.SetHeader("Content-Type", "application/json")
	s.apiTest.SetHeader("Authorization", "Basic "+s.BasicAuth(s.createdService.ServiceID.String(), "123456"))
	s.apiTest.SetHeader("x-subject", s.service.UserId)
	s.apiTest.SetHeader("x-action", "*")
	s.apiTest.SetHeader("x-resource", "*")
	s.apiTest.SetHeader("x-tenant", s.tenant)

	s.apiTest.SendRequest()
	return nil
}

func (s *systemAssignRoleTest) iRequestToAssignRoleToUserWhileFieldsAreMissing(userRole *godog.Table) error {
	body, err := s.apiTest.ReadRow(userRole, nil, false)
	if err != nil {
		return err
	}
	userTenantRole := userTenantRole{}

	if err = s.apiTest.UnmarshalJSON([]byte(body), &userTenantRole); err != nil {
		return err
	}
	if userTenantRole.TenantName != "" {
		s.userTenantRole.TenantName = userTenantRole.TenantName
	} else {
		s.userTenantRole.TenantName = ""
	}

	if userTenantRole.Role != "" {
		s.userTenantRole.RoleID = s.createdRoleResponseId
	} else {
		s.userTenantRole.RoleID = uuid.Nil
	}

	if userTenantRole.User == "" {
		userTenantRole.User = uuid.Nil.String()
	}

	reqbody, err := json.Marshal(s.userTenantRole)
	if err != nil {
		return err
	}
	s.apiTest.Body = string(reqbody)

	s.apiTest.URL = "/v1/system/users/" + userTenantRole.User + "/roles"
	s.apiTest.SetHeader("Content-Type", "application/json")
	s.apiTest.SetHeader("Authorization", "Basic "+s.BasicAuth(s.createdService.ServiceID.String(), "123456"))
	s.apiTest.SetHeader("x-subject", s.service.UserId)
	s.apiTest.SetHeader("x-action", "*")
	s.apiTest.SetHeader("x-resource", "*")
	s.apiTest.SetHeader("x-tenant", s.tenant)

	s.apiTest.SendRequest()
	return nil
}

func (s *systemAssignRoleTest) myRequestShouldFailWith(message string) error {

	if err := s.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := s.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (s *systemAssignRoleTest) theRoleShouldSuccessfullyBeAssigned() error {
	if err := s.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	if err := s.apiTest.AssertStringValueOnPathInResponse("ok", "true"); err != nil {
		return err
	}
	row, err := s.DB.IsRoleAssigned(context.Background(), db.IsRoleAssignedParams{
		TenantName: s.tenant,
		RoleID:     s.createdRoleResponseId,
		UserID:     s.createdUser.UserId,
	})
	if err != nil {
		return err
	}
	count, ok := row.(int64)
	if !ok {
		return fmt.Errorf("error getting user role")
	}
	if count == 0 {
		return fmt.Errorf("user role not found")
	}

	return nil
}

func (s *systemAssignRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		s.apiTest.Method = http.MethodPost
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = s.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^a permissions registered on the domain$`, s.aPermissionsRegisteredOnTheDomain)
	ctx.Step(`^A registered domain and tenant$`, s.aRegisteredDomainAndTenant)
	ctx.Step(`^i have role$`, s.iHaveRole)
	ctx.Step(`^I have service with$`, s.iHaveServiceWith)
	ctx.Step(`^I have user$`, s.iHaveUser)
	ctx.Step(`^I request to  assign role to user$`, s.iRequestToAssignRoleToUser)
	ctx.Step(`^I request to assign  role to user while fields are missing$`, s.iRequestToAssignRoleToUserWhileFieldsAreMissing)
	ctx.Step(`^my request should fail with "([^"]*)"$`, s.myRequestShouldFailWith)
	ctx.Step(`^the role should successfully be  assigned$`, s.theRoleShouldSuccessfullyBeAssigned)
}
