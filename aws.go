package main

import (
	"fmt"
    "context"

    "github.com/davecgh/go-spew/spew"
	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func main() {
	fmt.Println("welp")
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	serviceEC2 := ec2.New(config)
	req := serviceEC2.DescribeInstancesRequest(&ec2.DescribeInstancesInput{})
	resp, err := req.Send(context.Background()) // the background context is never canceled
	if err != nil {
		fmt.Println("error")
	}
    // fmt.Printf("%T:%#v", resp, resp)
    spew.Dump(resp.Reservations[0])
    // fmt.Println(resp.Reservations)
}
