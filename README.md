# Go-grep 

This repo implements grep-like behaviour using Golang. To install, make sure `task` is installed and then run 
```task build-cli```

Usage is 
```
go-grep <pattern> <path-to-file-or-directory>
```
with an option flag `--r` to indicate whether or not to recursivley search the path provided. If the path is a file, 
this flag is ignored. 

The pattern can contain the following regex patterns: 
- pattern anchoring at beginning of string with `^` and end of string with `$`
- exact string matching
- matching on digits with `\d`
- matching on alphanumeric characters with `\w`
- matching on character groups, e.g. `[abc]` will match to any character in `a`, `b`, `c`
- matching on negative character groups, e.g. `[^abc]` will match to any character not in `a`, `b`, `c`
- matching one or more times with `+`, e.g. `a+` will match to `a` one or more times
- matching zero or one times with `?`
- matching to a wildcard with `*`
- matching to an alternating group, e.g. `(cat|dog)` will match to either `cat` or `dog`
- backreferencing to previous alternating groups, e.g. `(cat|dog) and \1` will match to `cat and cat` or `dog and dog`