# `logs`

Attach logs from a task to the terminal. 

## Syntax
```
elk logs [task] [flags]
```

This command takes only one argument which is the name of the task that you want to attach the logs to the terminal. Optionally you can follow the args with the logs flags.

If the `task` do not have a `log` property it will throw an error.

## Examples

```
elk logs foo
elk logs foo -f ./elk.yml
elk logs foo --file ./elk.yml
elk logs foo -g
elk logs foo --global
```

## Flags
| Flag                                  | Short code | Description                                       | 
| -------                               | ------     | -------                                           | 
| [file](#file)                         | f          | Specify which file to use                         |
| [global](#global)                     | g          | Use global file                                   |

### file

This flag force `elk` to use a particular file.

Example:
```
elk logs test -f ./elk.yml
elk logs test --file ./elk.yml
```

### global

This force the task to run from the global file either declared at `ELK_FILE` or the default global path `~/elk.yml`.

Example:

```
elk logs test -g
elk logs test --global
```
