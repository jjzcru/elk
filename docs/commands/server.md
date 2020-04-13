server
==========

Start a graphql server ⚛️

## Syntax

```
elk server [flags]
```

This commands run `elk` as a `graphql` server, it enables user to run remote commands, either `sync` like the `run` 
command or `async` like `detached` mode. The tasks executed using the server are bound to the server process, meaning
that if the server process gets terminated, the tasks being executed will be terminated as well.

This command do not take any argument. Be default it will try to search for an `ox.yml` in the local directory, 
if not found I will search for the global file as a fallback. The server only keeps the file path of the configuration
in memory and not the actual content, the user can edit the file content on the fly without a need to restart the 
server for changes.

## Examples

```
elk server
elk server -q
elk server -p 9090 -q
elk server -p 9090 -q -d
elk server -g
elk server --global
elk server -q -f ./ox.yml
elk server -q -g
```

## Flags
| Flag                                  | Short code | Description                                          | 
| -------                               | ------     | -------                                              | 
| [detached](#detached)                 | d          | Run the server in detached mode and returns the PGID |
| [port](#port)                         | p          | Port where the server is going to run                |
| [query](#query)                       | q          | Enables graphql playground endpoint 🎮               |
| [file](#file)                         | f          | Specify the file to used                             |
| [global](#global)                     | g          | Use global file                                      |

### detached
Run the server in detached mode and returns the PGID

Example:
```
elk server -d
elk server --detached
```

### port
Specify the port that is going to be used by the server. If not set is going to use the port `8080` by default

Example:
```
elk server -p 3000
elk server --port 3000
```

### query
Enables a [GraphQL Playground][playground] endpoint to test the server and see the [GraphQL Schema][documentation].

Example:
```
elk server -q
elk server --query
```

### file

This flag force `elk` to use a particular file path to fetch the tasks.

Example:
```
elk server -f ./ox.yml
elk server -—file ./ox.yml
```

### global

This force `elk` to fetch the tasks from the `global` file.

Example:

```
elk server -g
elk server —-global
```

[playground]: https://github.com/prisma-labs/graphql-playground
[documentation]: ../../pkg/server/graph/schema.graphqls