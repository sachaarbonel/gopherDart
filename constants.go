package main

import "go/token"

// this is probably wrong...
var typesMap = map[string]string{
	"string":  "String",
	"float64": "double",
	"float32": "double",
	"int":     "int",
	"int16":   "int",
	"int32":   "int",
	"int64":   "int",
	"uint":    "int",
	"uint16":  "int",
	"uint32":  "int",
	"uint64":  "int",
	"byte":    "int",
	"rune":    "int",
	"nil":     "null",
}

var tokenMap = map[token.Token]string{
	token.ADD:        "+",
	token.ADD_ASSIGN: "+=",
	token.AND:        "&",
	token.AND_ASSIGN: "&=",
	token.AND_NOT:    "&!",
	token.ARROW:      "->",
	token.ASSIGN:     "=",
	token.BREAK:      "break",
	token.CASE:       "case",
	token.CHAN:       "chan",
	token.CHAR:       "String",
	token.COLON:      ":",
	token.CONTINUE:   "continue",
	token.DEC:        "--",
	token.DEFINE:     "=",
	token.EQL:        "==",
	token.GEQ:        ">=",
	token.GTR:        ">",
	token.INC:        "++",
	token.LAND:       "&&",
	token.LEQ:        "<=",
	token.LSS:        "<",
	token.LOR:        "||",
	token.MUL:        "*",
	token.NEQ:        "!=",
	token.NOT:        "!",
	token.OR:         "|",
	token.SUB:        "-",
	token.SUB_ASSIGN: "-=",
}

var sliceHeader = `// ListSlice is the emulator for go slices in dart.
class ListSlice {
List source;
int start = 0;
int end = 0;
int length = 0;

ListSlice([List this.source, this.start = 0, this.end = 0]) {
if (this.source == null) {
  this.source = new List();
}

if (start == end && source.length > 0) {
  end = source.length;
}
length = end - start;
}

ListSlice slice(int subStart, int subEnd) {
return new ListSlice(source, start + subStart, start + subEnd);
}

dynamic elementAt(int index) {
int sourceIndex = start + index;
if (sourceIndex < start || sourceIndex >= end || sourceIndex < 0 || sourceIndex >= source.length) {
  return null;
}
return source[sourceIndex];
}

setAt(int index, val) {
int sourceIndex = start + index;
if (sourceIndex < start || sourceIndex >= end || sourceIndex < 0 || sourceIndex >= source.length) {
  return; // throw?
}
source[sourceIndex] = val;
}

String toString() {
String result = "";
for (int i = start; i < end; i++) {
  result += source[i].toString();
}
return result;
}

void add(element) {
source.add(element);
end++;
length++;
}

void copy(src) {
this.source = new List.from(src);
this.start = 0;
this.end = this.source.length;
}
}

`

var getTypeName = `getTypeName(dynamic obj) {
  return reflect(obj).type.reflectedType.toString();
}`
