package main

import (
	"net"
	"testing"
)

// TestGetOutboundIP tests the getOutboundIP function, which is used to determine
// the IP address of the current machine as seen from outside of the network.
func TestGetOutboundIP(t *testing.T) {
	// Call the getOutboundIP function to obtain the outbound IP address.
	ip := getOutboundIP()

	// Check if the IP address is non-empty. If it is empty, throw an error.
	if ip == "" {
		t.Errorf("Expected non-empty outbound IP")
	}

	// Attempt to connect to Google's DNS server at 8.8.8.8:80.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		// If we cannot connect, throw a fatal error with the details.
		t.Fatalf("Error connecting to 8.8.8.8: %v", err)
	}
	defer conn.Close()

	// Obtain the local address of the connection.
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// Store the expected IP address, which should be the local address of
	// the connection we just established.
	expectedIP := localAddr.IP.String()

	// Compare the obtained IP address to the expected one. If they do not match,
	// throw an error with the details.
	if ip != expectedIP {
		t.Errorf("Expected IP %s, but got %s", expectedIP, ip)
	}
}

func TestGetLocalIp(t *testing.T) {
	ips := getLocalIp()
	if len(ips) == 0 {
		t.Fatal("Expected at least one IP address, but got none")
	}
	for _, ip := range ips {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			t.Errorf("Invalid IP Address %s", ip)
		}
	}
}
