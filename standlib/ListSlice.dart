// ListSlice is the emulator for go slices in dart.
class ListSlice {
  List source;

  ListSlice(this.source);

  ListSlice slice(int subStart, int subEnd) {
    ListSlice sliced;
    sliced = ListSlice(source.sublist(subStart, subEnd));
    return sliced;
  }

  dynamic elementAt(int index) {
    return source.elementAt(index);
  }

  void setAt(int index, val) {
    source[index] = val;
  }

  String toString() {
    String result = "";
    for (int i = 0; i < source.length; i++) {
      result += source[i].toString();
    }
    return result;
  }

  void add(element) {
    source.add(element);
  }

  void copy(src) {
    this.source = List.from(src);
  }
}
