package main

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

const (
	Protocol Option = iota
	Host
	Port
	Username
	Password
	Path
	Query
	Fragment
	Black Colour = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func main() {
	// TODO: read from STDIN (pipe operator)

	input := Input{}.From(os.Args[1:])
	uri, err := url.Parse(input.Url)

	if err != nil {
		fmt.Println("\n", Chalk(Red, " unable to parse url:"), input.Url)
		os.Exit(1)
	}

	if input.Argc == 2 {
		for option := range input.Options {
			fmt.Println(Extract(uri, option))
		}

		return
	}

	PrintSegments(uri, input.Options)

	if false == Selected(input.Options, Query) {
		os.Exit(0)
	}

	PrintQuery(uri)
}

// An option identifier
type Option = int

// A set of valid URL parsing options
type Options = map[Option]struct{}

// A 256-colour terminal value
type Colour = int

// Parsed input container
type Input struct {
	// Options set supplied by user
	Options

	// The URL being parsed
	Url string

	// Count of arguments supplied to program
	Argc int
}

// Generate a new input container from an array of args
func (_ Input) From(args []string) *Input {
	var subject string
	argc := len(args)
	options := make(Options, 8)
	flags := make(map[string]Option, 10)
	flags["--protocol"] = Protocol
	flags["--scheme"] = Protocol
	flags["--host"] = Host
	flags["--port"] = Port
	flags["--username"] = Username
	flags["--user"] = Username
	flags["--password"] = Password
	flags["--path"] = Path
	flags["--query"] = Query
	flags["--fragment"] = Fragment

	for index, flag := range args {
		if 1+index == argc {
			subject = flag

			continue
		}

		option, exists := flags[flag]

		if exists == false {
			fmt.Println("\n", Chalk(Red, " unknown option given:"), flag)
			os.Exit(1)
		}

		options[option] = struct{}{}
	}

	if len(options) == 0 {
		options = All()
	}

	this := new(Input)
	this.Argc = argc
	this.Options = options
	this.Url = subject

	return this
}

// Check if a given option is exists in the provided option set
func Selected(options Options, option Option) bool {
	_, selected := options[option]

	return selected
}

// Wrap the given string in escape sequence to colour the output
func Chalk(colour Colour, input string) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[0m", colour, input)
}

// Extract the given segment (option) from the URL instance
func Extract(uri *url.URL, segment Option) string {
	switch segment {
	case Protocol:
		return uri.Scheme
	case Host:
		return uri.Hostname()
	case Port:
		return uri.Port()
	case Username:
		return uri.User.Username()
	case Password:
		password, exists := uri.User.Password()

		if false == exists {
			return ""
		}

		return password
	case Path:
		return uri.Path
	case Query:
		return uri.RawQuery
	case Fragment:
		return uri.Fragment
	}

	return ""
}

// Create an option set containing all flags
func All() Options {
	options := make(Options, 8)

	for option := 0; option < 8; option++ {
		options[option] = struct{}{}
	}

	return options
}

// Display the URL instance segments in a pretty printed format
func PrintSegments(uri *url.URL, options Options) {
	var padding int
	labels := make([]string, 8)
	labels[Protocol] = Chalk(Green, "Protocol")
	labels[Host] = Chalk(Green, "Host")
	labels[Port] = Chalk(Green, "Port")
	labels[Username] = Chalk(Green, "Username")
	labels[Password] = Chalk(Green, "Password")
	labels[Path] = Chalk(Green, "Path")
	labels[Query] = Chalk(Green, "Query")
	labels[Fragment] = Chalk(Green, "Fragment")

	for _, label := range labels {
		length := len(label)

		if length <= padding {
			continue
		}

		padding = length
	}

	fmt.Println()
	fmt.Print(Chalk(Yellow, " URL Segments\n"))
	fmt.Print(Chalk(Yellow, " ------------\n\n"))

	for option, label := range labels {
		if false == Selected(options, option) {
			continue
		}

		value := Extract(uri, option)

		if 0 == len(value) {
			continue
		}

		fmt.Printf(fmt.Sprintf(" %%%ds: %%s\n", padding), label, value)
	}
}

// Display the URL instance query parameters in a pretty printed format
func PrintQuery(uri *url.URL) {
	query, err := url.ParseQuery(uri.RawQuery)

	if err != nil {
		panic(err)
	}

	parameters := len(query)

	if 0 == parameters {
		return
	}

	labels := make(map[string]string, parameters)
	fields := make([]string, parameters)
	index := 0
	padding := 0

	for field := range query {
		fields[index] = field
		labels[field] = Chalk(Green, field)
		length := len(labels[field])

		if length > padding {
			padding = length
		}

		index++
	}

	sort.Strings(fields)
	fmt.Println()
	fmt.Print(Chalk(Yellow, " Query Parameters\n"))
	fmt.Print(Chalk(Yellow, " ----------------\n\n"))

	for _, field := range fields {
		fmt.Println(fmt.Sprintf(
			fmt.Sprintf(" %%%ds: %%s", padding),
			labels[field],
			strings.Join(query[field], ", "),
		))
	}
}
