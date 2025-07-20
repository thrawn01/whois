/*
 * Copyright 2014-2024 Li Kexian
 * Copyright 2025 Derrick Wippler
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Go module for domain and ip whois information query
 * https://www.likexian.com/
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/likexian/gokit/xjson"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/thrawn01/whois"
	"golang.org/x/net/proxy"
)

func main() {
	server := flag.String("h", "", "specify the whois server")
	outJSON := flag.Bool("j", false, "output format as json")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Printf("Usage:\n\t%s [-j] [-h server] domain\n", os.Args[0])
		fmt.Printf(`
domain:
  domain or IPv4 or IPv6 or ASN for query

options:
  -h string specify the whois server
  -j        output format as json
  -v        show the whois version
`)
		os.Exit(1)
	}

	dialer := proxy.FromEnvironment()
	var text string
	var err error
	if contextDialer, ok := dialer.(proxy.ContextDialer); ok {
		text, err = whois.NewClient().
			SetDialer(contextDialer).Whois(flag.Args()[0], *server)
	} else {
		text, err = whois.NewClient().Whois(flag.Args()[0], *server)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *outJSON {
		info, err := whoisparser.Parse(text)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		data, err := xjson.PrettyDumps(info)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(data)
		os.Exit(0)
	}

	fmt.Println(text)
	os.Exit(0)
}
