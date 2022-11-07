Feature: Create Domain
    As a system user,
    i want to create a new Domains
    so that i assign different domains for my permissions and tenants.

  Background: 
    Given I have service with
      | name | user_id                                |
      | sso2 | "a93fab67-1c11-4cdc-b410-f6fc728f592a" |

  Scenario Outline: All required fields are given
    Given i am a system user
    When i send the request:
      | name   |
      | <name> |
    Then the result should be successfull "<message>"

    Examples: 
      | name   | message |
      | system | true    |

  Scenario Outline: Required fields are missing
    Given i am a system user
    When i send the request:
      | name   |
      | <name> |
    Then the result should be empty error "<message>"

    Examples: 
      | name | message                      |
      |      | domain name can not be blank |
