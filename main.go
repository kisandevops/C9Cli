package main

import "C"
import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"gopkg.in/yaml.v2"
)

type Orglist struct {
	Org []struct {
		Name     string `yaml:"Name"`
		Quota    string `yaml:"Quota"`
		OrgUsers struct {
			LDAP []struct {
				Name string `yaml:"Name"`
				Role string `yaml:"Role"`
			} `yaml:"LDAP"`
			SSO []struct {
				Name string `yaml:"Name"`
				Role string `yaml:"Role"`
			} `yaml:"SSO"`
			UAA []struct {
				Name string `yaml:"Name"`
				Role string `yaml:"Role"`
			} `yaml:"UAA"`
		} `yaml:"OrgUsers"`
		Spaces []struct {
			Name       string `yaml:"Name"`
			SpaceUsers struct {
				LDAP []struct {
					Name string `yaml:"Name"`
					Role string `yaml:"Role"`
				} `yaml:"LDAP"`
				SSO []struct {
					Name string `yaml:"Name"`
					Role string `yaml:"Role"`
				} `yaml:"SSO"`
				UAA []struct {
					Name string `yaml:"Name"`
					Role string `yaml:"Role"`
				} `yaml:"UAA"`
			} `yaml:"SpaceUsers"`
		} `yaml:"Spaces"`
	} `yaml:"Org"`
}
type Quotalist struct {
	Quota []struct {
		Name        string `yaml:"Name"`
		MemoryLimit string `yaml:"memory_limit"`
	} `yaml:"quota"`
}
type ProtectedList struct {
	Org   []string `yaml:"Org"`
	Quota []string `yaml:"quota"`
}
type InitClusterConfigVals struct {
	ClusterDetails struct {
		EndPoint  string `yaml:"EndPoint"`
		User      string `yaml:"User"`
		Pwd       string `yaml:"Pwd"`
		Org       string `yaml:"Org"`
		Space     string `yaml:"Space"`
		EnableASG bool   `yaml:"EnableASG"`
	} `yaml:"ClusterDetails"`
}

func main()  {

	var endpoint, user, pwd, org, space, asg, operation string

	flag.StringVar(&endpoint, "e", "api.sys-domain", "Use with init operation, Provide PCF Endpoint")
	flag.StringVar(&user, "u", "user", "Use with init operation, Provide UserName")
	flag.StringVar(&pwd, "p", "pwd", "Use with init operation, Provide Password")
	flag.StringVar(&org, "o", "org", "Use with init operation, Provide Org")
	flag.StringVar(&space, "s", "space", "Use with init operation, Provide Space")
	flag.StringVar(&asg, "a", "true", "Use with init operation, Enable ASGs ?.")
	flag.StringVar(&operation, "i", "init", "Provide Operation to be performed: init, create-{org,space,org-user,space-user,quota, ")
	flag.Parse()

	ClusterName := strings.ReplaceAll(endpoint, ".", "-")

	fmt.Printf("Operation: %v\n", operation)

	if operation == "init" {

		fmt.Println("Initializing C9Cli")

		fmt.Printf("ClusterName: %v\n", ClusterName)
		fmt.Printf("EndPoint: %v\n", endpoint)
		fmt.Printf("User: %v\n", user)
		fmt.Printf("Pwd: %v\n", pwd)
		fmt.Printf("Org: %v\n", org)
		fmt.Printf("Space: %v\n", space)
		fmt.Printf("EnableASG: %v\n", asg)

		Init(ClusterName, endpoint, user, pwd, org, space, asg)

	} else if operation == "create-org"{

		fmt.Printf("ClusterName: %v\n", ClusterName)
		SetupConnection (ClusterName)
		CreateOrUpdateOrgs (ClusterName)

	} else if operation == "create-quota" {

		fmt.Printf("ClusterName: %v\n", ClusterName)
		SetupConnection (ClusterName)
		CreateOrUpdateQuotas(ClusterName)

	} else if operation == "create-org-user" {

		fmt.Printf("ClusterName: %v\n", ClusterName)
		SetupConnection(ClusterName)
		CreateOrUpdateOrgUsers(ClusterName)
	} else if operation == "create-space"{

		fmt.Printf("ClusterName: %v\n", ClusterName)
		SetupConnection (ClusterName)
		CreateOrUpdateSpaces (ClusterName)

	} else if operation == "create-space-user"{

		fmt.Printf("ClusterName: %v\n", ClusterName)
		SetupConnection (ClusterName)
		CreateOrUpdateSpaceUsers (ClusterName)

	} else {
		fmt.Println("Provide Valid input operation")
	}
}
func SetupConnection(clustername string) error {

	var InitClusterConfigVals InitClusterConfigVals
	ConfigFile := "~/.C9Cli"+"/mgmt/"+clustername+"/config.yml"

	fileConfigYml, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileConfigYml), &InitClusterConfigVals)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Endpoint: %v\n", InitClusterConfigVals.ClusterDetails.EndPoint)
	fmt.Printf("User: %v\n", InitClusterConfigVals.ClusterDetails.User)
	fmt.Printf("Pwd: %v\n", InitClusterConfigVals.ClusterDetails.Pwd)
	fmt.Printf("Org: %v\n", InitClusterConfigVals.ClusterDetails.Org)
	fmt.Printf("Space: %v\n", InitClusterConfigVals.ClusterDetails.Space)
    //fmt.Println(InitClusterConfigVals.ClusterDetails.EndPoint)

	cmd := exec.Command("cf", "login", "-a", InitClusterConfigVals.ClusterDetails.EndPoint, "-u", InitClusterConfigVals.ClusterDetails.User, "-p", InitClusterConfigVals.ClusterDetails.Pwd, "-o", InitClusterConfigVals.ClusterDetails.Org, "-s", InitClusterConfigVals.ClusterDetails.Space)
	if _, err := cmd.Output(); err != nil{
		fmt.Println("Connection failed")
		fmt.Println("command: ", cmd)
		fmt.Println("Err Code: ", err)
		panic(err)
	} else {
		fmt.Println("Connection Passed")
		fmt.Println("command: ", cmd)
		fmt.Println(cmd.Stdout)
	}
	return err
}
func CreateOrUpdateOrgs(clustername string) error {

	var Orgs Orglist
	var ProtectedOrgs ProtectedList

	OrgsYml := "~/C9Cli/mgmt/"+clustername+"/Org.yml"
	fileOrgYml, err := ioutil.ReadFile(OrgsYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileOrgYml), &Orgs)
	if err != nil {
		panic(err)
	}


	ProtectedOrgsYml := "~/C9Cli/mgmt/"+clustername+"/ProtectedResources.yml"
	fileProtectedYml, err := ioutil.ReadFile(ProtectedOrgsYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileProtectedYml), &ProtectedOrgs)
	if err != nil {
		panic(err)
	}


	LenOrgs := len(Orgs.Org)
	LenProtectedOrgs := len(ProtectedOrgs.Org)

	for i := 0; i < LenOrgs; i++ {
		var count, totalcount int
		fmt.Println("Org: ", Orgs.Org[i].Name)
		for p := 0; p < LenProtectedOrgs; p++ {
			fmt.Println("Protected Org: ", ProtectedOrgs.Org[p])
			if ProtectedOrgs.Org[p] == Orgs.Org[i].Name {
				count = 1
			} else {
				count = 0
			}
		}
		totalcount = totalcount + count
		fmt.Println(totalcount)

		if totalcount == 0 {
			fmt.Println("This is not Protected Org")
			guid := exec.Command("cf", "org", Orgs.Org[i].Name, "--guid")
			if _, err := guid.Output(); err == nil{

				fmt.Println("command: ", guid)
				fmt.Println("Org exists: ", guid.Stdout)
				fmt.Println("Updating Org quota")
				SetQuota := exec.Command("cf", "set-quota", Orgs.Org[i].Name, Orgs.Org[i].Quota)
				if _, err := SetQuota.Output(); err != nil{
					fmt.Println("command: ", SetQuota)
					fmt.Println("Err Code: ", err)
				} else {
					fmt.Println("command: ", SetQuota)
					fmt.Println(SetQuota.Stdout)
				}
			} else {
				fmt.Println("command: ", guid)
				fmt.Println("Err Code: ", err)
				fmt.Println("Pulling Guid Id: ", guid.Stdout)
				fmt.Println("Org doesn't exists, Creating Org")
				createorg := exec.Command("cf", "create-org", Orgs.Org[i].Name)
				if _, err := createorg.Output(); err != nil{
					fmt.Println("command: ", createorg)
					fmt.Println("Err Code: ", err)
				} else {
					fmt.Println("command: ", createorg)
					fmt.Println(createorg.Stdout)
				}
				attachquota := exec.Command("cf", "set-quota", Orgs.Org[i].Name, Orgs.Org[i].Quota)
				if _, err := attachquota.Output(); err != nil{
					fmt.Println("command: ", attachquota)
					fmt.Println("Err Code: ", err)
				} else {
					fmt.Println("command: ", attachquota)
					fmt.Println(attachquota.Stdout)
				}
			}
		} else {
			fmt.Println("This is a protected Org")
		}
	}
	return err
}
func CreateOrUpdateSpaces(clustername string) error {

	var Orgs Orglist
	var ProtectedOrgs ProtectedList


	OrgsYml := "~/C9Cli/mgmt/"+clustername+"/Org.yml"
	fileOrgYml, err := ioutil.ReadFile(OrgsYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileOrgYml), &Orgs)
	if err != nil {
		panic(err)
	}

	ProtectedOrgsYml := "~/C9Cli/mgmt/"+clustername+"/ProtectedResources.yml"
	fileProtectedYml, err := ioutil.ReadFile(ProtectedOrgsYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileProtectedYml), &ProtectedOrgs)
	if err != nil {
		panic(err)
	}

	LenOrgs := len(Orgs.Org)
	LenProtectedOrgs := len(ProtectedOrgs.Org)


	for i := 0; i < LenOrgs; i++ {

		var count, totalcount int
		fmt.Println("Org: ", Orgs.Org[i].Name)
		for p := 0; p < LenProtectedOrgs; p++ {
			fmt.Println("Protected Org: ", ProtectedOrgs.Org[p])
			if ProtectedOrgs.Org[p] == Orgs.Org[i].Name {
				count = 1
			} else {
				count = 0
				}
			}
		totalcount = totalcount + count
		fmt.Println(totalcount)

		if totalcount == 0 {
			fmt.Println("This is not Protected Org")
			guid := exec.Command("cf", "org", Orgs.Org[i].Name, "--guid")

			if _, err := guid.Output(); err == nil {

				fmt.Println("command: ", guid)
				fmt.Println("Org exists: ", guid.Stdout)
				SpaceLen := len(Orgs.Org[i].Spaces)

				TargetOrg := exec.Command("cf", "t", "-o", Orgs.Org[i].Name)
				if _, err := TargetOrg.Output(); err == nil {
					fmt.Println("command: ", TargetOrg)
					fmt.Println("Targeting: ", TargetOrg.Stdout)
				} else {
					fmt.Println("command: ", TargetOrg)
					fmt.Println("Err Code: ", err)
				}

				for j := 0; j < SpaceLen; j++ {

					fmt.Println("Creating Spaces")

					guid = exec.Command("cf", "space", Orgs.Org[i].Spaces[j].Name, "--guid")

					if _, err := guid.Output(); err == nil{
						fmt.Println("command: ", guid)
						fmt.Println("Space exists: ", guid.Stdout)
					} else {
						fmt.Println("command: ", guid)
						fmt.Println("Pulling Space Guid ID: ", guid.Stdout )
						fmt.Println("Creating Space")

						CreateSpace := exec.Command("cf", "create-space", Orgs.Org[i].Spaces[j].Name, "-o", Orgs.Org[i].Name)

						if _, err := CreateSpace.Output(); err != nil {
							fmt.Println("command: ", CreateSpace)
							fmt.Println("Err Code: ", err)
						} else {
							fmt.Println("command: ", CreateSpace)
							fmt.Println(CreateSpace.Stdout)
						}
					}
				}
			} else {
				fmt.Println("command: ", guid )
				fmt.Println("Err Code: ", err)
				fmt.Println("Org doesn't exists, Please create Org")
			}
		} else {
			fmt.Println("This is a protected Org")
		}
	}
	return err
}
func CreateOrUpdateQuotas(clustername string) error {

	var Quotas Quotalist
	var ProtectedQuota ProtectedList


	QuotaYml := "~/C9Cli/mgmt/"+clustername+"/Quota.yml"
	fileQuotaYml, err := ioutil.ReadFile(QuotaYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileQuotaYml), &Quotas)
	if err != nil {
		panic(err)
	}

	ProtectedQuotasYml := "~/C9Cli/mgmt/"+clustername+"/ProtectedResources.yml"
	fileProtectedQYml, err := ioutil.ReadFile(ProtectedQuotasYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileProtectedQYml), &ProtectedQuota)
	if err != nil {
		panic(err)
	}

	LenQuota := len(Quotas.Quota)
	LenProtectedQuota := len(ProtectedQuota.Quota)

	for i := 0; i < LenQuota; i++ {

		var count, totalcount int
		fmt.Println("Quota: ", Quotas.Quota[i].Name)

		for p := 0; p < LenProtectedQuota; p++ {
			fmt.Println("Protected Quota: ", ProtectedQuota.Quota[p])
			if strings.Trim(ProtectedQuota.Quota[p], "") == strings.Trim(Quotas.Quota[i].Name, "") {
				count = 1
			} else {
				count = 0
			}
		}
		totalcount = totalcount + count
		fmt.Println(totalcount)

		if totalcount == 0 {

			fmt.Println("This is not Protected Quota")
			Quotadetails := exec.Command("cf", "quota", Quotas.Quota[i].Name)

			if _, err := Quotadetails.Output(); err != nil{
				fmt.Println("command: ", Quotadetails)
				fmt.Println("Err Code: ", err)
				//fmt.Println("Quota Doesn't exits: ", Quotadetails.Stdout)
				fmt.Println("Creating Quota")
				cmd := exec.Command("cf", "create-quota", Quotas.Quota[i].Name, "-m", Quotas.Quota[i].MemoryLimit, "-i", "-1", "-r", "-1", "-s", "-1", "-a", "-1", "--allow-paid-service-plans")
				if _, err := cmd.Output(); err != nil{
					fmt.Println("command: ", cmd)
					fmt.Println("Err Code: ", err)
				} else {
					fmt.Println("command: ", cmd)
					fmt.Println(cmd.Stdout)
				}
				QuotaGet := exec.Command("cf", "quota", Quotas.Quota[i].Name)
				if _, err := QuotaGet.Output(); err != nil{
					fmt.Println("command: ", QuotaGet)
					fmt.Println("Err Code: ", err)
				} else {
					fmt.Println("command: ", QuotaGet)
					fmt.Println(QuotaGet.Stdout)
				}
			} else {
				fmt.Println("command: ", Quotadetails)
				fmt.Println("Quota exists: ", Quotadetails.Stdout)
				fmt.Println("Updating Org")
				cmd := exec.Command("cf", "update-quota", Quotas.Quota[i].Name, "-m", Quotas.Quota[i].MemoryLimit, "-i", "-1", "-r", "-1", "-s", "-1", "-a", "-1", "--allow-paid-service-plans")
				if _, err := cmd.Output(); err != nil{
					fmt.Println("command: ", cmd)
					fmt.Println("Err Code: ", err)
				} else {
					fmt.Println("command: ", cmd)
					fmt.Println(cmd.Stdout)
				}
				QuotaGet := exec.Command("cf", "quota", Quotas.Quota[i].Name)
				if _, err := QuotaGet.Output(); err != nil{
					fmt.Println("command: ", QuotaGet)
					fmt.Println("Err Code: ", err)
				} else {
					fmt.Println("command: ", QuotaGet)
					fmt.Println(QuotaGet.Stdout)
				}
			}
		} else {
			fmt.Println("This is a protected Org")
		}
	}
	return err
}
func CreateOrUpdateOrgUsers(clustername string) error {

	var Orgs Orglist
	var ProtectedOrgs ProtectedList

	OrgsYml := "~/C9Cli/mgmt/"+clustername+"/Org.yml"
	fileOrgYml, err := ioutil.ReadFile(OrgsYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileOrgYml), &Orgs)
	if err != nil {
		panic(err)
	}

	ProtectedOrgsYml := "~/C9Cli/mgmt/"+clustername+"/ProtectedResources.yml"
	fileProtectedYml, err := ioutil.ReadFile(ProtectedOrgsYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileProtectedYml), &ProtectedOrgs)
	if err != nil {
		panic(err)
	}


	LenProtectedOrgs := len(ProtectedOrgs.Org)
	LenOrgs := len(Orgs.Org)

	for i := 0; i < LenOrgs; i++ {

		var count, totalcount int
		fmt.Println("Org: ", Orgs.Org[i].Name)
		for p := 0; p < LenProtectedOrgs; p++ {
			fmt.Println("Protected Org: ", ProtectedOrgs.Org[p])
			if ProtectedOrgs.Org[p] == Orgs.Org[i].Name {
				count = 1
			} else {
				count = 0
			}
		}
		totalcount = totalcount + count
		fmt.Println(totalcount)

		if totalcount == 0 {

			guid := exec.Command("cf", "org", Orgs.Org[i].Name, "--guid")
			if _, err := guid.Output(); err == nil{

				//fmt.Println("Err: ", err)

				fmt.Println("command: ", guid)
				fmt.Println("Org exists: ", guid.Stdout)
				fmt.Println("Updating Org Users")
				fmt.Println("Updating LDAP Users")

				LDAPOrgLen := len(Orgs.Org[i].OrgUsers.LDAP)

				for j := 0; j < LDAPOrgLen; j++ {

					//cf set-org-role USER ORGNAME ROLE(like OrgManager)

					cmd := exec.Command("cf", "set-org-role", Orgs.Org[i].OrgUsers.LDAP[j].Name, Orgs.Org[i].Name, Orgs.Org[i].OrgUsers.LDAP[j].Role)

					if _, err := cmd.Output(); err != nil{
						fmt.Println("command: ", cmd)
						fmt.Println("Err Code: ", err)
					} else {
						fmt.Println("command: ", cmd)
						fmt.Println(cmd.Stdout)
					}
				}

				fmt.Println("Updating UAA Users")

				UAAOrgLen := len(Orgs.Org[i].OrgUsers.UAA)

				for k := 0; k < UAAOrgLen; k++ {

					//cf set-org-role USER ORGNAME ROLE(like OrgManager)

					cmd := exec.Command("cf", "set-org-role", Orgs.Org[i].OrgUsers.UAA[k].Name, Orgs.Org[i].Name, Orgs.Org[i].OrgUsers.UAA[k].Role)
					if _, err := cmd.Output(); err != nil{
						fmt.Println("command: ", cmd)
						fmt.Println("Err Code: ", err)
					} else {
						fmt.Println("command: ", cmd)
						fmt.Println(cmd.Stdout)
					}
				}

				fmt.Println("Updating SSO Users")

				SSOOrgLen := len(Orgs.Org[i].OrgUsers.SSO)

				for l := 0; l < SSOOrgLen; l++ {
					//cf set-org-role USER ORGNAME ROLE(like OrgManager)
					cmd := exec.Command("cf", "set-org-role", Orgs.Org[i].OrgUsers.SSO[l].Name, Orgs.Org[i].Name, Orgs.Org[i].OrgUsers.SSO[l].Role)
					if _, err := cmd.Output(); err != nil{
						fmt.Println("command: ", cmd)
						fmt.Println("Err Code: ", err)
					} else {
						fmt.Println("command: ", cmd)
						fmt.Println(cmd.Stdout)
					}
				}
			} else {
				fmt.Println("command: ", guid)
				fmt.Println("Err Code: ", err)
				fmt.Println("Pulling Org Guid Id: ", guid.Stdout)
				fmt.Println("Please Create Org")
			}
		} else {
			fmt.Println("This is a protected Org")
		}
	}
	return err
}
func CreateOrUpdateSpaceUsers(clustername string) error {

	var Orgs Orglist
	var ProtectedOrgs ProtectedList

	OrgsYml := "~/C9Cli/mgmt/"+clustername+"/Org.yml"
	fileOrgYml, err := ioutil.ReadFile(OrgsYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileOrgYml), &Orgs)
	if err != nil {
		panic(err)
	}

	ProtectedOrgsYml := "~/C9Cli/mgmt/"+clustername+"/ProtectedResources.yml"
	fileProtectedYml, err := ioutil.ReadFile(ProtectedOrgsYml)

	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal([]byte(fileProtectedYml), &ProtectedOrgs)
	if err != nil {
		panic(err)
	}

	LenProtectedOrgs := len(ProtectedOrgs.Org)
	LenOrgs := len(Orgs.Org)

	for i := 0; i < LenOrgs; i++ {

		var count, totalcount int
		fmt.Println("Org: ", Orgs.Org[i].Name)
		for p := 0; p < LenProtectedOrgs; p++ {
			fmt.Println("Protected Org: ", ProtectedOrgs.Org[p])
			if ProtectedOrgs.Org[p] == Orgs.Org[i].Name {
				count = 1
			} else {
				count = 0
			}
		}
		totalcount = totalcount + count
		fmt.Println(totalcount)

		if totalcount == 0 {
			guid := exec.Command("cf", "org", Orgs.Org[i].Name, "--guid")

			if _, err := guid.Output(); err == nil {

				fmt.Println("command: ", guid)
				fmt.Println("Org exists: ", guid.Stdout)
				targetOrg := exec.Command("cf", "t", "-o", Orgs.Org[i].Name)
				if _, err := targetOrg.Output(); err == nil {
					fmt.Println("command: ", targetOrg)
					fmt.Println("Targeted Org: ", targetOrg.Stdout)
				} else {
					fmt.Println("command: ", targetOrg)
					fmt.Println("Err Code: ", targetOrg.Stderr)
				}
				SpaceLen := len(Orgs.Org[i].Spaces)

				for j := 0; j < SpaceLen; j++ {

					guid = exec.Command("cf", "space", Orgs.Org[i].Spaces[j].Name, "--guid")
					if _, err := guid.Output(); err == nil {
						fmt.Println("command: ", guid)
						fmt.Println("Space exists: ", guid.Stdout)
						fmt.Println("Creating Space users")

						fmt.Println("Updating LDAP Users")

						LDAPSpaceUsersLen := len(Orgs.Org[i].Spaces[j].SpaceUsers.LDAP)

						for k := 0; k < LDAPSpaceUsersLen; k++ {
							cmd := exec.Command("cf", "set-space-role", Orgs.Org[i].Spaces[j].SpaceUsers.LDAP[k].Name, Orgs.Org[i].Name, Orgs.Org[i].Spaces[j].Name, Orgs.Org[i].Spaces[j].SpaceUsers.LDAP[k].Role)
							if _, err := cmd.Output(); err != nil{
								fmt.Println("command: ", cmd)
								fmt.Println("Err Code: ", err)
							} else {
								fmt.Println("command: ", cmd)
								fmt.Println(cmd.Stdout)
							}
						}

						fmt.Println("Updating UAA Users")

						UAASpaceUsersLen := len(Orgs.Org[i].Spaces[j].SpaceUsers.UAA)

						for l := 0; l < UAASpaceUsersLen; l++ {
							cmd := exec.Command("cf", "set-space-role", Orgs.Org[i].Spaces[j].SpaceUsers.UAA[l].Name, Orgs.Org[i].Name, Orgs.Org[i].Spaces[j].Name, Orgs.Org[i].Spaces[j].SpaceUsers.UAA[l].Role)
							if _, err := cmd.Output(); err != nil{
								fmt.Println("command: ", cmd)
								fmt.Println("Err Code: ", err)
							} else {
								fmt.Println("command: ", cmd)
								fmt.Println(cmd.Stdout)
							}
						}

						fmt.Println("Updating SSO Users")

						SSOSpaceUsersLen := len(Orgs.Org[i].Spaces[j].SpaceUsers.SSO)
						for m := 0; m < SSOSpaceUsersLen; m++ {
							cmd := exec.Command("cf", "set-space-role", Orgs.Org[i].Spaces[j].SpaceUsers.SSO[m].Name, Orgs.Org[i].Name, Orgs.Org[i].Spaces[j].Name, Orgs.Org[i].Spaces[j].SpaceUsers.SSO[m].Role)
							if _, err := cmd.Output(); err != nil{
								fmt.Println("command: ", cmd)
								fmt.Println("Err Code: ", err)
							} else {
								fmt.Println("command: ", cmd)
								fmt.Println(cmd.Stdout)
							}
						}

					} else {
						fmt.Println("command: ",guid)
						fmt.Println("Err Code: ", err)
						fmt.Println("Space doesn't exists, Please create Space")
					}
				}
			} else {
				fmt.Println("command: ", guid)
				fmt.Println("Err Code: ", err)
				fmt.Println("Org doesn't exists, Please create Org")
			}
		}
	}
	return err
}
func Init(clustername string, endpoint string, user string, pwd string, org string, space string, asg string) (err error) {

	type ClusterDetails struct {
		EndPoint         string `yaml:"EndPoint"`
		User         string `yaml:"User"`
		Pwd        string `yaml:"Pwd"`
		Org            string `yaml:"Org"`
		Space string  `yaml:"Space"`
		EnableASG     string `yaml:"EnableASG"`
	}

	// Variables - host, namespace
	ConfigPath := "~/.C9Cli"+"/mgmt/"+clustername
	mgmtpath := "~/C9Cli/mgmt/"+clustername
	ASGPath := "~/C9Cli/mgmt/"+clustername+"/ASGs"
	//ConfigFile := "~/C9Cli/mgmt/"+clustername+"/config.yml"
	OrgsYml := "~/C9Cli/mgmt/"+clustername+"/Org.yml"
	QuotasYml := "~/C9Cli/mgmt/"+clustername+"/Quota.yml"
	ProtectedResourcesYml := "~/C9Cli/mgmt/"+clustername+"/ProtectedResources.yml"

	//println(data)

	_, err = os.Stat(ConfigPath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(ConfigPath, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

		var data = `---
ClusterDetails:
  EndPoint: {{ .EndPoint }}
  User: {{ .User }}
  Pwd: {{ .Pwd }}
  Org: {{ .Org }}
  Space: {{ .Space }}
  EnableASG: {{ .EnableASG }}`

		// Create the file:
		err = ioutil.WriteFile(ConfigPath+"/config.tmpl", []byte(data), 0644)
		check(err)

		values := ClusterDetails{EndPoint: endpoint, User: user, Pwd: pwd, Org: org, Space: space, EnableASG: asg}

		var templates *template.Template
		var allFiles []string

		if err != nil {
			fmt.Println(err)
		}

		filename := "config.tmpl"
		fullPath := ConfigPath + "/config.tmpl"
		if strings.HasSuffix(filename, ".tmpl") {
			allFiles = append(allFiles, fullPath)
		}

		fmt.Println(allFiles)
		templates, err = template.ParseFiles(allFiles...)
		if err != nil {
			fmt.Println(err)
		}

		s1 := templates.Lookup("config.tmpl")
		f, err := os.Create(ConfigPath + "/config.yml")
		if err != nil {
			panic(err)
		}

		fmt.Println("Creating .C9Cli folder and config files")

		err = s1.Execute(f, values)
		defer f.Close() // don't forget to close the file when finished.
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("~/.C9Cli/mgmt/<cluster> exists, please manually edit file to make changes or provide new cluster name")
	}

	_, err = os.Stat(mgmtpath)
	if os.IsNotExist(err) {

		fmt.Println("Creating C9Cli/mgmt/<cluster> folder")
		errDir := os.MkdirAll(mgmtpath, 0755)

		var OrgTmp = `---
Org:
  - Name: "org-test"
    Quota: "default"
    OrgUsers:
      LDAP:
        - Name: User1
          Role: OrgManager
        - Name: User2
          Role: OrgManager
        - Name: User3
          Role: OrgAudidor
      SSO:
        - Name: User1
          Role: OrgManager
        - Name: User2
          Role: OrgManager
        - Name: User3
          Role: OrgAudidor
      UAA:
        - Name: User1
          Role: OrgManager
        - Name: User2
          Role: OrgManager
        - Name: User3
          Role: OrgAudidor
    Spaces:
      - Name: Space1
        SpaceUsers:
          LDAP:
            - Name: User1
              Role: SpaceManager
            - Name: User2
              Role: SpaceManager
            - Name: User3
              Role: SpaceAudidor
          SSO:
            - Name: User1
              Role: SpaceManager
            - Name: User2
              Role: SpaceManager
            - Name: User3
              Role: SpaceAudidor
          UAA:
            - Name: User1
              Role: SpaceManager
            - Name: User2
              Role: SpaceManager
            - Name: User3
              Role: SpaceAudidor
      - Name: Space2
        SpaceUsers:
          LDAP:
            - Name: User1
              Role: SpaceManager
            - Name: User2
              Role: SpaceManager
            - Name: User3
              Role: SpaceAudidor
          SSO:
            - Name: User1
              Role: SpaceManager
            - Name: User2
              Role: SpaceManager
            - Name: User3
              Role: SpaceAudidor
          UAA:
            - Name: User1
              Role: SpaceManager
            - Name: User2
              Role: SpaceManager
            - Name: User3
              Role: SpaceAudidor`

		var QuotasTmp = `---
quota:
  - Name: default
    memory_limit: 1024M
  - Name: small_quota
    memory_limit: 2048M
  - Name: medium_quota
    memory_limit: 2048M
  - Name: large_quota
    memory_limit: 2048M`

		var ProtectedListTmp = `---
  Org:
    - system
    - healthwatch
    - dynatrace
  quota:
    - default`

		fmt.Println("Creating C9Cli/mgmt/<cluster>/ sample yaml files")
		err = ioutil.WriteFile(OrgsYml, []byte(OrgTmp), 0644)
		check(err)
		err = ioutil.WriteFile(QuotasYml, []byte(QuotasTmp), 0644)
		check(err)
		err = ioutil.WriteFile(ProtectedResourcesYml, []byte(ProtectedListTmp), 0644)
		check(err)

		if errDir != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("C9Cli/mgmt/<cluster> exists, please manually edit file to make changes or provide new cluster name")
	}

	_, err = os.Stat(ASGPath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(ASGPath, 0755)
		if errDir != nil {
			log.Fatal(err)
			fmt.Println("C9Cli/mgmt/<cluster>/ASGs exist, please manually edit file to make changes or provide new cluster name")
		} else {
			fmt.Println("Creating C9Cli/mgmt/<cluster>/ASGs")
		}
	}

	return
}
func check(e error) {
	if e != nil {
		fmt.Println("C9Cli/mgmt/<cluster>/ yamls exists, please manually edit file to make changes or provide new cluster name")
		panic(e)
	}
}
