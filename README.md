# polygen

Simple-minded Go code generator

How many times have you written the same loop or function with slight type variations? Wouldn't it be nice to have a simple way to use the templating system, hassle-free to generate all that similar code?

## polygen to the rescue!

polygen lets you create templates for functions in files, and then use `go generate` to create the version of the function you need. The generated code will be created inside the package where `go generate` is invoked, so you can have a library of functions you reuse often with different types and generate their specific variations when you need them.

A simple case of transforming a slice of any type to a slice of `interface{}` would look like this.
 
 ```Go
{{ $a := T "a" }}
{{ $NA := Name "a" }}

func convert{{$NA}}slice(xs []{{$a}}) []interface{} {
	ys := make([]interface{}, len(xs))
	for i, x := range xs {
		ys[i] = x
	}

	return ys
}
```
Here, `a` is just a placeholder name for a type. If you save this as `convertslice.tgo` inside package `github.com/youruser/awesome/utils/`, you can create a version of this function that works for `int` just by adding this line to a file in the package that needs it.

```Go
//go:generate polygen github.com/youruser/polythings/convertslice.tgo -t"a:int"`
```

The first parameter is the package + template path. The `-t` flag lets you map your placeholder types to actual types. Of course, you can map as many types as you want. You can also add a second parameter to manually set the name of the output file (although you souldn't usually need it).

polygen provides two template functions to help you create the "generic code" you need.

The `T` function provides the actual type as set by the command line flag `-t`.

The `Name` function provides a representation of the type that can be used as part of an idenfifier. The rules to convert types to names are very simple:

* `[]` is changed to `List`
* `*` is changed to `PtrTo`
* Any remaining `[` or `]` are removed

For example, `[]*Person` is changed to `ListPtrToPerson`, and `map[string]Foo` is changed to `mapstringFoo`.

Use `Name` as part of your generated functions and you will get functions with different names for different types. This lets you generate different variants into the same package.

Currently, imported types from other packages are not supported, as this would mean to find the correct package and import it in the generated code. The easy workaround is to create a custom type in your package.

```Go

type foo otherpackage.Foo

//go:generate path/to/template.tgo outfile -t"a:foo"
```
