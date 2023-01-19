// Copyright 2022 vflaux

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 		http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vflaux/semver-tagger/pkg/semvertagger"
)

func init() { Root.AddCommand(NewCmdTag(&options)) }

// NewCmdTag creates a new cobra.Command for the tag subcommand.
func NewCmdTag(options *[]semvertagger.Option) *cobra.Command {
	var versionLabelKey string

	cmd := &cobra.Command{
		Use:   "tag IMAGE [TAG...]",
		Short: "Tag a remote image if its version is greater than the version of the remote tag",
		Args:  cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ref := args[0]
			tags := args[1:]

			if len(tags) == 0 {
				tags = []string{"latest"}
			}

			err := semvertagger.Tag(ref, tags, versionLabelKey, *options...)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&versionLabelKey, "version-label", "org.opencontainers.image.version", "Specifies the image label containing the version.")

	return cmd
}
