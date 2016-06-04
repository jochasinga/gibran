gibran
======
A minimal web framework	in Go that promises nothing but love and abstraction.

Why Gibran
----------
Go web frameworks often feel like _extra horns on a rhinoceros_ or _extra sets
of teeth for a shark_. They feel excessive, considering how [Go was already designed with excellent
idioms to work on an MVC web application](https://medium.com/code-zen/why-i-don-t-use-go-web-frameworks-1087e1facfa4#.chb3c27ad). No one sums this up best than Albert Einstein:


> "Everything should be made as simple as possible, but not simpler."
> —Albert Einstein


A Go web framework, IMHO, should
+ bring structure just enough to be a clean canvas for creativity of the builder.
+ allow the builder to adopt any packages.
+ save time, and most importantly,
+ be extensible and feels like the language.


What Gibran does
---------------
A few things gibran will do for you are:
+ [Inversion of Control](https://medium.com/code-zen/wtf-is-dependency-injection-1c599231d95c#.yoai7vj6i) -- It helps you write loose-coupling code and encourage the use of sub-packages.
+ Structure -- It generates the standard project structure for you, but it won't
enforce you to use it. Restructure as desired, but make sure you understand its use of [interfaces](https://medium.com/code-zen/go-interfaces-and-delegation-pattern-f962c138dc1e#.wlmzvpnfo) beforehand.
+ That's it. Trust me, you don't need anything more in Go.

Getting Started
---------------
To get started, first install with
```bash

$ go get github.com/jochasinga/gibran

```

This installed `gibran` commandline program to the `$GOPATH/bin/gibran`, if your `$GOPATH` is set correctly. In the current directory, usually `$GOPATH/src/<yourvcs>/<yourname>` Start your new project with
```bash

$ gibran startproject myapp

```

or to start a project in any directory, simple provide the directory path as the second argument.

```bash

$ gibran startproject myapp $HOME/code/jochasinga

```

The project directory should be created as shown below:

```bash

myapp/
├── brokers
├── controllers
├── main.go
├── models
├── routers
├── tests
└── views

```

If you look into `main.go` file, you will see a very minimal scaffolding code to
let you run the server right away. You may notice something is different:

```go

package main

import (
        "log"
        "net/http"

        "yourvcs/yourname/myapp/brokers"
)

var (
    Routers = brokers.NewRouters()
)

func main() {
        router := Routers.HandleRoutes()
        log.Fatal(http.ListenAndServe(":8080", router))
}

```

In the root of the project directory, Run the server with `gibran run` and browser to `localhost:8080` to see the welcome page.

Brokers and Delegates
---------------------
gibran generates a `broker` and a `delegate` for every single sub-package.
A broker is an interface representing all the methods in a package and act as
the project's global agency for import. A delegate is a struct which implements
the corresponding broker, thus allowing the package's methods to be useable from
other packages via the broker.
By importing brokers instead of the actual package, you are always using abstractions
instead of directly using the package's namespace to refer to its functions.
This has a couple of benefits, such as avoid import cycles and easy unit testing.

More to come
------------
