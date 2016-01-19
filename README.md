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

* Function pointers
* no support for go standard libraries
* Cannot input functions as arguments
* Labels 

## Status

* Untested
* 


#TODO
Figure out how to deal with conditional compilation commands.
Fix the millions of nil dereferences.

Probably should do major refactor for posterity.