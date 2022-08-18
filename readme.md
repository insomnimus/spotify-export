# spotify-export

Export your Spotify playlists into json files!

# Installation
Simply run:

`go install github.com/insomnimus/spotify-export@latest`

# Usage
> Note: You need to register an application from the Spotify developer portal and save your id, secret and the redirect URI you entered somewhere.

```text
spotify-export [OPTIONS] [--] <PLAYLIST_NAME...>
Export spotify playlists into json files

OPTIONS:
  -i, --id <SPOTIFY_ID>: The spotify app ID [env: SPOTIFY_ID]
  -s, --secret <SPOTIFY_SECRET>: The spotify app secret [env: SPOTIFY_SECRET]
  -r, --redirect <REDIRECT_URI>: The spotify app redirect URI [env: SPOTIFY_REDIRECT_URI]
  -o, --out <DIRECTORY>: Save exported files to <directory> or print to stdout if <directory> is "-" (default: ".")
  -B, --no-browser: Do not launch the default browser for authentication
  --: Stop processing options
  -V, --version: Show version
  -h, --help: Show this message

ARGS:
  <PLAYLIST_NAME...>: One or more playlist names or UNIX-style glob patterns
```
