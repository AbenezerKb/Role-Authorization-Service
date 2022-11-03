package updatestatus

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

type updateUserStatusTest struct {
	test.TestInstance
	apiTest        src.ApiTest
	OK             bool `json:"ok"`
	service        db.CreateServiceParams
	userStatus     dto.UpdateUserStatus
	createdService db.CreateServiceRow
	createdUser    dto.RegisterUser
}

func TestUpdateServiceStatus(t *testing.T) {
	u := &updateUserStatusTest{}
	u.TestInstance = test.Initiate(context.Background(), "../../../../")
	u.apiTest.InitializeServer(u.Server)
	u.apiTest.InitializeTest(t, "update user status test", "feature/update_user_status.feature", u.InitializeScenario)
}

func (u *updateUserStatusTest) iHaveServiceWith(service *godog.Table) error {
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

func (u *updateUserStatusTest) iSendTheRequestToUpdateTheStatus() error {
	u.apiTest.SetHeader("Authorization", "Basic "+u.BasicAuth(u.createdService.ServiceID.String(), "123456"))
	u.apiTest.SetHeader("x-subject", u.service.UserID.String())
	u.apiTest.SetHeader("x-action", "*")
	u.apiTest.SetHeader("x-resource", "*")
	u.apiTest.SetHeader("x-tenant", "administrator")
	u.apiTest.SendRequest()
	return nil
}

func (u *updateUserStatusTest) iWantToUpdateTheUsersStatusTo(status string) error {
	u.userStatus.Status = status
	u.userStatus.ServiceID = u.createdService.ServiceID
	u.userStatus.UserID = u.createdUser.UserId
	body, err := json.Marshal(&u.userStatus)
	if err != nil {
		return err
	}
	u.apiTest.Body = string(body)
	return nil
}

func (u *updateUserStatusTest) theUserStatusShouldUpdateTo(status string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	user, err := u.DB.GetUserWithUserIdAndServiceId(context.Background(), db.GetUserWithUserIdAndServiceIdParams{
		UserID:    u.userStatus.UserID,
		ServiceID: u.userStatus.ServiceID,
	})
	if err != nil {
		return err
	}
	if u.apiTest.AssertEqual(user.Status, status); err != nil {
		return err
	}
	return nil
}

func (u *updateUserStatusTest) thereIsAUserWithTheFollowingDetails(user *godog.Table) error {
	body, err := u.apiTest.ReadRow(user, nil, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body), &u.createdUser)
	if err != nil {
		return err
	}

	u.createdUser.ServiceID = u.createdService.ServiceID

	err = u.DB.RegisterUser(context.Background(), db.RegisterUserParams{
		UserID:    u.createdUser.UserId,
		ServiceID: u.createdUser.ServiceID,
	})
	if err != nil {
		return err
	}
	return nil
}
func (u *updateUserStatusTest) theUserIsNotRegisteredOnTheSystem(user *godog.Table) error {
	body, err := u.apiTest.ReadRow(user, nil, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body), &u.createdUser)
	if err != nil {
		return err
	}

	u.userStatus.UserID = u.createdUser.UserId
	return nil
}

func (u *updateUserStatusTest) thenIShouldGetAnErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}

func (u *updateUserStatusTest) thenIShouldGetAFieldErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}
func (u *updateUserStatusTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		u.apiTest.URL = "/v1/users/status"
		u.apiTest.Method = http.MethodPatch
		u.apiTest.SetHeader("Content-Type", "application/json")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.Pool.Exec(ctx, "truncate table services,tenants,users,roles,permissions,role_permissions,tenant_users_roles,domains,permission_domains,permissions_hierarchy cascade;")
		return ctx, nil
	})
	ctx.Step(`^Then I should get a field error with message "([^"]*)"$`, u.thenIShouldGetAFieldErrorWithMessage)
	ctx.Step(`^the user is not registered on the system$`, u.theUserIsNotRegisteredOnTheSystem)
	ctx.Step(`^Then I should get an error with message "([^"]*)"$`, u.thenIShouldGetAnErrorWithMessage)
	ctx.Step(`^I have service with$`, u.iHaveServiceWith)
	ctx.Step(`^I send the request to update the status$`, u.iSendTheRequestToUpdateTheStatus)
	ctx.Step(`^I want to update the user\'s status to "([^"]*)"$`, u.iWantToUpdateTheUsersStatusTo)
	ctx.Step(`^the user status should update to "([^"]*)"$`, u.theUserStatusShouldUpdateTo)
	ctx.Step(`^there is a user with the following details:$`, u.thereIsAUserWithTheFollowingDetails)
}
