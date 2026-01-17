package main

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
  IP        string
  Message   string
  Server    string
  SourceURL string
}

const (
 URI  = "https://www.sslshopper.com/ssl-checker.html#hostname"
 PORT = ":443"
)

func getSSLShopperMessage(ip string, err error) string {
 if err == nil {
  return "Congratulations! Your certificate is installed correctly."
 }
 errStr := err.Error()
 if strings.Contains(errStr, "deadline exceeded") || strings.Contains(errStr, "timeout") || strings.Contains(errStr, "refused") {
  return "No SSL certificates were found. Check firewall/port 443."
 }
return "Issue detected: " + errStr
}

func auditIP(ip string, wg *sync.WaitGroup, results chan<- Result) {
 defer wg.Done()

 address := ip + PORT
 timeout := 10 * time.Second
 sourceURL := fmt.Sprintf(URI+"=%s", ip)

dialer := &net.Dialer{Timeout: timeout}
conn, err := tls.DialWithDialer(dialer, "tcp", address, &tls.Config{InsecureSkipVerify: true})
	
msg := getSSLShopperMessage(ip, err)
serverHeader := "Hidden"

	if err == nil {
	defer conn.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	resp, hErr := client.Get("https://" + ip)
	if hErr == nil {
	serverHeader = resp.Header.Get("Server")
		if serverHeader == "" {
		serverHeader = "Protected / Not Informed"
		}
		resp.Body.Close()
	}
	fmt.Printf("[%s] Done\n", ip)
	} else {
	fmt.Printf("[%s] SSL port failed or closed.\n", ip)
	}

	results <- Result{IP: ip, Message: msg, Server: serverHeader, SourceURL: sourceURL}
}

func main() {
 var ips []string
  if len(os.Args) > 1 {
	ips = append(ips, os.Args[1])
  } else {
	file, err := os.Open("ips.txt")
	 if err != nil {
		fmt.Println("Error: Provide an IP address as an argument or create a file. 'ips.txt'.")
		return
	 }
	defer file.Close()
	scanner := bufio.NewScanner(file)
	 for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		 if ip != "" && !strings.HasPrefix(ip, "#") {
			ips = append(ips, ip)
		 }
	 }
  }

  if len(ips) == 0 {
   fmt.Println("No IP addresses found to process.")
   return
  }

  resultsChan := make(chan Result, len(ips))
  var wg sync.WaitGroup

  fmt.Printf("\n--- Initiating Audit on %d targets ---\n\n", len(ips))

  for _, ip := range ips {
   wg.Add(1)
	go auditIP(ip, &wg, resultsChan)
  }

  go func() {
   wg.Wait()
	close(resultsChan)
  }()

  var finalResults []Result
  for res := range resultsChan {
   finalResults = append(finalResults, res)
  }

	
  fmt.Println("\nAudit Completed! Choose the report format:")
  fmt.Println("1) Display summary in terminal")
  fmt.Println("2) Export to CSV file")
  fmt.Print("\nOpção: ")

  var choice string
  fmt.Scanln(&choice)

  if choice == "1" {
   fmt.Printf("\n%-15s | %-25s | %s\n", "IP", "Server", "Status SSL")
   fmt.Println(strings.Repeat("-", 80))
	for _, r := range finalResults {
	 fmt.Printf("%-15s | %-25s | %s\n", r.IP, r.Server, r.Message)
	}
  } else if choice == "2" {
   filename := "report.csv"
   f, _ := os.Create(filename)
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write([]string{"IP", "Result", "Serverr", "font SSL Checker"})
	 for _, r := range finalResults {
	   w.Write([]string{r.IP, r.Message, r.Server, r.SourceURL})
	  }
	  fmt.Printf("\nReport saved successfully: %s\n", filename)
	} else {
	  fmt.Println("Leaving without generating a report.")
	}
}

