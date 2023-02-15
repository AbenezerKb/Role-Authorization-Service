Feature: Get Roles List
    As an authorized user,
    I want to get the list of all roles
    So that I can manage roles

    Background:
        Given I have service with
            | name | user_id                              |
            | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |
        And The role "Admin" is registered with the following permission in the tenant
            | name           | description    | effect | action                 | resource                      | domains |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |

    Scenario Outline: I get all the roles
        When I request to get all roles under my tenant
        Then I should get all roles in my tenant
            | name  |
            | Admin |

