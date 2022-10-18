Feature: Create role
    As an Admin,
    I want to create a role
    So that I can assign it to a user

    Background:
        Given I have service with
            | name | user_id                              |
            | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |
        And a permissions registered on the domain
            | name           | description    | effect | action                 | resource                      | domains |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |

    Scenario Outline: successfully creates a new role
        When I request to create a role with the "delete service" permissions
            | name    |
            | my_role |
        Then the role should successfully be created

    Scenario Outline: Required fields are missing
        When I request to create a role with the "delete service" permissions
            | name   |
            | <name> |
        Then my request should fail with "<message>"
        Examples:
            | name | message               |
            |      | role name is required |