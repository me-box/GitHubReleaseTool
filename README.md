# GitHub Release Tool

A tool to coordinate a release across multiple GitHub repositories.

Using the `config.json` it will:

  * Based on tags on the designated `MainRepo`, calculate the next version using semantic versioning (Major, Minor, Patch)
  * Collate all PR messages since the last release across the `MainRepo` and any `CoreRepos` into a `CHANGELOG`
  * Present a list of proposed changes so the user can double check
  * If the user agrees to go ahead it will use the GitHub API to:
    * Tag the `MainRepo` and all `CoreRepos` with the new version
    * Create or update the `Version` file in the main Databox repo
    * Create a release on the `MainRepo`  with the generated `CHANGELOG`
    * If the `-docs` flag is given:
      * Build a documentation file from the `README.md` files in the  repos indicated in the `Docs` section
      * Upload the generated documentation under a release [TODO]

## Building

```
go get github.com/me-box/GitHubReleaseTool
cd [to go path]/src/github.com/me-box/GitHubReleaseTool
go build ghrelease.go
```

or to build in a container
```
make build # produces me-box/ghrelease
```

## Usage

```
  -config string
        path of the config file (default "./config.json")
  -docs
        build docs from README.md files in main and core repos
  -docsOutFile string
        Where should the docs be output (default "./Documentation.md")
  -major
        Major release
  -minor
        Minor release
  -patch
        Patch/Bugfix release (default: true)
  -release
        set this to `false` to disable releasing (for example if you just want to rebuild the docs) (default: true)
```

## Acknowledgements

Development of Databox has been supported by the following EPSRC funding:

```
EP/N028260/1, Databox: Privacy-Aware Infrastructure for Managing Personal Data

EP/N028260/2, Databox: Privacy-Aware Infrastructure for Managing Personal Data

EP/N014243/1, Future Everyday Interaction with the Autonomous Internet of Things

EP/M001636/1, Privacy-by-Design: Building Accountability into the Internet of Things (IoTDatabox)

EP/M02315X/1, From Human Data to Personal Experience
```
