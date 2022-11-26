Feature: Update Permission Status

    As a admin
    I want to update the permission's status
    So that I can activate or deactivate permission as appropriate

    Background:
        Given I have service with
            | name | user_id                              |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |

    @success
    Scenario Outline: Successfully Update permission Status
        Given i have the following permission in tenant "vendor_1"
            | name   | description   | effect   | action   | resource   | fields   |
            | <name> | <description> | <effect> | <action> | <resource> | <fields> |
        And I want to update the permission's status to "<status>"
        When I send the request to update the status
        Then the permission status should update to "<status>"
        Examples:
            | name           | description    | effect | action                 | resource                      | fields | status   |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | *      | ACTIVE   |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | *      | INACTIVE |
    @failure
    Scenario Outline: Failed permission Status Update
        Given the permission is not on the system "<id>"
        And I want to update the permission's status to "<status>"
        When I send the request to update the status
        Then Then I should get an error with message "<message>"

        Examples:
            | id                                   | status | message                    |
            | 2302f310-e60f-4434-b54b-d133eaa63a2c | ACTIVE | permission does not exists |

    @failure
    Scenario Outline: Missing required field values
        Given i have the following permission in tenant "vendor_1"
            | name   | description   | effect   | action   | resource   | fields   |
            | <name> | <description> | <effect> | <action> | <resource> | <fields> |
        And I want to update the permission's status to "<status>"
        When I send the request to update the status
        Then Then I should get a field error with message "<message>"
        Examples:
            | name           | description    | effect | action                 | resource                      | fields | status              | message            |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | *      | something not valid | invalid status     |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | *      |                     | status is required |
