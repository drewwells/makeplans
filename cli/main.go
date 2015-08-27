package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/drewwells/makeplans"
)

type account struct {
	Name  string
	Token string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must pass key file")
	}
	bs, err := ioutil.ReadFile(os.Args[1])
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
	for _, svc := range svcs {
		fmt.Println(svc)
	}

	svc := svcs[0]
	svc.Price = "20.0"
	err = client.ServiceSave(svc)
	if err != nil {
		log.Fatal(err)
	}

	// get next cross fit
	slots, err := client.SlotNextDate("393")
	if err != nil {
		// log.Fatal(err)
	}
	for _, slot := range slots {
		fmt.Println("Next running", slot)
	}

	// get next other cross fit
	slots, err = client.SlotNextDate("394")
	if err != nil {
		// log.Fatal(err)
	}

	// get next running
	slots, err = client.SlotNextDate("395")
	if err != nil {
		// log.Fatal(err)
	}
	for _, slot := range slots {
		fmt.Println("Next running", slot)
	}
}
