Feature: Create service

    As a system user,
    i want to create a new service
    so that i can user the service for authorization.

    Scenario Outline: All required fields are given
        Given i am a system user
        When i send the request:
            | name   |
            | <name> |
        Then the result should be successfull "<message>"
        Examples:
            | name | message |
            | sso  | true    |

    Scenario Outline: Required fields are missing
        Given i am a system user
        When i send the request:
            | name   |
            | <name> |
        Then the request should fail with error message "<message>"
        Examples:
            | name | message                  |
            |      | service name is required |

    Scenario Outline: Field inputs are invalid
        Given i am a system user
        When i send the request:
            | name   |
            | <name> |
        Then the request should fail with error message "<message>"
        Examples:
            | name | message                                  |
            | ss   | name must be between 3 and 32 characters |
