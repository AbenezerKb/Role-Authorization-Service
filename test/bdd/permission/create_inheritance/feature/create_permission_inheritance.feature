Feature: Create Permission Inheritance
    As an authorized system,
    I want to create dependency betweeen permissions
    So that I can have a relations betweeen permissions

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
            | get service    | get service    | allow  | "admin:service:get"    | "admin:service:getservice"    | vendor  |

    Scenario Outline: Successfully create a dependency between permission
        Given I want to have a relation between "delete service" permission as a parent and "get service" permission as a child
        When I send request to create the inheritance
        Then the request should be successfull

    Scenario Outline: Required fields are missing
        Given I want to have a relation between permissions
            | permission   | inherited_permissions   |
            | <permission> | <inherited_permissions> |
        When I send request to create the inheritance
        Then i should get error with message "<message>"
        Examples:
            | permission     | inherited_permissions | message                            |
            |                | get service           | permission is required             |
            | delete service |                       | inherited permissions are required |