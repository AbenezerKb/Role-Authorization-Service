Feature: Delete Service

    As a system user,
    I want to delete the service
    So that I can remove remove my service from the system

    # @success
    # Scenario: Successful Delete Service
    #     Given I have a registered service
    #         | name | user_id                              |
    #         | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |
    #     When I delete the service:
    #         | action | subject                              | resource |
    #         | *      | a93fab67-1c11-4cdc-b410-f6fc728f592a | *        |
    #     Then The service should be deleted

    # @failure
    # Scenario Outline: no service
    #     When I delete the service with id "<id>"
    #     Then the request should fail with error message "<message>"
    #     Examples:
    #         | id                                   | message                    |
    #         | 60d56419-c2e9-4ee4-951f-04644d245ee3 | no record of service found |

    Scenario Outline: Required fields are missing
        Given I have a registered service
            | name      | user_id                              |
            | ride-plus | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        When I delete the service with input:
            | action   | subject   | resource   |
            | <action> | <subject> | <resource> |
        Then the request should fail with field error message "<message>"
        Examples:
            | action | subject                              | resource | message              |
            |        | a93fab67-1c11-4cdc-b410-f6fc728f592a | *        | action is required   |
            | *      |                                      | *        | subject is required  |
            | *      | a93fab67-1c11-4cdc-b410-f6fc728f592a |          | resource is required |
