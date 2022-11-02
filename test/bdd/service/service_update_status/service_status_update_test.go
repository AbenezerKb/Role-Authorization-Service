package serviceupdatestatus

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
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type updateServiceStatusTest struct {
	test.TestInstance
	apiTest        src.ApiTest
	OK             bool `json:"ok"`
	service        db.CreateServiceParams
	serviceStatus  dto.UpdateServiceStatus
	createdService db.CreateServiceRow
}

func TestDeleteService(t *testing.T) {
	u := &updateServiceStatusTest{}
	u.TestInstance = test.Initiate(context.Background(), "../../../../")
	u.apiTest.InitializeServer(u.Server)
	u.apiTest.InitializeTest(t, "update service status test", "feature/update_service_status.feature", u.InitializeScenario)
}

func (u *updateServiceStatusTest) iUpdateTheServicesStatusTo(status string) error {
	u.serviceStatus.Status = status
	u.serviceStatus.ServiceID = u.createdService.ServiceID
	body, err := json.Marshal(&u.serviceStatus)
	if err != nil {
		return err
	}
	u.apiTest.Body = string(body)
	u.apiTest.SendRequest()
	return nil
}

func (u *updateServiceStatusTest) theServiceStatusShouldUpdateTo(status string) error {
	service, err := u.DB.GetServiceById(context.Background(), u.createdService.ServiceID)
	if err != nil {
		return err
	}
	if u.apiTest.AssertEqual(service.Status, status); err != nil {
		return err
	}
	return nil
}

func (u *updateServiceStatusTest) thenIShouldGetErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (u *updateServiceStatusTest) theServiceIsRegisteredOnTheSystem(service *godog.Table) error {
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

	if u.createdService, err = u.DB.CreateService(context.Background(), u.service); err != nil {
		return err
	}

	if _, err := u.DB.Pool.Exec(context.Background(), "UPDATE services set status = 'ACTIVE'  where id = $1", u.createdService.ServiceID); err != nil {
		return err
	}

	if err := u.Opa.Refresh(context.Background(), fmt.Sprintf("Created service with name - [%v]", u.service.Name)); err != nil {
		return err
	}

	return nil
}
func (u *updateServiceStatusTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		u.apiTest.URL = "/v1/services/status"
		u.apiTest.Method = http.MethodPatch
		u.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^the service is registered on the system$`, u.theServiceIsRegisteredOnTheSystem)

	ctx.Step(`^I update the service\'s status to "([^"]*)"$`, u.iUpdateTheServicesStatusTo)
	ctx.Step(`^the service status should update to "([^"]*)"$`, u.theServiceStatusShouldUpdateTo)
	ctx.Step(`^Then I should get error with message "([^"]*)"$`, u.thenIShouldGetErrorWithMessage)
}
