# cako.io cli
A command line interface for archiving and serving pages from cako.io.

## Usage:
```
  -checkLinks
    	Only check links for reachability (do not archive files).
  -outDir string
    	Output directory to save files (default "./saved/")
  -page string
    	Crawl specified page name only (default "all/")
  -password string
    	Password for private site access
  -serve
    	Serve locally saved files
  -skipAssets
    	Only crawl html files
  -skipExisting
    	Skip crawling pages already in output directory
  -youtubeDataApiKey string
    	YouTube Data API key for verifying YouTube video availability.
```

## Example:

```
cako -outDir cako.io-archive
```

After archiving, run with `-serve` to start a local server for the given output directory.

```
cako -serve -outDir cako.io-archive
```

`-checkLinks` may also be used to verify reachability of external links.  You must create a YouTube Data API key and set it locally via environment variable `YOUTUBE_DATA_API_KEY` or command line option `-youtubeDataApiKey`.

```
cako -checkLinks -youtubeDataApiKey {API_KEY}
```

## Installation:
Download the latest [Release](https://github.com/david-cako/cako.io-cli/releases) and install.
```
sudo install -o root -g wheel -m 755 cako-darwin-arm64 /usr/local/bin/cako
```

### Enabling binary execution on macOS:

When you run a compiled binary from the Releases page, you may see an error that "Apple cannot check for malicious software".  Run this command to allow it through Gatekeeper.

```
sudo xattr -d com.apple.quarantine /usr/local/bin/cako
```
