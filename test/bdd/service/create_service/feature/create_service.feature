Feature: Create service

    As a system user,
    i want to create a new service
    so that i can user the service for authorization.

    Scenario Outline: All required fields are given
        Given i am a system user
        When i send the request:
            | name   | user_id   |
            | <name> | <user_id> |
        Then the result should be successfull "<message>"
        Examples:
            | name | user_id                              | message |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a | true    |

    Scenario Outline: Required fields are missing
        Given i am a system user
        When i send the request:
            | name   | user_id   |
            | <name> | <user_id> |
        Then the request should fail with error message "<message>"
        Examples:
            | name | user_id                              | message                  |
            |      | a93fab67-1c11-4cdc-b410-f6fc728f592a | service name is required |
            | sso  |                                      | user id is required      |

    Scenario Outline: Field inputs are invalid
        Given i am a system user
        When i send the request:
            | name   | user_id  |
            | <name> | <user_id |
        Then the request should fail with error message "<message>"
        Examples:
            | name | user_id                              | message                                  |
            | ss   | a93fab67-1c11-4cdc-b410-f6fc728f592a | name must be between 3 and 32 characters |
