alias b := build
alias t := test

run:
  go run .

build:
  go build .

test:
  go test ./test/ -v
