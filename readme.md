# cako.io cli
A command line interface for archiving and serving pages from cako.io.

## Usage:
```
  -outDir string
        Output directory to save files (default "./saved/")
  -page string
        Crawl specified page name only
  -password string
        Password for private site access
  -serve
        Serve locally saved files
  -skipAssets
        Only crawl html files
  -skipExisting
        Skip crawling pages already in output directory
```

## Example:

```
cako -skipExisting
```

After archiving, run with `-serve` to start a local server for the given output directory.

```
cako -serve
```

### Enabling binary execution on macOS:

When you run a compiled binary from the Releases page, you may see an error that "Apple cannot check for malicious software".  Run this command to allow it through Gatekeeper.

```
xattr -d com.apple.quarantine cako.io-cli-darwin-arm64
```
