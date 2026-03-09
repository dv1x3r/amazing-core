# Startup Messages

On game start, before the main menu becomes interactive, `AkamaiServerCheck` opens a temporary session to perform a version check, then closes it.

| # | Message                                                | Invoked at                |
| - | ------------------------------------------------------ | ------------------------- |
| 1 | [`GetClientVersionInfo`](./get-client-version-info.md) | `AkamaiServerCheck.cs:70` |
