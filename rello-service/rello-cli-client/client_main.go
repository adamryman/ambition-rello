package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/pkg/errors"

	// This Service
	pb "github.com/adamryman/ambition-rello/rello-service"
	clientHandler "github.com/adamryman/ambition-rello/rello-service/generated/cli/handlers"
	grpcclient "github.com/adamryman/ambition-rello/rello-service/generated/client/grpc"
	httpclient "github.com/adamryman/ambition-rello/rello-service/generated/client/http"
)

var (
	_ = strconv.ParseInt
	_ = strings.Split
	_ = json.Compact
	_ = errors.Wrapf
	_ = pb.RegisterRelloServer
)

func main() {
	// The addcli presumes no service discovery system, and expects users to
	// provide the direct address of an service. This presumption is reflected in
	// the cli binary and the the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings.

	var (
		httpAddr = flag.String("http.addr", "", "HTTP address of addsvc")
		grpcAddr = flag.String("grpc.addr", ":5040", "gRPC (HTTP) address of addsvc")
		method   = flag.String("method", "checklistwebhook", "checklistwebhook")
	)

	var (
		flagModelCheckListWebhook  = flag.String("checklistwebhook.model", "", "")
		flagActionCheckListWebhook = flag.String("checklistwebhook.action", "", "")
	)
	flag.Parse()

	var (
		service pb.RelloServer
		err     error
	)
	if *httpAddr != "" {
		service, err = httpclient.New(*httpAddr)
	} else if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while dialing grpc connection: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		service, err = grpcclient.New(conn)
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch *method {

	case "checklistwebhook":
		var err error

		var ModelCheckListWebhook pb.Model
		if flagModelCheckListWebhook != nil && len(*flagModelCheckListWebhook) > 0 {
			err = json.Unmarshal([]byte(*flagModelCheckListWebhook), &ModelCheckListWebhook)
			if err != nil {
				panic(errors.Wrapf(err, "unmarshalling ModelCheckListWebhook from %v:", flagModelCheckListWebhook))
			}
		}

		var ActionCheckListWebhook pb.Action
		if flagActionCheckListWebhook != nil && len(*flagActionCheckListWebhook) > 0 {
			err = json.Unmarshal([]byte(*flagActionCheckListWebhook), &ActionCheckListWebhook)
			if err != nil {
				panic(errors.Wrapf(err, "unmarshalling ActionCheckListWebhook from %v:", flagActionCheckListWebhook))
			}
		}

		request, err := clientHandler.CheckListWebhook(ModelCheckListWebhook, ActionCheckListWebhook)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling clientHandler.CheckListWebhook: %v\n", err)
			os.Exit(1)
		}

		v, err := service.CheckListWebhook(context.Background(), request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling service.CheckListWebhook: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Client Requested with:")
		fmt.Println(ModelCheckListWebhook, ActionCheckListWebhook)
		fmt.Println("Server Responded with:")
		fmt.Println(v)
	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}
}
