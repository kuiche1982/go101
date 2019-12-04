// github.com/bradleyjkemp/grpc-tools
// for http win: fiddler, mac: Charles
// protoc -I ./ -I ./extproto/ --go_out=plugins=grpc:./ test.proto
// a good place to tryout the decoding, encoding,etc. unfortunately, it doesn't support hpack2 and protobuf https://gchq.github.io/CyberChef/#recipe=From_Hexdump(/disabled)From_Hex('Auto')Gunzip()Raw_Deflate('None%20(Store)'/disabled)Zlib_Deflate('Dynamic%20Huffman%20Coding'/disabled)Bzip2_Compress(9,31/disabled)From_Base64('A-Za-z0-9%2B/%3D',true/disabled)To_Hexdump(16,false,false)Decode_text('UTF-8%20(65001)'/disabled)&input=MWYgOGIgMDggMDAgMDAgMDAgMDAgMDAgMDIgZmYgZTIgZTIgMmEgNDkgMmQgMmUKZDEgMmIgMjggY2EgMmYgYzkgMTcgNjIgYzkgNGQgY2MgY2MgOTMgOTIgNGUgY2YKY2YgNGYgY2YgNDkgZDUgMDcgOGIgMjUgOTUgYTYgZTkgYTcgZTYgMTYgOTQgNTQKNDIgOTQgNDggYzkgYTMgNGIgOTYgNjQgZTYgYTYgMTYgOTcgMjQgZTYgMTYgNDAKMTUgYzggNDAgMTUgMjQgMTYgNjQgZWEgMjcgZTYgZTUgZTUgOTcgMjQgOTYgNjQKZTYgZTcgMTUgNDMgNjQgOTUgOTQgYjggNzggM2MgNTIgNzMgNzIgZjIgODMgNTIKMGIgNGIgNTMgOGIgNGIgODQgODQgYjggNTggZmMgMTIgNzMgNTMgMjUgMTggMTUKMTggMzUgMzggODMgYzAgNmMgMjUgNGQgMmUgNWUgYTggOWEgZTIgODIgZmMgYmMKZTIgNTQgMjEgMDkgMmUgNzYgZGYgZDQgZTIgZTIgYzQgNzQgOTggM2EgMTggZDcKMjggMTIgNmEgNWMgNzAgNmEgNTEgNTkgNjYgNzIgYWEgOTAgMjcgMTcgMmIgOTgKMmYgMjQgYTQgMDcgZjIgOGEgMWUgYjIgNWQgNTIgYzIgMjggNjIgMTAgYjMgOTUKYzQgOWIgMmUgM2YgOTkgY2MgYzQgYTcgYzQgYTkgNWYgNjYgYTggOWYgMDEgOTIKYjIgNjIgZDQgOWEgYzAgYzQgOTggYzQgMDYgNzYgYjAgMzEgMjAgMDAgMDAgZmYKZmYgNTcgZTcgZmIgYjMgMjAgMDEgMDAgMDA

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func main() {
	fmt.Println("starting")
	conn, err := grpc.Dial("127.0.0.1:8087", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := NewHelloServiceClient(conn)
	ip, _ := net.ResolveIPAddr("TCP", "127.0.0.1")
	for i := 0; i < 1; i++ {
		resp, err := client.Hello(context.TODO(), &HelloRequest{
			Name: "kkkk",
			I64:  64,
			Hr: &HelloResponse{
				Message: "somemsg",
			},
		}, grpc.PeerCallOption{
			PeerAddr: &peer.Peer{
				Addr: ip,
			},
		})
		fmt.Println(resp, err)
	}
}
