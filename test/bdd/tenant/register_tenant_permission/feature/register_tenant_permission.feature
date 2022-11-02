Feature: Register tenant role
    As an Admin,
    I want to create a new permission
    So that I am able to add new permissions based on my preference

    Background:
        Given I have service with
            | name | user_id                              |
            | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |

    Scenario Outline: successfully add new permission in my tenant
        Given I want to register the following permission under the tenant
            | name           | description    | effect | action                 | resource                      |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" |
        When I send request to add the permission
        Then the request should successfull
        And the permission should be accessible through the tenant

    Scenario Outline: Required fields are missing
        Given I want to register the following permission under the tenant
            | name   | description   | effect   | action   | resource   |
            | <name> | <description> | <effect> | <action> | <resource> |
        When I send request to add the permission
        Then The request should fail with error "<message>"
        Examples:
            | name           | description    | effect | action               | resource                    | message                                   |
            |                | delete service | allow  | admin:service:delete | admin:service:deleteservice | permission name is required               |
            | delete service |                | allow  | admin:service:delete | admin:service:deleteservice | permission description is required        |
            | delete service | delete service |        | admin:service:delete | admin:service:deleteservice | effect: statement effect is required.     |
            | delete service | delete service | allow  |                      | admin:service:deleteservice | action: statement action is required.     |
            | delete service | delete service | allow  | admin:service:delete |                             | resource: statement resource is required. |
