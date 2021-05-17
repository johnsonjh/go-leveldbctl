# jleveldbctl

JLevelDB control command.

This command provides easy way to CRUD operation on LevelDB.

```sh
$ jleveldbctl put foo bar
put foo: bar into ./.
$ jleveldbctl get foo
bar
```

## Features

* Initialize JLevelDB
* Put key/value into JLevelDB
* Get value with key
* Delete key
* Dump all key/values in JLevelDB
* Print all keys

## Install

```sh
$ go install -a -v github.com/johnsonjh/jleveldbctl/cmd/jleveldbctl
```

## Usage

```sh
jleveldbctl [global options] command [command options] [arguments...]

COMMANDS:
     init, i    Initialize a JLevelDB
     walk, w    Walk in a JLevelDB
     keys, k    Search all keys in a JLevelDB
     put, p     Put a value into a JLevelDB
     get, g     Gut a value from a JLevelDB
     delete, d  Delete a value from a JLevelDB
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dbdir value, -d value  JLevelDB Directory (default: "./") [$JLEVELDB_DIR]
   --hexkey, --xk           get / put hexadecimal keys
   --hexvalue, --xv         get / put hexadecimal values
   --help, -h               show help
   --version, -v            print the version
```

For hexadecimal keys and values:

```sh
$ export JLEVELDB_DIR=${HOME}/.bitcoin/index
$ leveldbctl -xk g 62f2a1f90489f1f74e441f325ec6f532df8286847d7c7a14000000000000000000|xxd -p
89fe04a3db1d801d92188d350880fec55300008020edd4d15faba7c63dd7
c83961bf6783a691fb8f5f6887120000000000000000009f413c1df7e296
4af9babb54e46d4414eaad550b27b409e29ab80a832ac64ce9966ab95ddf
8e1417baf3db320a
```

Search for a key / value using a key prefix:

```sh
$ leveldbctl --dbdir=testdb s foo
foo: bar
foo2: bar2
foo3: bar3
```
