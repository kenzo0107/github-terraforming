package github

import (
	"context"
	"os"
	"sync"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Iface : github interface
type Iface interface {
	GetAllTeams(string) ([]*github.Team, error)
	GetAllRepositories(string) ([]*github.Repository, error)
}

// Instance : github instance
type Instance struct {
	client *github.Client
}

var instance *Instance
var once sync.Once

var (
	githubToken        string
	githubOrganization string
)

// Repository : github repository
type Repository struct {
	github.Repository
	ShapedName string
}

func init() {
	githubToken = os.Getenv("GITHUB_TOKEN")
	githubOrganization = os.Getenv("GITHUB_ORGANIZATION")
}

// New : new github client
func New(token string) Iface {
	once.Do(func() {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		ctx := context.Background()
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)
		instance = &Instance{
			client: client,
		}
	})
	return instance
}

// GetAllTeams : Get all teams of out organization in Github
func (d *Instance) GetAllTeams(org string) (allTeams []*github.Team, err error) {
	page := 1
	for {
		opt := &github.ListOptions{
			Page: page,
		}
		teams, _, err := d.client.Teams.ListTeams(context.Background(), org, opt)
		if err != nil {
			return nil, errors.Wrap(err, "GetAllTeamOfGithub error occured")
		}
		allTeams = append(allTeams, teams...)

		if len(teams) == 0 {
			break
		}
		page++
	}
	return
}

// GetAllRepositories ... Get all repositories of out organization in Github
func (d *Instance) GetAllRepositories(org string) (all []*github.Repository, err error) {
	page := 1
	for {
		opts := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{
				Page: page,
			},
		}
		repos, _, err := d.client.Repositories.ListByOrg(context.Background(), org, opts)
		if err != nil {
			return nil, errors.Wrap(err, "GetAllTeamOfGithub error occured")
		}
		all = append(all, repos...)

		if len(repos) == 0 {
			break
		}
		page++
	}
	return
}
