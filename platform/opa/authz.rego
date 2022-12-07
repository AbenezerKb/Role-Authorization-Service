package authz

import future.keywords.every

default allow := false

allow {
	authorize
}

allow {
	IsAdmin
}

authorize {
	check_user_permissions[permission]
	match(permission)
}

IsAdmin {
	not authorize
	check_admin_permissions[permission]
	match(permission)
}

match(permission) {
	has_action(permission)
	has_resource(permission)
	has_fields(permission)
}

match(permission) {
	pc := permission.child[_]
	pc.status = "ACTIVE"
	pc.statement.effect = "allow"
	has_resource(pc)
	has_action(pc)
}

check_service = service {
	service := data.services[input.service]
	service.status = "ACTIVE"
}

check_tenant = tenant {
	tenant := check_service.tenants[input.tenant]
	tenant.status = "ACTIVE"
}

check_admin_tenant = tenant {
	tenant := check_service.tenants.administrator
	tenant.status = "ACTIVE"
}

check_user = user {
	user := check_tenant.users[input.subject]
	user_check(user)
}

check_admin_user = user {
	user := check_admin_tenant.users[input.subject]
	user_check(user)
}

user_check(user) {
	user.status = "ACTIVE"
	user.user_role_status
}

check_user_role[role] {
	role := check_user.role[_]
	role.status = "ACTIVE"
}

check_admin_role[role] {
	role := check_admin_user.role[_]
	role.status = "ACTIVE"
}

check_user_permissions[permission] {
	permission := check_user_role[role].permissions[_]
	check_permission(permission)
}

check_admin_permissions[permission] {
	permission := check_admin_role[role].permissions[_]
	check_permission(permission)
}

check_permission(permission) {
	permission.status = "ACTIVE"
	permission.statement.effect == "allow"
}

#
# Variable expansion
#
no_variables(a) {
	contains(a, "${") == false
}

variables(a) {
	indexof(a, "${") < indexof(a, "}")
}

# Defines the funservice := data.services[input.service]
# 'expanded'.
# Note: currently, only expands the one variable we know: ${a2:username}.
expand(orig) = expanded {
	# split(input.subjects[_], ":", ["user", _, username])
	expanded := replace(orig, "${username}", input.subject)
}

no_wildcard(a) {
	contains(a, "*") == false
}

wildcard(a) {
	endswith(a, ":*")
}

# Check that it does not end with ":*" AND that it is not a solitary "*".
# Note: The latter is done so that we don't end up with 'input.resource = *'
# rules in our partial results.
# Note that we avoid "not", which hinders partial result optimizations, see
# https://github.com/open-policy-agent/opa/issues/709.
not_wildcard(a) {
	endswith(a, ":*") == false
	a != "*"
}

# This supports these business rules:
# (a) A wildcard may only occur in the last section.
# (b) A wildcard may not be combined with a prefix (e.g. cannot say "x:y:foo*").
# (c) A wildcard applies to the current section and any deeper sections
#     (e.g. "a:*" matches "a:b" and "a:b:c", etc.).
wildcard_match(a, b) {
	startswith(a, trim(b, "*"))
}

#
# Resource matching
#
resource_matches(inp, stored) {
	no_variables(stored)
	not_wildcard(stored)
	inp == stored
}

resource_matches(inp, stored) {
	no_variables(stored)
	wildcard(stored)
	wildcard_match(inp, stored)
}

resource_matches(inp, stored) {
	variables(stored)
	not_wildcard(stored)
	inp == expand(stored)
}

resource_matches(inp, stored) {
	variables(stored)
	wildcard(stored)
	wildcard_match(inp, expand(stored))
}

resource_matches(_, "*") = true

has_resource(permission) {
	statement_resource := permission.statement.resource
	resource_matches(input.resource, statement_resource)
}

fields_matches(_, ["*"]) = true

fields_matches(inp, stored) {
	every i in inp {
		i = stored[_]
	}
}

has_fields(permission) {
	statement_fields := permission.statement.fields
	fields_matches(input.fields, statement_fields)
}

action_matches(inp, stored) {
	no_wildcard(stored)
	inp == stored
}

action_matches(inp, stored) = action_match(split(stored, ":"), split(inp, ":"))

action_match([service, "*"], [service, _, _]) = true

action_match([service, type, "*"], [service, type, _]) = true

action_match([service, "*", verb], [service, _, verb]) = true

action_match(["*", verb], [_, _, verb]) = true

action_match(["*"], _) = true

has_action(permission) {
	statement_action := permission.statement.action
	action_matches(input.action, statement_action)
}

allowedPermissions[rm] {
	matchedPermissions[permission]
	rm := json.remove(permission, ["child"])
}

matchedPermissions[permission] {
	check_user_permissions[permission]
}

matchedPermissions[child] {
	child := check_user_permissions[_].child[_]
	child.statement.effect == "allow"
}
