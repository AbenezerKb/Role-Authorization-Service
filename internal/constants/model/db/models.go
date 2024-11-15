// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type Status string

const (
	StatusPENDING  Status = "PENDING"
	StatusACTIVE   Status = "ACTIVE"
	StatusINACTIVE Status = "INACTIVE"
)

func (e *Status) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Status(s)
	case string:
		*e = Status(s)
	default:
		return fmt.Errorf("unsupported scan type for Status: %T", src)
	}
	return nil
}

type NullStatus struct {
	Status Status
	Valid  bool // Valid is true if Status is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatus) Scan(value interface{}) error {
	if value == nil {
		ns.Status, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Status.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Status, nil
}

type Domain struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Status    Status       `json:"status"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	ServiceID uuid.UUID    `json:"service_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Permission struct {
	ID             uuid.UUID     `json:"id"`
	Status         Status        `json:"status"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Statement      pgtype.JSON   `json:"statement"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	ServiceID      uuid.UUID     `json:"service_id"`
	TenantID       uuid.NullUUID `json:"tenant_id"`
	DeletedAt      sql.NullTime  `json:"deleted_at"`
	DeleteOrUpdate bool          `json:"delete_or_update"`
}

type PermissionDomain struct {
	ID           uuid.UUID `json:"id"`
	PermissionID uuid.UUID `json:"permission_id"`
	DomainID     uuid.UUID `json:"domain_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PermissionsHierarchy struct {
	ID        uuid.UUID `json:"id"`
	Parent    uuid.UUID `json:"parent"`
	Child     uuid.UUID `json:"child"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role struct {
	ID        uuid.UUID    `json:"id"`
	Status    Status       `json:"status"`
	Name      string       `json:"name"`
	TenantID  uuid.UUID    `json:"tenant_id"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type RolePermission struct {
	ID           uuid.UUID `json:"id"`
	PermissionID uuid.UUID `json:"permission_id"`
	RoleID       uuid.UUID `json:"role_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RoleTenant struct {
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ID         uuid.UUID `json:"id"`
	Status     Status    `json:"status"`
	UpdatedAt  time.Time `json:"updated_at"`
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

type Service struct {
	ID        uuid.UUID    `json:"id"`
	Status    Status       `json:"status"`
	Name      string       `json:"name"`
	Password  string       `json:"password"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Tenant struct {
	ID         uuid.UUID    `json:"id"`
	Status     Status       `json:"status"`
	TenantName string       `json:"tenant_name"`
	ServiceID  uuid.UUID    `json:"service_id"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DomainID   uuid.UUID    `json:"domain_id"`
	Inherit    bool         `json:"inherit"`
}

type TenantUsersRole struct {
	ID        uuid.UUID    `json:"id"`
	TenantID  uuid.UUID    `json:"tenant_id"`
	UserID    uuid.UUID    `json:"user_id"`
	RoleID    uuid.UUID    `json:"role_id"`
	Status    Status       `json:"status"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type User struct {
	ID        uuid.UUID    `json:"id"`
	UserID    uuid.UUID    `json:"user_id"`
	Status    Status       `json:"status"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	ServiceID uuid.UUID    `json:"service_id"`
}
