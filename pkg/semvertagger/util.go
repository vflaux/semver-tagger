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
	"fmt"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type Options struct {
	name   []name.Option
	remote []remote.Option
}

func MakeOptions(opts ...Option) Options {
	opt := Options{
		remote: []remote.Option{
			remote.WithAuthFromKeychain(authn.DefaultKeychain),
		},
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Option is a functional option for crane.
type Option func(*Options)

// WithTransport is a functional option for overriding the default transport
// for remote operations.
func WithTransport(t http.RoundTripper) Option {
	return func(o *Options) {
		o.remote = append(o.remote, remote.WithTransport(t))
	}
}

// Insecure is an Option that allows image references to be fetched without TLS.
func Insecure(o *Options) {
	o.name = append(o.name, name.Insecure)
}

func ParseReference(r string, opt ...Option) (name.Reference, error) {
	o := MakeOptions(opt...)
	ref, err := name.ParseReference(r, o.name...)
	if err != nil {
		return nil, fmt.Errorf("parsing reference %q: %v", r, err)
	}
	return ref, nil
}
