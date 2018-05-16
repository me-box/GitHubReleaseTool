# GitHub Release Tool

A tool to coordinate a release across multiple GitHub repositories. 

## It Will

- Look at the tags on the main repo (in config.json) and calculate the next
  version using semvar (Major, Minor, Patch)
- Collate All PRs across main and core repos (in config.json) into a changelog
- Present a List of proposed changes so that the user can double check and
  cancel if necessary
- If the user agrees
   - Tag the main repo and all core components with the new version using the GitHub API
   - Create/Update the Version file in the main databox repo
   - Create the release on the main repo using the GitHub API with the generated changelog
- if the -docs tag is enabled
  - build a documents files from the repos listed in the Docs section of the config file
  - Todo: upload the docs with a release

## Building

```
go get github.com/Toshbrown/GHR
cd [to go path]/src/github.com/Toshbrown/GHR
go build ghrelease.go
```

## Usage

```
  -config string
        path of the config file (default "./config.json")
  -docs
        build docs from Readme.md files in main and core repos
  -docsOutFile string
        Where should the docs be output (default "./Documtation.md")
  -major
        Major release
  -minor
        Minor release
  -patch
        Patch/Bugfix release (default true)
  -release
        set this to false (-release=false) disable releasing (for example if you just want to rebuild the docs) (default true)
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
