// ListSlice is the emulator for go slices in dart.
class ListSlice {
List source;
int start = 0;
int end = 0;
int length = 0;

    ListSlice([int this.length=-1, List this.source, this.start = 0, this.end = 0]) {
        if (this.source == null) {
          this.source = new List();
        }


        if (start == end && source.length > 0) {
          end = source.length;
        }
        if (length == -1) {
            length = end - start;
        }
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