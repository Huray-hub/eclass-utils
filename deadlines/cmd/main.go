package main

import (
	"fmt"
	"log"

	in "github.com/Huray-hub/eclass-utils/deadlines/internal"
	dl "github.com/Huray-hub/eclass-utils/deadlines/pkg"
)

func main() {
	opts, err := in.GetOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	creds, err := in.GetCreds()
	if err != nil {
		log.Fatal(err.Error())
	}

	assignments, err := dl.Deadlines(opts, creds)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(assignments)
}
