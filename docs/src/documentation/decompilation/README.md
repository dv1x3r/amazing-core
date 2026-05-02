# Game decompilation

## Assembly-CSharp

The game is built using Unity engine. The latest version of the game was compiled with Unity 5.3.6f1.

Most of the game’s core logic is contained in the `Assembly-CSharp.dll` file.

This DLL may be decompiled, modified, and reassembled using common .NET tools such as:

- `ILSpy`: primary decompiler, successfully decompiles nearly all code without issues;
- `dnSpy`: used for edge cases where ILSpy fails (in this project, only a single file required it);

**This should be done for educational purposes only.**
