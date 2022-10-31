Feature: Delete Role
    As an admin
    I want to delete a role
    So that I can remove a role from my system

    Background:
        Given I have service with
            | name | user_id                              |
            | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |

    @success
    Scenario: I successfully delete role
        Given i have a role "admin" in tenant "vendor_1" with the permissions below
            | name           | description    | effect | action                 | resource                      | domains |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |
        And I  want to delete the role "admin"
        When i send a request to deleted the role
        Then the role should be delete
        And the role should not be found in the tenant

    @failer
    Scenario: role does not exists
        Given the role does not exists in the system
        When i send a request to deleted the role
        Then the request should fail with error "role does not exists"