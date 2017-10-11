package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	config "./config"
	"./githubProvider"
	tu "./tutils"
)

func main() {

	//Process options
	configPath := flag.String("config", "./config.json", "path of the config file")
	configMajor := flag.Bool("major", false, "Major release")
	configMinor := flag.Bool("minor", false, "Minor release")
	configPatch := flag.Bool("patch", true, "Patch/Bugfix release")
	flag.Parse()

	//load config file
	cfg, err := config.ConfigFromFile(*configPath)
	tu.CheckExit(err)

	ghp := githubProvider.New(cfg.AccessToken)
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
