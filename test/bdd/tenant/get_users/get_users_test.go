package getusers

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/platform/argon"
	"2f-authorization/test"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"

	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type Service struct {
	Name      string `json:"name"`
	ServiceID string `json:"service_id"`
}
type DomainTenantsRoles struct {
	Domain     string `json:"domain"`
	TenantName string `json:"tenant_name"`
	Role       string `json:"role"`
}
type getUserTenantUsersWithRoles struct {
	test.TestInstance
	apiTest         src.ApiTest
	serviceName     string
	servicePassword string
	serviceUserId   uuid.UUID
	serviceID       uuid.UUID
	Subject         uuid.UUID
	Domain          uuid.UUID
	tenant          uuid.UUID
	users           map[string]uuid.UUID
	roles           map[string]uuid.UUID
	userids         map[string]uuid.UUID
	userToRoles     map[string][]string
}
type TenantUserRoles struct {
	Name  string `json:"name"`
	Roles string `json:"roles"`
}
type TenanUserROlesResponse struct {
	Name  string `json:"name"`
	Roles string `json:"roles"`
}

func TestGetTenantUsersRoles(t *testing.T) {

	r := &getUserTenantUsersWithRoles{}
	r.users = make(map[string]uuid.UUID)
	r.roles = make(map[string]uuid.UUID)
	r.userids = make(map[string]uuid.UUID)
	r.userToRoles = map[string][]string{}
	r.TestInstance = test.Initiate(context.Background(), "../../../../")
	r.apiTest.InitializeServer(r.Server)
	r.apiTest.InitializeTest(t, "get tenant users with roles", "feature/get_users.feature", r.InitializeScenario)
}
func (r *getUserTenantUsersWithRoles) aRegisteredDomainAndTenantAndRole(domaintenantrole *godog.Table) error {
	body, _ := r.apiTest.ReadRow(domaintenantrole, nil, false)
	domainTenantRole := DomainTenantsRoles{}
	if err := json.Unmarshal([]byte(body), &domainTenantRole); err != nil {
		return err
	}
	tenantId := r.DB.Pool.QueryRow(context.Background(), "INSERT INTO tenants(tenant_name,service_id,domain_id) VALUES ($1,$2,$3) RETURNING  id", domainTenantRole.TenantName, r.serviceID, r.Domain)
	var id *uuid.UUID
	tenantId.Scan(&id)
	r.tenant = *id

	return nil
}

func (r *getUserTenantUsersWithRoles) iAmSystemUser() error {

	return nil
}

func (r *getUserTenantUsersWithRoles) iHaveServiceWith(service *godog.Table) error {

	var err error
	body, _ := r.apiTest.ReadRow(service, nil, false)
	serviceInstance := Service{}
	if err := json.Unmarshal([]byte(body), &serviceInstance); err != nil {
		return err
	}
	r.serviceName = serviceInstance.Name
	r.serviceUserId = uuid.New()
	if r.servicePassword, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}
	serviceCreated, err := r.DB.CreateService(context.Background(), db.CreateServiceParams{
		Name:     serviceInstance.Name,
		Password: r.servicePassword,
		UserID:   r.serviceUserId,
	})
	if err != nil {
		return err
	}

	r.serviceID = serviceCreated.ServiceID
	if _, err := r.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE' where id = $1", r.serviceID); err != nil {
		return err
	}
	if err := r.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", r.serviceName)); err != nil {
		return err
	}
	sub := r.DB.Pool.QueryRow(context.Background(), "select ur.id from (select * from tenants where tenant_name='administrator')tn INNER JOIN (select id,user_id,tenant_id from tenant_users_roles ) tur ON tn.id=tur.tenant_id Inner join (select id from users)ur on ur.id=tur.user_id ")

	var i *uuid.UUID
	sub.Scan(&i)
	r.Subject = *i

	dom := r.DB.Pool.QueryRow(context.Background(), "select id from domains where name='administrator'")
	dom.Scan(&i)
	r.Domain = *i

	return nil
}

func (r *getUserTenantUsersWithRoles) iSendRequest(tenantname string) error {
	r.apiTest.SetHeader("Authorization", "Basic "+r.BasicAuth(r.serviceID.String(), "123456"))
	r.apiTest.SetHeader("x-subject", r.Subject.String())
	r.apiTest.SetHeader("x-action", "*")
	r.apiTest.SetHeader("x-resource", "*")
	r.apiTest.SetHeader("x-tenant", r.tenant.String())
	r.apiTest.SendRequest()
	return nil
}

func (r *getUserTenantUsersWithRoles) tenantUsersAndRole(tenantUserRole *godog.Table) error {
	rolesnames := make(map[string]string)
	body, _ := r.apiTest.ReadRows(tenantUserRole, nil, false)
	tenantUserRoles := []TenantUserRoles{}

	if err := json.Unmarshal([]byte(body), &tenantUserRoles); err != nil {
		return err
	}
	for _, rols := range tenantUserRoles {
		assgnedRoles := []string{}
		roles := strings.Split(rols.Roles, ",")
		for _, role := range roles {
			assgnedRoles = append(assgnedRoles, role)
			rolesnames[role] = role
		}
		r.userToRoles[rols.Name] = assgnedRoles
	}
	for _, rname := range rolesnames {

		var rid *uuid.UUID

		roleId := r.DB.Pool.QueryRow(context.Background(), "insert into roles ( name,tenant_id) values ($1,$2) RETURNING id ", rname, r.tenant)
		if roleId != nil {
			roleId.Scan(&rid)
			r.roles[rname] = *rid
		}

	}

	for _, tur := range tenantUserRoles {
		roles := strings.Split(tur.Roles, ",")
		userid := uuid.New()
		usr := r.DB.Pool.QueryRow(context.Background(), "insert into users (user_id,service_id) values ($1,$2) RETURNING id", userid, r.serviceID)
		var id *uuid.UUID
		usr.Scan(&id)
		r.users[tur.Name] = userid
		r.userids[tur.Name] = *id

		for i := 0; i < len(roles); i++ {

			if _, err := r.DB.Pool.Exec(context.Background(), "INSERT INTO tenant_users_roles (tenant_id,user_id,role_id) values ($1,$2,$3) ", r.tenant, id, r.roles[roles[i]]); err != nil {
				return err
			}

		}
	}

	return nil
}

func (r *getUserTenantUsersWithRoles) theResultShouldBeNameAndRoles(results *godog.Table) error {
	body, _ := r.apiTest.ReadRows(results, nil, false)
	turs := []TenanUserROlesResponse{}
	if err := json.Unmarshal([]byte(body), &turs); err != nil {
		return err
	}
	for _, tur := range turs {

		roles := strings.Split(tur.Roles, ",")
		storedRole := r.userToRoles[tur.Name]
		for _, role := range roles {
			isRoleExist := false
			for _, rol := range storedRole {
				if role == rol {
					isRoleExist = true
				}
			}
			if !isRoleExist {
				return fmt.Errorf("user and  the user role dose not matched")
			}
		}
	}

	return nil
}

func (r *getUserTenantUsersWithRoles) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.URL = "/v1/tenants/users"
		r.apiTest.Method = http.MethodGet
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^A registered domain and tenant and role$`, r.aRegisteredDomainAndTenantAndRole)
	ctx.Step(`^i am system user$`, r.iAmSystemUser)
	ctx.Step(`^I have service with$`, r.iHaveServiceWith)
	ctx.Step(`^I send request  "([^"]*)"$`, r.iSendRequest)
	ctx.Step(`^tenant users and role$`, r.tenantUsersAndRole)
	ctx.Step(`^the result should be <name> and <roles>$`, r.theResultShouldBeNameAndRoles)
}
