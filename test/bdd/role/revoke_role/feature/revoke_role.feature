Feature: Assign  role
  As a user
  I want to revoke a user's role
  So that I can remove a user’s permission for executing a specific task on a protected resource

  Background:
    Given I have service with
      | name | user_id                              |
      | sso2 | a93fab67-1c11-4cdc-b410-f6fc728f592a |
    And A registered domain and tenant
      | domain | tenant_name |
      | vendor | vendor_1    |
    And I have user
      | user_id                              |
      | a93fab67-1c11-4cdc-b410-f6fc728f5921 |
    And a permissions registered on the domain
      | name           | description    | effect | action                 | resource                      | domains |
      | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | vendor  |
    And i have role
      | name  |
      | Admin |
    And I have user with role
      | user_id                              | role_name |
      | a93fab67-1c11-4cdc-b410-f6fc728f5921 | Admin     |

  Scenario Outline: successfully revoked  user role
    When I request to revoke the role of the user
      | user_id                              | role_name |
      | a93fab67-1c11-4cdc-b410-f6fc728f592a | Admin     |
    Then the role should successfully revoked

  Scenario Outline: Required fields are missing
    When I request to assign  role to user while fields are missing
      | user_id   | role_name   |
      | <user_id> | <role_name> |
    Then my request should fail with "<message>"

    Examples:
      | user_id                              | role_name | message                     |
      | 00000000-0000-0000-0000-000000000000 | Admin     | User id required            |
      | a93fab67-1c11-4cdc-b410-f6fc728f592a |           | Role id or name is require |
