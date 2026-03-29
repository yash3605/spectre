package domain

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/yash3605/spectre/internal/models"
)

func whoisServer(tld string) string {
	switch tld {
	case "com", "net":
		return "whois.verisign-grs.com"
	case "org":
		return "whois.pir.org"
	case "io":
		return "whois.nic.io"
	case "dev":
		return "whois.nic.google"
	default:
		return "whois.iana.org"
	}
}

func Lookup(domain string) models.Result {
	parts := strings.Split(domain, ".")
	tld := parts[len(parts)-1]

	addrStr := whoisServer(tld)
	conn, err := net.Dial("tcp", addrStr+":43")
	if err != nil {
		fmt.Println("Error Connecting TCP", err)
		return models.Result{
			Title:  domain,
			Data:   map[string]string{},
			Status: models.StateError,
		}
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "%s\r\n", domain)
	if err != nil {
		fmt.Println("Error getting response from connection", err)
		return models.Result{
			Title:  domain,
			Data:   map[string]string{},
			Status: models.StateError,
		}
	}

	raw, err := io.ReadAll(conn)
	if err != nil {
		fmt.Println("Error reading response from connection", err)
		return models.Result{
			Title:  domain,
			Data:   map[string]string{},
			Status: models.StateError,
		}
	}

	data := make(map[string]string)
	var order []string
	lines := strings.Split(string(raw), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			if _, ok := data[parts[0]]; !ok {
				data[parts[0]] = parts[1]
				order = append(order, parts[0])
			}
		}
	}

	return models.Result{
		Title:  domain,
		Data:   data,
		Order:  order,
		Status: models.StateSuccess,
	}
}
