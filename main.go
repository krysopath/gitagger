package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

var (
	GitlabToken = os.Getenv("CI_BUILD_TOKEN")
	GitlabAddr  = os.Getenv("CI_API_V4_URL")
	Name        = os.Getenv("USER")
	Host        = os.Getenv("HOSTNAME")
	Email       = fmt.Sprintf("%s@%s", Name, Host)
	Message     = "automated"
)

func toContent(contentOrPath string) string {
	var content string

	switch strings.HasPrefix(contentOrPath, "<@") {
	case true:
		path := strings.TrimLeft(contentOrPath, "<@")
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("cant read file at", path, err)
		}
		content = string(data)
	case false:
		content = contentOrPath
	default:
		content = contentOrPath
	}

	return content
}

//App is main cli interface
type App struct {
	client *gitlab.Client
}

func (a *App) createTag(project string, ref string, tag string, msg string) {
	message := toContent(msg)
	p := &gitlab.CreateTagOptions{
		TagName: &tag,
		Ref:     &ref,
		Message: &message,
	}
	_, _, err := a.client.Tags.CreateTag(project, p)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (a *App) createFile(project string, ref string, filename string, contentOrPath string) {
	content := toContent(contentOrPath)
	p := &gitlab.CreateFileOptions{
		Branch:        &ref,
		AuthorName:    &Name,
		AuthorEmail:   &Email,
		Content:       &content,
		CommitMessage: &Message,
	}
	_, _, err := a.client.RepositoryFiles.CreateFile(project, filename, p)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (a *App) updateFile(project string, ref string, filename string, contentOrPath string) {
	content := toContent(contentOrPath)
	p := &gitlab.UpdateFileOptions{
		Branch:        &ref,
		AuthorName:    &Name,
		AuthorEmail:   &Email,
		Content:       &content,
		CommitMessage: &Message,
	}
	val, resp, err := a.client.RepositoryFiles.UpdateFile(project, filename, p)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(val, resp)
}

func (a *App) createUpdateFile(project string, ref string, filename string, contentOrPath string) {
	content := toContent(contentOrPath)
	update := &gitlab.UpdateFileOptions{
		Branch:        &ref,
		AuthorName:    &Name,
		AuthorEmail:   &Email,
		Content:       &content,
		CommitMessage: &Message,
	}
	val, resp, err := a.client.RepositoryFiles.UpdateFile(project, filename, update)
	if err != nil {
		if resp.Response.StatusCode == 400 {
			create := &gitlab.CreateFileOptions{
				Branch:        &ref,
				AuthorName:    &Name,
				AuthorEmail:   &Email,
				Content:       &content,
				CommitMessage: &Message,
			}

			_, _, err := a.client.RepositoryFiles.CreateFile(project, filename, create)
			if err != nil {
				log.Fatalf(err.Error())
			}

		} else {
			log.Fatalf(err.Error())
		}
	}
	fmt.Println(val, resp)
}

//Act decides the subcommand for the cli app
func (a *App) Act() {
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Help")
		os.Exit(1)

	}
	action, extraArgs := args[0], args[1:]
	switch action {
	case "createTag":
		if len(extraArgs) > 0 {
			for i, s := range extraArgs {
				data := strings.Split(s, ":")
				pid, ref, tag, msg := data[0], data[1], data[2], data[3]
				a.createTag(pid, ref, tag, msg)
				fmt.Println(i, pid, ref, tag, msg)
			}
		} else {
			fmt.Println("Help")
			os.Exit(1)
		}

	case "createFile":
		if len(extraArgs) > 0 {
			for i, s := range extraArgs {
				data := strings.Split(s, ":")
				pid, ref, path, content := data[0], data[1], data[2], data[3]
				a.createFile(pid, ref, path, content)
				fmt.Println(i, pid, ref, path, content)
			}
		} else {
			fmt.Println("Help")
			os.Exit(1)
		}
	case "updateFile":
		if len(extraArgs) > 0 {
			for i, s := range extraArgs {
				data := strings.Split(s, ":")
				pid, ref, path, content := data[0], data[1], data[2], data[3]
				a.updateFile(pid, ref, path, content)
				fmt.Println(i, pid, ref, path, content)
			}
		} else {
			fmt.Println("Help")
			os.Exit(1)
		}
	case "createUpdateFile":
		if len(extraArgs) > 0 {
			for i, s := range extraArgs {
				data := strings.Split(s, ":")
				pid, ref, path, content := data[0], data[1], data[2], data[3]
				a.createUpdateFile(pid, ref, path, content)
				fmt.Println(i, pid, ref, path, content)
			}
		} else {
			fmt.Println("Help")
			os.Exit(1)
		}
	default:
		fmt.Println("Help")
	}
}

func init() {
	flag.Parse()
}

func main() {
	git, err := gitlab.NewClient(
		GitlabToken,
		gitlab.WithBaseURL(GitlabAddr))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	app := App{
		client: git,
	}
	app.Act()
}
