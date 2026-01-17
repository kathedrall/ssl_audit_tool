# SSL-Audit & Recon Tool

A high-performance, concurrent CLI utility built in Go for SSL/TLS validation and infrastructure reconnaissance. This tool is designed to identify misconfigured certificates and exposed server signatures across large networks, helping to prevent information disclosure and assist in infrastructure hardening.

---

### Installation and Compilation

Since the program is written in Go, it can be compiled into a standalone binary for any operating system.

**1. Standard Compilation:**
To compile for your current operating system, use:
```bash
go build -o ssl-audit main.go
```

**2. Static Compilation (Recommended for Linux/Termux):**
To avoid segmentation faults or issues with dynamic libraries in environments like Termux or Proot-Debian, compile with CGO disabled:
```bash
CGO_ENABLED=0 go build -o ssl-audit main.go
```

**3. Cross-compilation:**
* Linux: `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ssl-audit-linux main.go`
* Windows: `GOOS=windows GOARCH=amd64 go build -o ssl-audit.exe main.go`

---

### Usage Modes

The tool identifies your intent based on how you run it in your terminal.

#### A) Single Target Audit
If you need to check a specific IP or domain quickly, pass it directly as an argument.
**Command:**
```bash
./ssl-audit 142.250.190.40
```

#### B) Bulk Audit (Using ips.txt)
To audit a list of targets simultaneously, follow these steps:

1. **Create the file:** In the same folder as the program, create a file named `ips.txt`.
2. **Add your targets:** List one IP or domain per line.
3. **Run the program:** Execute the tool without any arguments. It will automatically detect and load your list.

**Example `ips.txt` content:**
```text
bbc.com
200.185.59.31
mysecure-api.com
142.250.190.46
# 140.82.121.3 (This line is a comment and will be skipped)
104.16.124.96
```

**Command:**
```bash
./ssl-audit
```

---

### Understanding the Outputs (Returns)

The tool provides distinct ways to visualize the gathered intelligence after the scan completes.

#### 1. Real-time Terminal Feedback
While the audit is running, you will see the status of each target being processed by the concurrent Goroutines.

#### 2. Terminal Summary Table
If you choose option **1** after the scan, the tool generates a formatted table showing the IP, the Server type (via Banner Grabbing), and the SSL Status.

**Example Return:**
```text
IP              | Server                    | Status SSL
--------------------------------------------------------------------------------
bbc.com         | BBC-GTM                   | Congratulations! Your certificate is installed correctly.
104.16.124.96   | cloudflare                | Congratulations! Your certificate is installed correctly.
200.185.59.31   | Hidden                    | No SSL certificates were found. Check firewall/port 443.
```

#### 3. CSV Report
If you choose option **2**, a file named `report.csv` is generated. This report includes a "source SSL Checker" column with direct links to the SSLShopper web checker for every target.

---

### Key Technical Features

1. **Concurrent Scanning:** Uses Go routines to audit all targets simultaneously. This ensures that the total execution time is minimal, regardless of the number of IPs. 


2. **Banner Grabbing:** Attempts to capture the `Server` header to detect exposed infrastructure information (e.g., Apache/Nginx versions). This is critical for reducing the attack surface. 


3. **SSLShopper Style Diagnostics:** Emulates professional diagnostic messages for common network issues, such as connection timeouts, refused connections, or missing certificates.

---

### License
This project is licensed under the MIT License.

