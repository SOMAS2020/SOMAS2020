## Requirements

- [Go 1.15.x](https://golang.org/dl/) (1.15.5 preferred)
    - `go version` should produce `go version go1.15....`

(Note that the reference system as in CI is `Ubuntu 20.04`. It is preferable to run this OS, but if not, that's fine.)

## Text Editor & Language Server

- You're encouraged to use [VSCode](https://code.visualstudio.com/), with the [Go extension](https://code.visualstudio.com/docs/languages/go).
- You _must_ use [gopls](https://godoc.org/golang.org/x/tools/gopls). This is a language server that ensures an even code style among members, and it helps catches mistakes as you type them. This should be available in most LSP-enabled editors. If you don't understand this, use VSCode and the Go extension as in 1.

## Coding rules

1. Trust the language server. Red lines == death. Yellow lines == close to death. An example where it might be very tempting to let yellow lines pass are in `struct`s:
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

4. Do not use `panic` or `die`--return an `error` instead!

5. Do not use system-specific packages (e.g. `internal/syscall/unix`).

6. Keep your team repo base up to date with the main repo.

7. Use the superior `errors.Errorf` to create your errors so that we have a stack trace.

8. If you define a new `iota`, please implement the `String()`, `GoString()`, `MarshalText()`, and `MarshalJSON()` functions. See [`#150`](https://github.com/SOMAS2020/SOMAS2020/pull/150).

## Repo

1. Each team will work off a fork of the main repo. Your team is responsible for all development happening in the fork, and are responsible to keep your own fork up-to-date, as well as to pull in changes to the main repo periodically. (Remember to give your teammates _write access_ to the fork!)

2. Your team's fork must pass CI + infrastructure team code reviews before it can be merged into the main repo. Make sure you detail your changes succinctly in your PR, and _KEEP DIFFS SMALL_. Infra might not need to read every line, but having small reviews to do is helpful. Again, DO NOT TOUCH CODE YOU'RE NOT SUPPOSED TO. 

3. Make sure your fork is up-to-date with the main repo's `main` branch before submitting a PR. You should install [https://github.com/wei/pull](https://github.com/wei/pull) on your fork to automate this.

4. Your fork should inherit the Github Actions for CI as well, this means PRs into your `main` branch should run CI.

5. Feel free to link your PR in the Discord infrastructure channel to request a review.


## Dependencies

1. The usual way of getting dependencies for Golang should work, i.e. `go get <MODULE_LINK>`. Try to refrain from including rarely-used or dodgy-looking dependencies: everyone running the code needs to get the code on their computer (Golang does this automatically). If in doubt, ask.
