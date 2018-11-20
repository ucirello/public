# cirello.io/oversight

[![travis-ci](https://api.travis-ci.org/ucirello/oversight.svg?branch=master)](https://travis-ci.org/ucirello/oversight)
[![GoDoc](https://godoc.org/cirello.io/oversight?status.svg)](https://godoc.org/cirello.io/oversight)


Package oversight makes a nearly complete implementation of the Erlang
supervision trees.

Refer to: http://erlang.org/doc/design_principles/sup_princ.html

go get [-f -u] cirello.io/oversight

http://godoc.org/cirello.io/oversight


## Quickstart
```
supervise := oversight.Oversight(
	oversight.Processes(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
			log.Println(1)
		}
		return nil
	}),
)

ctx, cancel := context.WithCancel(context.Background())
defer cancel()
if err := supervise(ctx); err != nil {
	log.Fatal(err)
}
```