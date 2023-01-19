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

package semvertagger

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
)

func Tag(refStr string, tags []string, versionLabelKey string, opt ...Option) error {

	o := MakeOptions(opt...)
	ref, err := name.ParseReference(refStr, o.name...)
	if err != nil {
		return fmt.Errorf("parsing reference %q: %v", refStr, err)
	}

	imageVersionStr, err := getVersion(ref, versionLabelKey)
	if err != nil {
		return err
	}

	var desc *remote.Descriptor

	imageVersion, err := semver.NewVersion(imageVersionStr)
	if err != nil {
		return fmt.Errorf("parsing image version %q: %w: %q", refStr, err, imageVersionStr)
	}

	err = parseTags(tags, imageVersion)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		tagRef := ref.Context().Tag(tag)
		tagVersionStr, err := getVersion(tagRef, versionLabelKey, o.remote...)
		if err != nil {
			transportErr := errors.Unwrap(err)
			if transportErr != nil {
				if transportErr.(*transport.Error).Errors[0].Code != transport.ManifestUnknownErrorCode {
					log.Printf("getting version %q: %v", ref, err)
					continue
				}
			} else {
				log.Printf("getting version %q: %v", ref, err)
				continue
			}
		}

		var tagVersion *semver.Version
		if tagVersionStr != "" {
			tagVersion, err = semver.NewVersion(tagVersionStr)
			if err != nil {
				log.Printf("parsing tag version %q: %v: %q", tag, err, tagVersionStr)
				continue
			}
		}

		if tagVersion == nil || imageVersion.GreaterThan(tagVersion) {
			if desc == nil {
				desc, err = remote.Get(ref, o.remote...)
				if err != nil {
					log.Printf("fetching %q: %v", ref, err)
					continue
				}
			}
			err := remote.Tag(tagRef, desc, o.remote...)
			if err != nil {
				log.Printf("failed to tag %q as %q: %v", ref, tagRef, err)
				continue
			}
			log.Printf("tagged %q as %q", ref.Identifier(), tagRef.Identifier())
		} else {
			log.Printf("tag %q is already the latest: %q >= %q", tag, tagVersionStr, imageVersionStr)
		}
	}

	return nil
}

func getVersion(ref name.Reference, versionLabelKey string, options ...remote.Option) (string, error) {
	img, err := remote.Image(ref, options...)
	if err != nil {
		return "", fmt.Errorf("reading image: %w", err)
	}

	config, err := img.ConfigFile()
	if err != nil {
		return "", fmt.Errorf("getting config file: %w", err)
	}

	version, found := config.Config.Labels[versionLabelKey]
	if !found {
		return "", fmt.Errorf("no version found in image labels")
	}

	return version, nil
}

func parseTags(tags []string, version *semver.Version) error {
	var err error
	for i := range tags {
		tags[i], err = parseTag(tags[i], version)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseTag(tag string, version *semver.Version) (string, error) {
	buf := bytes.Buffer{}
	t := template.New("")

	t.Funcs(template.FuncMap{
		"major":      func() uint64 { return version.Major() },
		"minor":      func() uint64 { return version.Minor() },
		"patch":      func() uint64 { return version.Patch() },
		"prerelease": func() string { return version.Prerelease() },
		"metadata":   func() string { return version.Metadata() },
	})

	t, err := t.Parse(tag)
	if err != nil {
		return "", fmt.Errorf("parsing tag %q: %w", tag, err)
	}
	if err := t.Execute(&buf, version); err != nil {
		return "", fmt.Errorf("parsing tag %q: %w", tag, err)
	}
	return buf.String(), nil
}
