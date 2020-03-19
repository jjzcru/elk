cron
==========

Run one or more task as a `cron job`.

## Syntax
```
elk cron [crontab] [tasks] [flags]
```

This command takes at least two arguments. The first one is going to be `crontab` which is the syntax used to describe 
a `cron job`.

The rest of the arguments are the names of the `task` that are going to be executed follow by the flags.

## Examples

```
elk cron "*/1 * * * *" foo
elk cron "*/1 * * * *" foo bar
elk cron "*/1 * * * *" foo -d
elk cron "*/2 * * * *" foo -t 1s
elk cron "*/2 * * * *" foo --delay 1s
elk cron "*/2 * * * *" foo -e FOO=BAR --env HELLO=WORLD
elk cron "*/6 * * * *" foo -l ./foo.log -d
elk cron "*/1 * * * *" foo --ignore-logfile
elk cron "*/2 * * * *" foo --ignore-error
elk cron "*/5 * * * *" foo --deadline 09:41AM
elk cron "*/1 * * * *" foo --start 09:41PM
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


### detached

This will group all the tasks under the same `PGID` and then it will detach from the process, and returns the `PGID` so 
the user can kill the process later.

Example:

```
elk cron "* * * * *" test -d
elk cron "* * * * *" test --detached
```

### env

This flag will overwrite whatever env variable already declared in the file. You can call this flag multiple times.

Example:
```
elk cron "* * * * *" test -e HELLO=WORLD --env FOO=BAR
```

### file

This flag force `elk` to use a particular file path to run the commands.

Example:
```
elk cron "* * * * *" test -f ./elk.yml
elk cron "* * * * *" test --file ./elk.yml
```

### global

This force the task to run from the global file either declared at `ELK_FILE` or the default global path `~/elk.yml`.

Example:

```
elk cron "* * * * *" test -g
elk cron "* * * * *" test --global
```

### ignore-log

Force task to output to stdout.

Example:

```
elk cron "* * * * *" test --ignore-logfile
```

### ignore-error

Ignore errors that happened during a `task`.

Example:

```
elk cron "* * * * *" test --ignore-error
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
elk cron "* * * * *" test --delay 1s
elk cron "* * * * *" test --delay 500ms
elk cron "* * * * *" test --delay 2h
elk cron "* * * * *" test --delay 2h45m
```

### log

This saves the output to a specific file.

Example:

```
elk cron "* * * * *" test -l ./test.log
elk cron "* * * * *" test --log ./test.log
```

### watch

This requires that the task has a property `sources` already setup, otherwise it will throw an error. When this flag is 
enable it will kill the existing process and create a new one every time a file that match the regex is changed.

The property `sources` uses a `go` regex to search for all the paths, inside the `dir` property, that matches the 
criteria and adds a `watcher` to all the files.

Example:

```
elk cron "* * * * *" test -w
elk cron "* * * * *" test --watch
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
elk cron "* * * * *" test -t 1s
elk cron "* * * * *" test --timeout 500ms
elk cron "* * * * *" test --timeout 2h
elk cron "* * * * *" test --timeout 2h45m
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
elk cron "* * * * *" test --deadline 11:00PM
elk cron "* * * * *" test --deadline 2020-12-12T09:41:00Z00:00
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
elk cron "* * * * *" test --deadline 11:00PM
elk cron "* * * * *" test --deadline 2020-12-12T09:41:00Z00:00
```