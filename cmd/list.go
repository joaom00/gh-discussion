package cmd

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	gh "github.com/cli/go-gh"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/joaom00/gh-discussions/ui"

	"github.com/spf13/cobra"
)

type listOptions struct {
	Repository string
	Limit      int
}

func NewListCmd() *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List of discussions of a repository",
		RunE: func(_ *cobra.Command, _ []string) error {
			if opts.Repository == "" {
				// repo, err := resolveRepository()
				repo, err := gh.CurrentRepository()
				if err != nil {
					return err
				}
				opts.Repository = fmt.Sprintf("%s/%s", repo.Owner(), repo.Name())
			}

			return listRun(opts)

		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository to get discussions")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of items to fetch")

	return cmd
}

func listRun(opts listOptions) error {
	client, err := gh.GQLClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	var query struct {
		Repository struct {
			Discussions struct {
				Nodes []ui.Discussion
			} `graphql:"discussions(first: $first)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	s := strings.Split(opts.Repository, "/")
	owner, name := s[0], s[1]

	variables := map[string]interface{}{
		"first": graphql.Int(opts.Limit),
		"owner": graphql.String(owner),
		"name":  graphql.String(name),
	}

	err = client.Query("Discussions", &query, variables)
	if err != nil {
		log.Fatal(err)
	}

	nodes := query.Repository.Discussions.Nodes
	if len(nodes) == 0 {
		fmt.Printf("There are no discussions in %s\n", opts.Repository)
		return nil
	}

	model := ui.NewModel(&nodes)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	// items := []list.Item{}

	// for _, v := range nodes {
	// 	items = append(items, item{title: v.Title, desc: v.Author.Login})
	// }

	// d := list.NewDefaultDelegate()

	// m := ui.Model{List: list.New(items, d, 0, 0)}

	// p := tea.NewProgram(m, tea.WithAltScreen())

	// if err := p.Start(); err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }

	return nil
}
