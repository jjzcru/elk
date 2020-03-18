## [v0.3.0](https://github.com/jjzcru/elk/tree/v0.3.0) (2020-03-18)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.3.0)

**Commands:**
- **[cron]** Create command

**Documentation:**
- **[cron]** Create documentation
- Specify `units` for `--delay` and `--timeout` flags

**Flags:**
- **--interval** Add flag to **[run]** command
- **--ignore-error** Add flag to **[run]** command


**Syntax:**
- **ignore_error** Add property at `task` level

**Bug Fix:**
- Fix build binary for `windows` in CI
- Fix build binary for `macOS` in CI
- Fix `context` error when running task in `watch` mode

**Misc:**
- Increase test code coverage
- Use the same `os`, via CI, to compile binaries instead of using `go` cross compile feature

## [v0.2.1](https://github.com/jjzcru/elk/tree/v0.2.1) (2020-03-12)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.2.1)

**Bug Fix:**
- Fix issue when context was call before the program finish

## [v0.2.0](https://github.com/jjzcru/elk/tree/v0.2.0) (2020-03-12)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.2.0)

**Misc:**
- Use **ELK_FILE** env variable as a global file
- Change the way `dep` is declared in `elk.yml`

**Commands:**
- **[config]** Remove command
- **[install]** Remove command

**Flags:**
- **--timeout** Add flag to **[run]** command
- **--deadline** Add flag to **[run]** command
- **--start** Add flag to **[run]** command
- **--delay** Add flag to **[run]** command

## [v0.1.0](https://github.com/jjzcru/elk/tree/v0.1.0) (2020-03-04)
[Release](https://github.com/jjzcru/elk/releases/tag/v0.1.0)

**Commands:**
- **[config]** Create command
- **[init]** Create command
- **[install]** Create command
- **[logs]** Create command
- **[ls]** Create command
- **[run]** Create command

**Documentation:**
- **[config]** Create documentation
- **[init]** Create documentation
- **[install]** Create documentation
- **[logs]** Create documentation
- **[ls]** Create documentation
- **[run]** Create documentation
