## ghopen

Simple tool to open files in the GitHub viewer. You'll be linked to the
canonical url, with support for line numbers.

### Usage

```
# will open current repository
$ ghopen

# if you are already in a subdirectory (relative to the git root), it will open
# that directory
$ ghopen

# open specific directory
$ ghopen app/

# open file
$ ghopen app/models/file.rb

# open file at line number 10
$ ghopen app/models/file.rb 10
```

### Vim integration

This will open the file and line number when pressing F2:

```
function! GHOpen()
  silent! :call system('ghopen ' . expand('%') . ' ' . line('.'))
endfunction
map <silent><F2> :call GHOpen()<return>
```
