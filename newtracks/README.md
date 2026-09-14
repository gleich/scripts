# newtracks

Run `go run ./newtracks`, or install with `make` and run `newtracks`.
Requires macOS, Music, and permission for your terminal to automate Music.

Add source playlist names to the `playlists` slice in `newtracks.go`:

```go
var playlists = []string{
	"chill",
}
```

Each source gets a `new <name>` playlist containing its last 50 entries in
reverse order. The bottom of the source's stored playlist order is treated as
newest. View the destination in playlist order to see newest first.
Generated playlists live in the `NEW` playlist folder, which is created if
needed. Existing destinations are moved into it even when their tracks are
unchanged.

The command runs once, preserving duplicate entries and using all entries when
there are fewer than 50. An empty source produces an empty destination.
Unchanged playlists are skipped. Changed playlists have their contents rebuilt,
so logged counts include all removed and added entries. Source playlists and
library songs are preserved. A failed update stops the command; rerun it to
repair a partially filled destination.

Run `go test ./newtracks` for unit tests. To also exercise Music using temporary
playlists that are removed afterward, run
`NEWTRACKS_MUSIC_TEST=1 go test ./newtracks -run TestMusic -v`.
The integration test requires at least two tracks in the Music library.
