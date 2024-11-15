package domain

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/test"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type deleteDomainTest struct {
	domainrequest dto.Domain
	test.TestInstance
	apiTest        src.ApiTest
	servicemodel   db.CreateServiceParams
	createdService db.CreateServiceRow
}

func TestDeleteDomain(t *testing.T) {
	d := &deleteDomainTest{}
	d.TestInstance = test.Initiate(context.Background(), "../../../../")
	d.apiTest.InitializeServer(d.Server)
	d.apiTest.InitializeTest(t, "create domain test", "feature/delete_domain.feature", d.InitializeScenario)
}

func (d *deleteDomainTest) iAmASystemUser() error {
	return nil
}

func (d *deleteDomainTest) iHaveDomain(domainName *godog.Table) error {

	body, err := d.apiTest.ReadRow(domainName, nil, false)
	if err != nil {
		return err
	}
	if err := d.apiTest.UnmarshalJSON([]byte(body), &d.domainrequest); err != nil {
		return err
	}
	d.servicemodel.Password = "password"

	if err != nil {
		return err
	}
	if d.createdService, err = d.DB.CreateService(context.Background(), d.servicemodel); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	_, err = d.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", d.createdService.ServiceID)
	if err != nil {

		return err

	}
	if err := d.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", d.servicemodel.Name)); err != nil {
		return err
	}
	_, err = d.DB.CreateDomain(context.Background(), db.CreateDomainParams{Name: d.domainrequest.Name, ServiceID: d.createdService.ServiceID})
	if err != nil {
		return err
	}

	return nil
}

func (d *deleteDomainTest) iSendTheRequest(domain *godog.Table) error {
	body, err := d.apiTest.ReadRow(domain, nil, false)
	if err != nil {
		return err
	}

	d.apiTest.Body = body
	d.apiTest.SetHeader("Authorization", "Basic "+basicAuth(d.createdService.ServiceID.String(), "password"))
	d.apiTest.SetHeader("x-subject", d.servicemodel.UserID.String())
	d.apiTest.SetHeader("x-action", "*")
	d.apiTest.SetHeader("x-resource", "*")
	d.apiTest.SetHeader("x-tenant", "administrator")

	d.apiTest.SendRequest()
	return nil
}

func (d *deleteDomainTest) theResultShouldBeEmptyError(message string) error {

	if err := d.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := d.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (d *deleteDomainTest) theResultShouldBeNotFoundError(message string) error {
	if err := d.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := d.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}

func (d *deleteDomainTest) theResultShouldBeSuccessfull(message string) error {

	if err := d.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	if err := d.apiTest.AssertStringValueOnPathInResponse("ok", message); err != nil {
		return err
	}

	return nil
}

func (d *deleteDomainTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		d.apiTest.URL = "/v1/domains"
		d.apiTest.Method = http.MethodDelete
		d.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = d.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^i am a system user$`, d.iAmASystemUser)
	ctx.Step(`^I have domain$`, d.iHaveDomain)
	ctx.Step(`^i send the request:$`, d.iSendTheRequest)
	ctx.Step(`^the result should be empty error "([^"]*)"$`, d.theResultShouldBeEmptyError)
	ctx.Step(`^the result should be not found error "([^"]*)"$`, d.theResultShouldBeNotFoundError)
	ctx.Step(`^the result should be successfull "([^"]*)"$`, d.theResultShouldBeSuccessfull)
}
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
