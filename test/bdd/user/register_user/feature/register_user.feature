Feature: Register User
    As a system user,
    i want to Register a new user
    so that i will add miltiple users in my service.

    Background:
        Given I have service with
            | tenant_name | user_id                              |
            | sso2        | a93fab67-1c11-4cdc-b410-f6fc728f592a |

    Scenario Outline: successfully register the new user
        Given I have the user id:
            | user_id                              |
            | 4a2650b5-1c1c-437b-adc5-7d5d24a91126 |
        When i send the request to add the user
        Then the request should be successfull

    Scenario Outline: Required fields are missing
        Given I have the user id:
            | user_id   |
            | <user_id> |
        When i send the request to add the user
        Then the request should fail with error "<message>"
        Examples:
            | user_id | message       |
            |         | invalid input |
    Scenario Outline: user already exists in the service
        Given I have a user registered on my service:
            | user_id   |
            | <user_id> |
        And I want to add the same user again:
            | user_id   |
            | <user_id> |
        When i send the request to add the user
        Then the request should fail with error "<message>"
        Examples:
            | user_id                              | message                          |
            | 4a2650b5-1c1c-437b-adc5-7d5d24a91126 | user with this id already exists |

