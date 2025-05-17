## gh subrepo extension

![GitHub Release](https://img.shields.io/github/v/release/teaglebuilt/gh-subrepo)
[![License](https://img.shields.io/badge/License-MIT-default.svg)](./LICENSE.md)

This is an extension to enable management of git submodules with the [github cli](https://cli.github.com/).

## Instructions

### Installation

1. Install github cli extension. `gh extension install teaglebuilt/gh-subrepo`

2. Run `gh subrepo init` to setup the tool. A folder will be created in your current working directory.

3. Your good to go..`gh subrepo -h`

### Commands
```
Commands:
      clone     Clone a remote repository into a local subdirectory
      init      Turn a current subdirectory into a subrepo
      pull      Pull upstream changes to the subrepo
      push      Push local subrepo changes upstream

      fetch     Fetch a subrepo's remote branch (and create a ref for it)
      branch    Create a branch containing the local subrepo commits
      commit    Commit a merged subrepo branch into the mainline

      status    Get status of a subrepo (or all of them)
      clean     Remove branches, remotes and refs for a subrepo
      config    Set subrepo configuration properties

      help      Documentation for git-subrepo (or specific command)
      version   Display git-subrepo version info
      upgrade   Upgrade the git-subrepo software itself
```

### Options
```
Show the command summary
    --[no-]help           Help overview
    --[no-]version        Print the git-subrepo version number
    -a, --[no-]all        Perform command on all current subrepos
    -A, --[no-]ALL        Perform command on all subrepos and subsubrepos
    -b, --[no-]branch ... Specify the upstream branch to push/pull/fetch
    -e, --[no-]edit       Edit commit message
    -f, --[no-]force      Force certain operations
    -F, --[no-]fetch      Fetch the upstream content first
    -M, --[no-]method ... Join method: 'merge' (default) or 'rebase'
    -m, --[no-]message ...
                          Specify a commit message
    --[no-]file ...       Specify a commit message file
    -r, --[no-]remote ... Specify the upstream remote to push/pull/fetch
    -s, --[no-]squash     Squash commits on push
    -u, --[no-]update     Add the --branch and/or --remote overrides to .gitrepo
    -q, --[no-]quiet      Show minimal output
    -v, --[no-]verbose    Show verbose output
    -d, --[no-]debug      Show the actual commands used
    -x, --[no-]DEBUG      Turn on -x Bash debugging
```
### Inspiration

⭐ Inspired by [git-subrepo](https://github.com/ingydotnet/git-subrepo)
