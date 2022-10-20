Feature: Assign  role
    As an Admin,
    I want to assign role to my users
    So that I can assign it to a user

  Background: 
    Given I have service with
      | name | user_id                              |
      | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
    And A registered domain and tenant
      | domain | tenant_name |
      | vendor | vendor_1    |
    And I have user
      | user_id                              |
      | a93fab67-1c11-4cdc-b410-f6fc728f592a |
    And a permissions registered on the domain
      | name           | description    | effect | action                 | resource                      | domains |
      | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |
    And i have role
      | name  |
      | admin |

  Scenario Outline: successfully assign a new role to the user
    When I request to  assign role to user
      | user_id                              | role_name |
      | a93fab67-1c11-4cdc-b410-f6fc728f592a | admin     |
    Then the role should successfully be  assigned

  Scenario Outline: Required fields are missing
    When I request to assign  role to user while fields are missing
      | user_id   | role_name   |
      | <user_id> | <role_name> |
    Then my request should fail with "<message>"

    Examples: 
      | user_id                              | role_name | message          |
      | 00000000-0000-0000-0000-000000000000 | admin     | User id required |
      | a93fab67-1c11-4cdc-b410-f6fc728f592a |           | Role id required |
