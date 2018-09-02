# golang_errors_benchmark

Purpose of this repository is to propose a "good enough" and opinionated way of dealing with adding some context to golang error handling.

### Background
So we all know that golang on itself has very limited usage of errors. Quite a lot of people use [pkg/errors](https://github.com/pkg/errors) to add additional context of the errors as well as stacktraces to where the error happened.


Most used functions from `pkg/errors` are `errors.WithMessage` and `errors.Wrap`, they both have different usage, and this repository shows simple way of using them.

### Requirements

So requirements for "good enough" error handling strategy are:

* capture place of where error happened
* add additional context for errors, as on different code levels, errors has different purpose (error saying that you failed to fetch User vs error saying your http request failed)

### Tests
`similar_users` show a very stupid usage of very popular code structuring, where calls are: Service -> Repository -> DBDriver.

3 different strategies are used when adding additional errors messages:

* use `errors.WithMessage` only
* use `errors.Wrap` only
* mixed approach

### Results

Obviously using only `errors.WithMessage` doesn't fulfill our requirements, as stacktrace is not gathered.

On the other hand, using `errors.Wrap` everywhere leads to capturing not necessary and duplicated stacktrace frames.

You can also see performance differences using all 3 strategies when calling `service.Find(...)`


```
all-with wrapper 5000000	       247 ns/op
all-wrap wrapper 1000000	      1865 ns/op
mixed approach   1000000	      1138 ns/op
```

In case of out-of-context error wrapping those are benchmarking results:

```
errors.WithMessage 20000000	        59.7 ns/op
errors.Wrap        2000000	       649 ns/op
```

So using errors.WithMessage is 10 times faster.

This seems like a not important matter, errors doesn't happen to often in our application right?

Yeah, right, but imagine a sudden flow of 100% error calls together with increased traffic in your service. Then this can have real impact on your response times if you have like 10 calls per request to `errors.something` this will make a much bigger difference.

### Proposal

How to use this package in "good enough" way? 
My proposal is to :

* Use `errors.Wrap` whenever calling code that is not yours (for examply stdlib or 3rd party code)
* Use `errors.WithMessage` in your code elsewhere to add some more information for error
* Use `errors.New` from `dpg/errors` to produce new errors with stacktrace from the beginning

tried this approach and it works fine. Of course isn't perfect when 3rd party library already returns `dpg/errors`, but this has been covered already on many different issues at github (like [this one](https://github.com/pkg/errors/issues/144#issuecomment-354554983))