logs
==========

Attach logs from a task to the terminal. 

## Syntax
```
elk logs [tasks] [flags]
```
This command takes one or more `tasks` as arguments, and attach the `log` content to `stdout`. If the `task` do 
not have a `log` property it will throw an error.

## Examples

```
elk logs foo bar
elk logs foo -f ./ox.yml
elk logs foo bar -f ./ox.yml
elk logs foo --file ./ox.yml
elk logs foo -g
elk logs foo bar -g
elk logs foo --global
```

## Flags
| Flag                                  | Short code | Description                                       | 
| -------                               | ------     | -------                                           | 
| [file](#file)                         | f          | Specify which file to use to get the tasks        |
| [global](#global)                     | g          | Use global file                                   |

### file

This flag force `elk` to use a particular file.

Example:
```
elk logs test -f ./ox.yml
elk logs test --file ./ox.yml
elk logs test bar -f ./ox.yml
elk logs test bar --file ./ox.yml
```

### global

This force the task to run from the global file either declared at `ELK_FILE` or the default global path `~/ox.yml`.

Example:

```
elk logs test -g
elk logs test --global
elk logs test bar -g
elk logs test bar --global
```
