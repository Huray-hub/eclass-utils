# Assignments

Assignments is a utility that will allow you to seamlessly be up-to-date with the state of your assignments.

## Features 

Prints the assignments of your enrolled courses to the terminal. These are the following options available for convenience:

- **Inclusion of expired assignments**: usually there are expired assignments from previous years.
(default = false)

- **Include only favorite courses**: fetch only courses marked as 'favorites'. 
(default = false)

- **Exclusion of selected courses**: there are courses that do not have assignments.
(default = empty)

- **Exclusion of selected assignments**: professors tend to divide assignments in a non-common pattern like per lab classes.
(default = empty)

- **Export an ICS file**: produces a calendar file that can be imported from any calendar app. See [here](https://support.google.com/calendar/answer/37118?hl=en&co=GENIE.Platform%3DDesktop). 
(default = false)

- **Plain text**: The output will be printed in csv format instead of a table.
(default = false)


## Installation

1. See [releases](https://github.com/Huray-hub/eclass-utils/releases) for pre-built binaries.
2. Install as go package (requires [Go](https://go.dev))
```sh
go install github.com/Huray-hub/eclass-utils/assignment/cmd/assignments@latest
```
2. Build from source requires [Go](https://go.dev).
- clone repo
- cd into the repo 
- run
```sh
make install-assignments-local
# or 
make install-assignments-gopath
```

## Configuration
### Location
After the installation, a yaml file will be stored in one of the following locations:
- Unix
    - XDG_CONFIG_HOME/.config/eclass-utils/config.yaml
    - ~/.config/eclass-utils/config.yaml
- Windows
    - %AppData%/eclass-utils/config.yaml
- MacOS
    - no idea I'm poor

### File structure
Configuration file looks like this: 
```yaml
credentials:
    username: your-eclass-username
    password: your-eclass-password-ENCRYPTED
options:
    plainText: false # print output in csv format instead of pretty table
    includeExpired: false # include expired assignments
    exportICS: false # export calendar file (for import to calendar apps) 
    excludedCourses: # the following course codes are excluded
        CS179: 
        ICE257: 
        ICE290: 
    excludedAssignments:
        # for course code ICE245, assignments containing these strings are excluded
        ICE245: 
            - ΗΛΕΚ01
            - ΗΛΕΚ02
            - ΗΛΕΚ03
            - ΗΛΕΚ05
        CS157:
            - ΓΙΑ ΟΣΟΥΣ ΔΕΝ ΑΝΗΚΟΥΝ ΣΕ ΚΑΠΟΙΟ ΤΜΗΜΑ
        ICE325:
            - '15η εργασία : Εκπρόθεσμες υποβολές ασκήσεων.'
            - 6. Απαλλακτική εργασία
        ICE326:
            - Άσκηση (project) Εργαστηρίου - 2023
    baseDomain: eclass.uniwa.gr # domain of your university
    onlyFavoriteCourses: false # include only favorite courses
secretKey: generated-secret-key # always ignore this one
```


## Command-line-flags
- Use -h for help

*Note that in the runtime, cmd-flags override the config file (will not modify the file)*

### Authentication

College password can be cached, **encrypted** in the config file. It can't be update it manually from config file. 

Use -h for help
