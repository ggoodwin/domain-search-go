package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/joho/godotenv"
	domainavailability "github.com/whois-api-llc/domain-availability-go"
)

type DomainAvailability interface {
	IsAvailable(domain string) (bool, error)
}

type DomainAvailabilityImpl struct {}

func (d *DomainAvailabilityImpl) IsAvailable(domain string) (bool, error) {
	// Load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading .env file")
			return false, err
	}

	// Get the value of an environment variable
	apiKey := os.Getenv("API_KEY")
	client := domainavailability.NewBasicClient(apiKey)

	domainAvailabilityResp, _, err := client.Get(context.Background(), domain)
	if err != nil {
			log.Fatal(err)
			return false, err
	}

	if domainAvailabilityResp.IsAvailable != nil {
			if *domainAvailabilityResp.IsAvailable {
				return true, nil
			} else {
				return false, nil
			}
	}
	return false, nil
}

func main() {
	myService := new(DomainAvailabilityImpl)
	rpc.Register(myService)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	defer listener.Close()

	log.Printf("Serving RPC on port %d", 1234)
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Serve error: ", err)
	}
}
