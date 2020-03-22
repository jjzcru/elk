![Build Status](https://github.com/jjzcru/elk/workflows/Build%20Status/badge.svg?branch=master)

Elk
==========

Elk 🦌 is a minimalist, [YAML][yaml] based task runner that aims to help developers to focus on building cool stuff,
instead of remembering to perform tedious tasks.

Since it's written in [Go][go], most of the commands runs across multiple operating systems (`Linux`, `macOS`, 
`Windows`) and use the same syntax between them thanks to this [library][sh].

*Why should i use this?* You can watch some [Use Cases]

## Table of contents
  * [Getting Started](#getting-started)
    + [Installation](#installation)
  * [Syntax](#syntax)
  * [Use Cases](#use-cases)
  * [Commands](#commands)
  * [Changelog][changelog]
  * [Releases][releases]

## Getting Started
The main use case for `elk` is that you are able to run any command/s in a declarative way in any path. 

By default the global file that is going to be used is `~/ox.yml`. You can change this path if you wish to use another 
file by setting the `env` variable `ELK_FILE`.

`elk` will first search if there is a `ox.yml` file in the current directory and use that first, if the file is not 
found it will use the `global` file. 

This enables the user to have multiples `ox.yml` one per project while also having one for the system itself.

### Installation
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

## Use Cases
The goal of `elk` is to run `tasks` in a declarative way, anything that you could run on your terminal, you can run 
behind `elk`. If you handle multiple projects, languages, task or you want to automate your workflow you can use `elk`
to achieve that, just declare you workflow and `elk` will take care of the rest.

To learn about some use cases for `elk` go to [Use Cases][use-cases] to figure out 😉.

## Commands

| Command           | Description                                            | Syntax                               |
| -------           | ------                                                 | -------                              |
| [cron][cron]      | Run one or more task as a `cron job` ⏱                | `elk cron [crontab] [tasks] [flags]` |
| [exec][exec]      | Execute ad-hoc commands ⚡                              | `elk exec [commands] [flags]`        |
| [init][init]      | This command creates a dummy file in current directory | `elk init [flags]`                   |
| [logs][logs]      | Attach logs from a task to the terminal 📝             | `elk logs [task] [flags]`            |
| [ls][ls]          | List tasks                                             | `elk ls [flags]`                     |
| [run][run]        | Run one or more tasks 🤖                               | `elk run [tasks] [flags]`            |
| [version][version]| Display version number                                 | `elk version [flags]`                |


[go]: https://golang.org/
[yaml]: https://yaml.org/
[sh]: https://github.com/mvdan/sh

[releases]: https://github.com/jjzcru/elk/releases
[changelog]: https://github.com/jjzcru/elk/blob/master/CHANGELOG.md

[syntax]: docs/syntax/syntax.md
[use-cases]: docs/syntax/use-cases.md

[cron]: docs/commands/cron.md
[init]: docs/commands/init.md
[logs]: docs/commands/logs.md
[ls]: docs/commands/ls.md
[run]: docs/commands/run.md
[version]: docs/commands/version.md
[exec]: docs/commands/exec.md
