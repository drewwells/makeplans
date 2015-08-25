package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/drewwells/makeplans"
)

type account struct {
	Name  string
	Token string
}

func main() {
	bs, err := ioutil.ReadFile("account.json")
	if err != nil {
		log.Fatal(err)
	}

	ac := account{}
	err = json.Unmarshal(bs, &ac)
	if err != nil {
		log.Fatal(err)
	}

	client := makeplans.New(ac.Name, ac.Token)

	svcs, err := client.Services()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("% #v\n", svcs)

	svc := svcs[0]
	svc.Price = "20.0"
	err = client.ServiceSave(svc)
	if err != nil {
		log.Fatal(err)
	}

	svcs, _ = client.Services()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("% #v\n", svcs[0])
}
