package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"

	"github.com/golang-devops/go-psexec/shared"
)

var (
	serverFlag   = flag.String("server", "http://localhost:62677", "The endpoint server address")
	executorFlag = flag.String("executor", "winshell", "The executor to use")
)

func handleRecovery() {
	if r := recover(); r != nil {
		log.Printf("ERROR: %s\n", getErrorStringFromRecovery(r))
	}
}

func main() {
	defer handleRecovery()

	flag.Parse()

	exeAndArgs := flag.Args()
	if len(exeAndArgs) == 0 {
		panic("Need at least one additional argument")
	}

	session, err := createNewSession()
	checkError(err)

	fmt.Printf("Using session id: %d\n", session.SessionId)

	var exe string
	var args []string = []string{}

	exe = exeAndArgs[0]
	if len(exeAndArgs) > 1 {
		args = exeAndArgs[1:]
	}

	encryptedJson, err := session.EncryptAsJson(&shared.ExecDto{*executorFlag, exe, args})
	checkError(err)

	req := session.NewRequest()
	req.Json = shared.EncryptedJsonContainer{encryptedJson}

	url := combineServerUrl("/auth/exec")

	resp, err := req.Post(url)
	checkError(err)
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}