exec
==========

Execute ad-hoc commands ⚡

## Syntax
```
elk exec [commands] [flags]
```
This commands enables you to run commands without declaring them in a `.yml` file. 

You can execute multiple commands at the same time like `elk exec "clear" "curl -s http://localhost:8080/health"`.

You can specify what the behavior of the commands with the available `flags`.

### Examples
```
elk exec "echo Hello World"
elk exec "clear" "curl -s http://localhost:8080/health"
elk exec "clear" "curl -s http://localhost:8080/health" -i 2s
elk exec "echo This is: {{.foo}}" -v foo=bar
elk exec "echo This is $bar" -e bar=foo
elk exec "echo $foo $bar" --env-file ./example.env
elk exec "exit 1" "echo hello world" --ignore-error
elk exec "echo Hello World" --delay 1s
elk exec "echo Hello World" --start 09:41AM
elk exec "echo Hello World" --deadline 09:41AM
elk exec "echo Hello World" --timeout 5s
```

## Flags

| Flag                                  | Short code | Description                                       | 
| -------                               | ------     | -------                                           | 
| [detached](#detached)                 | d          | Run the task in detached mode and returns the PGID|
| [env](#env)                           | e          | Set `env` variable to the command/s               |
| [env-file](#env-file)                 |            | Set `env` variable to the command/s with a file   |
| [var](#var)                           | v          | Set `var` variable to the command/s               |
| [delay](#delay)                       |            | Set a delay to the commands                       |
| [dir](#dir)                           |            | Set a directory to the commands                   |
| [log](#log)                           | l          | Log output to a file                              |
| [ignore-error](#ignore-error)         |            | Ignore errors from the commands                   |
| [timeout](#timeout)                   | t          | Set a timeout to the commands                     |
| [deadline](#deadline)                 |            | Set a deadline to the commands                    |
| [start](#start)                       |            | Set a date/datetime to the commands               |
| [interval](#interval)                 | i          | Set a duration for an interval                    | 

### detached

This will group all the commands under the same `PGID`, detach from the process and returns the `PGID` so the user can 
kill the process later.

Example:

```
elk exec "echo Hello World" -d
elk exec "echo Hello World" --detached
```

### env

This flag will set `env` variables for the commands. This flag can be called multiple times.

Example:
```
elk exec "echo This is $bar $foo" -e bar=foo --env foo=bar
elk exec "curl $url/health" -e url="http://localhost:8080"
```

### env-file

This flag will let the user load `env` variables from a file.

Example:
```
elk exec "echo This is $bar $foo" --env-file ./example.env
elk exec "curl $url/health" --env-file ./example.env
```

### var

This flag will set `var` variable in all the commands. You can call this flag multiple times.

Example:
```
elk exec "curl {{.url}}/health" -v url="http://localhost:8080"
elk exec "curl {{.url}}/health" --var url="http://localhost:8080"
```

### delay

This flag will run the commands after some duration.

This flag supports the following duration units:
- `ns`: Nanoseconds
- `ms`: Milliseconds
- `s`: Seconds
- `m`: Minutes
- `h`: Hours

Example:

```
elk exec "echo Hello world" --delay 500ms
```

### dir

This flag specify the `directory` where the command is going to run. Be default the commands run in the current 
directory.

Example:

```
elk exec "touch example.txt" --dir /home/developer/Desktop
```

### log

This saves the output to a file.

Example:

```
elk exec "echo Hello world" -l ./test.log
elk exec "echo Hello world" --log ./test.log
```

### ignore-error

Ignore errors that happened during the commands.

Example:

```
elk exec "exit 1" "echo Hello World" --ignore-error
```

### timeout

This flag with kill the commands after some duration since the program was started.

This flag supports the following duration units:
- `ns`: Nanoseconds
- `ms`: Milliseconds
- `s`: Seconds
- `m`: Minutes
- `h`: Hours

Example:

```
elk exec "echo Hello world" -t 500ms
elk exec "echo Hello world" --timeout 500ms
```

### deadline

This flag with kill the commands at a particular datetime.

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
elk exec "echo Hello world" --deadline 09:41AM
elk exec "echo Hello world" --deadline 2007-01-09T09:41:00Z00:00
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
elk exec "echo Hello world" --start 09:41AM
elk exec "echo Hello world" --start 2007-01-09T09:41:00Z00:00
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
elk exec "echo Hello world" -i 2s
elk exec "echo Hello world" --interval 2s
```