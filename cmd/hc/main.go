package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	lipgloss "charm.land/lipgloss/v2"
)

var version = "dev"

type statusCode struct {
	Code    int
	Message string
	Explain string
}

var statusCodes = []statusCode{
	{100, "Continue", "Request received, please continue"},
	{101, "Switching Protocols", "Switching to new protocol; obey Upgrade header"},
	{102, "Processing", "WebDAV; RFC 2518"},
	{103, "Early Hints", "RFC 8297"},

	{200, "OK", "Request fulfilled, document follows"},
	{201, "Created", "Document created, URL follows"},
	{202, "Accepted", "Request accepted, processing continues off-line"},
	{203, "Non-Authoritative Information", "Request fulfilled from cache"},
	{204, "No Content", "Request fulfilled, nothing follows"},
	{205, "Reset Content", "Clear input form for further input"},
	{206, "Partial Content", "Partial content follows"},
	{207, "Multi-Status", "WebDAV; RFC 4918"},
	{208, "Already Reported", "WebDAV; RFC 5842"},
	{226, "IM Used", "RFC 3229"},

	{300, "Multiple Choices", "Object has several resources -- see URI list"},
	{301, "Moved Permanently", "Object moved permanently -- see URI list"},
	{302, "Found", "Object moved temporarily -- see URI list"},
	{303, "See Other", "Object moved -- see Method and URL list"},
	{304, "Not Modified", "Document has not changed since given time"},
	{305, "Use Proxy", "You must use proxy specified in Location to access this resource"},
	{306, "Switch Proxy", "Subsequent requests should use the specified proxy"},
	{307, "Temporary Redirect", "Object moved temporarily -- see URI list"},
	{308, "Permanent Redirect", "Object moved permanently"},

	{400, "Bad Request", "Bad request syntax or unsupported method"},
	{401, "Unauthorized", "No permission -- see authorization schemes"},
	{402, "Payment Required", "No payment -- see charging schemes"},
	{403, "Forbidden", "Request forbidden -- authorization will not help"},
	{404, "Not Found", "Nothing matches the given URI"},
	{405, "Method Not Allowed", "Specified method is invalid for this resource"},
	{406, "Not Acceptable", "URI not available in preferred format"},
	{407, "Proxy Authentication Required", "You must authenticate with this proxy before proceeding"},
	{408, "Request Timeout", "Request timed out; try again later"},
	{409, "Conflict", "Request conflict"},
	{410, "Gone", "URI no longer exists and has been permanently removed"},
	{411, "Length Required", "Client must specify Content-Length"},
	{412, "Precondition Failed", "Precondition in headers is false"},
	{413, "Payload Too Large", "Payload is too large"},
	{414, "Request-URI Too Long", "URI is too long"},
	{415, "Unsupported Media Type", "Entity body in unsupported format"},
	{416, "Requested Range Not Satisfiable", "Cannot satisfy request range"},
	{417, "Expectation Failed", "Expect condition could not be satisfied"},
	{418, "I'm a teapot", "The HTCPCP server is a teapot"},
	{419, "Authentication Timeout", "Previously valid authentication has expired"},
	{420, "Method Failure / Enhance Your Calm", "Spring Framework / Twitter"},
	{422, "Unprocessable Entity", "WebDAV; RFC 4918"},
	{423, "Locked", "WebDAV; RFC 4918"},
	{424, "Failed Dependency / Method Failure", "WebDAV; RFC 4918"},
	{425, "Unordered Collection", "Internet draft"},
	{426, "Upgrade Required", "Client should switch to a different protocol"},
	{428, "Precondition Required", "RFC 6585"},
	{429, "Too Many Requests", "RFC 6585"},
	{431, "Request Header Fields Too Large", "RFC 6585"},
	{440, "Login Timeout", "Microsoft"},
	{444, "No Response", "Nginx"},
	{449, "Retry With", "Microsoft"},
	{450, "Blocked by Windows Parental Controls", "Microsoft"},
	{451, "Unavailable For Legal Reasons", "RFC 7725"},
	{494, "Request Header Too Large", "Nginx"},
	{495, "Cert Error", "Nginx"},
	{496, "No Cert", "Nginx"},
	{497, "HTTP to HTTPS", "Nginx"},
	{498, "Token expired/invalid", "Esri"},
	{499, "Client Closed Request", "Nginx"},

	{500, "Internal Server Error", "Server got itself in trouble"},
	{501, "Not Implemented", "Server does not support this operation"},
	{502, "Bad Gateway", "Invalid responses from another server/proxy"},
	{503, "Service Unavailable", "The server cannot process the request due to a high load"},
	{504, "Gateway Timeout", "The gateway server did not receive a timely response"},
	{505, "HTTP Version Not Supported", "Cannot fulfill request"},
	{506, "Variant Also Negotiates", "RFC 2295"},
	{507, "Insufficient Storage", "WebDAV; RFC 4918"},
	{508, "Loop Detected", "WebDAV; RFC 5842"},
	{509, "Bandwidth Limit Exceeded", "Apache bw/limited extension"},
	{510, "Not Extended", "RFC 2774"},
	{511, "Network Authentication Required", "RFC 6585"},
	{598, "Network read timeout error", "Unknown"},
	{599, "Network connect timeout error", "Unknown"},
}

// Styles

var lightDark lipgloss.LightDarkFunc

type colorPair struct {
	light, dark string
}

var (
	codeColorPairs = [5]colorPair{
		{"27", "75"},   // 1xx blue
		{"28", "78"},   // 2xx green
		{"136", "220"}, // 3xx yellow
		{"166", "208"}, // 4xx orange
		{"160", "196"}, // 5xx red
	}
	explainColors = colorPair{"242", "248"}
)

func initStyles() {
	lightDark = lipgloss.LightDark(lipgloss.HasDarkBackground(os.Stdin, os.Stdout))
}

func codeStyle(code int) lipgloss.Style {
	i := code/100 - 1
	if i < 0 || i >= len(codeColorPairs) {
		i = len(codeColorPairs) - 1
	}
	c := codeColorPairs[i]
	return lipgloss.NewStyle().Bold(true).Foreground(lightDark(lipgloss.Color(c.light), lipgloss.Color(c.dark)))
}

func msgStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true)
}

func explainStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color(explainColors.light), lipgloss.Color(explainColors.dark)))
}

// Output

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func printCode(sc statusCode) {
	code := codeStyle(sc.Code).Render(strconv.Itoa(sc.Code))
	msg := msgStyle().Render(sc.Message)
	explain := explainStyle().Render(sc.Explain)
	fmt.Fprintf(lipgloss.Writer, "%s  %s\n     %s\n\n", code, msg, explain)
}

func printCodes(codes []statusCode) {
	sort.Slice(codes, func(i, j int) bool { return codes[i].Code < codes[j].Code })
	for _, sc := range codes {
		printCode(sc)
	}
}

// Filtering

func matchByPattern(codes []statusCode, pattern string) ([]statusCode, error) {
	re, err := regexp.Compile(strings.ReplaceAll(pattern, "x", `\d`) + "$")
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %s", pattern)
	}
	var matched []statusCode
	for _, sc := range codes {
		if re.MatchString(strconv.Itoa(sc.Code)) {
			matched = append(matched, sc)
		}
	}
	return matched, nil
}

func matchByText(codes []statusCode, text string) ([]statusCode, error) {
	re, err := regexp.Compile("(?i)" + text)
	if err != nil {
		return nil, fmt.Errorf("invalid search pattern: %s", text)
	}
	var matched []statusCode
	for _, sc := range codes {
		if re.MatchString(sc.Message + sc.Explain) {
			matched = append(matched, sc)
		}
	}
	return matched, nil
}

// CLI

func usage() {
	fmt.Print(`hc - Look up HTTP status codes

Usage:
  hc <code>              Look up a code (supports regex, x = any digit)
  hc -s <text>           Search by description (case-insensitive)
  hc                     List all codes

Examples:
  hc 418                 I'm a teapot
  hc 1xx                 All 1xx informational codes
  hc '30[12]'            301 and 302 only
  hc -s timeout          Codes mentioning "timeout"

Options:
  -s, --search <text>    Search by name or description
  -h, --help             Show this help
  -v, --version          Show version
`)
}

func main() {
	args := os.Args[1:]

	var (
		search string
		code   string
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h", "--help":
			usage()
			return
		case "-v", "--version":
			fmt.Printf("hc version %s\n", version)
			return
		case "-s", "--search":
			if i+1 >= len(args) {
				fatal("--search requires an argument")
			}
			i++
			search = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Fprintf(os.Stderr, "Unknown option: %s\n", args[i])
				usage()
				os.Exit(1)
			}
			code = args[i]
		}
	}

	initStyles()

	var (
		results []statusCode
		err     error
	)

	switch {
	case search != "":
		results, err = matchByText(statusCodes, search)
	case code != "":
		results, err = matchByPattern(statusCodes, code)
	default:
		results = statusCodes
	}
	if err != nil {
		fatal("%s", err)
	}
	if len(results) == 0 {
		fatal("No code found matching: %s", search+code)
	}
	printCodes(results)
}
