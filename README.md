# gopherDart
Go to Dart transpiler

To build
```
go build
```

To Run
```
./gopherDart /path/to/go/package
```

##Goal
Take our go, serverside, code and transpile it into dart that can be run client side. 
Basically GopherJs, except with Dart.



## Supported

* Generated code looks like it works :thumbsup:
* Very preliminary dependency resolution
    * Recursive transpilation.
    * Outputs empty files for most Go builtin packages, due to transpilation failures.
* 

## Todo

* defer
* recover
* Various control statements
* support for go standard libraries

## Maybe todo
* First class functions
* Function pointers

## Unsupported (and not todo) 
* goroutines

## In progress

* Working on dependency resolution
* Alex is learning Dart
* Starting real testing.

