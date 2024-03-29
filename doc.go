/*
Package applicator provides an interface to run applicator functions based on
struct tags

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
    if err := applicator.ApplyRun(&a); err != nil {
        // error
    }

Builtin functions

    trim
            Calls strings.TrimSpace on the string

    lower
            Calls strings.ToLower on the string

    fillNil
            Sets the pointer to the zero value, if nil

*/
package applicator
