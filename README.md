# C9Cli
Cloud Foundry - cf management

C9Cli helps maaintaing Cloud Foundry - Orgs/Spaces/Users/ASGs

C9Cli support following operations:
- init ( To initialize config )
- create-org
- create-space
- create-org-users
- create-space-users
- create-quotas

To Initialize configurations: 
C9Cli.exe -i init -e api.sys-domain -u <cf-login-user> -p <cf-login-pwd> -o <org> -s <space> -a true (to enable ASGs)
This will create config and sample ymls for further operation
  
  cf login config:
  ~/.C9Cli/mgmt/<endpoint>/config.yml - holds all cf login details created during init operation
  
  config.yml:
  ---
  ClusterDetails:
    EndPoint: api.sys-domain
    User: <user>
    Pwd: <pwd>
    Org: system
    Space: system
    EnableASG: true
  
  Yamls for Org/Space management:
  ~/C9Cli/mgmt/<endpoint>/Org.yml
  ~/C9Cli/mgmt/<endpoint>/Quota.yml
  ~/C9Cli/mgmt/<endpoint>/ProtectedResources.yml
  
Operators/Users are required to fill the the above Yaml files before executing create operations. The Orgs and Quotas listed ProtectedResources.yml are not touched during the  cli execution even if they are added to Org/Quota.yml
  
Org.yml:
---
Org:
  - Name: "org-test"
    Quota: "default"
    OrgUsers:
      LDAP:
        - Name: User1
          Role: OrgManager
        - Name: kxv6921
          Role: OrgAuditor
      SSO:
        - Name: User1
          Role: OrgManager
         UAA:
        - Name: User3
          Role: OrgAudidor
    Spaces:
      - Name: Space1
        SpaceUsers:
          LDAP:
            - Name: kxv6921
              Role: SpaceManager
          SSO:
            - Name: User3
              Role: SpaceAudidor
          UAA:
             - Name: User3
               Role: SpaceAudidor
      - Name: Space2
        SpaceUsers:
          LDAP:
            - Name: User1
              Role: SpaceManager
          SSO:
            - Name: User3
              Role: SpaceAudidor
          UAA:
            - Name: User3
              Role: SpaceAudidor

ProtectedResources.yml
---
  Org:
    - system
    - healthwatch
    - dynatrace
    - p-spring-cloud-services
    - splunk-nozzle-org
  quota:
    - default
 
Quota.yml
---
quota:
  - Name: default
    memory_limit: 10G
  - Name: small_quota
    memory_limit: 2048M
  - Name: medium_quota
    memory_limit: 2048M
  - Name: large_quota
    memory_limit: 2048M
    
    
Pending Items:
- Delete Orgs/Spaces/Users
- Assiging ASGs
