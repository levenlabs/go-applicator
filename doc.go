/*
Package helper provides an interface to run helper functions based on
struct tags

For example, if you accept strings over an API endpoint, you might want
to trim all strings when receiving them.

    type MethodArgs struct {
        Username string `helper:"trim,lowercase"`
        Name     string `helper:"trim"`
    }

Once you have an instance of MethodArgs, you can call helper.Run() to
run all the helpers.

    a := MethodArgs{"myusername "}
    if err := helper.Run(&a); err != nil {
        // error
    }

Builtin helpers

    trim
            Calls strings.TrimSpace on the string

    lower
            Calls strings.ToLower on the string

*/

package helper
