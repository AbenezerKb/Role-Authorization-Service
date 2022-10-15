Feature: Create Tenant
    As a system user,
    i want to create a new Tenant
    so that i will have miltiple tenants in my service.

  Background: 
    Given I have service with
      | tenant_name | user_id                                |
      | sso2        | "a93fab67-1c11-4cdc-b410-f6fc728f592a" |

  Scenario Outline: All required fields are given
    Given i am a system user
    When i send the request:
      | tenant_name   |
      | <tenant_name> |
    Then the result should be successfull "<message>"

    Examples: 
      | tenant_name | message |
      | shop1       | true    |

  Scenario Outline: Required fields are missing
    Given i am a system user
    When i send the request:
      | tenant_name   |
      | <tenant_name> |
    Then the result should be empty error "<message>"

    Examples: 
      | tenant_name | message                      |
      |             | tenant name can not be blank |
