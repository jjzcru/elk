![Build Status](https://github.com/jjzcru/elk/workflows/Build%20Status/badge.svg?branch=develop)

Elk
==========

Elk is a minimalist, [YAML][yaml] based task runner that aims to help developers 
to focus on building cool stuff, instead of remembering to perform tedious 
tasks.

Since it's written in [Go][go], most of the commands runs across multiple 
operating systems (`linux`, `macOS`, `Windows`) and use the same syntax between
them thanks to this [library][sh].

Once installed, you will have an example file under `~/elk.yml`. 
To run a task run `elk run [task]`.

```yml
version: '1'
tasks:
  example:
    cmds:
      - echo "Hello world"
```

To run the `example` task will be:

```
elk run example
```

This will print `Hello world`


## Installation

### Download Binary
- Grab the latest binary of your platform from the [Releases][releases] page
- Add the binary to PATH
- Give executable permisions with `chmod`
- Run `elk version` to make sure that the binary is installed
- Run `elk install`

## Usage

### Getting started
The main use case for `elk` is that you are able to run any command/s in a declarative way 
in any path. 

By default the global file that is going to be used is `~/elk.yml`. You can change this 
path if you wish to use another file by setting the env variable `ELK_FILE`.

`elk` will first search if there is a `elk.yml` file in the current directory and use that 
first, if the file is not found it will use the global file. This enables the user to have 
multiples `elk.yml` one per project while also having one for the system itself.

### Syntax

```yml
version: '1'
env_file: /tmp/test.env
env:
  HELLO: WORLD
  MACHINE: WALL-E
tasks:
  # This prints HELLO WORLD
  hello:
    description: "Print hello world"
    env:
      HELLO: HELLO
    cmds:
      - echo $HELLO WORLD 

  # This puts WORLD in the file ./test.log
  test-log:
    description: "Print WORLD"
    log: ./test.log
    cmds:
      - echo $HELLO 
    
  restart:
    description: 'Restart the machine'
    cmds:
      - reboot
  
  shutdown:
    description: 'Command to shutdown the machine'
    cmds:
      - echo "$(hostname) is going to shutdown"
      - shutdown

  cra-example:
    description: "Compile and runs a CRA app"
    dir: /tmp/create-react-app-example
    watch: "[a-zA-Z]*.jsx$" # All .jsx files
    dir: /tmp/create-react-app-example
    deps:
      - name: build
    cmds:
      - lite-server --baseDir="build"

  build:
    dir: /tmp/create-react-app-example
    cmds:
      - npm run build
```

#### Global
Anything that is declared on this level is inherit by the tasks.

`version`

Identifies what is the current version syntax that `elk` is going to interpret

`env_file`

This is a path to a file that declares the `env` variables as `ENV_NAME=ENV_VALUE` 
where each line is a different `env` variable. This overwrites the existing `env`
variable.

`env`

In here you declare all the `env` variable that you wish that all the task inherit
this property overwrites the existing `env` variables, also the ones declared in
the `env_file` property

`tasks`

In here you have a list of all the tasks that you wish to perform. The name 
of the task is going to be used to know which task is going to perform.


#### Task
`env_file`

This is a path to a file that declares the `env` variables as `ENV_NAME=ENV_VALUE` 
where each line is a different `env` variable. This overwrites the existing `env`
variable already declared on global.

`env`

In here you declare all the `env` variable that you wish that all the task inherit
this property overwrites the existing `env` variables, also the ones declared in
the `env_file` property and global.

`description`

In here you describe what is the purpose of the task, this is also display by the
`ls` command.

`dir`

This specifies what is the directory in which the commands are going to run. If not
set is going to use the current directory.

`log`

Saves the output of the command to a file

`watch` 

This is a regex for the files that are going to activate the re-run of the tasks

`deps`

This is a list of all the dependencies that the task requires to run. The `dep` declaration
takes 2 properties:
- `name`: It takes a `string` which is the name of the task that you which to 
run as a dependency
- `detached`: It takes a `boolean` which tells if the dependency should run in `detached` 
mode, is `false` as default.

```yml
test:
    deps:
      - name: build
      - name: hello
        detached: true

```

If a `dep` is run as `detached` it will run without waiting the result of the previous command.
If you are going to run a long running task is recommended to run in detached mode because
the main won't run until all the task that are not detached finish running.

`cmds`

This is a list of all the command that are required to run to perform the `task`. If at least one 
of them fail the entire task fails.

```yml
hello:
    description: "Print hello world"
    env:
      HELLO: HELLO
    cmds:
      - echo $HELLO WORLD 
```

### Commands
- [init](docs/init.md)
- [ls](docs/ls.md)
- [run](docs/run.md)
- [logs](docs/logs.md)


[go]: https://golang.org/
[yaml]: https://yaml.org/
[sh]: https://github.com/mvdan/sh
[releases]: https://github.com/jjzcru/elk/releases