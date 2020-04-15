# Applicator

[![GoDoc](https://godoc.org/github.com/levenlabs/go-applicator?status.svg)](https://godoc.org/github.com/levenlabs/go-applicator)

Package applicator provides an interface to run functions based on struct tags

For example, if you accept strings over an API endpoint, you might want to trim
all strings when receiving them. Additionally, you want to lowercase the
Username field.

    type MethodArgs struct {
        Username string `apply:"trim,lower"`
        Name     string `apply:"trim"`
    }

Once you have an instance of MethodArgs, you can call applicator.Apply() to
apply the functions for each field.

    a := MethodArgs{"myusername "}
    if err := applicator.Apply(&a); err != nil {
        // error
    }

Builtin functions

    trim
            Calls strings.TrimSpace on the string

    lower
            Calls strings.ToLower on the string

    fillNil
            Ensures that the field is not nil by filling it with the default value
