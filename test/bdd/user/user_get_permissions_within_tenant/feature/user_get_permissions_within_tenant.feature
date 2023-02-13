Feature: Get User Permissions Within Tenant
    As an user,
    I want to get the list of all permissions i have within a tenant
    So that I know what i am allowed and not

    Background:
        Given I have service with
            | name | user_id                              |
            | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |
        And I have a role "Admin" in tenant "vendor_1" with the following permissions
            | name           | description    | effect | action | resource                    | fields | domains |
            | delete service | delete service | allow  | Post   | admin:service:deleteservice | *      | vendor  |
        And I am registered on the system
            | user_id                              |
            | 4a2650b5-1c1c-437b-adc5-7d5d24a91126 |
        And I am granted an "Admin" role

    @success
    Scenario: I get my permissions within the tenant
        Given I want to get my permissions
        When I request to get my permissions
        Then The Request should be successfull
