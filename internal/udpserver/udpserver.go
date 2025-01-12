package udpserver

import (
	"fmt"
	"net"

	"github.com/exploitz0169/flipdns/pkg/parser"
)

// Temp just to test the parser
func Run(conn net.PacketConn) {
	buf := make([]byte, 512)

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s\n", n, addr)

		header, err := parser.ParseDNSHeader(buf[:12])
		if err != nil {
			fmt.Println("Error parsing DNS header:", err)
			continue
		}

		// Print the parsed DNS header
		fmt.Printf("ID: %d\n", header.ID)
		fmt.Printf("QR: %d\n", header.Flags.QR)
		fmt.Printf("OPCODE: %d\n", header.Flags.OPCODE)
		fmt.Printf("AA: %d\n", header.Flags.AA)
		fmt.Printf("TC: %d\n", header.Flags.TC)
		fmt.Printf("RD: %d\n", header.Flags.RD)
		fmt.Printf("RA: %d\n", header.Flags.RA)
		fmt.Printf("Z: %d\n", header.Flags.Z)
		fmt.Printf("RCODE: %d\n", header.Flags.RCODE)
		fmt.Printf("QDCOUNT: %d\n", header.QDCOUNT)
		fmt.Printf("ANCOUNT: %d\n", header.ANCOUNT)
		fmt.Printf("NSCOUNT: %d\n", header.NSCOUNT)
		fmt.Printf("ARCOUNT: %d\n", header.ARCOUNT)

		questions, err := parser.ParseDNSQuestions(buf[12:], header.QDCOUNT)
		if err != nil {
			fmt.Println("Error parsing DNS questions:", err)
			continue
		}

		question := questions[0]

		fmt.Printf("QTYPE: %d\n", question.QTYPE)
		fmt.Printf("QCLASS: %d\n", question.QCLASS)
		fmt.Printf("DOMAIN: %s\n", question.DOMAIN)

	}
}
