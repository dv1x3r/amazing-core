# Source reference

## pkg.go.dev

Amazing Core keeps its detailed handler, message, and GSF object reference in Godoc comments.

[**Message handlers**](https://pkg.go.dev/github.com/dv1x3r/amazing-core@{{#include ../../vars/reference-version.md}}/internal/game#Handler)
document the `internal/game` handler methods. Use this page to understand when the client calls a handler and what server workflow the handler coordinates.

[**Message payloads**](https://pkg.go.dev/github.com/dv1x3r/amazing-core@{{#include ../../vars/reference-version.md}}/internal/network/gsf/messages#pkg-types)
document `internal/network/gsf/messages`. Use this page for request, response, and notification payload structs. These comments should focus on the wire message and confirmed client behavior.

[**GSF object types**](https://pkg.go.dev/github.com/dv1x3r/amazing-core@{{#include ../../vars/reference-version.md}}/internal/network/gsf/types#pkg-types)
document `internal/network/gsf/types`. Use this page for reusable objects nested inside messages, such as players, avatars, assets, items, zones, OIDs, and enum values.

## Unreleased versions

Public pkg.go.dev can show an exact unreleased commit only through a Go pseudo-version.

Resolve a branch to a pseudo-version with:

```sh
go list -m github.com/dv1x3r/amazing-core@dev
```

Then open:

```text
https://pkg.go.dev/github.com/dv1x3r/amazing-core@<pseudo-version>
```
