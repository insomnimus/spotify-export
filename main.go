package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/zmb3/spotify"
)

const VERSION = "0.1.0"

func getPlaylists(c *Creds) ([]spotify.SimplePlaylist, error) {
	limit := 50
	offset := 0
	pls := make([]spotify.SimplePlaylist, 0, 50)

	for {
		resp, err := c.Client.CurrentUsersPlaylistsOpt(&spotify.Options{
			Limit:  &limit,
			Offset: &offset,
		})
		if err != nil {
			return nil, err
		}
		pls = append(pls, resp.Playlists...)
		if len(pls) < 50 {
			break
		}
		offset += 50
	}

	return pls, nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	args := parseArgs()

	if args.out != "-" {
		fi, err := os.Stat(args.out)
		if err != nil {
			log.Fatalln(err)
		}
		if !fi.IsDir() {
			log.Fatalf("%s is not a directory\n", args.out)
		}
	}

	creds, err := login(&args)
	if err != nil {
		log.Fatalln(err)
	}

	pls, err := getPlaylists(creds)
	if err != nil {
		log.Fatalln("error retreiving playlists", err)
	}

	tasks := make([]spotify.SimplePlaylist, 0, len(args.names))
	tasksfound := make(map[spotify.ID]bool)

	for _, s := range args.names {
		g, _ := glob.Compile(strings.ToLower(s), 0)
		found := false

		for _, p := range pls {
			if _, ok := tasksfound[p.ID]; ok {
				found = true
				continue
			}

			if (g != nil && g.Match(strings.ToLower(p.Name))) || strings.EqualFold(s, p.Name) {
				found = true
				tasks = append(tasks, p)
				tasksfound[p.ID] = true
			}
		}

		if !found {
			log.Fatalf("no playlist matched the pattern `%s`\n", s)
		}
	}

	// Export playlists
	for _, p := range tasks {
		full, err := creds.Client.GetPlaylist(p.ID)
		if err != nil {
			log.Fatalf("failed to get full playlist %s: %s\n", p.Name, err)
		}
		j, err := json.MarshalIndent(full, "", "\t")
		if err != nil {
			log.Fatalf("failed to serialize %s: %s\n", p.Name, err)
		}

		if args.out == "-" {
			os.Stdout.Write(j)
			os.Stdout.Write([]byte("\n\n"))
		} else {
			out := filepath.Join(args.out, fmt.Sprintf("%s (%s).json", p.Name, string(p.ID)))
			if err := os.WriteFile(out, j, 0o644); err != nil {
				log.Fatalf("error writing to %s: %s\n", out, err)
			}
			fmt.Printf("saved %s to %s\n", p.Name, out)
		}
	}
}
