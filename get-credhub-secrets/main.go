package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
)

var (
	apiEndpoint               string
	username                  string
	password                  string
	get                       string
	insecureSkipTLSValidation string
	ca                        string
)

func main() {
	flag.StringVar(&apiEndpoint, "api", "", "Set the CredHub API target where commands are sent")
	flag.StringVar(&username, "username", "", "Authentication username")
	flag.StringVar(&password, "password", "", "Authentication password")
	flag.StringVar(&get, "get", "", "Newline-separated list of secrets to fetch.\nSecrets must be of the format SECRET_NAME=/<path>/<to>/<secret>")
	flag.StringVar(&insecureSkipTLSValidation, "insecureSkipTLSValidation", "false", "Disable TLS validation (not recommended)")
	flag.StringVar(&ca, "ca", "", "Trusted CA certificate (x509)")

	flag.Parse()

	if apiEndpoint == "" {
		log.Fatal("Missing required \"api\" input")
	}
	if username == "" {
		log.Fatal("Missing required \"username\" input")
	}
	if password == "" {
		log.Fatal("Missing required \"password\" input")
	}
	if get == "" {
		log.Fatal("Missing required \"get\" input")
	}

	options := []credhub.Option{
		credhub.SkipTLSValidation(insecureSkipTLSValidation == "true"),
		credhub.Auth(auth.UaaPassword("credhub_cli", "", username, password)),
	}
	if ca != "" {
		options = append(options, credhub.CaCerts(ca))
	}
	ch, err := credhub.New(apiEndpoint,
		options...,
	)
	if err != nil {
		log.Fatalf("unable to authenticate: %q", err)
	}

	fmt.Println("Parsing Tokens")
	pattern := regexp.MustCompile(`(?m)(?P<name>\w+):\s*(?P<path>.+?)?(?:\.(?P<key>.*))?$`)
	matches := pattern.FindAllStringSubmatch(get, -1)
	if len(matches) == 0 {
		log.Fatalln("No credentials to fetch")
	}
	for _, line := range matches {
		name := line[1]
		path := line[2]

		fmt.Printf("Fetching %q from %q\n", name, path)
		cred, err := ch.GetLatestVersion(path)
		if err != nil {
			fmt.Errorf("failed while fetching credential: %q", err)
		}

		var s string
		key := line[3]
		if key == "" {
			switch cred.Value.(type) {
			case string:
				s = cred.Value.(string)
			default:
				b, err := json.Marshal(cred.Value)
				if err != nil {
					fmt.Errorf("failed while parsing credential: %q", err)
				}
				s = string(b)
			}

			fmt.Printf("::set-output name=%s::%v\n", name, s)
			continue
		}

		c, ok := cred.Value.(map[string]interface{})
		if !ok {
			log.Fatalf("could not find key %q in credential %q", key, path)
		}

		fmt.Printf("Getting %q from %q\n", key, name)
		switch c[key].(type) {
		case string:
			s = c[key].(string)
		default:
			b, err := json.Marshal(c[key])
			if err != nil {
				fmt.Errorf("failed while parsing credential: %q", err)
			}
			s = string(b)
		}

		fmt.Printf("::set-output name=%s::%v\n", name, escape(s))
	}
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\n", "%0A")
	s = strings.ReplaceAll(s, "\r", "%0D")
	return s
}
