# Release

- No race conditions in testing allowed, as this is automation and using clis. Too problematic.

## Run Full Tests

```shell
export GOTESTS='superslow'

packages=(
    "docgen"
    "docker"
    "fancy"
    "gotools"
    "tooling"
    )
for pkg in $packages
do
    #gotestsum --format dots-v2 --packages "./$pkg"  -- -shuffle=on -tags integration
    GOTESTS='slow' go test ./... -json -v -shuffle=on -race -tags integration | tparse -notests -smallscreen #-pulse 1s
done
```

Alternative using tparse (easier to read)

```shell
GOTESTS='slow' go test ./... -json -v -shuffle=on -tags integration | tparse -notests -smallscreen
```
