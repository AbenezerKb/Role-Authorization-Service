Feature: Get Permission List
    As an authorized user,
    I want to get the list of all possible permissions
    So that I can use them to create roles

    Background:
        Given I have service with
            | name | user_id                              |
            | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain and tenant
            | domain | tenant_name |
            | vendor | vendor_1    |
        And a permissions registered on the domain
            | name           | description    | effect | action                 | resource                      | domains |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |

    Scenario Outline: I get all permissions
        When I request to get all permissions under my tenant
        Then I should get all permissions in my tenant
            | name           | description    | effect | action                 | resource                      | 
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | 


