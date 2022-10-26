Feature: Authorize user
    As an system
    I want to authorize a user
    So that I can protect the resource from unprotected users

    Background:
        Given I have service with
            | name | user_id                              |
            | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |
        And i have a role "admin" in tenant "vendor_1" with the permissions below
            | name           | description    | effect | action | resource                    | domains |
            | delete service | delete service | allow  | Post   | admin:service:deleteservice | vendor  |
        And There is a user registered on the system:
            | user_id                              |
            | 4a2650b5-1c1c-437b-adc5-7d5d24a91126 |

    @success
    Scenario: I successfully authorize the user
        Given The user is granted with the "admin" role
        And I  want to authorize the user to perform the belowe action on the resource:
            | subject                              | action | resource                    | tenant   |
            | 4a2650b5-1c1c-437b-adc5-7d5d24a91126 | Post   | admin:service:deleteservice | vendor_1 |
        When i send a request to authorize the user
        Then the user should be allowed
    @failer
    Scenario: The user should be denied
        Given the user is not granted any role
        And I  want to authorize the user to perform the belowe action on the resource:
            | subject                              | action | resource                    | tenant   |
            | 765650b5-1c1c-437b-adc5-7d5d24a91126 | Post   | admin:service:deleteservice | vendor_1 |
        When i send a request to authorize the user
        Then the user should be denied