package main

import (
	"flag"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/proto/gnmi"

	"ribeiro/gnmi/helper"

	"github.com/kylelemons/godebug/pretty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	targetAddress = flag.String("target_address", "localhost:32123", "The target address:port.")
	targetName    = flag.String("target_name", "", "Will use this hostname to verify server certificate during TLS handshake.")
	timeOut       = flag.Duration("time_out", 10*time.Second, "Timeout for the Get request, 10 seconds by default.")
	query         = flag.String("query", "", "XPath query or queries. Example: system/openflow/controllers/controller[main]/connections/connection[0]/state/address")
)

func main() {
	flag.Parse()

	if *query == "" {
		log.Exit("-query must be set")
	}
	queries := helper.ParseQuery(*query)
	getRequest := helper.ToGetRequest(queries)

	creds := helper.ClientCertificates(*targetName)

	conn, err := grpc.Dial(*targetAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Exitf("did not connect: %v", err)
	}
	defer conn.Close()
	c := gnmi.NewGNMIClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
	defer cancel()

	getResponse, err := c.Get(ctx, &getRequest)
	if err != nil {
		log.Exitf("could not get: %v", err)
	}

	pretty.Print(getResponse)
}