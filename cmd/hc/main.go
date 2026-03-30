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
	{102, "Processing", "Server is processing the request but no response is available yet"},
	{103, "Early Hints", "Used to return response headers before the final response"},

	{200, "OK", "Request fulfilled, document follows"},
	{201, "Created", "Document created, URL follows"},
	{202, "Accepted", "Request accepted, processing continues off-line"},
	{203, "Non-Authoritative Information", "Response payload has been modified by a transforming proxy"},
	{204, "No Content", "Request fulfilled, nothing follows"},
	{205, "Reset Content", "Clear input form for further input"},
	{206, "Partial Content", "Partial content follows"},
	{207, "Multi-Status", "Response contains multiple status codes for multiple sub-requests"},
	{208, "Already Reported", "Members already enumerated in a previous reply"},
	{226, "IM Used", "Server fulfilled a GET request with instance-manipulations applied"},

	{300, "Multiple Choices", "Object has several resources -- see URI list"},
	{301, "Moved Permanently", "Object moved permanently -- see URI list"},
	{302, "Found", "Object moved temporarily -- see URI list"},
	{303, "See Other", "Object moved -- see Method and URL list"},
	{304, "Not Modified", "Document has not changed since given time"},
	{305, "Use Proxy", "You must use proxy specified in Location to access this resource"},
	{306, "(Unused)", "Reserved; previously proposed as Switch Proxy"},
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
	{409, "Conflict", "Request conflicts with the current state of the target resource"},
	{410, "Gone", "URI no longer exists and has been permanently removed"},
	{411, "Length Required", "Client must specify Content-Length"},
	{412, "Precondition Failed", "Precondition in headers is false"},
	{413, "Content Too Large", "Request content is larger than the server is willing or able to process"},
	{414, "URI Too Long", "URI is too long"},
	{415, "Unsupported Media Type", "Entity body in unsupported format"},
	{416, "Range Not Satisfiable", "Cannot satisfy request range"},
	{417, "Expectation Failed", "Expect condition could not be satisfied"},
	{418, "I'm a teapot", "The HTCPCP server is a teapot"},
	{419, "Authentication Timeout", "Previously valid authentication has expired"},
	{420, "Method Failure / Enhance Your Calm", "Spring Framework: method failure; Twitter: rate limiting"},
	{421, "Misdirected Request", "Request was directed at a server unable to produce a response"},
	{422, "Unprocessable Content", "Server understands the content type but cannot process the contained instructions"},
	{423, "Locked", "The resource is currently locked"},
	{424, "Failed Dependency", "Request failed because it depended on another request that failed"},
	{425, "Too Early", "Server is unwilling to risk processing a request that might be replayed"},
	{426, "Upgrade Required", "Client should switch to a different protocol"},
	{428, "Precondition Required", "Server requires the request to be conditional"},
	{429, "Too Many Requests", "User has sent too many requests in a given amount of time"},
	{431, "Request Header Fields Too Large", "Server refuses to process because header fields are too large"},
	{440, "Login Timeout", "The client's session has expired and must log in again"},
	{444, "No Response", "Server returned no information and closed the connection"},
	{449, "Retry With", "Request should be retried after performing the appropriate action"},
	{450, "Blocked by Windows Parental Controls", "Parental controls are blocking access to the requested resource"},
	{451, "Unavailable For Legal Reasons", "Denied access due to a legal demand"},
	{494, "Request Header Too Large", "Client sent too large of a request or header line"},
	{495, "Cert Error", "Client certificate error prevented the connection"},
	{496, "No Cert", "Client did not provide a required certificate"},
	{497, "HTTP to HTTPS", "HTTP request was sent to an HTTPS port"},
	{498, "Token expired/invalid", "Expired or invalid token in the request"},
	{499, "Client Closed Request", "Client closed the connection before the server could respond"},

	{500, "Internal Server Error", "Server got itself in trouble"},
	{501, "Not Implemented", "Server does not support this operation"},
	{502, "Bad Gateway", "Invalid responses from another server/proxy"},
	{503, "Service Unavailable", "The server cannot process the request due to a high load"},
	{504, "Gateway Timeout", "The gateway server did not receive a timely response"},
	{505, "HTTP Version Not Supported", "Cannot fulfill request"},
	{506, "Variant Also Negotiates", "Server has an internal configuration error during content negotiation"},
	{507, "Insufficient Storage", "Not enough storage to complete the request"},
	{508, "Loop Detected", "Server detected an infinite loop while processing the request"},
	{509, "Bandwidth Limit Exceeded", "Server has exceeded the bandwidth limit set by the administrator"},
	{510, "Not Extended", "Further extensions to the request are required for the server to fulfill it"},
	{511, "Network Authentication Required", "Client needs to authenticate to gain network access"},
	{598, "Network read timeout error", "Informal; used by some HTTP proxies"},
	{599, "Network connect timeout error", "Informal; used by some HTTP proxies"},
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
