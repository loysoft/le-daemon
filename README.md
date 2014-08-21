le-daemon
=========

Command line tool to dump stderr and stdout of executed command to [Logentries](http://logentries.com)


Installation
------------

```shell
go get github.com/loysoft/le-daemon
go install github.com/loysoft/le-daemon
```

Usage
-----
Add a new manual TCP token log at [logentries.com](https://logentries.com/quick-start/) and copy the [token](https://logentries.com/doc/input-token/).


Usage examples:

```shell
le-daemon -token <logentries_token> ls /
```
