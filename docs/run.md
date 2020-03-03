# `run`

Fetch a logs of a task to terminal

```
Usage: 
  elk run [task] [flags]

Flags:
    -d, --detached          Run the task in detached mode and return the PGID
    -e, --env               Overwrites declared env variables
    -f, --file              Specify which file to use 
    -g, --global            Force to use the global file
    -l, --log               File path that log the output from the task
    -w, --watch             Enable watch mode
        --ignore-log        Force the output of the task to be on the terminal
```

## Flags
`detached`

This will group all the tasks under the same `PGID` and then it will detach from the process, and returns the `PGID` so the user can kill the process later.

Example:

```
elk run test -d
```

`env`

This flag will overwrite whatever env variable already declared in the file. You can call this flag multiple times.

Example:
```
elk run test -e HELLO=WORLD -e FOO=BAR
```

`log`

This saves the output to a specific file.

Example:

```
elk run test -l ./test.log
```

`watch`

This requires that the task has a property `watch` already setup, otherwise it will throw an error. When this flag is enable it will kill the existing process and create a new one everytime that a file that match the regex is changed.

Example:

```
elk run test -w
```