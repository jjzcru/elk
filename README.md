![Build Status](https://github.com/jjzcru/elk/workflows/Build%20Status/badge.svg?branch=master)

Elk
==========

Elk is a minimalist, [YAML][yaml] based task runner that aims to help developers to focus on building cool stuff, 
instead of remembering to perform tedious tasks.

Since it's written in [Go][go], most of the commands runs across multiple operating systems (`Linux`, `macOS`, 
`Windows`) and use the same syntax between them thanks to this [library][sh].

## Table of contents
  * [Getting Started](#getting-started)
    + [Installing](#installing)
  * [Syntax](#syntax)
  * [Commands](#commands)
  * [Changelog][changelog]
  * [Releases][releases]

## Getting Started
The main use case for `elk` is that you are able to run any command/s in a declarative way in any path. 

By default the global file that is going to be used is `~/elk.yml`. You can change this path if you wish to use another 
file by setting the `env` variable `ELK_FILE`.

`elk` will first search if there is a `elk.yml` file in the current directory and use that first, if the file is not 
found it will use the `global` file. 

This enables the user to have multiples `elk.yml` one per project while also having one for the system itself.

### Installing
1. Grab the latest binary of your platform from the [Releases](https://github.com/jjzcru/elk/releases) page.
2. If you are running on `macOS` or `Linux`, run `chmod +x elk` to give `executable` permissions to the binary. If you
are on `windows` you can ignore this step.
3. Add the binary to `$PATH`.
4. Run `elk version` to make sure that the binary is installed.

## Syntax
The syntax consists on two main section one is `global` which sets defaults for all the tasks and the other is `tasks` 
which defines the behavior for each of the task.

To learn about the properties go to [Syntax Documentation][syntax].

### Example

```yml
version: ‘1’
env_file: /tmp/test.env
env:
  FOO: BAR
tasks:
  hello:
    description: “Print Hello World”
    env:
      FOO: Hello
      BAR: World
    cmds:
      - echo $FOO $BAR
```

## Commands

| Command       | Description                                      | Syntax                               | 
| -------       | ------                                           | -------                              | 
| [cron][cron]  | Run one or more task as a `cron job`             | `elk cron [crontab] [tasks] [flags]` |
| [init][init]  | This command creates a dummy file in current dir | `elk init`                           |
| [logs][logs]  | Attach logs from a task to the terminal          | `elk logs [task] [flags]`            |
| [ls][ls]      | List tasks                                       | `elk ls [flags]`                     |
| [run][run]    | Run one or more task                             | `elk run [tasks] [flags]`            |

[go]: https://golang.org/
[yaml]: https://yaml.org/
[sh]: https://github.com/mvdan/sh

[releases]: https://github.com/jjzcru/elk/releases
[changelog]: https://github.com/jjzcru/elk/blob/master/CHANGELOG.md

[cron]: docs/commands/cron.md
[init]: docs/commands/init.md
[logs]: docs/commands/logs.md
[ls]: docs/commands/ls.md
[run]: docs/commands/run.md
[syntax]: docs/syntax.md