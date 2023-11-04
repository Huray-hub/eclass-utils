# Assignments

Assignments is a utility that will allow you seamlessly be up-to-date with the state of your assignments.

## Features 

Prints the assignments of your enrolled courses to the terminal (stdin). There
are the following options for convenience:

- **Inclusion of expired assignments**: usually there are expired assignments from previous years.
(default = false)

- Include only favortie courses: fetch only courses marked as 'favorites'. 
(default = false)

- **Exclusion of selected courses**: there are courses that do not have assignments.
(default = empty)

- **Exclusion of selected assignments**: professors tend to divide assignments in a non-common pattern like per lab classes.
(default = empty)

- **Export an ICS file**: produces a calendar file that can be imported from any calendar app. See [here](https://support.google.com/calendar/answer/37118?hl=en&co=GENIE.Platform%3DDesktop). 
(default = false)

- **Plain text**: The output will be printed in csv format instead of a table.
(default = false)

- **Manual add assignments**: Some professors put the assignments on other sections/platforms or nowhere at all. (TODO)
(default = empty)

## Installation Options

1. See releases for pre-built binaries.
2. TODO: Build from source requires [Go](https://go.dev).
3. TODO: Package managers

## Configuration

Two ways: through **config file** and **cmd-flags**

*Note that in the runtime, cmd-flags override the config file (will not modify the file)*

### Config File
After the installation, a yaml file will be stored in one of the following locations:

- Unix
    - XDG_CONFIG_HOME/.config/eclass-utils/config.yaml
    - ~/.config/eclass-utils/config.yaml
- Windows
    - %AppData%/eclass-utils/config.yaml
- Mac
    - no idea I'm poor


### Command-line-flags
- Use -h for help
TODO


## Disclaimers
If you choose to cache college credentials during the installation, please make sure to not give read access to the config file. Currently, creds are be stored there. There was not much time to deal with OS secret API storages. In the future, probably will either store the credentials encrypted to the file or provide support for the secret storage. 
