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

* goroutines
* defer
* recover
* Various control statements

* Function pointers
* no support for go standard libraries
* Cannot input functions as arguments
* Labels 
* Pointers (StarExpr)


## Not todo
* goroutines

## In progress

* Working on dependency resolution
* Alex is learning Dart
* Starting real testing.

* Untested
* 


#TODO
Figure out how to deal with conditional compilation commands.
Fix the millions of nil dereferences.
Probably should do major refactor for posterity.
There are issues with double for loop, for range and switch
