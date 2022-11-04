Feature: Update Role Status

    As a admin
    I want to update the role's status
    So that I can activate or deactivate role as appropriate

    Background:
        Given I have service with
            | name | user_id                              |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |

    @success
    Scenario Outline: Successfully Update Role Status
        Given i have a role "admin" in tenant "vendor_1" with the permissions below
            | name   | description   | effect   | action   | resource   | domains   |
            | <name> | <description> | <effect> | <action> | <resource> | <domains> |
        And I want to update the role's status to "<status>"
        When I send the request to update the status
        Then the role status should update to "<status>"
        Examples:
            | name           | description    | effect | action                 | resource                      | domains | status   |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  | ACTIVE   |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  | INACTIVE |
    @failure
    Scenario Outline: Failed Role Status Update
        Given the role is not on the system "<id>"
        And I want to update the role's status to "<status>"
        When I send the request to update the status
        Then Then I should get an error with message "<message>"

        Examples:
            | id                                   | status | message        |
            | 2302f310-e60f-4434-b54b-d133eaa63a2c | ACTIVE | role not found |

    @failure
    Scenario Outline: Missing required field values
        Given i have a role "admin" in tenant "vendor_1" with the permissions below
            | name   | description   | effect   | action   | resource   | domains   |
            | <name> | <description> | <effect> | <action> | <resource> | <domains> |
        And I want to update the role's status to "<status>"
        When I send the request to update the status
        Then Then I should get a field error with message "<message>"
        Examples:
            | name           | description    | effect | action                 | resource                      | domains | status              | message            |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  | something not valid | invalid status     |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |                     | status is required |
