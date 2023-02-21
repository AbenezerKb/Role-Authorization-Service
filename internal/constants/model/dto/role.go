package dto

import (
	"2f-authorization/internal/constants"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type UpdateRole struct {
	// RoleID is the name of the role.
	RoleID uuid.UUID `json:"role_id"`
	// PermissionID is the list of permissions id's.
	PermissionsID []uuid.UUID `json:"permissions_id"`
}

func (r UpdateRole) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RoleID, validation.NotIn(uuid.Nil.String()).Error("Role id is required")),
		validation.Field(&r.PermissionsID, validation.By(validatePermissions)),
	)
}

type CreateRole struct {
	// Name is the name of the role.
	Name string `json:"name"`
	// TenantName is the name of the tenant the role is in.
	TenantName string `json:"tenant_name"`
	// PermissionID is the list of permissions id's.
	PermissionID []uuid.UUID `json:"permissions_id"`
	// ServiceID is the id of the service the tenant belongs to.
	ServiceID uuid.UUID `json:"service_id"`
}

func (r CreateRole) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("role name is required")),
		validation.Field(&r.TenantName, validation.Required.Error("tenant name is required")),
		validation.Field(&r.PermissionID, validation.By(validatePermissions)),
	)
}
func validatePermissions(id interface{}) error {
	p_id := id.([]uuid.UUID)
	if len(p_id) == 0 {
		return fmt.Errorf("atleast one permission is required")
	}

	for _, id := range p_id {
		if err := id == uuid.Nil; err {
			return fmt.Errorf("permission-id is required")
		}
	}

	return nil
}

type Role struct {
	// ID is the unique identifier for the service.
	// It is automatically generated when the role is created.
	ID uuid.UUID `json:"id,omitempty"`
	// Name is the name of the role.
	Name string `json:"name,omitempty"`
	// Permissions are the list of permissions names this role contains
	Permissions []string `json:"permissions,omitempty"`
	// Status is the status of the role.
	Status string `json:"status,omitempty"`
	// DeletedAt is the time this service was created.
	DeletedAt time.Time `json:"-"`
	// CreatedAt is the time this service was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// UpdatedAt is the time this service was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type TenantUsersRole struct {
	RoleTenant
	//UserID is the user identifier which going to get the the role
	UserID uuid.UUID `json:"user_id"`
}

type RoleTenant struct {
	//RoleID is id of the role which is going to be assigned to the user.
	RoleID uuid.UUID `json:"role_id"`
	//RoleName is the name of the role which is going to be assigned to the user.
	RoleName string `json:"role_name"`
	//TenantName The Name of the tenante which is given when the tenant is created
	TenantName string `json:"tenant_name"`
}

func (t TenantUsersRole) Validate() error {

	return validation.ValidateStruct(
		&t,
		validation.Field(&t.TenantName, validation.Required.Error("tenant is required")),
		validation.Field(&t.UserID, is.UUID, validation.NotIn(uuid.Nil.String()).Error("User id required")),
		validation.Field(&t.RoleID, validation.When(t.RoleName == "",
			validation.NotIn(uuid.Nil.String()).Error("Role id or name is required"))),
		validation.Field(&t.RoleName, validation.When(t.RoleID == uuid.Nil,
			validation.Required.Error("Role id or name is required"),
		)),
	)
}

type GetAllRolesReq struct {
	//TenantName is the name of the tenant
	TenantName string `json:"tenant_name"`
	//ServiceID  is the service id of service.
	ServiceID uuid.UUID `json:"service_id"`
}

func (g GetAllRolesReq) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.TenantName, validation.Required.Error("tenant name can not be blank")),
		validation.Field(&g.ServiceID, validation.NotIn(uuid.Nil.String()).Error("service id is required")))
}

type UpdateRoleStatus struct {
	// Status is new status of the role that is going to replace the old one
	Status string `json:"status"`
}

func (u UpdateRoleStatus) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Status, validation.Required.Error("status is required"), validation.In(constants.Active, constants.InActive).Error("invalid status")),
	)
}
