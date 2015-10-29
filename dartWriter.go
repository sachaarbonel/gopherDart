package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"

	"github.com/lologarithm/gopherDart/dart"
)

//  TODO: Make indent part of the context so that we can have 'correct' indentation.
//        One thing is to separate indentation level with how much indentation to draw with any single call.

// LibraryContext lets you pass around values that span the library.
type LibraryContext struct {
	Name        string
	Indentation string
	Class       *ClassContext
	NextIota    int
}

// ClassContext lets you pass around class specific values
type ClassContext struct {
	Name     string
	Recv     string
	RecvDecl interface{}

	Fields map[string]string // Fields maps a field name in the class to its type.
}

// LoadToLibrary accepts an AST of a go file and adds it to the library passed in.
// Due to go file structure we have to first group up functions wih structs before printing.
func LoadToLibrary(f *ast.File, lib *dart.Library) string {
	lib.Name = f.Name.Name
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil {
				lib.Funcs = append(lib.Funcs, d)
			} else {
				for _, rt := range d.Recv.List {
					switch rtt := rt.Type.(type) {
					case *ast.StarExpr:
						id, ok := rtt.X.(*ast.Ident)
						if !ok {
							fmt.Printf("Function %s receiver X incorrect type!?", d.Name.Name)
						}
						if _, ok := lib.Classes[id.Name]; !ok {
							lib.Classes[id.Name] = &dart.Class{
								Name:    id.Name,
								Fields:  []*ast.Field{},
								Methods: []*ast.FuncDecl{},
							}
						}

						lib.Classes[id.Name].Methods = append(lib.Classes[id.Name].Methods, d)
					default:
						fmt.Printf("Func declaration not being handled: %s", reflect.TypeOf(rt.Type))
					}
				}
				//d.Recv.List[0].Type
			}
		case *ast.GenDecl:
			switch d.Tok {
			case token.TYPE:
				for _, s := range d.Specs {
					ts := s.(*ast.TypeSpec)
					switch tsType := ts.Type.(type) {
					case *ast.StructType:
						if _, ok := lib.Classes[ts.Name.Name]; !ok {
							lib.Classes[ts.Name.Name] = &dart.Class{
								Name:    ts.Name.Name,
								Fields:  []*ast.Field{},
								Methods: []*ast.FuncDecl{},
							}
						}
						lib.Classes[ts.Name.Name].Fields = tsType.Fields.List
					case *ast.FuncType:
						lib.FuncTypes = append(lib.FuncTypes, d)
					case *ast.InterfaceType:
						lib.Interfaces = append(lib.Interfaces, d)
					default:
						fmt.Printf("Unknown type lib declaration: %s, %v\n", reflect.TypeOf(ts.Type), ts.Type)
					}
				}
			case token.VAR, token.CONST:
				lib.Vars = append(lib.Vars, d)
			}
		default:
			fmt.Printf("Other declaration in file? %s, %v", reflect.TypeOf(d), d)
		}
	}

	return ""
}

// Print returns the string representation of this library
func Print(lib *dart.Library) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString("library ")
	buf.WriteString(lib.Name)
	buf.WriteString(";")
	buf.WriteString("\n\n")
	buf.WriteString(sliceHeader) // Write the slice class at start of dart file.

	ctx := &LibraryContext{
		Name:        lib.Name,
		Indentation: "",
		Class:       nil,
	}
	for _, v := range lib.Vars {
		printDecl(v, buf, "", ctx)
		buf.WriteString(";\n")
		ctx.NextIota = 0
	}
	buf.WriteString("\n")
	for _, f := range lib.FuncTypes {
		printDecl(f, buf, "", ctx)
		buf.WriteString(";\n")
	}
	buf.WriteString("\n")
	for _, f := range lib.Interfaces {
		printDecl(f, buf, "", ctx)
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
	for _, c := range lib.Classes {
		printClass(c, buf, ctx)
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
	for _, f := range lib.Funcs {
		printFunc(f, buf, "", ctx)
		buf.WriteString("\n")
	}
	return buf.Bytes()
}

func printClass(cl *dart.Class, buf *bytes.Buffer, ctx *LibraryContext) {
	buf.WriteString("class ")
	buf.WriteString(cl.Name)
	buf.WriteString(" {")
	buf.WriteString("\n")
	ctx.Class = &ClassContext{
		Name:   cl.Name,
		Fields: map[string]string{},
	}

	newDefBuf := &bytes.Buffer{} // Buffer a new function up.
	newBodyBuf := &bytes.Buffer{}

	newDefBuf.WriteString("  ") // indent constructor
	newDefBuf.WriteString(cl.Name)
	newDefBuf.WriteString("({")
	indent := "  "
	for idx, f := range cl.Fields {
		printExpr(f.Type, buf, "  ", ctx)
		buf.WriteString(" ")
		// if !f.Names[0].IsExported() {
		// 	buf.WriteString("_")
		// } // TODO: Before we can do this we need to be able to generalize this behavior.
		buf.WriteString(f.Names[0].Name)
		buf.WriteString(";\n")

		newDefBuf.WriteString(f.Names[0].Name)
		newBodyBuf.WriteString(indent + "  ")
		newBodyBuf.WriteString("this.")
		newBodyBuf.WriteString(f.Names[0].Name)
		newBodyBuf.WriteString(" = ")
		newBodyBuf.WriteString(f.Names[0].Name)
		newBodyBuf.WriteString(";\n")
		if idx < len(cl.Fields)-1 {
			newDefBuf.WriteString(", ")
		}
	}
	newDefBuf.WriteString("}) {\n")

	buf.WriteString("\n")
	buf.Write(newDefBuf.Bytes())
	buf.Write(newBodyBuf.Bytes())
	buf.WriteString("  }\n\n") // Close constructor

	for idx, m := range cl.Methods {
		ctx.Class.RecvDecl = m.Recv.List[0]
		printFunc(m, buf, indent, ctx)
		if idx < len(cl.Methods)-1 {
			buf.WriteString("\n")
		}
	}
	buf.WriteString("}\n")
}

func printFunc(f *ast.FuncDecl, buf *bytes.Buffer, indent string, ctx *LibraryContext) {

	if f.Type.Results == nil {
		buf.WriteString(indent)
		buf.WriteString("void ")
	} else if len(f.Type.Results.List) > 1 {
		buf.WriteString(indent)
		// Bundle multiple returns into lists
		buf.WriteString("List ")
	} else {
		printExpr(f.Type.Results.List[0].Type, buf, indent, ctx)
		buf.WriteString(" ")
	}

	if f.Name.Name == "String" {
		// Specially handle the go 'String' function to the dart equivolent.
		buf.WriteString("toString")
	} else {
		buf.WriteString(f.Name.Name)
	}
	buf.WriteString("(")
	printParams(f.Type.Params, buf, "", ctx)
	buf.WriteString(") {\n")

	for _, stmt := range f.Body.List {
		printStmt(stmt, buf, indent+"  ", ctx)
	}
	buf.WriteString(indent)
	buf.WriteString("}\n")
}

func printParams(ps *ast.FieldList, buf *bytes.Buffer, indent string, ctx *LibraryContext) {
	for idx, p := range ps.List {
		for nIdx, n := range p.Names {
			printExpr(p.Type, buf, "", ctx)
			buf.WriteString(" ")
			printExpr(n, buf, "", ctx)
			if nIdx < len(p.Names)-1 {
				buf.WriteString(", ")
			}
		}
		if idx < len(ps.List)-1 {
			buf.WriteString(", ")
		}
	}
}

func printExpr(e ast.Expr, buf *bytes.Buffer, indent string, ctx *LibraryContext) {
	buf.WriteString(indent)
	switch et := e.(type) {
	case *ast.StarExpr:
		//buf.WriteString("*") // TODO: Dart pointer vs struct?
		printExpr(et.X, buf, "", ctx)
	case *ast.Ident:
		if et.Obj != nil && ctx.Class != nil && et.Obj.Decl == ctx.Class.RecvDecl {
			buf.WriteString("this")
		} else if tv, ok := typesMap[et.Name]; ok {
			buf.WriteString(tv)
		} else if et.Name == "iota" {
			buf.WriteString(strconv.Itoa(ctx.NextIota))
			ctx.NextIota++
		} else {
			buf.WriteString(et.Name)
		}
	case *ast.UnaryExpr:
		// For now dont print &
		if et.Op != token.AND {
			buf.WriteString(tokenMap[et.Op])
		}
		printExpr(et.X, buf, "", ctx)
	case *ast.CompositeLit:
		buf.WriteString("new ")
		printExpr(et.Type, buf, "", ctx)
		buf.WriteString("(")
		// TODO: How to deal with construction of []rune from a string
		// If array type we should be adding [] around the parameters being passed in.
		for idx, subExp := range et.Elts {
			printExpr(subExp, buf, "", ctx)
			if idx < len(et.Elts)-1 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(")")
	case *ast.KeyValueExpr:
		printExpr(et.Key, buf, "", ctx)
		buf.WriteString(":")
		printExpr(et.Value, buf, "", ctx)
	case *ast.BasicLit:
		if et.Kind == token.STRING && et.Value[0] == '`' {
			// Make sure we replace any go string literals
			buf.WriteString("\"")
			buf.WriteString(et.Value[1 : len(et.Value)-1])
			buf.WriteString("\"")
		} else {
			buf.WriteString(et.Value)
		}
	case *ast.BinaryExpr:
		printExpr(et.X, buf, "", ctx)
		buf.WriteString(" ")
		buf.WriteString(tokenMap[et.Op])
		buf.WriteString(" ")
		printExpr(et.Y, buf, "", ctx)
	case *ast.SelectorExpr:
		printExpr(et.X, buf, "", ctx)
		buf.WriteString(".")
		printExpr(et.Sel, buf, "", ctx)
	case *ast.CallExpr:
		ident, ok := et.Fun.(*ast.Ident)
		if ok {
			if ident.Name == "len" {
				printExpr(et.Args[0], buf, "", ctx)
				buf.WriteString(".length")
				break
			} else if ident.Name == "copy" {
				printExpr(et.Args[0], buf, "", ctx)
				buf.WriteString(".copy(")
				printExpr(et.Args[1], buf, "", ctx)
				buf.WriteString(")")
				break
			} else if ident.Name == "String" {
				// Handle the conversion of function 'String' to 'toString'
				buf.WriteString("toString(")
				for idx, arg := range et.Args {
					printExpr(arg, buf, "", ctx)
					if idx < len(et.Args)-1 {
						buf.WriteString(", ")
					}
				}
				buf.WriteString(")")
			}
		}
		doCont := false
		switch subT := et.Fun.(type) {
		case *ast.Ident:
			if _, ok := typesMap[subT.Name]; ok {
				// We should only print the args
				for idx, arg := range et.Args {
					printExpr(arg, buf, "", ctx)
					if idx < len(et.Args)-1 {
						buf.WriteString(", ")
					}
				}
				doCont = true
			}
		}
		if doCont {
			break
		}
		printExpr(et.Fun, buf, "", ctx)
		buf.WriteString("(")
		for idx, arg := range et.Args {
			printExpr(arg, buf, "", ctx)
			if idx < len(et.Args)-1 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(")")
	case *ast.ParenExpr:
		buf.WriteString("(")
		printExpr(et.X, buf, "", ctx)
		buf.WriteString(")")
	case *ast.IndexExpr:
		// TODO: Is this the expr for accessing values in a map?
		// TODO: check out if we have a map, if so use the [],if slice use the .elementAt() function.
		printExpr(et.X, buf, "", ctx)
		buf.WriteString(".elementAt(")
		printExpr(et.Index, buf, "", ctx)
		buf.WriteString(")")
		// buf.WriteString("[")
		// buf.WriteString("]")
	case *ast.ArrayType:
		buf.WriteString("ListSlice")
	case *ast.SliceExpr:
		printExpr(et.X, buf, "", ctx)
		buf.WriteString(".slice(")
		printExpr(et.Low, buf, "", ctx)
		buf.WriteString(",")
		printExpr(et.High, buf, "", ctx)
		buf.WriteString(")")
	case *ast.MapType:
		buf.WriteString("Map")
	case *ast.InterfaceType:
		buf.WriteString("interface")
	case *ast.TypeAssertExpr:
		// We don't handle this here at all!
	case nil:
		// Do nothing i guess?
	default:
		fmt.Printf("%sUnknown expr type: %s\n", indent, reflect.TypeOf(e))
	}
}

func printStmt(e ast.Stmt, buf *bytes.Buffer, indent string, ctx *LibraryContext) {
	buf.WriteString(indent)
	switch st := e.(type) {
	case *ast.ReturnStmt:
		buf.WriteString("return")
		if st.Results != nil {
			buf.WriteString(" ")
			if len(st.Results) > 1 {
				buf.WriteString("[")
			}
			for idx, e := range st.Results {
				printExpr(e, buf, "", ctx)
				if idx < len(st.Results)-1 {
					buf.WriteString(", ")
				}
			}
			if len(st.Results) > 1 {
				buf.WriteString("]")
			}
		}
		if indent != "" {
			buf.WriteString(";\n")
		}
	case *ast.ExprStmt:
		printExpr(st.X, buf, "", ctx)
		if indent != "" {
			buf.WriteString(";\n")
		}
	case *ast.DeclStmt:
		printDecl(st.Decl, buf, "", ctx)
		if indent != "" {
			buf.WriteString(";\n")
		}
	case *ast.AssignStmt:
		printAssignStmt(st, buf, indent, ctx)
	case *ast.IfStmt:
		if st.Init != nil {
			// TODO: This doesn't correctly scope the variables defined here.
			// Perhaps the entire if could just be inside of an 'if (true)' block?
			printStmt(st.Init, buf, indent, ctx)
			buf.WriteString(";\n")
			buf.WriteString(indent)
		}
		buf.WriteString("if (")
		printExpr(st.Cond, buf, "", ctx)
		buf.WriteString(") {\n")
		for _, subStmt := range st.Body.List {
			printStmt(subStmt, buf, indent+"  ", ctx)
		}
		buf.WriteString(indent)
		buf.WriteString("}")
		if st.Else != nil {
			buf.WriteString(" else ")
			switch st.Else.(type) {
			case *ast.IfStmt:
				printStmt(st.Else, buf, indent, ctx)
			default:
				buf.WriteString("{")
				printStmt(st.Else, buf, indent, ctx)
				buf.WriteString("}")
			}
		}
		buf.WriteString("\n")
	case *ast.ForStmt:
		if st.Init == nil && st.Cond == nil && st.Post == nil {
			buf.WriteString("while (true) {\n")
		} else {
			buf.WriteString("for (")
			printStmt(st.Init, buf, "", ctx)
			buf.WriteString(";")
			printExpr(st.Cond, buf, "", ctx)
			buf.WriteString(";")
			printStmt(st.Post, buf, "", ctx)
			buf.WriteString(") {\n")
		}
		for _, stmt := range st.Body.List {
			printStmt(stmt, buf, indent+"  ", ctx)
		}
		buf.WriteString(indent)
		buf.WriteString("}\n")
	case *ast.SwitchStmt:
		buf.WriteString("switch (")
		printStmt(st.Init, buf, "", ctx)
		buf.WriteString(") {\n")
		for _, stmt := range st.Body.List {
			printStmt(stmt, buf, indent+"  ", ctx)
		}
		buf.WriteString(indent)
		buf.WriteString("}\n")
	case *ast.RangeStmt:
		printRangeStmt(st, buf, indent, ctx)
	case *ast.CaseClause:
		for _, caseExpr := range st.List {
			buf.WriteString("  case ")
			printExpr(caseExpr, buf, "", ctx)
			buf.WriteString(":\n")
			buf.WriteString(indent + "  ")
		}
		for _, stmt := range st.Body {
			printStmt(stmt, buf, indent+"    ", ctx)
		}
		if len(st.Body) > 0 {
			// If last statement isn't a branch statement then we have to add in an explicit break.
			switch st.Body[len(st.Body)-1].(type) {
			case *ast.BranchStmt:
			default:
				buf.WriteString(indent + "    break;\n")
			}
		} else {
			buf.WriteString("break;\n")
		}
	case *ast.BranchStmt:
		buf.WriteString(tokenMap[st.Tok])
		buf.WriteString(";\n")
	case *ast.IncDecStmt:
		printExpr(st.X, buf, "", ctx)
		buf.WriteString(tokenMap[st.Tok])
		if indent != "" {
			buf.WriteString(";\n")
		}
	case *ast.GoStmt:
		fmt.Printf("Go statements not yet handled.\n")
		// foo() => new Future(() { FUNC() };
	case *ast.BlockStmt:
		for _, stmt := range st.List {
			printStmt(stmt, buf, indent, ctx)
		}
	case nil, *ast.EmptyStmt:
		// Do nothing?
	default:
		fmt.Printf("%sUnhandled function statement: %v\n", indent, reflect.TypeOf(e))
	}
}

func printAssignStmt(st *ast.AssignStmt, buf *bytes.Buffer, indent string, ctx *LibraryContext) {
	if len(st.Lhs) > 1 && len(st.Rhs) == 1 {
		isAssign := st.Tok == token.DEFINE

		switch ms := st.Rhs[0].(type) {
		case *ast.TypeAssertExpr:
			if isAssign {
				buf.WriteString("var ")
			}
			printExpr(st.Lhs[1], buf, "", ctx)
			buf.WriteString(" = ")
			printExpr(ms.X, buf, "", ctx)
			buf.WriteString(" is ")
			printExpr(ms.Type, buf, "", ctx)
			buf.WriteString(";\n")

			buf.WriteString(indent)
			if isAssign {
				buf.WriteString("var ")
			}
			printExpr(st.Lhs[0], buf, "", ctx)
			buf.WriteString(" = ")
			printExpr(ms.X, buf, "", ctx)
		case *ast.IndexExpr:
			// Probably a map?
			if isAssign {
				buf.WriteString("var ")
			}
			printExpr(st.Lhs[1], buf, "", ctx)
			buf.WriteString(" = ")
			printExpr(ms.X, buf, "", ctx)
			buf.WriteString(".containsKey(")
			printExpr(ms.Index, buf, "", ctx)
			buf.WriteString(");\n")
			buf.WriteString(indent)
			if isAssign {
				buf.WriteString("var ")
			}
			printExpr(st.Lhs[0], buf, "", ctx)
			buf.WriteString(" = ")
			printExpr(ms.X, buf, "", ctx)
			buf.WriteString("[")
			printExpr(ms.Index, buf, "", ctx)
			buf.WriteString("]")
		case *ast.CallExpr:
			buf.WriteString("var tmpList ")
			buf.WriteString(tokenMap[st.Tok])
			buf.WriteString(" ")
			printExpr(st.Rhs[0], buf, "", ctx)
			buf.WriteString(";\n")
			for idx, lh := range st.Lhs {
				buf.WriteString(indent)
				if isAssign {
					buf.WriteString("var ")
				}
				printExpr(lh, buf, "", ctx)
				buf.WriteString(" = tmpList[")
				buf.WriteString(strconv.Itoa(idx))
				buf.WriteString("];\n")
			}
		default:
			fmt.Printf("Unknown multi-assign from type: %s\n", reflect.TypeOf(st.Rhs[0]))
		}
		// Unpack all returns from Rhs into Lhs
	} else if len(st.Lhs) == len(st.Rhs) {
		for idx, lh := range st.Lhs {
			if st.Tok == token.DEFINE {
				buf.WriteString("var ")
			}
			// Handle assigning to index here
			lhTyped, ok := st.Lhs[idx].(*ast.IndexExpr)
			if ok {
				// If we have an index expr we need to specially handle
				// choosing between map and std array expressions.
				if isMapIndex(lhTyped) {
					printExpr(lhTyped.X, buf, "", ctx)
					buf.WriteString("[")
					printExpr(lhTyped.Index, buf, "", ctx)
					buf.WriteString("] = ")
					printExpr(st.Rhs[idx], buf, "", ctx)
				} else {
					printExpr(lhTyped.X, buf, "", ctx)
					buf.WriteString(".setAt(")
					printExpr(lhTyped.Index, buf, "", ctx)
					buf.WriteString(",")
					printExpr(st.Rhs[idx], buf, "", ctx)
					buf.WriteString(")")
				}
				continue
			}
			// If assign rhs is func 'append' we specially handle it.
			switch rhTyped := st.Rhs[idx].(type) {
			case *ast.CallExpr:
				ident, ok := rhTyped.Fun.(*ast.Ident)
				if ok {
					if ident.Name == "append" {
						printExpr(lh, buf, "", ctx)
						buf.WriteString(".add(")
						printExpr(rhTyped.Args[1], buf, "", ctx)
						buf.WriteString(")")
						// TODO: deal with ..., use 'addAll'
						// TODO: deal with args[0] not being the RHS.
						//   We would need to create an add inside of an add.
						continue
					} else if ident.Name == "make" {
						printExpr(lh, buf, "", ctx)
						buf.WriteString(tokenMap[st.Tok])
						buf.WriteString("new ")
						printExpr(rhTyped.Args[0], buf, "", ctx)
						buf.WriteString("()")
						continue
					}
				}
			}
			// Default is to print the left, token, right, newline if there are more assignments.
			printExpr(lh, buf, "", ctx)
			buf.WriteString(" ")
			buf.WriteString(tokenMap[st.Tok])
			buf.WriteString(" ")
			printExpr(st.Rhs[idx], buf, "", ctx)

			if idx < len(st.Lhs)-1 {
				buf.WriteString(";\n")
			}
		}
	} else {
		fmt.Printf("Assign statement with more elements on the right hand side than left!? %v", st.Tok)
	}
	// no indent means that this expression is part of another stmt. If it is indented then we need to add a newline.
	if indent != "" {
		buf.WriteString(";\n")
	}
}

func isMapIndex(lhTyped *ast.IndexExpr) bool {
	typX, ok := lhTyped.X.(*ast.Ident)
	if ok {
		decX, ok := typX.Obj.Decl.(*ast.AssignStmt)
		if ok {
			for _, rh := range decX.Rhs {
				complt, ok := rh.(*ast.CompositeLit)
				if ok {
					_, ok := complt.Type.(*ast.MapType)
					if ok {
						return true
					}
				}
			}
		}
	}
	return false
}

func printRangeStmt(r *ast.RangeStmt, buf *bytes.Buffer, indent string, ctx *LibraryContext) {
	// TODO: Figure out if you are ranging
	// over a map, slice, channel

	buf.WriteString("for (int ")
	printExpr(r.Key, buf, "", ctx)
	buf.WriteString(" = 0; ")
	printExpr(r.Key, buf, "", ctx)
	buf.WriteString(" < ")
	printExpr(r.X, buf, "", ctx)
	buf.WriteString(".length; ")
	printExpr(r.Key, buf, "", ctx)
	buf.WriteString("++) {\n")
	buf.WriteString(indent + "  ")
	buf.WriteString("var ")
	printExpr(r.Value, buf, "", ctx)
	buf.WriteString(" = ")
	printExpr(r.X, buf, "", ctx)
	buf.WriteString("[")
	printExpr(r.Key, buf, "", ctx)
	buf.WriteString("];\n")
	for _, stmt := range r.Body.List {
		printStmt(stmt, buf, indent+"  ", ctx)
	}
	buf.WriteString(indent)
	buf.WriteString("}\n")
}

func printDecl(d ast.Decl, buf *bytes.Buffer, indent string, ctx *LibraryContext) {
	switch dt := d.(type) {
	case *ast.GenDecl:
		switch dt.Tok {
		case token.TYPE:
			for _, s := range dt.Specs {
				ts := s.(*ast.TypeSpec)
				switch tsType := ts.Type.(type) {
				case *ast.FuncType:
					buf.WriteString("typedef ")
					if tsType.Results == nil {
						buf.WriteString(indent)
						buf.WriteString("void ")
					} else if len(tsType.Results.List) > 1 {
						buf.WriteString(indent)
						buf.WriteString("List ")
					} else {
						printExpr(tsType.Results.List[0].Type, buf, indent, ctx)
						buf.WriteString(" ")
					}
					printExpr(ts.Name, buf, "", ctx)
					buf.WriteString("(")
					printParams(tsType.Params, buf, "", ctx)
					buf.WriteString(")")
				case *ast.InterfaceType:
					buf.WriteString("abstract class ")
					printExpr(ts.Name, buf, "", ctx)
					buf.WriteString("{\n")
					buf.WriteString("}")
				default:
					fmt.Printf("Unknown type in generic declaration printing: %s\n", reflect.TypeOf(ts.Type))
				}
			}
		case token.VAR, token.CONST:
			var lastVal *ast.ValueSpec
			for sIdx, s := range dt.Specs {
				ts := s.(*ast.ValueSpec)
				for nIdx, n := range ts.Names {
					buf.WriteString("var ")
					buf.WriteString(n.Name)
					buf.WriteString(" = ")
					if len(ts.Values) > nIdx {
						printExpr(ts.Values[nIdx], buf, "", ctx)
						lastVal = ts
					} else if lastVal != nil {
						printExpr(lastVal.Values[len(lastVal.Values)-1], buf, "", ctx)
					} else {
						buf.WriteString("null")
					}
					if nIdx < len(ts.Names)-1 {
						buf.WriteString(";\n")
					}
				}
				if sIdx < len(dt.Specs)-1 {
					buf.WriteString(";\n")
				}
			}
		default:
			fmt.Printf("Unknown token in gen decl: %s", dt.Tok)
		}
	default:
		fmt.Printf("%sUnhandled declaration: %v\n", indent, reflect.TypeOf(d))
	}
}
