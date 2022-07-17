# audioslave

[![Tests](https://github.com/bevzzz/audioslave/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/bevzzz/audioslave/actions/workflows/tests.yml)

## Description
Listening to music while working is func. Sometimes, however, while writing code, emails, a message to your colleague, or a piece of documentation, you might wish that it was just a bit quiter.  
Here's where `audioslave` comes into play: it's a command line utility that _controls the output volume_ on your device based on your current typing speed.

## Install
```bash
go install github.com/bevzzz/audioslave@latest
```  

Compiled binary for those who want to use it without installing the `go tool` is coming soon!

## Run
```bash
./audioslave [OPTIONS]
```  
A number of options can be specified to customize your experience :
- **Min volume**: set the volume that should be set at "peak activity"
- **Average CPM** (characters-per-minute): make sure the volume changes relative to your typing speed
- **Window** and **interval**: current typing speed is calculated as an average of the last N observation. How many observations should be taken into account and how often should the typing speed be measured?    

Explore the available options with:
```bash
./audioslave --help
```

## Contribute  
Contributions are welcome!  
If you found a bug or an improvement, feel free to open an issue.  
For submitting PRs, please follow the [open source contribution guide](https://opensource.guide/how-to-contribute/#opening-a-pull-request]).

### Known issues
- The tool does not capture global keystrokes, which means it will only alternate volume when you're focused on the same window in which the terminal is running. (no issue for it yet)
