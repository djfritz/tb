# Tagebuch

Tagebuch (`tb`) is a simple command-line tool for organizing daily journals ("Tagebuch" is German for "diary") and todo lists, written in Go.

It organizes notes by day in a directory structure (`year/month/day`), supports multiple journals, and optionally synchronizes across machines via Git.

## Installation

```bash
go install github.com/djfritz/tb@latest
```

## Quick Start

```bash
# Initialize a new journal called "work"
tb work init

# Edit today's entry (opens $EDITOR)
tb work edit today

# Print today's entry
tb work print today

# Add a todo item
tb work todo add "Review pull requests"

# List all todos
tb work todo
```

## Commands

Commands support prefix matching (e.g., `tb work e tod` expands to `tb work edit today`). Invalid commands display help at the current level.

```
tb <journal>
    init                    Initialize a new journal
    edit
        today               Edit today's entry
        yesterday           Edit yesterday's entry
        tomorrow            Edit tomorrow's entry
        <year/month/day>    Edit a specific date (e.g., 2026/1/6)
    print
        today               Print today's entry (also lists attached files)
        yesterday           Print yesterday's entry
        tomorrow            Print tomorrow's entry
        <year/month/day>    Print a specific date
    todo                    List all todo items
        add <text>          Add a todo item
        complete <number>   Complete a todo item by its number
    search <term>           Search entries using grep-style regular expressions
    list                    List all days with entries (for scripting)
    sync                    Manually sync with git remote (pull then push)
    alias                   List all aliases
        add <name> <date>   Create an alias to a date (e.g., alias add "great thoughts" 2026/1/6)
        remove <name>       Remove an alias by name
    files
        add <date> <path>   Copy a file into a day's directory
        list <date>         List files attached to a day
        remove <date> <name> Remove a file from a day
        copy <date> <name> <dest> Copy a file out to a destination path
    calendar                Show this month's calendar (days with entries highlighted in green, files marked with *)
        last                Show last month's calendar
        next                Show next month's calendar
        <year/month>        Show a specific month (e.g., 2026/1)
    serve <host:port>       Start a web server for managing todos (e.g., serve localhost:8080)
```

## Configuration

### Base Directory

By default, journals are stored in `~/.tb/`. Override with the `-b` flag:

```bash
tb -b /path/to/journals work edit today
```

### Editor

`tb` uses the `$EDITOR` environment variable to open entries for editing. Ensure this is set:

```bash
export EDITOR=vim  # or emacs, etc.
```

## Multiple Journals

`tb` supports working with multiple journals, each in a separate directory under the base path:

```bash
tb personal init
tb work init
tb ideas init
```

Each journal maintains its own entries and todo list.

## Git Synchronization

Enable Git sync to automatically pull before reading and push after writing entries.

1. Initialize a Git repository in your journal directory:
   ```bash
   cd ~/.tb/work
   git init
   git remote add origin <your-remote-url>
   ```

2. Enable sync by adding to the journal's config file (`.tagebuch`):
   ```
   git=true
   ```

When enabled, `tb` will:
- `git pull` before reading entries or todos
- `git add -A && git commit && git push` after writing

Git errors are printed to stderr but don't prevent the operation from completing.

## Directory Structure

```
~/.tb/
└── work/                   # Journal name
    ├── .tagebuch           # Config file (presence marks valid journal)
    ├── todo                # Todo list (one item per line)
    ├── aliases             # Named aliases to dates (name=year/month/day)
    └── 2026/
        └── 1/
            └── 6/
                ├── entry       # Daily entry file
                └── photo.jpg   # Attached files
```
