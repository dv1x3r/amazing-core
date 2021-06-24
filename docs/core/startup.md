# Starting the game

| Action                   | Where                                 | What happens                                                                  |
| ------------------------ | ------------------------------------- | ----------------------------------------------------------------------------- |
| 1. Server check request  | AkamaiServerCheck.OnConnect()         | Sending [GSFGetClientVersionInfoSvc](../messages/get_client_version_info.md)  |
| 2. Server check response | GetClientVersionInfoResponseHandler() | Comparing clientVersionInfo and GameSettings.clientLocalVersion               |
|                          |                                       | Setting GameSettings.needsUpdate flag and blocking the game if local < server |
