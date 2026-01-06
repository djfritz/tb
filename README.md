# Tagebuch

Tagebuch (tb) is a simple wrapper for organizing a daily journal ("Tagebuch") and todo list, written in Go. 

It organizes notes by day, supports multiple journals, and supports synchronizing with multiple installations via `git`. 

## Commands

`tb` organized as follows. A malformed command will present help at the current command level. 

```
tb <name of journal>
    init                ( initialize a new journal)
    edit
        yesterday       ( edit yesterday's entry )
        today           ( edit today's entry )
        tomorrow        ( edit tomorrow's entry )
        year/month/day  ( edit an entry given by year/month/day )
    print
        yesterday       ( print yesterday's entry )
        today           ( print today's entry )
        tomorrow        ( print tomorrow's entry )
        year/month/day  ( print an entry given by year/month/day )
    todo                ( print todo items )
        add             ( add a todo list item )
        complete        ( remove a todo list item by number )
    search <term>       ( print matching lines and entries to the given search term. Supports `grep` styled regular expressions )
```

## Installing

`tb` is still a work in progress. For now, `go install github.com/djfritz/tb@latest` is the best way to install and use `tb`.

## Todo lists

Each journal keeps a simple todo list (see above commands). Adding new items with the same content will be ignored (deduplicated).

## Multiple journals



## Git synchronization

## My usage as an example

## Planned features

- sync on demand instead of on every invocation
- aliases for named documents
- support for adding files to days 
