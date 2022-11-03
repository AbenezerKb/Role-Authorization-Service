Feature: update User Status

    As a admin
    I want to update user's status
    So that I can activate or deactivate user as appropriate

    Background:
        Given I have service with
            | name | user_id                              |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |

    @success
    Scenario Outline: Successfully Update User Status
        Given there is a user with the following details:
            | user_id   | service   |
            | <user_id> | <service> |
        And I want to update the user's status to "<status>"
        When I send the request to update the status
        Then the user status should update to "<status>"
        Examples:
            | user_id                              | service | status   |
            | 8e70df10-7957-426c-912a-9ed00f887e46 | sso     | ACTIVE   |
            | 8e70df10-7957-426c-912a-9ed00f887e46 | sso     | INACTIVE |

    @failure
    Scenario Outline: Failed User Status Update
        Given the user is not registered on the system
            | user_id   |
            | <user_id> |
        And I want to update the user's status to "<status>"
        When I send the request to update the status
        Then Then I should get an error with message "<message>"

        Examples:
            | user_id                              | status | message        |
            | 2302f310-e60f-4434-b54b-d133eaa63a2c | ACTIVE | user not found |

    @failure
    Scenario Outline: Missing required field values
        Given there is a user with the following details:
            | user_id   | service   |
            | <user_id> | <service> |
        And I want to update the user's status to "<status>"
        When I send the request to update the status
        Then Then I should get a field error with message "<message>"
        Examples:
            | user_id                              | status              | message             |
            | 00000000-0000-0000-0000-000000000000 | ACTIVE              | user id is required |
            | 8e70df10-7957-426c-912a-9ed00f887e46 | something not valid | invalid status      |
            | 8e70df10-7957-426c-912a-9ed00f887e46 |                     | status is required  |
