package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	proto "git.amocrm.ru/ilnasertdinov/http-server-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := flag.String("addr", "localhost:8091", "grpc address")
	id := flag.Uint64("id", 0, "account id")
	flag.Parse()

	if *id == 0 {
		log.Fatal("id is required")
	}

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := proto.NewAccountServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.DisableAccount(ctx, &proto.DisableAccountRequest{AccountId: *id})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ok=%v info=%s\n", res.Ok, res.Info)
}
