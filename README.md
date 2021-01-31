# C9Cli
Cloud Foundry - cf management

C9Cli helps maaintaing Cloud Foundry - Orgs/Spaces/Users/ASGs

C9Cli support following operations:
- init ( To initialize config )
- create-org
- create-space
- create-org-user
- create-space-user
- create-quota

To Initialize configurations: 
C9Cli.exe -i init -e api.sys-domain -u <cf-login-user> -p <cf-login-pwd> -o <org> -s <space> -a true (to enable ASGs)
This will create config and sample ymls for further operation
  
  cf login config:
  ~/.C9Cli/mgmt/<endpoint>/config.yml - holds all cf login details created during init operation
  
  Yamls for Org/Space management:
  ~/C9Cli/mgmt/<endpoint>/Org.yml
  ~/C9Cli/mgmt/<endpoint>/Quota.yml
  ~/C9Cli/mgmt/<endpoint>/ProtectedResources.yml
  
Operators/Users are required to fill the the above Yaml files before executing create operations. The Orgs and Quotas listed ProtectedResources.yml are not touched during the  cli execution even if they are added to Org/Quota.yml
    
Pending Items:
- Delete Orgs/Spaces/Users
- Assiging ASGs
