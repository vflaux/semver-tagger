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
	"net/http"
	"os"

	"github.com/docker/cli/cli/config"
	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/spf13/cobra"
	"github.com/vflaux/semver-tagger/pkg/semvertagger"
)

func init() {
	Root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable debug logs")
	Root.PersistentFlags().BoolVar(&insecure, "insecure", false, "Allow image references to be fetched without TLS")
	Root.PersistentFlags().Var(platform, "platform", "Specifies the platform in the form os/arch[/variant] (e.g. linux/amd64).")
}

const commandName = "semver-tagger"

var (
	verbose  = false
	insecure = false
	platform = &platformValue{}

	// Crane options for this invocation.
	options = []semvertagger.Option{}

	// Root is the top-level cobra.Command for semvertagger.
	Root = &cobra.Command{
		Use:               commandName,
		Short:             commandName + " is a tool to remotely tag an image if it is the latest version",
		Run:               func(cmd *cobra.Command, _ []string) { _ = cmd.Usage() },
		DisableAutoGenTag: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				logs.Debug.SetOutput(os.Stderr)
			}
			if insecure {
				options = append(options, semvertagger.Insecure)
			}

			options = append(options, semvertagger.WithPlatform(platform.platform))

			// Add any http headers if they are set in the config file.
			cf, err := config.Load(os.Getenv("DOCKER_CONFIG"))
			if err != nil {
				logs.Debug.Printf("failed to read config file: %v", err)
			} else if len(cf.HTTPHeaders) != 0 {
				options = append(options, semvertagger.WithTransport(&headerTransport{
					inner:       http.DefaultTransport,
					httpHeaders: cf.HTTPHeaders,
				}))
			}
		},
	}
)

// headerTransport sets headers on outgoing requests.
type headerTransport struct {
	httpHeaders map[string]string
	inner       http.RoundTripper
}

// RoundTrip implements http.RoundTripper.
func (ht *headerTransport) RoundTrip(in *http.Request) (*http.Response, error) {
	for k, v := range ht.httpHeaders {
		if http.CanonicalHeaderKey(k) == "User-Agent" {
			// Docker sets this, which is annoying, since we're not docker.
			// We might want to revisit completely ignoring this.
			continue
		}
		in.Header.Set(k, v)
	}
	return ht.inner.RoundTrip(in)
}
