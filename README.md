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

<<<<<<< HEAD


* goroutines
* defer
* recover
* Various control statements

* Function pointers
* no support for go standard libraries
* Cannot input functions as arguments
* Labels 
* Pointers (StarExpr)
=======
* defer
* recover
* Various control statements
* support for go standard libraries
* Fix assignment statements not being concurrent
```a,b := b,a``` is not transpiled properly
* We lose the length of slices when we transpile to ListSlice

## Maybe todo
* First class functions
* Function pointers
>>>>>>> origin/master

## Not todo
* goroutines

## In progress

* Working on dependency resolution
* Alex is learning Dart
* Starting real testing.

<<<<<<< HEAD
* Untested
* 


#TODO
Figure out how to deal with conditional compilation commands.
Fix the millions of nil dereferences.

Probably should do major refactor for posterity.
=======
>>>>>>> origin/master
