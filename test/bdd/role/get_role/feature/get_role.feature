Feature: Get Role Details
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

    @success
    Scenario Outline: Successfully Get Role Detail
        Given i have a role "admin" in tenant "vendor_1" with the following permissions:
            | name   | description   | effect   | action   | resource   | domains   |
            | <name> | <description> | <effect> | <action> | <resource> | <domains> |
        When I send the request to get the role details
        Then The request should be successfull
        Examples:
            | name           | description    | effect | action                 | resource                      | domains | status   |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  | ACTIVE   |
            | create service | create service | allow  | "admin:service:create" | "admin:service:createservice" | vendor  | INACTIVE |
    @failure
    Scenario Outline: Getting the role failed
        Given the role does not exists under the tenant "<id>"
        When I send the request to get the role details
        Then I should get an error with message "<message>"

        Examples:
            | id                                   | message        |
            | 2302f310-e60f-4434-b54b-d133eaa63a2c | role not found |

