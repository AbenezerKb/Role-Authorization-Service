Feature: Get Tenant users with their role
    As a system user 
    I want  to get tenant users with their roles
    So that I can get tenant users with their roles

  Background: 
    Given I have service with
      | name | service_id                           |
      | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |
    And A registered domain and tenant and role
      | domain | tenant_name | role          |
      | vendor | vendor_1    | manager,Admin |
    And tenant users and role
      | name  | roles         |
      | name1 | manager,Admin |
      | name2 | Admin         |

  @success
  Scenario Outline: Successfully get tenant users
    Given i am system user
    When I send request  "vendor_1"
    Then the result should be <name> and <roles>
      | name  | roles         |
      | name1 | manager,Admin |
      | name2 | Admin         |
