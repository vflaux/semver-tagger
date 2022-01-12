# TL;DR

```
echo "${REGISTRY_PASSWORD}" | semver-tagger auth login --password-stdin "${REGISTRY_URL}" -u "${REGISTRY_USER}"
semver-tagger tag "${IMAGE}" latest '{{major}}' '{{major}}.{{minor}}'
```

# Description

semver-tagger is a tool to add tags to a repository only if the provided image (first argument) is the latest version available on the registry.
The [Semantic Versions](http://semver.org) is read from images metadata `org.opencontainers.image.version` so your images must be correctly labeled.
The tool interact directly with the registry and do not require the image locally (no need to pull/push).

This is usefull in ci if you want indempotent jobs and "latest" tags on your registry.
If you run again a ci job of an old version of your image (you may want to rebuild a lost image from sources), semver-tagger will ensure that you don't override a more recent version of your image.
It is also usefull if you maintains different major versions of your image as you can push them in any order without risking to tag to a wrong version.

For example, with a registry that contains the following images/tags:

- foo/bar:1: v1.0.0
- foo/bar:1.0: v1.0.0
- foo/bar:1.0.0: v1.0.0

Then you push the image:
- foo/bar:1.1.0

And run  `semver-tagger tag foo/bar:1.1.0 foo/bar:latest foo/bar:1 foo/bar:1.1`

This will set the following tags:
- foo/bar:latest: v1.1.0
- foo/bar:1: v1.1.0
- foo/bar:1.1: v1.1.0

But if you (re) run `semver-tagger tag foo/bar:1.0.0 foo/bar:latest foo/bar:1 foo/bar:1.0`
No tag is changed because each tags versions are superior or equals to the version of `foo/bar:1.1.0` (`foo/bar:latest`, `foo/bar:1` & `foo/bar:1.0` ar already at v1.1.0).

# Usage

```
Usage:
  semver-tagger [flags]
  semver-tagger [command]

Available Commands:
  auth        Log in or access credentials
  help        Help about any command
  tag         Tag a remote image if its version is greater than the version of the remote tag

Flags:
  -h, --help                help for semver-tagger
      --insecure            Allow image references to be fetched without TLS
      --platform platform   Specifies the platform in the form os/arch[/variant] (e.g. linux/amd64). (default all)
  -v, --verbose             Enable debug logs

Use "semver-tagger [command] --help" for more information about a command.
```

## Tag

```
Tag a remote image if its version is greater than the version of the remote tag

Usage:
  semver-tagger tag IMAGE [TAG...] [flags]

Flags:
  -h, --help   help for tag

Global Flags:
      --insecure            Allow image references to be fetched without TLS
      --platform platform   Specifies the platform in the form os/arch[/variant] (e.g. linux/amd64). (default all)
  -v, --verbose             Enable debug logs
```

### Templates

Specified tags are interpreted as go templates with the following functions defined:

- `major`: major version of the image semantic version
- `minor`: minor version of the image semantic version
- `patch`: patch version of the image semantic version
- `prerelease`: pre-release version of the image semantic version
- `metadata`: metadata on the version of the image semantic version

example:
```
'{{major}}' '{{major}}.{{minor}}' '{{major}}.{{minor}}.{{patch}}'
```

## Auth

```
Log in or access credentials

Usage:
  semver-tagger auth [flags]
  semver-tagger auth [command]

Available Commands:
  get         Implements a credential helper
  login       Log in to a registry

Flags:
  -h, --help   help for auth

Global Flags:
      --insecure            Allow image references to be fetched without TLS
      --platform platform   Specifies the platform in the form os/arch[/variant] (e.g. linux/amd64). (default all)
  -v, --verbose             Enable debug logs

Use "semver-tagger auth [command] --help" for more information about a command.
```

### get
```
Implements a credential helper

Usage:
  semver-tagger auth get [flags]

Examples:
  # Read configured credentials for reg.example.com
  echo "reg.example.com" | semver-tagger auth get
  {"username":"AzureDiamond","password":"hunter2"}

Flags:
  -h, --help   help for get

Global Flags:
      --insecure            Allow image references to be fetched without TLS
      --platform platform   Specifies the platform in the form os/arch[/variant] (e.g. linux/amd64). (default all)
  -v, --verbose             Enable debug logs
```
### login
```
Log in to a registry

Usage:
  semver-tagger auth login [OPTIONS] [SERVER] [flags]

Examples:
  # Log in to reg.example.com
  crane auth login reg.example.com -u AzureDiamond -p hunter2

Flags:
  -h, --help              help for login
  -p, --password string   Password
      --password-stdin    Take the password from stdin
  -u, --username string   Username

Global Flags:
      --insecure            Allow image references to be fetched without TLS
      --platform platform   Specifies the platform in the form os/arch[/variant] (e.g. linux/amd64). (default all)
  -v, --verbose             Enable debug logs
```

