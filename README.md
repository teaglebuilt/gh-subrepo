## gh subrepo extension

[![CI](https://github.com/teaglebuilt/gh-subrepo/actions/workflows/ci.yaml/badge.svg)](https://github.com/teaglebuilt/gh-subrepo/actions/workflows/ci.yaml)
[![License](https://img.shields.io/badge/License-MIT-default.svg)](./LICENSE.md)

This is an extension to enable management of git submodules with the [github cli](https://cli.github.com/).

## Instructions

### Installation

```
gh extension install teaglebuilt/gh-subrepo
```

### Commands

- **clone**: `gh subrepo clone {repo_url.git} {path}`
- **pull**: `gh subrepo pull {path}`
- **push**: `gh subrepo push {path}`
- **fetch**: `gh subrepo fetch {path}`
- **status**: `gh subrepo status {path}`
- **branch**: `gh subrepo branch {path}`

### Inspiration

⭐ Inspired by [git-subrepo](https://github.com/ingydotnet/git-subrepo)
