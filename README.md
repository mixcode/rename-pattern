
# rename-pattern

a file rename tool with simple \*? style matchmatch or regexp match.

## install
```
go install github.com/mixcode/rename-pattern@latest
```


## Example

* substitute all 'hello' to 'world' in ZIP files
```
# '*' and '?' matches filenames
# the '-d' flag rename files
rename-pattern -d hello world *.zip
```

* substitute file\_1.zip to name\_001.zip
```
# ':' matches digits
rename-pattern -d file_: name_%03d *.zip
```

* substitute file\_1.zip to 001\_name.zip
```
# $POS or %[POS] referes a match at the position
rename-pattern -d *|_: %[2]03d_$1 *.zip
```

* '-r' flag: use regexp for match
```
rename-pattern -r '(\d+)' %03d *.zip
```

* use '-s' to feed filenames from STDIN
```
ls -1 \*.zip | rename-pattern -d : %05d
```


