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

`tasks`
In here you have a list of all the tasks that you wish to perform. The name of the task is going to be used to know 
which task is going to perform.

### Task
In the `task` level you can overwrite the values set at the `global` level to this particular task.

`env_file`
This is a path to a file that declares the `env` variables as `ENV_NAME=ENV_VALUE` where each line is a different 
`env` variable. This overwrites the existing `env` variable already declared on global.

`env`
In here you declare all the `env` variable that you wish that all the task inherit this property overwrites the 
existing `env` variables, also the ones declared in the `env_file` property and global.

`description`
In here you describe what is the purpose of the task, this is also display by the `ls` command.

`dir`
This specifies what is the directory in which the commands are going to run. If not set is going to use the current 
directory.

`log`
Saves the output of the command to a file.

`ignore_error`
Ignore errors that happened during a `task`.

`sources` 
This is a regex for the files that are going to activate the re-run of the tasks

`deps`
This is a list of all the dependencies that the task requires to run. The `dep` declaration takes 2 properties:

- `name`: It takes a `string` which is the name of the task that you which to run as a dependency.

- `detached`: It takes a `boolean` which tells if the dependency should run in `detached` mode, is `false` as default.

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

## Use Cases
- [Typescript](#typescript)
- [Back-end](#back-end)
- [CI/CD](#ci/cd)
- [Automation](#automation)

### Typescript

**Example File:** [Typescript Example][typescript-example]

If you are a `typescript` developer you need to compile your project and start it over and over again. You can use
`elk` to automate this process.

With this file we have three tasks `build`, `serve` and `health`. With this document we are setting the property 
`sources` in the `serve` task and also we are telling that `serve` depends on the `build` task. We could run `serve`
in `watch` mode by running the command `elk run serve -w` and now everytime we update a `.ts` file the program will
recompile the project and execute it.

There is no need on building and running manually or the need to create a custom script to do that, you just need to 
declare the behavior and `elk` will take care of the rest.

Now that you have your service running lets imagine that you want to health check to make sure that the service is up
and running, you could do a `curl http://localhost:8080/health` to make sure that this is happening, but how about if 
we don't want to do that, and we just want to see if the service is running and no need to run the command anymore.

You could use the `interval` flag in the `run` command to achieve this, just run `elk run health -i 2s` and now we
are going to health check the service each 2 seconds.

### Back-end

**Example File:** [Back-end Example][back-end-example]

In this example we have two different projects in [NodeJS](https://nodejs.org), `service_1` and `service_2`, each with 
their own dependencies and env variables, instead of using the same `env` as the system, each `task` runs on it's own
so we can use the `env` property to specify the `PORT` on which we want them to run, no need to update the env variable
in a `.zshrc` or `.bashrc`, or setting the `env` variable on the terminal.

Now lets say that we want that every single time our application gets saves and we don't want to have a terminal open
because we are working on the two services at the same time, we can run both services in `detached` and `watch` mode
and save the output of those services to a file with the property `log`, we could run the command 
`elk run service_1 service_2 -d -w`.

Now let make it a little bit more complicated, how about if we want to know if our services are alive or not. We could 
run a health check every seconds to check on the services and print them green if they are alive and red if they are 
dead. 

We can run `elk run health -i 1s`, this command will clear the terminal and output the state of the services.

### CI/CD

**Example File:** [CI/CD Example][ci-cd-example]

We can use `elk` as a `CI/CD` build system too. Let's say we are doing a self host `CI/CD` pipeline with
[Jenkins](https://jenkins.io/). We already set up our `jobs` inside `Jenkins` and now we are configuring our `Build`
step using a command.

We have two projects that we need to deploy one is `service_1` and the other is `service_2`. To deploy this service
we first need to make sure that all the test are working, then we need to build the application and then we need to
deploy it. The deployment for this applications is moving the project from `/home/example/ci` to `/home/example/deploy`
and run a script with [pm2](https://pm2.io/).

For the `service_1` by using the `deps` property we just need to run a single command `elk run service_1_deploy` to 
deploy our application. But what is happening under the hood?.

`service_1_deploy` has a dependency on `service_1_build` that at the same time it has a dependency on `service_1_test`
so the first step, the application will run the task `service_1_test` if the test are successfull, it will run 
`service_1_build` which will compile the application and then move the content of the build to `/home/example/deploy`
now `service_1_deploy` gets executed and this will run `app.js` with `pm2` and saves the current `pm2` configuration.

Now for the second service, `service_2_deploy`, we have two direct dependencies, first we are going to run 
`service_2_test` and if everything works fine it will execute `service_2_build` and after that it will run 
`service_2_deploy`.

In both cases if a `dependency` fail the program will stop and `Jenkins` will check the job as failed, instead of 
running complex `.sh` or `.bat` files, you just need to declare how would you like the application to run. In this
examples the process for deploying `service_1` and `service_2` is practically the same, for `service_1` we have three
levels depth of dependency and `service_2` has two levels, but with two dependencies at the same level. This allow us
to create either simple or complex dependency tree dependending on our use case.

### Automation

**Example File:** [Automation Example][automation-example]

You don't need to use `elk` only for your job, you can also use it to automate tasks in the real world, with the rise
of `IoT` devices, let's imagine that you create an `http` server that talk with your `IoT` devices and has an endpoint 
that takes a query param called `command` which receives a text with the command that you want to execute.

Now lets imagine that you are starting your day at `9:00AM` in the morning, and you now that you are a workaholic so
you want to make sure to turn down your computer when you end you working shift, let's say it finish at `5:00PM`, you 
could use the flag `start` to delay the execution of some commands to a particular time and also run it as `detached`
so we don't have a terminal hanging out with the process and we kill it by accident

To acomplish our task to shutdown our machine we could run `elk run shutdown --start 5:00PM -d`. Going on with the day
you know that you want to eat at `1:00PM` so we want to set up an alarm that reminds us that is time to out, but we
are feeling lazy for cooking so we will go out to buy some food and we are going to leave at `1:10PM` and probably by 
back at `1:30PM`. We can program all of that bu running:

- `elk run alarm --start 1:00PM`: To setup the alarm that is going to remind us that is time to eat.
- `elk run open_the_door 1:10PM`: So the garage door get open while we are getting ready.
- `elk run close_the_door 1:11PM`: To close the garage door once we leave.
- `elk run open_the_door 1:30PM`: So the garage door get open while we are heading back.
- `elk run close_the_door 1:31PM`: To close the garage door when we enter.

We can program all that in the morning and just going on with our day. If you never turn down your computer you can
use the command `cron` to creates more complex automate scenarios.

[typescript-example]: ./examples/typescript.yml
[back-end-example]: ./examples/back-end.yml
[cra-example]: ./examples/create-react-app.yml
[ci-cd-example]: ./examples/ci_cd.yml
[automation-example]: ./examples/automation.yml