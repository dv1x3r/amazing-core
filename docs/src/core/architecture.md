# Architecture

## Folder Structure

```
cmd/           - entry point
data/          - sql migrations
internal/      
├── api/       - http server for admin dashboard and asset streaming
├── game/      - game server and message handling
├── network/   - tcp server protocol implementation
├── services/  - business logic and database interaction
├── config/    - configuration variables
├── lib/       - shared libraries (e.g. logging, helpers)
tools/         - development tools (e.g. asset importers)
web/           - embedded frontend for admin dashboard
```
