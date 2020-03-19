`run`
==========

Run one or more task

## Syntax
```
elk run [tasks] [flags]
```

This command takes at least one argument which is the name of the `task`. You can run multiple `task` in a single command.

You can overwrite properties declared in the `syntax` with `flags`.

### Examples
```
elk run foo
elk run foo bar
elk run foo -d
elk run foo -d -w
elk run foo -t 1s
elk run foo --delay 1s
elk run foo -e FOO=BAR --env HELLO=WORLD
elk run foo -l ./foo.log -d
elk run foo --ignore-logfile
elk run foo --ignore-error
elk run foo --deadline 09:41AM
elk run foo --start 09:41PM
elk run foo -i 2s
elk run foo --interval 2s
```

## Flags

| Flag                                  | Short code | Description                                       | 
| -------                               | ------     | -------                                           | 
| [detached](#detached)                 | d          | Run the task in detached mode and returns the PGID|
| [env](#env)                           | e          | Set env variable to the task/s                    |
| [file](#file)                         | f          | Run task from a file                              |
| [global](#global)                     | g          | Run task from global file                         |
| [help](#help)                         | h          | Help for run                                      |
| [ignore-logfile](#ignore-logfile)     |            | Ignores task log property                         |
| [ignore-error](#ignore-error)         |            | Ignore errors from task                           |
| [delay](#delay)                       |            | Set a delay to a task                             |
| [log](#log)                           | l          | Log output from a task to a file                  |
| [watch](#watch)                       | w          | Enable watch mode                                 |
| [timeout](#timeout)                   | t          | Set a timeout to a task                           |
| [deadline](#deadline)                 |            | Set a deadline to a task                          |
| [start](#start)                       |            | Set a date/datetime to a task                     |
| [interval](#interval)                 | i          | Set a duration for an interval                    | 

### detached

This will group all the tasks under the same `PGID` and then it will detach from the process, and returns the `PGID` so 
the user can kill the process later.

Example:

```
elk run test -d
elk run test --detached
```

### env

This flag will overwrite whatever env variable already declared in the file. You can call this flag multiple times.

Example:
```
elk run test -e HELLO=WORLD --env FOO=BAR
```

### file

This flag force `elk` to use a particular file path to run the commands.

Example:
```
elk run test -f ./elk.yml
elk run test --file ./elk.yml
```

### global

This force the task to run from the global file either declared at `ELK_FILE` or the default global path `~/elk.yml`.

Example:

```
elk run test -g
elk run test --global
```

### ignore-log

Force task to output to stdout.

Example:

```
elk run test --ignore-logfile
```

### ignore-error

Ignore errors that happened during a `task`.

Example:

```
elk run test --ignore-error
```

### delay

This flag will run the task after some duration.

This commands supports the following duration units:
- `ns`: Nanoseconds
- `ms`: Milliseconds
- `s`: Seconds
- `m`: Minutes
- `h`: Hours

Example:

```
elk run test --delay 1s
elk run test --delay 500ms
elk run test --delay 2h
elk run test --delay 2h45m
```

### log

This saves the output to a specific file.

Example:

```
elk run test -l ./test.log
elk run test --log ./test.log
```

### watch

This requires that the task has a property `sources` already setup, otherwise it will throw an error. When this flag is 
enable it will kill the existing process and create a new one every time a file that match the regex is changed.

The property `sources` uses a `go` regex to search for all the paths, inside the `dir` property, that matches the 
criteria and adds a `watcher` to all the files.

Example:

```
elk run test -w
elk run test --watch
```

### timeout

This flag with kill the task after some duration since the program was started.

This commands supports the following duration units:
- `ns`: Nanoseconds
- `ms`: Milliseconds
- `s`: Seconds
- `m`: Minutes
- `h`: Hours

Example:

```
elk run test -t 1s
elk run test --timeout 500ms
elk run test --timeout 2h
elk run test --timeout 2h45m
```

### deadline

This flag with kill the task at a particular datetime.

It supports the following datetime standards:
- `ANSIC`: `Mon Jan _2 15:04:05 2006`
- `UnixDate`: `Mon Jan _2 15:04:05 MST 2006`
- `RubyDate`: `Mon Jan 02 15:04:05 -0700 2006`
- `RFC822`: `02 Jan 06 15:04 MST`
- `RFC822Z`: `02 Jan 06 15:04 -0700`
- `RFC850`: `Monday, 02-Jan-06 15:04:05 MST`
- `RFC1123`: `Mon, 02 Jan 2006 15:04:05 MST`
- `RFC1123Z`: `Mon, 02 Jan 2006 15:04:05 -0700`
- `RFC3339`: `2006-01-02T15:04:05Z07:00`
- `RFC3339Nano`: `2006-01-02T15:04:05.999999999Z07:00`
- `Kitchen`: `3:04PM`

If the `Kitchen` format is used and the time is before the current time it will run at the same time in the following 
day.

Example:

```
elk run test --deadline 11:00PM
elk run test --deadline 2020-12-12T09:41:00Z00:00
```

### start

This flag with run the task at a particular datetime.

It supports the following datetime standards:
- `ANSIC`: `Mon Jan _2 15:04:05 2006`
- `UnixDate`: `Mon Jan _2 15:04:05 MST 2006`
- `RubyDate`: `Mon Jan 02 15:04:05 -0700 2006`
- `RFC822`: `02 Jan 06 15:04 MST`
- `RFC822Z`: `02 Jan 06 15:04 -0700`
- `RFC850`: `Monday, 02-Jan-06 15:04:05 MST`
- `RFC1123`: `Mon, 02 Jan 2006 15:04:05 MST`
- `RFC1123Z`: `Mon, 02 Jan 2006 15:04:05 -0700`
- `RFC3339`: `2006-01-02T15:04:05Z07:00`
- `RFC3339Nano`: `2006-01-02T15:04:05.999999999Z07:00`
- `Kitchen`: `3:04PM`

If the `Kitchen` format is used and the time is before the current time it will run at the same time in the following 
day.

Example:

```
elk run test --start 11:00PM
elk run test --start 2020-12-12T09:41:00Z00:00
```

### interval

This flag will run a task in a new process every time the interval ticks. Enabling `interval` disables the `watch` mode.

This commands supports the following duration units:
- `ns`: Nanoseconds
- `ms`: Milliseconds
- `s`: Seconds
- `m`: Minutes
- `h`: Hours

Example:

```
elk run test -i 1s
elk run test --interval 500ms
elk run test --interval 2h
elk run test --interval 2h45m
```