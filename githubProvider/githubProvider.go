package githubProvider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tu "github.com/toshbrown/GHR/tutils"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//GithubProvider provides authenticated access to GitHub APIs
type GithubProvider struct {
	ctx    context.Context
	client *github.Client
}

//New create a new GithubProvider that with authentication
func New(accessToken string) *GithubProvider {
	ghp := GithubProvider{}
	ghp.ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ghp.ctx, ts)
	ghp.client = github.NewClient(tc)
	return &ghp
}

//GenerateReleaseText gets the last release form mainRepo then looks for PR on mainRepo, coreRepos,  otherRepos since the last version and genorates a change log
func (ghp *GithubProvider) GenerateReleaseText(mainRepo []string, coreRepos [][]string, otherRepos [][]string, major *bool, minor *bool, patch *bool) ([]string, string) {

	changes, lastRelease, GhpErr := ghp.getChangesSinceLastRelease(mainRepo[0], mainRepo[1])
	tu.CheckWarn(GhpErr)

	currentVersion := lastRelease.TagName
	nextVersion := ghp.calculateNextVersion(currentVersion, major, minor, patch)

	var releaseText []string

	fmt.Println("nextVersion :: " + nextVersion)
	fmt.Println("currentVersion :: " + *currentVersion)
	releaseText = append(releaseText, "# Changes since last version:\n")
	releaseText = append(releaseText, "## Changes to "+mainRepo[0]+"/"+mainRepo[1]+":\n")
	for _, change := range changes {
		releaseText = append(releaseText, fmt.Sprintf("  - %s see %s/%s/pull/%d \n", change.Title, mainRepo[0], mainRepo[1], change.PrNum))
	}
	releaseText = append(releaseText, "# Changes to core repositories:\n")
	for _, coreRepo := range coreRepos {
		coreChanges, changeErr := ghp.getChangesSinceRelease(lastRelease, coreRepo[0], coreRepo[1])
		tu.CheckWarn(changeErr)
		releaseText = append(releaseText, "## Changes to "+coreRepo[0]+"/"+coreRepo[1]+":\n")
		for _, coreChange := range coreChanges {
			releaseText = append(releaseText, fmt.Sprintf("  - %s see %s/%s/pull/%d \n", coreChange.Title, coreRepo[0], coreRepo[1], coreChange.PrNum))
		}
	}
	fmt.Println("--------------------------------------------------")
	fmt.Println("The Folowing repos will also be tagged with " + nextVersion)
	fmt.Println("--------------------------------------------------")
	fmt.Println(mainRepo[0] + "/" + mainRepo[1])
	for _, repo := range coreRepos {
		fmt.Println(repo[0] + "/" + repo[1])
	}
	for _, repo := range otherRepos {
		fmt.Println(repo[0] + "/" + repo[1])
	}

	return releaseText, nextVersion
}

func (ghp *GithubProvider) DoRelease(mainRepo []string, coreRepos [][]string, otherRepos [][]string, releaseText []string, nextversion string) {

	//create version file in main repo
	getOpts := github.RepositoryContentGetOptions{}
	fileCont, _, _, getErr := ghp.client.Repositories.GetContents(ghp.ctx, mainRepo[0], mainRepo[1], "Version", &getOpts)
	tu.CheckWarn(getErr)

	if fileCont != nil {
		//update
		msg := "Updating version to " + nextversion
		sha := fileCont.GetSHA()
		fileOpts := github.RepositoryContentFileOptions{
			Message: &msg,
			Content: []byte(nextversion),
			SHA:     &sha,
		}
		_, _, updateError := ghp.client.Repositories.UpdateFile(ghp.ctx, mainRepo[0], mainRepo[1], "Version", &fileOpts)
		tu.CheckWarn(updateError)

	} else {
		//create
		msg := "Creating version file for version " + nextversion
		fileOpts := github.RepositoryContentFileOptions{
			Message: &msg,
			Content: []byte(nextversion),
		}
		_, _, updateError := ghp.client.Repositories.CreateFile(ghp.ctx, mainRepo[0], mainRepo[1], "Version", &fileOpts)
		tu.CheckWarn(updateError)
	}

	//tag core and other repos
	for _, repo := range append(coreRepos, otherRepos...) {
		//msg := "Tagging release " + nextversion

		//get the last commit
		commitOpts := github.CommitsListOptions{}
		commits, _, listCommitErr := ghp.client.Repositories.ListCommits(ghp.ctx, repo[0], repo[1], &commitOpts)
		tu.CheckWarn(listCommitErr)

		//
		objType := "commit"
		objSha := commits[0].GetSHA()
		gitObj := github.GitObject{
			Type: &objType,
			SHA:  &objSha,
		}
		refStr := "refs/tags/" + nextversion
		ref := github.Reference{
			Ref:    &refStr,
			Object: &gitObj,
		}
		_, _, tagErr := ghp.client.Git.CreateRef(ghp.ctx, repo[0], repo[1], &ref)
		tu.CheckWarn(tagErr)
	}

	//create release
	tagName := nextversion
	relName := strings.Title(mainRepo[1]) + " Version: " + nextversion
	relText := strings.Join(releaseText, "")
	release := github.RepositoryRelease{
		TagName: &tagName,
		Name:    &relName,
		Body:    &relText,
	}
	ghp.client.Repositories.CreateRelease(ghp.ctx, mainRepo[0], mainRepo[1], &release)

}

type changes struct {
	Title string
	URL   string
	PrNum int
}

//getChangesSinceLastRelease gets all merged Pull requests on the repo at owner/repo since the last tagged release
func (ghp *GithubProvider) getChangesSinceLastRelease(owner string, repo string) ([]changes, *github.RepositoryRelease, error) {

	tu.Log("getChangesSinceLastRelease " + owner + "/" + repo)

	lastRelease, _, gitErr := ghp.client.Repositories.GetLatestRelease(ghp.ctx, owner, repo)
	if gitErr != nil {
		return nil, nil, gitErr
	}

	res, changesErr := ghp.getChangesSinceRelease(lastRelease, owner, repo)

	return res, lastRelease, changesErr
}

//getChangesSinceRelease gets all merged Pull requests on the repo at owner/repo since the provided release
func (ghp *GithubProvider) getChangesSinceRelease(release *github.RepositoryRelease, owner string, repo string) ([]changes, error) {

	result := []changes{}

	prOpts := github.PullRequestListOptions{
		State: "all",
		Base:  "",
	}
	prs, _, prErr := ghp.client.PullRequests.List(ghp.ctx, owner, repo, &prOpts)
	if prErr != nil {
		return nil, prErr
	}

	for _, pr := range prs {
		if pr.MergedAt != nil && pr.MergedAt.After(release.CreatedAt.Time) {
			PrNumTmp := strings.Split(pr.GetURL(), "/")
			PrNum, _ := strconv.Atoi(PrNumTmp[len(PrNumTmp)-1])
			result = append(result, changes{
				pr.GetTitle(),
				pr.GetURL(),
				PrNum,
			})
		}
	}

	return result, nil

}

func (ghp *GithubProvider) calculateNextVersion(currentVersion *string, major *bool, minor *bool, patch *bool) string {

	var nextVersion string
	sp := strings.Split(*currentVersion, ".")
	if len(sp) < 3 {
		sp = []string{"0", "0", "0"}
	}
	if *major {
		verPart, _ := strconv.Atoi(sp[0])
		sp[0] = strconv.Itoa(verPart + 1)
		sp[1] = strconv.Itoa(0)
		sp[2] = strconv.Itoa(0)
		nextVersion = strings.Join(sp, ".")

	} else if *minor {
		verPart, _ := strconv.Atoi(sp[1])
		sp[1] = strconv.Itoa(verPart + 1)
		sp[2] = strconv.Itoa(0)
		nextVersion = strings.Join(sp, ".")
	} else if *patch {
		verPart, _ := strconv.Atoi(sp[2])
		sp[2] = strconv.Itoa(verPart + 1)
		nextVersion = strings.Join(sp, ".")
	}
	return nextVersion
}
