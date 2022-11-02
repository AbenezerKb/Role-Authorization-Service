Feature: update service Status

    As a admin
    I want to update services's status
    So that I can activate or deactivate service's
    Background: the service is registered on the system
        Given the service is registered on the system
            | name | user_id                              |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |

    @success
    Scenario Outline: Successfully update service status
        When I update the service's status to "<status>"
        Then the service status should update to "<status>"
        Examples:
            | name | status   |
            | sso  | INACTIVE |
            | sso  | ACTIVE   |

    @failure
    Scenario Outline: Invalid Status
        When I update the service's status to "<status>"
        Then Then I should get error with message "<message>"
        Examples:
            | status              | message            |
            | something not valid | invalid status     |
            |                     | status is required |

