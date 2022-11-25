Feature: Get Pemrission Details
    As an admin
    I want to fetch particular role
    So that I can see the details and manage that particular role

    Background:
        Given I have service with
            | name | user_id                              |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |
        And A permission registered on the domain
            | name           | description    | effect | action                 | fields | resource                      |
            | delete service | delete service | allow  | "admin:service:delete" | *      | "admin:service:deleteservice" |

    @success
    Scenario Outline: Successfully Get Permission Detail
        Given I want to get the permission detail
        And I send the request
        Then The request should be successfull
        And I should get the permission detail
            | name           | description    | effect | action                 | fields | resource                      | status |
            | delete service | delete service | allow  | "admin:service:delete" | *      | "admin:service:deleteservice" | ACTIVE |
    @failure
    Scenario Outline: Getting the permission detail failed
        Given the permission does not exists "<id>"
        When I send the request to get the permission details
        Then I should get an error with message "<message>"

        Examples:
            | id                                   | message                    |
            | 2302f310-e60f-4434-b54b-d133eaa63a2c | permission does not exists |

