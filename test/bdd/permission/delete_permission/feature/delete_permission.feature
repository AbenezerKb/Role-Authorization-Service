Feature: Delete permission
    As a user,
    i want to delete the permission
    so that i can remove permissions from my tenant.
    Background: I have a registered service,
        Given I have a registered service
            | name | user_id                              |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |
        And A permission registered on the tenant
            | name           | description    | effect | action                 | fields | resource                      |
            | delete service | delete service | allow  | "admin:service:delete" | *      | "admin:service:deleteservice" |

    @success
    Scenario Outline: Successful Delete Permmission
        Given I want to delete the permission
        And I send a request to delete the permission
        Then The request should be successfull
        And the permission should be deleted
