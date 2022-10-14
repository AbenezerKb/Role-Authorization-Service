Feature: Create Domain
    As a system user,
    i want to create a new domains
    so that i assign diffrent domains for my service.
    Background: Create service and domain
        Given I have a registered service
            | name | user_id                              |
            | sso  | a93fab67-1c11-4cdc-b410-f6fc728f592a |
        And A registered domain
            | name   |
            | system |
    @success
    Scenario Outline: Successful Create Permmission
        When I create a permmission in the domain:
            | name   | description   | effect   | action   | resource   | domains   |
            | <name> | <description> | <effect> | <action> | <resource> | <domains> |
        Then The service should be deleted
        Examples:
            | name           | description    | effect | action                 | resource                      | domains      |
            | delete service | delete service | allow  | "admin:service:delete" | "admin:service:deleteservice" | system |

    Scenario Outline: Required fields are missing
        When I create a permmission in the domain:
            | name   | description   | effect   | action   | resource   | domains   |
            | <name> | <description> | <effect> | <action> | <resource> | <domains> |
        Then The request should fail with error "<message>"
        Examples:
            | name           | description    | effect | action               | resource                    | domains | message                            |
            |                | delete service | allow  | admin:service:delete | admin:service:deleteservice | system  | permission name is required        |
            | delete service |                | allow  | admin:service:delete | admin:service:deleteservice | system  | permission description is required |
            | delete service | delete service |        | admin:service:delete | admin:service:deleteservice | system  | statement effect is required       |
            | delete service | delete service | allow  |                      | admin:service:deleteservice | system  | statement action is required       |
            | delete service | delete service | allow  | admin:service:delete |                             | system  | statement resource is required     |
