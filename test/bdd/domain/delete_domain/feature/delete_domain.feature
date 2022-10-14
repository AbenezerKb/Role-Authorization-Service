Feature: Create Domain
    As a system user,
    i want to delete a Domain which is in may service
    so that i i will have more flexibility on my service domains.

  Background: 
    Given I have domain
      | name   |
      | system |

  Scenario Outline: Domain name is given
    Given i am a system user
    When i send the request:
      | name   |
      | <name> |
    Then the result should be successfull "<message>"

    Examples: 
      | name   | message |
      | system | true    |

  Scenario Outline: Required field are missing
    Given i am a system user
    When i send the request:
      | name   |
      | <name> |
    Then the result should be empty error "<message>"

    Examples: 
      | name | message                      |
      |      | domain name can not be blank |

  Scenario Outline: Domain Name not found
    Given i am a system user
    When i send the request:
      | name   |
      | <name> |
    Then the result should be not found error "<message>"

    Examples: 
      | name       | message                   |
      | testDomain | no record of domain found |
