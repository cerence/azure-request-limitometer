package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hetalsonavane/azure-request-limitometer/internal/config"
	"github.com/hetalsonavane/azure-request-limitometer/pkg/common"
	"github.com/hetalsonavane/azure-request-limitometer/pkg/outputs"

	flag "github.com/spf13/pflag"
)

var azureClient = common.Client

const (
	cliName        = "limitometer"
	cliDescription = "Collects the number of remaining requests in Azure Resource Manager"
	cliVersion     = "2.0.0"
)

var (
	nodename = flag.String("node", "", "Valid node in the resource group to create compute queries. Environment Variable: NODE_NAME")
	target   = flag.String("output", "pushgateway", "Target output for the limitometer, supported values are: [influxdb|pushgateway]")
)

func printUsage() {
	if flag.Args()[0] == "help" {
		fmt.Printf("%s\n\n", cliName)
		fmt.Println(cliDescription)
		flag.PrintDefaults()
		os.Exit(2)
	}
}

func printHelp() {
	if flag.Args()[0] == "version" {
		fmt.Printf("%s version %s\n", cliName, cliVersion)
		os.Exit(0)
	}
}

func main() {
	flag.Parse()

	if len(flag.Args()) > 0 {
		printHelp()
		printUsage()
	}

	if err := config.ParseEnvironment(); err != nil {
		log.Fatalf("failed to parse environment: %s\n", err)
	}

	env, exists := os.LookupEnv("NODE_NAME")
	if exists {
		*nodename = env
	}

	log.Printf("Starting limitometer with %s as target VM", *nodename)
	requestsRemaining := getRequestsRemaining(*nodename)
	//cjson, _ := json.Marshal(requestsRemaining)
	//log.Printf("%s\n", cjson)

	log.Printf("Writing to database: %s", *target)
	if strings.ToLower(*target) == "influxdb" {
		outputs.WriteOutputInflux(requestsRemaining, "requestRemaining")
	} else if strings.ToLower(*target) == "pushgateway" {
		outputs.WriteOutputPushGateway(requestsRemaining)
	} else {
		log.Printf("Did not provide a output through -output flag. Exiting.")
	}

}
