package deleteservice

import (
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/platform/argon"
	"2f-authorization/test"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type deleteServiceTest struct {
	test.TestInstance
	apiTest        src.ApiTest
	OK             bool `json:"ok"`
	service        db.CreateServiceParams
	createdService db.CreateServiceRow
}

func TestDeleteService(t *testing.T) {
	c := &deleteServiceTest{}
	c.TestInstance = test.Initiate(context.Background(), "../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "delete service test", "feature/delete_service.feature", c.InitializeScenario)
}

func (d *deleteServiceTest) iDeleteTheService(service *godog.Table) error {
	body, err := d.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	d.apiTest.Body = body
	d.apiTest.SetHeader("Authorization", "Basic "+basicAuth(d.createdService.ServiceID.String(), "123456"))

	d.apiTest.SendRequest()
	return nil
}

func (d *deleteServiceTest) iDeleteTheServiceWithId(id string) error {
	d.apiTest.SetHeader("Authorization", "Basic "+basicAuth(id, "123456"))
	d.apiTest.SendRequest()
	return nil
}

func (d *deleteServiceTest) iHaveARegisteredService(service *godog.Table) error {
	body, err := d.apiTest.ReadRow(service, nil, false)
	if err != nil {
		return err
	}
	if err = d.apiTest.UnmarshalJSON([]byte(body), &d.service); err != nil {
		return err
	}

	if d.service.Password, err = argon.CreateHash("123456", argon.DefaultParams); err != nil {
		return err
	}

	if d.createdService, err = d.DB.CreateService(context.Background(), d.service); err != nil {
		return err
	}

	if _, err := d.DB.Pool.Exec(context.Background(), "UPDATE services set status = true where id = $1", d.createdService.ServiceID); err != nil {
		return err
	}

	if err := d.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", d.service.Name)); err != nil {
		return err
	}

	return nil
}

func (d *deleteServiceTest) theRequestShouldFailWithErrorMessage(message string) error {
	if err := d.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := d.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}

func (d *deleteServiceTest) theServiceShouldBeDeleted() error {
	if err := d.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	return nil
}

func (d *deleteServiceTest) theRequestShouldFailWithFieldErrorMessage(message string) error {
	if err := d.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := d.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (d *deleteServiceTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		d.apiTest.URL = "/v1/services"
		d.apiTest.Method = http.MethodDelete
		d.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = d.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})

	ctx.Step(`^I delete the service:$`, d.iDeleteTheService)
	ctx.Step(`^I delete the service with id "([^"]*)"$`, d.iDeleteTheServiceWithId)
	ctx.Step(`^I have a registered service$`, d.iHaveARegisteredService)
	ctx.Step(`^the request should fail with field error message "([^"]*)"$`, d.theRequestShouldFailWithFieldErrorMessage)
	ctx.Step(`^the request should fail with error message "([^"]*)"$`, d.theRequestShouldFailWithErrorMessage)
	ctx.Step(`^The service should be deleted$`, d.theServiceShouldBeDeleted)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
