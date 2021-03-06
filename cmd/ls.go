// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/iancmcc/jig/config"
	"github.com/iancmcc/jig/fs"
	"github.com/iancmcc/jig/match"
	"github.com/iancmcc/jig/vcs"
	"github.com/spf13/cobra"
)

var (
	limit int
	all   bool
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List repositories",
	Long:  `List repositories below the current directory, optionally sorted by similarity to a search string`,
	Run: func(cmd *cobra.Command, args []string) {
		here, _ := filepath.Abs("")
		root, err := config.FindClosestJigRoot("")
		if err != nil {
			logrus.Fatal("No jig root found. Use 'jig init' to create one.")
		}
		var repos <-chan string
		if all {
			repos = fs.DefaultFinder().FindBelowWithChildrenNamed(root, ".git", 1)
		} else {
			ch := make(chan string)
			repos = ch
			go func() {
				defer close(ch)
				manifest, err := config.DefaultManifest("")
				if err != nil {
					return
				}
				for _, r := range manifest.Repos {
					path, err := vcs.RepoToPath(r.Repo)
					if err != nil {
						continue
					}
					ch <- filepath.Join(root, path)
				}
			}()
		}
		if len(args) == 0 {
			var i int
			for repo := range repos {
				if limit > 0 && i >= limit {
					break
				}
				rel, _ := filepath.Rel(here, repo)
				fmt.Println(rel)
				i++
			}
			return
		}
		matcher := match.DefaultMatcher(args[0])
		for repo := range repos {
			matcher.Add(strings.TrimPrefix(repo, root))
		}
		for i, repo := range matcher.Match() {
			if limit > 0 && i >= limit {
				break
			}
			rel, _ := filepath.Rel(here, filepath.Join(root, repo))
			fmt.Println(rel)
		}
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().IntVarP(&limit, "limit", "n", 0, "Limit the number of results returned (default is no limit)")
	lsCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "Show all repositories, not just those in the manifest")
}
