# Concurrency Control Algorithms

Implementations of classic mutual exclusion algorithms: Bakery, Dekker, and Peterson.

## Project structure

- `ada/` – Ada implementations of the algorithms
- `go/` – Go implementations with process simulation and execution tracing
- `Makefile` – builds the Ada programs
- `display-travel-2.bash` – script for visualizing Ada and Go program output

## Building

For Ada:

```bash
make
```

This builds the binaries: `bakery`, `dekker`, and `peterson`.

## Running

Ada:

```bash
./bakery > out
./display-travel-2.bash out
```

Or:

```bash
make run-bakery
make run-dekker
make run-peterson
```

Go:

```bash
go run ./go/bakery.go > out
./display-travel-2.bash out

go run ./go/dekker.go > out
./display-travel-2.bash out

go run ./go/peterson.go > out
./display-travel-2.bash out
```
