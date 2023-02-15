Feature: Update User's Role Status

    As a user
    I want to update the user's role status
    So that I can activate or deactivate the role as appropriate

    Background:
        Given I have service with
            | name | user_id                              |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |
        And a role "Admin" in tenant "vendor_1" with the following permissions
            | name           | description    | effect | action                 | resource                      | domains |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |
    @success
    Scenario Outline: Successfully Update Role Status
        Given the user has the following role in the following tenant
            | user_id   | tenant   | role_id   |
            | <user_id> | <tenant> | <role_id> |
        And I want to update the user's role status to "<status>"
        When I send the request to update the status
        Then the role status should update to "<status>"
        Examples:
            | user_id                              | tenant   | role  | status   |
            | 8e70df10-7957-426c-912a-9ed00f887e46 | vendor_1 | Admin | INACTIVE |
            | 8e70df10-7957-426c-912a-9ed00f887e46 | vendor_1 | Admin | ACTIVE   |
    @failure
    Scenario Outline: Missing required field values
        Given the user has the following role in the following tenant
            | user_id   | tenant   | role_id   |
            | <user_id> | <tenant> | <role_id> |
        And I want to update the user's role status to "<status>"
        When I send the request to update the status
        Then Then I should get a field error with message "<message>"
        Examples:
            | user_id                              | tenant   | role  | status         | message            |
            | 8e70df10-7957-426c-912a-9ed00f887e46 | vendor_1 | Admin |                | status is required |
            | 8e70df10-7957-426c-912a-9ed00f887e46 | vendor_1 | Admin | invalid status | invalid status     |
