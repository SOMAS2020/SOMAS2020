## Requirements

- [Go 1.15.x](https://golang.org/dl/) (1.15.5 preferred)
    - `go version` should produce `go version go1.15....`

## Text Editor & Language Server

- You're encouraged to use [VSCode](https://code.visualstudio.com/), with the [Go extension](https://code.visualstudio.com/docs/languages/go).
- You _must_ use [gopls](https://godoc.org/golang.org/x/tools/gopls). This is a language server that ensures an even code style among members, and it helps catches mistakes as you type them. This should be available in most LSP-enabled editors. If you don't understand this, use VSCode.

## Coding rules

1. Trust the language server. Red lines == death. Yellow lines == close to death. An example where it might be very tempting to let yellow lines pass are in structs:
```golang
type S struct {
    name string
    age  int
}

s1 := S {"pitt", 42} // NO
s2 := S {"pitt"} // NO (even though it initialises age to 0)

s3 := S{name: "pitt", age: 42} // OK, if we add fields into S, this will still be correct
s4 := S{name: "pittson"} // OK if `pittson`'s age is 0
```

2. Write tests where required. There is ample guide online on how to do this in Golang. Tests will be run alongside CI when you pull into the main repo. If anyone breaks anything, it's easy to observe that if you have tests. Otherwise, your code will be broken unknowingly.

3. DO NOT TOUCH code you don't own unless you have a good reason to. If you have a good reason to, do it in a separate PR and notify the owners of the code.

4. Keep your team repo base up to date with the main repo.

## Repo

1. Each team will work off a fork of the main repo. Your team is responsible for all development happening in the fork, and are responsible to keep your own fork up-to-date, as well as to pull in changes to the main repo periodically.

2. Your team's fork must pass CI + infrastructure team code reviews before it can be merged into the main repo. Make sure you detail your changes succinctly in your PR, and _KEEP DIFFS SMALL_. Infra might not need to read every line, but having small reviews to do is helpful. Again, DO NOT TOUCH CODE YOU'RE NOT SUPPOSED TO. 

