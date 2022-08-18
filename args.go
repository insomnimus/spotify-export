package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Args struct {
	out         string
	names       []string
	openBrowser bool

	redirect, id, secret string
}

func usage() string {
	exec := filepath.Base(os.Args[0])
	return exec + "[OPTIONS] [--] <PLAYLIST_NAME...>"
}

func help() {
	fmt.Println("USAGE: ", usage())
	fmt.Println(`Export spotify playlists into json files

OPTIONS:
  -i, --id <SPOTIFY_ID>: The spotify app ID [env: SPOTIFY_ID]
  -s, --secret <SPOTIFY_SECRET>: The spotify app secret [env: SPOTIFY_SECRET]
  -r, --redirect <REDIRECT_URI>: The spotify app redirect URI [env: SPOTIFY_REDIRECT_URI]
  -o, --out <DIRECTORY>: Save exported files to <DIRECTORY> or print to stdout if <DIRECTORY> is "-" (default: ".")
  -B, --no-browser: Do not launch the default browser for authentication
  --: Stop processing options
  -V, --version: Show version
  -h, --help: Show this message

ARGS:
  <PLAYLIST_NAME...>: One or more playlist names or UNIX-style glob patterns`)
}

func preprocess(args []string) []string {
	processed := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			processed = append(processed, args[i:]...)
			break
		}

		if strings.HasPrefix(a, "--") {
			split := strings.SplitN(a, "=", 2)
			if len(split) <= 1 {
				processed = append(processed, a)
			} else {
				processed = append(processed, split[0], split[1])
			}
			continue
		}

		if strings.HasPrefix(a, "-") {
			if a == "-" {
				processed = append(processed, "-")
				continue
			}

		chars:
			for j, c := range a[1:] {
				switch c {
				case 'o', 's', 'i', 'r':
					processed = append(processed, "-"+string(c))
					if j+2 < len(a) {
						processed = append(processed, a[j+2:])
						break chars
					}
				default:
					processed = append(processed, "-"+string(c))
				}
			}

			continue
		}

		processed = append(processed, a)
	}

	return processed
}

func parseArgs() Args {
	args := Args{openBrowser: true}
	vals := preprocess(os.Args[1:])
	getval := func(i int, prev, argName string) string {
		if i+1 >= len(vals) || vals[i+1] == "" {
			log.Fatalf("the option %s requires a value but none was provided\n", argName)
		}
		if prev != "" {
			log.Fatalf("the option %s is specified more than once\n", argName)
		}
		return vals[i+1]
	}

loop:
	for i := 0; i < len(vals); i++ {
		a := vals[i]
		switch a {
		case "-h", "--help":
			help()
			os.Exit(0)

		case "-V", "--version":
			fmt.Printf("spotify-export V%s\n", VERSION)
			os.Exit(0)

		case "-B", "--no-browser":
			if !args.openBrowser {
				log.Fatalln("the option -B --no-browser can be used only once")
			}
			args.openBrowser = false

		case "-o", "--out":
			args.out = getval(i, args.out, "-o --out <DIRECTORY>")
			i++

		case "-i", "--id":
			args.id = getval(i, args.id, "-i --id <SPOTIFY_ID>")
			i++

		case "-s", "--secret":
			args.secret = getval(i, args.secret, "-s --secret <SPOTIFY_SECRET>")
			i++

		case "-r", "--redirect":
			args.redirect = getval(i, args.redirect, "-r --redirect <REDIRECT_URI>")
			i++

		case "--":
			for _, s := range vals[i+1:] {
				if s != "" {
					args.names = append(args.names, s)
				}
			}
			break loop
		default:
			if strings.HasPrefix(a, "-") {
				log.Fatalf("unknown option %s\nusage: %s\n", a, usage())
			}
			args.names = append(args.names, a)
		}
	}

	// validate args
	if len(args.names) == 0 {
		log.Fatalln("you must provide at least one playlist name")
	}

	envs := func(target *string, env, argName string) {
		if *target == "" {
			*target = os.Getenv(env)
			if *target == "" {
				log.Fatalf("missing required argument %s\n", argName)
			}
		}
	}

	if args.out == "" {
		args.out = "."
	}

	envs(&args.id, "SPOTIFY_ID", "-i --id <SPOTIFY_ID>")
	envs(&args.secret, "SPOTIFY_SECRET", "-s --secret <SPOTIFY_SECRET>")
	envs(&args.redirect, "SPOTIFY_REDIRECT_URI", "-r --redirect <REDIRECT_URI>")

	return args
}
