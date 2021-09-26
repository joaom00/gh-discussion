package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	// "github.com/cli/browser"
	"github.com/joaom00/gh-discussions/gh"
	// "github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

type Node struct {
	Title  string
	URL    string
	Author struct {
		Login string
	}
}

type discussionsResponse struct {
	Data struct {
		Repository struct {
			Discussions struct {
				Nodes []Node
			}
		}
	}
}

type listOptions struct {
	Repository string
	Limit      int
}

func NewListCmd() *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List of discussions of a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Repository == "" {
				repo, err := resolveRepository()
				if err != nil {
					return err
				}
				opts.Repository = repo
			}
			return listRun(opts)

		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository to get discussions")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of items to fetch")

	return cmd
}

func resolveRepository() (string, error) {
	cmdArgs := []string{
		"repo", "view",
	}

	out, eout, err := gh.Command(cmdArgs...)
	if err != nil {
		if strings.Contains(eout.String(), "not a git repository") {
			return "", errors.New("Try running this command from inside a git repository or with the -R flag")
		}
		return "", err
	}

	viewOut := strings.Split(out.String(), "\n")[0]
	repo := strings.TrimSpace(strings.Split(viewOut, ":")[1])

	return repo, nil

}

func listRun(opts listOptions) error {
	query := `query($owner: String!, $name: String!, $first: Int!) { 
	repository(owner: $owner, name: $name) {
        discussions(first: $first) {
            nodes {
                title
                url
                author {
                    login
                }
            }
        }
    }
}`

	ownerName := strings.Split(opts.Repository, "/")

	if opts.Limit > 100 {
		opts.Limit = 100
	}

	cmdArgs := []string{
		"api", "graphql",
		"-F", fmt.Sprintf("owner=%s", ownerName[0]),
		"-F", fmt.Sprintf("name=%s", ownerName[1]),
		"-F", fmt.Sprintf("first=%d", opts.Limit),
		"-f", fmt.Sprintf("query=%s", query),
	}

	out, _, err := gh.Command(cmdArgs...)
	if err != nil {
		return err
	}

	var resp discussionsResponse
	err = json.Unmarshal(out.Bytes(), &resp)
	if err != nil {
		return fmt.Errorf("failed to deserialize JSON: %w", err)
	}

	nodes := resp.Data.Repository.Discussions.Nodes
	if len(nodes) == 0 {
		fmt.Printf("There are no discussions in %s\n", opts.Repository)
		return nil
	}

	for _, v := range nodes {
		fmt.Printf("%s\t\t%s\n", v.Author.Login, v.Title)
	}

	// idx, err := fuzzyfinder.Find(
	// 	nodes,
	// 	func(i int) string {
	// 		return nodes[i].Title
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("Opening %s in your browser\n", nodes[idx].URL)
	// err = browser.OpenURL(nodes[idx].URL)
	// if err != nil {
	// 	return err
	// }

	return nil
}
