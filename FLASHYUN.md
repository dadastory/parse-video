# FlashYun library integration

Canonical upstream: https://github.com/wujunwei928/parse-video
Integration branch: `flashyun-downloads`; baseline `57590df94b49aafc827eb98e85c845eaa1cb3daf`.
The parent repository pins the tested submodule commit.

`parser.NewClient(httpClient).Parse(ctx, sharingText)` is the reusable library
entry point. Configure the standard HTTP client before construction, including
its transport, explicit proxy and timeout. Each invocation creates isolated
platform receivers and request clients, propagates cancellation, and preserves
native redirect policies without modifying another caller's HTTP client.
The existing package-level CLI/server entry points remain compatible.

The Xiaohongshu adapter requests public mobile pages, recognizes current
`noteData.data.noteData` and legacy `note.noteDetailMap` state, preserves image,
video and live-photo resources, and returns `ErrNoMedia` for missing media or a
platform error. Fixtures contain synthetic data rather than private share tokens.

Run `go test ./...` and `go test -race ./parser`. FlashYun imports only the parser
package; no Gin web server, upstream rate limiter, downloader or additional
container enters its runtime. External platform access remains best-effort.
