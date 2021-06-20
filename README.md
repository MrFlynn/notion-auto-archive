# Notion Auto-Archive

[![Tests](https://github.com/MrFlynn/notion-auto-archive/actions/workflows/test.yml/badge.svg)](https://github.com/MrFlynn/notion-auto-archive/actions/workflows/test.yml)

Automatically archives old tasks on a Notion task board. This is to prevent your
column of completed tasks from getting filled with old tasks.

## Getting Started
First, you will need to get an API key and the ID(s) of the task board(s) you wish
to manage with this application. Notion has a
[pretty good guide](https://developers.notion.com/docs) on how to get that
information.

Next, you will need to create a configuration file. You can name if whatever you
want, but in this example I'll call it `config.yml`. The configuration file has
the following format.

```yaml
apiKey: <api-key> # Can set using CLI or environment variables.
boards:
  - id: <board-id>
    archiveAfter: 24h
    selectors:
      columnName: Status
      sourceColumn: Completed
      targetColumn: Archived
  # More boards go here.
```

See the section titled **Configuration Options** for more details on what each
option does.

## Downloading the Application
Head to the [releases page](https://github.com/MrFlynn/notion-auto-archive/releases)
and grab a copy of the application for your platform, or you can use one of the
following alternative methods. Extract the downloaded binary and add it to your path.

### Using Go
You can run `go get github.com/mrflynn/notion-auto-archive` and it will download
and build the latest release of the application.

### Using Docker
This application is also available through Docker. The Docker image is available 
on all of the same platforms that the binaries are available on. Just run
`docker pull ghcr.io/mrflynn/notion-auto-archive:latest` to get the latest 
Docker image.

## Running the Application
### Using a Binary
Simply run the executable you downloaded with the following flags.

```bash
$ notion-auto-archive -config=/path/to/config.yml
```

You can specify the `-key` flag if your API key is not in your configuration file,
or by setting the `NOTION_AUTO_ARCHIVE_API_KEY` environment variable.

### Using Docker
Similarly, you can run the Dockerfile as follows.

```bash
$ docker run --rm --name notion-auto-archive \
    -v /path/to/config.yml:/config.yml:ro \
    ghcr.io/mrflynn/notion-auto-archive:latest
```

Refer to [Docker's docs](https://docs.docker.com/engine/reference/commandline/run/)
for how to set environment variables or change CLI flags.

## Configuration Options
The table below describes all of the configuration fields found in `config.yml`.

| Option                       | Default Value | Description |
| ---------------------------- | ------------- | ----------- |
| apiKey                       | Empty string  | Notion.so API key. Can also be set using CLI or environment variable. |
| boards                       | Empty array   | List of task board this application will manage. |
| boards.id                    | Empty string  | ID of board. |
| boards.archiveAfter          | 24h           | After how much time of inactivity should the task be moved. Must be in seconds, minutes, or hours. |
| boards.selector.columnName   | Status        | Name of selector field to group tasks (i.e. not started, in progress, complete, etc.) |
| boards.selector.sourceColumn | Completed     | Name of column from which to move tasks. |
| boards.selector.targetColumn | Archived      | Name of column to move tasks to. |