[![GitHub Workflow Status](https://github.com/bdreece/herobrian/actions/workflows/build.yml/badge.svg)](https://github.com/bdreece/herobrian/actions/workflows/build.yml)

# herobrian

A Minecraft server management platform.

## Getting Started

### Dependencies

Building `herobrian` requires the following build tools:

- go (1.26.1)
- node (^24.14.0 || >=25.0.0)
- pnpm (10.28.2)

### Building

This program is built using [Taskfile](https://taskfile.dev). Run the following
command to list the available tasks:

```
$ go tool task --list
task: Available tasks for this project:
* build:              Build application                        (aliases: b)
* clean:              Clean all build artifacts                (aliases: c)
* restore:            Restore all dependencies                 (aliases: r)
* test:               Run test suites                          (aliases: t)
* watch:              Launch application with live reload      (aliases: w)
* go:build:           Build Go app                             (aliases: go:b)
* go:clean:           Clean Go build artifacts                 (aliases: go:c)
* go:generate:        Generate Go sources                      (aliases: go:g)
* go:restore:         Restore Go dependencies                  (aliases: go:r)
* go:test:            Run Go unit tests                        (aliases: go:t)
* go:watch:           Watch Go app with Air                    (aliases: go:w)
* vite:build:         Build Vue application                    (aliases: vite:b)
* vite:clean:         Remove Vite build artifacts              (aliases: vite:c)
* vite:restore:       Restore Vite dependencies                (aliases: vite:r)
* vite:test:          Run Vitest testing suite                 (aliases: vite:t)
* vite:watch:         Launch Vite dev server                   (aliases: vite:w)
```
