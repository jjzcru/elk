`ls`
==========

List tasks

## Syntax

```
elk ls [flags]
```

This command do not take any argument. Be default it will try to search for an `elk.yml` in the local directory, 
if not found I will search for the global file as a fallback.

## Examples

```
elk ls
elk ls -f ./elk.yml
elk ls --file ./elk.yml
elk ls -g
elk ls --global
elk ls -a
elk ls --all
elk ls -a -f ./elk.yml
elk ls -a -g
```

## Flags
| Flag                                  | Short code | Description                                       | 
| -------                               | ------     | -------                                           | 
| [all](#all)                           | a          | Display all the properties from a task            |
| [file](#file)                         | f          | Specify which file to use                         |
| [global](#global)                     | g          | Use global file                                   |

### all
Display all the columns

Example:
```
elk ls -a
elk ls —-all
```

### file

This flag force `elk` to use a particular file path to fetch the tasks.

Example:
```
elk ls -f ./elk.yml
elk ls -—file ./elk.yml
```

### global

This force `elk` to fetch the tasks from the `global` file.

Example:

```
elk ls -g
elk ls —-global
```