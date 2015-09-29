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

For now it generates a file 'lib.dart' that contains the entire Go package in a single dart library.

In Progress

  1. Go slice support via dart class 'ListSlice' to emulate go slice behavior.
  2. Working on making Go interfaces into dart abstract classes.

Future TODO:

  1. Read in imports and parse them
  2. Decide what standard packages should be hand made in dart vs translated.
