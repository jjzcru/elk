Syntax
==========

The syntax consists on two main section one is `global` which serves to set defaults for all the tasks and the other is 
`tasks` which defines the behavior for each of the task.

## Properties
### Global
In the `global` level anything that is declared is inherit by the tasks.

`version`

Identifies what is the current version syntax that `elk` is going to interpret.

`env_file`

This is a path to a file that declares the `env` variables as `ENV_NAME=ENV_VALUE` where each line is a different `env` 
variable. This overwrites the existing `env` variable.

`env`

In here you declare all the `env` variable that you wish that all the task inherit this property overwrites the 
existing `env` variables, also the ones declared in the `env_file` property.

`vars`

It takes a map with all the variables that you wish to include in your program. Once you declared your `vars` you 
can write your `cmds` in [Go Template][go-template] syntax.

`tasks`

In here you have a list of all the tasks that you wish to perform. The name of the task is going to be used to know 
which task is going to perform.

### Task
In the `task` level you can overwrite the values set at the `global` level to this particular task.

`title`

This properties defines what is the title of the task. If not set is going to use the `name` of the task as a default.

`tags`

This propertie is a list of tags that is used to group tasks.

`env_file`

This is a path to a file that declares the `env` variables as `ENV_NAME=ENV_VALUE` where each line is a different 
`env` variable. This overwrites the existing `env` variable already declared on global.

`env`

In here you declare all the `env` variables that you wish that the task uses, `env` declared in here overwrites the 
ones written in the `env_file` property and global.

`vars`

It takes a `map` with all the variables that you wish to include in your program. `vars` declared in here overwrites
the ones that were declared at `global`. Once you declared your `vars` you can write your `cmds` in 
[Go Template][go-template] syntax.

Example: 
```yml
test:
  vars:
    hello: "hello"
  cmds:
    - "echo {{.hello}} world" # This will print "hello world"
```

`description`

In here you describe what is the purpose of the task, this is also display by the `ls` command.

`dir`

This specifies what is the directory in which the commands are going to run. If not set is going to use the current 
directory.

`log`

This properties sets where the output of the command is going to be stores. It has the following properties:
- `out` **Required**: This is the path where the `stdout` is going to be stored.
- `error` *optional*: This is the path where the `stderr` is going to be stored. If not set, is going to save `sterr`
in the same path as the `out` property.
- `format` *optional*: This is a *timestamp* prefix for all the outputs. This is the list of all the available formats:
  - `ANSIC`: *Mon Jan _2 15:04:05 2006*
  - `UnixDate`: *Mon Jan _2 15:04:05 MST 2006*
  - `RubyDate`: *Mon Jan 02 15:04:05 -0700 2006*
  - `RFC822`: *02 Jan 06 15:04 MST*
  - `RFC822Z`: *02 Jan 06 15:04 -0700*
  - `RFC850`: *Monday, 02-Jan-06 15:04:05 MST*
  - `RFC1123`: *Mon, 02 Jan 2006 15:04:05 MST*
  - `RFC1123Z`: *Mon, 02 Jan 2006 15:04:05 -0700*
  - `RFC3339`: *2006-01-02T15:04:05Z07:00*
  - `RFC3339Nano`: *2006-01-02T15:04:05.999999999Z07:00*
  - `Kitchen`: *3:04PM*

Example: 
```yml
test:
  log:
    out: ./hello.log
    error: ./hello-error.log
    format: RFC3339
  cmds:
    - "echo Hello world"
```

`ignore_error`

Ignore errors that happened during a `task`.

`sources` 

This is a regex for the files that are going to activate the re-run of the tasks

`deps`

This is a list of all the dependencies that the task requires to run. The `dep` declaration takes 2 properties:

- `name` **Required**: It takes a `string` which is the name of the task that you which to run as a dependency.

- `detached` *optional*: It takes a `boolean` which tells if the dependency should run in `detached` mode, is `false` 
as default.

- `ignore_error` *optional*: It takes a `boolean` which tells if the program should keep running if an error happens in
the dependency.

Example: 
```yml
test:
  deps:
    - name: build
    - name: hello
      detached: true
```

If a `dep` is run as `detached` it will run without waiting the result of the previous command. If you are going to run 
a long running task is recommended to run in detached mode because the main won’t run until all the task that are not 
detached finish running.

`cmds`

This is a list of all the command that are required to run to perform the `task`. If at least one of them fail the 
entire task fails.

Example:
```yml
hello:
  description: “Print hello world”
  env:
    HELLO: HELLO
  cmds:
    - echo $HELLO WORLD 
```

[go-template]: https://golang.org/pkg/text/template/