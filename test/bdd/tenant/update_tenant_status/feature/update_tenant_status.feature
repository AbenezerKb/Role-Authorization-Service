Feature: Update Tenant Status

    As a admin
    I want to update the tenant's status
    So that I can activate or deactivate tenant as appropriate

    Background:
        Given I have service with
            | name | user_id                              |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |

    @success
    Scenario Outline: Successfully Update Tenant Status
        Given I want to update the tenant's status to "<status>"
        When I send the request to update the status
        Then the tenant status should update to "<status>"
        Examples:
            | status   |
            | ACTIVE   |
            | INACTIVE |
    @failure
    Scenario Outline:Tenant Status Update Failed
        Given the tenant is not on the system "<id>"
        And I want to update the tenant's status to "<status>"
        When I send the request to update the status
        Then Then I should get an error with message "<message>"

        Examples:
            | id                                   | status | message          |
            | 2302f310-e60f-4434-b54b-d133eaa63a2c | ACTIVE | tenant not found |

    @failure
    Scenario Outline: Missing required field values
        Given I want to update the tenant's status to "<status>"
        When I send the request to update the status
        Then Then I should get a field error with message "<message>"
        Examples:
            | status              | message            |
            | something not valid | invalid status     |
            |                     | status is required |