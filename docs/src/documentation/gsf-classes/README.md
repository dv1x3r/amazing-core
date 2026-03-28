# GSF Classes

This section documents game requests and responses grouped by flow.

Each request is a GSF message sent by the client to the server.

## Assembly-CSharp

The game is built using Unity engine. The latest version of the game was compiled with 5.3.6f1.

Most of the game’s core logic is contained in the `Assembly-CSharp.dll` file.

This DLL may be decompiled, modified, and recompiled using common .NET C# tools such as:

- `ILSpy`: primary decompiler, successfully decompiles nearly all code without issues;
- `dnSpy`: used for edge cases where ILSpy fails (in this project, only a single file required it);
