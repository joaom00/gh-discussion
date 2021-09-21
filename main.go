package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/cli/browser"
	"github.com/cli/safeexec"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

type Node struct {
	Title string
	Url   string
}

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use: "discussions",
	}
}

func viewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "List of discussions of a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runView()

		},
	}
	return cmd
}

func runView() error {
	query := `query { 
	repository(owner: "vercel", name: "next.js") {
        discussions(first: 10) {
            nodes {
                title
                url
            }
        }
    }
}`

	cmdArgs := []string{
		"api", "graphql", "--paginate",
		"-f", fmt.Sprintf("query=%s", query),
	}

	out, _, err := gh(cmdArgs...)

	resp := map[string]map[string]map[string]map[string][]Node{}

	err = json.Unmarshal(out.Bytes(), &resp)
	if err != nil {
		return fmt.Errorf("failed to deserialize JSON: %w", err)
	}

	nodes := resp["data"]["repository"]["discussions"]["nodes"]

	idx, err := fuzzyfinder.Find(nodes, func(i int) string {
		return nodes[i].Title
	})

	if err != nil {
		log.Fatal(err)
	}

	err = browser.OpenURL(nodes[idx].Url)
	if err != nil {
		fmt.Printf("Cannot open url: %s", nodes[idx].Url)
	}

	return nil
}

func gh(args ...string) (sout, eout bytes.Buffer, err error) {
	ghBin, err := safeexec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("could not find gh. Is it installed? error: %w", err)
		return
	}

	cmd := exec.Command(ghBin, args...)
	cmd.Stderr = &eout
	cmd.Stdout = &sout

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run gh. error: %w, stderr: %s", err, eout.String())
		return
	}

	return

}

func main() {
	rc := rootCmd()
	rc.AddCommand(viewCmd())

	if err := rc.Execute(); err != nil {
		os.Exit(1)
	}
}
