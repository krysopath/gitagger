package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

var (
	GitlabtToken = os.Getenv("CI_BUILD_TOKEN")
	GitlabAddr   = os.Getenv("CI_API_V4_URL")
)

type GitlabData struct {
	users []*gitlab.User
}

type App struct {
	client *gitlab.Client
	data   GitlabData
}

func (a *App) listUsers() {
	users, _, err := a.client.Users.ListUsers(&gitlab.ListUsersOptions{})
	if err != nil {
		fmt.Println(err)
	}
	a.data.users = users

}

func (a *App) createTag(project string, ref string, tag string, msg string) {
	p := &gitlab.CreateTagOptions{
		TagName: &tag,
		Ref:     &ref,
		Message: &msg,
	}
	_, _, err := a.client.Tags.CreateTag(project, p)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (a *App) Act() {
	for i, s := range flag.Args() {
		data := strings.Split(s, ":")
		pid, ref, tag, msg := data[0], data[1], data[2], data[3]
		a.createTag(pid, ref, tag, msg)
		fmt.Println(i, pid, ref, tag, msg)
	}
}

func init() {
	flag.Parse()
}

func main() {
	git, err := gitlab.NewClient(
		GitlabtToken,
		gitlab.WithBaseURL(GitlabAddr))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	app := App{
		client: git,
	}
	app.Act()
}
