Feature: Update Role
    As an admin
    I want to update a role
    So that I can change the permissions given to that role

    Background:
        Given I have service with
            | name | user_id                              |
            | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |
        And i have a role "admin" in tenant "vendor_1" with the permissions below
            | name           | description    | effect | action                 | resource                      | domains |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |

    @success
    Scenario: I successfully update role
        Given I  want to update the role "admin" with the following permissions:
            | name           | description    | effect | action                 | resource                      | domains |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |
            | create service | create service | allow  | "admin:service:create" | "admin:service:createservice" | vendor  |
        When i send a request to update the role
        Then the role should be updated
