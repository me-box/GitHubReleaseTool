package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	config "github.com/me-box/GitHubReleaseTool/config"
	"github.com/me-box/GitHubReleaseTool/githubProvider"
	tu "github.com/me-box/GitHubReleaseTool/tutils"
)

func main() {

	//Process options
	configPath := flag.String("config", "./config.json", "path of the config file")
	configMajor := flag.Bool("major", false, "Major release")
	configMinor := flag.Bool("minor", false, "Minor release")
	configPatch := flag.Bool("patch", true, "Patch/Bugfix release")
	configDocs := flag.Bool("docs", false, "build docs from Readme.md files in main and core repos")
	configRelease := flag.Bool("release", true, "set this to false (-release=false) disable releasing (for example if you just want to rebuild the docs)")
	configDocsFile := flag.String("docsOutFile", "./Documtation.md", "Where should the docs be output")
	flag.Parse()

	//load config file
	cfg, err := config.ConfigFromFile(*configPath)
	tu.CheckExit(err)

	ghp := githubProvider.New(cfg.AccessToken)

	if *configDocs {

		fmt.Println("Building docs")
		docs, genDocsErr := ghp.GenerateDocs(cfg.Docs)
		tu.CheckExit(genDocsErr)

		err := ioutil.WriteFile(*configDocsFile, []byte(strings.Join(docs, "\n")), 0644)
		tu.CheckExit(err)

	}

	if *configRelease {
		releaseText, nextVersion := ghp.GenerateReleaseText(cfg.MainRepo, cfg.CoreRepos, cfg.OtherRepos, configMajor, configMinor, configPatch)
		fmt.Println("")
		fmt.Println("---------------------")
		fmt.Println("Proposed Release Text")
		fmt.Println("---------------------")
		fmt.Println(strings.Trim(fmt.Sprint(releaseText[:]), "[]"))
		fmt.Println("--------------------------------------------------")
		fmt.Println("Do you want to publish the release above? [no/yes]")
		fmt.Println("--------------------------------------------------")
		reader := bufio.NewReader(os.Stdin)
		ans, _ := reader.ReadString('\n')

		if ans != "yes\n" {
			tu.Log("No changes made you did not answer yes")
			return
		}

		fmt.Println("---------------------")
		fmt.Println("Releasing " + nextVersion + "........ ")
		fmt.Println("---------------------")
		ghp.DoRelease(cfg.MainRepo, cfg.CoreRepos, cfg.OtherRepos, releaseText, nextVersion)
	}
}
