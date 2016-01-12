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

## Supported

* Generated code looks like it works :thumbsup:
* Very preliminary dependency resolution
    * Recursive transpilation.
    * Outputs empty files for most Go builtin packages, due to transpilation failures.
* 

## Unsupported

* goroutines
* defer
* recover
* Various control statements
* First class functions
* Function pointers
* no support for go standard libraries

## Status

* Untested
* 
