package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type credentials struct {
	Credentials tokens `json:"Credentials"`
}
type tokens struct {
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Expiration      string `json:"Expiration"`
}

func main() {

	var serialNumber string
	fmt.Printf("Enter token code...\n")
	fmt.Scan(&serialNumber)

	var mfaCode string
	fmt.Printf("Enter token code...\n")
	fmt.Scan(&mfaCode)
	var execCommand string = "aws sts get-session-token --serial-number " + serialNumber + " --token-code " + mfaCode
	fmt.Printf(execCommand)
	cmd := exec.Command("bash", "-c", execCommand)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("error =>", err)
	}
	fmt.Println(string(stdout))

	var convertedCredentials credentials
	err = json.Unmarshal(stdout, &convertedCredentials)
	if err != nil {
		log.Fatalf("Credentials Error: %s", err.Error())
	}

	var accessToken string = "echo -e 'export AWS_ACCESS_KEY_ID=" + convertedCredentials.Credentials.AccessKeyId + "' >> ~/.bash_profile"
	var secretAccessToken string = "echo 'export AWS_SECRET_ACCESS_KEY=" + convertedCredentials.Credentials.SecretAccessKey + "' >> ~/.bash_profile"
	var sessionToken string = "echo 'export AWS_SESSION_TOKEN=" + convertedCredentials.Credentials.SessionToken + "' >> ~/.bash_profile"
	var checkSession string = "aws s3 ls"

	exec.Command("bash", "-c", "sed -i '' '/export AWS_ACCESS_KEY_ID=/d' ~/.bash_profile").Run()
	fmt.Println("deleted previous credentials", accessToken)

	exec.Command("bash", "-c", "sed -i '' '/export AWS_SECRET_ACCESS_KEY=/d' ~/.bash_profile").Run()
	fmt.Println("deleted previous credentials", secretAccessToken)

	exec.Command("bash", "-c", "sed -i '' '/export AWS_SESSION_TOKEN=/d' ~/.bash_profile").Run()
	fmt.Println("deleted previous credentials", sessionToken)

	exec.Command("bash", "-c", accessToken).Run()
	exec.Command("bash", "-c", secretAccessToken).Run()
	exec.Command("bash", "-c", sessionToken).Run()
	fmt.Println("append new credentials")

	fmt.Println("Updating ~/.bash_profile")
	exec.Command("bash", "-c", "source ~/.bash_profile").Run()
	fmt.Println("Updated ~/.bash_profile")

	stdout, err = exec.Command("bash", "-c", checkSession).Output()
	if err != nil {
		log.Fatalf("Session Error: %s", err.Error())
	}
	fmt.Println("aws s3 ls => ", string(stdout))
	os.Exit(1)
}
