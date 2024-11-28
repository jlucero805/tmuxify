# Tmuxify

## Getting Started

This follows XDG (Linux) standards.

### Prerequisites

Create the following directory for configurations:

```
~/.config/tmuxify/
```

### Configuring Tmux Sessions

In order to configure sessions, under the tmuxify configuration
directory, create any toml file. For example:

```
# ~/.config/tmuxify/application.toml

# This is the tmux session configuration
Root = "~/path/to/dir"
Name = "session-name"

# These configure tmux windows in order of declaration
[[Win]]
Root = "~/path/to/dir/subdir"

# This will create a tmux window with Nvim opened
[[Win]]
Root = "~/path/to/dir/subdir"
Nvim = true
```



