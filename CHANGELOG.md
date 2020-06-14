## [v0.7.3](https://github.com/jjzcru/elk/tree/v0.7.3) (2020-06-14)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.7.3)

**Server ⚛️:**
- Search detached task in graphql by status
- Search detached task in graphql by array of ids
- Ignores task `start` property if the is set in the past
- Search tasks `ids` using regex instead of complete `id`

**Misc 👾:**
- Add builds for `ARM`

## [v0.7.2](https://github.com/jjzcru/elk/tree/v0.7.2) (2020-04-20)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.7.2)

**Bug Fix 🐛:**
- Fix issue in `server` where the the elk file from the request was using the old syntax

## [v0.7.1](https://github.com/jjzcru/elk/tree/v0.7.1) (2020-04-19)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.7.1)

**Documentation 📖:**
- Instructions on how to install `elk` from terminal

**Misc 👾:**
- Installation from terminal with `Go Binaries`
- Move `main.go` to root directory

## [v0.7.0](https://github.com/jjzcru/elk/tree/v0.7.0) (2020-04-19)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.7.0)

**Server ⚛️:**
- **[health]** Add health check endpoint

**Commands 🤖:**
- **[server]** Enable authorization

**Documentation 📖:**
- **[deps]** Update documentation to enable `ignore_error` at `deps` level
- **[task]** Update documentation to enable `title` and `tags` at `task` level

**Flags 🚩:**
- **--auth** Add flag to **[server]** command
- **--token** Add flag to **[server]** command

**Syntax:**
- **task** Add property `title` and `tags`
- **deps** Add property `ignore_error`

**Misc 👾:**
- Integration with `Travis-CI`
- Integration with `Coverall`
- Integration with `Go Release`
- Add support for `golangci-lint`
- Add `go vet` to build pipeline

## [v0.6.0](https://github.com/jjzcru/elk/tree/v0.6.0) (2020-04-12)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.6.0)

**Commands 🤖:**
- **[server]** Create command

**Documentation 📖:**
- **[server]** Create documentation

**Bug Fix 🐛:**
- Fix issue where `logger` where being created to all the tasks instead of just the one that were being executed
  
## [v0.5.0](https://github.com/jjzcru/elk/tree/v0.5.0) (2020-03-31)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.5.0)

**Commands 🤖:**
- **[logs]** Able to display logs for multiple `tasks` in the same command

**Documentation 📖:**
- **[logs]** Add documentation for the updated `log` property

**Flags 🚩:**
- **--ignore-dep** Add flag to **[run]** command
- **--ignore-dep** Add flag to **[cron]** command
- **--ignore-log-format** Add flag to **[run]** command
- **--ignore-log-format** Add flag to **[cron]** command

**Syntax:**
- **log** Add properties `out`, `error` and `format` to log object

**Bug Fix 🐛:**
- Fix issue where `version` was not being display on `macOS amd64`

**Misc 👾:**
- `engine` now uses a `logger` per task, instead of one for the entire `engine`
- If you display multiple tasks with the command `logs` it display which output belong to which task
- Add which go version was used to built the binary in the `version` command
  
## [v0.4.0](https://github.com/jjzcru/elk/tree/v0.4.0) (2020-03-22)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.4.0)

**Commands 🤖:**
- **[exec]** Create command
- **[version]** Add built information 

**Documentation 📖:**
- **[exec]** Create documentation

**Flags 🚩:**
- **--var** Add flag to **[run]** command
- **--var** Add flag to **[cron]** command

**Syntax:**
- **vars** Add property at `task` level
- **vars** Add property at `global` level

## [v0.3.1](https://github.com/jjzcru/elk/tree/v0.3.1) (2020-03-19)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.3.1)

**Commands 🤖:**
- Remove the display of `help` when a `run`, `cron` and `ls` throw an error
- Remove `Examples` section from `--help`

**Documentation 📖:**
- Improve documentation in `README.md`
- Add `Syntax` section
- Add `Use Cases` section
- Add `Commands` section

**Flags 🚩:**
- Rename the flag `--ignore-log` to `--ignore-log-file`

**Syntax:**
- Rename the property `watch` to `sources`

**Misc 👾:**
- Rename the file `elk.yml` to `ox.yml`

## [v0.3.0](https://github.com/jjzcru/elk/tree/v0.3.0) (2020-03-18)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.3.0)

**Commands 🤖:**
- **[cron]** Create command

**Documentation 📖:**
- **[cron]** Create documentation
- Specify `units` for `--delay` and `--timeout` flags

**Flags 🚩:**
- **--interval** Add flag to **[run]** command
- **--ignore-error** Add flag to **[run]** command

**Syntax:**
- **ignore_error** Add property at `task` level

**Bug Fix 🐛:**
- Fix build binary for `windows` in CI
- Fix build binary for `macOS` in CI
- Fix `context` error when running task in `watch` mode

**Misc 👾:**
- Increase test code coverage
- Use the same `os`, via CI, to compile binaries instead of using `go` cross compile feature

## [v0.2.1](https://github.com/jjzcru/elk/tree/v0.2.1) (2020-03-12)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.2.1)

**Bug Fix 🐛:**
- Fix issue when context was call before the program finish

## [v0.2.0](https://github.com/jjzcru/elk/tree/v0.2.0) (2020-03-12)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.2.0)

**Misc 👾:**
- Use **ELK_FILE** env variable as a global file
- Change the way `dep` is declared in `elk.yml`

**Commands 🤖:**
- **[config]** Remove command
- **[install]** Remove command

**Flags 🚩:**
- **--timeout** Add flag to **[run]** command
- **--deadline** Add flag to **[run]** command
- **--start** Add flag to **[run]** command
- **--delay** Add flag to **[run]** command

## [v0.1.0](https://github.com/jjzcru/elk/tree/v0.1.0) (2020-03-04)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.1.0)

**Commands 🤖:**
- **[config]** Create command
- **[init]** Create command
- **[install]** Create command
- **[logs]** Create command
- **[ls]** Create command
- **[run]** Create command

**Documentation 📖:**
- **[config]** Create documentation
- **[init]** Create documentation
- **[install]** Create documentation
- **[logs]** Create documentation
- **[ls]** Create documentation
- **[run]** Create documentation
