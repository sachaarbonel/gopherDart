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
  3. Finish 'else' cases in if/elseif/else groups.

Future TODO:

  1. Use Streams in place of channels
  2. Read in imports and parse them
  3. Decide what standard packages should be hand made in dart vs translated.
  4. TypedArray vs normal list slice for number types.
  5. Figure out range statements for map and channel types.
  6. Use a future async function for 'go' commands.
