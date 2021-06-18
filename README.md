# golang-project-template
An opinionated repository structure for projects written in Go. This repository
contains the common configurations I use across my Go projects. Each of these
configuration files contain what I would consider to be "sane" defaults.

This project template is also very Github-oriented. It uses Github Actions and
the Github Container Registry (GHCR). If you wish to use this template elsewhere
just note that you will have to make fairly significant changes to make it work.

## Getting Started
1. Click _Use this template_ above to create a repository using this template.
Follow Github's prompts.
2. Clone the new repository you just made.
3. Run the following shell script to substitute all of the placeholder variables
in all of the configuration files. When prompted, enter the name of the
project's license.
```bash
$ ./substitute.sh
```
4. Delete the substitution script and edit this file.
```
$ git rm substitute.sh
$ $EDITOR README.md
```