# Introduction to game files

**WARNING: For educational purposes only!**

## How this can help

- You can run parts of the decompiled code with breakpoints to understand how it works (helped to implement codec)
- By modifying original game dll you can log stuff to the text file (helped to understand exceptions and implement login responses)

## How to decompile

The Amazing World is made on [Unity](https://unity.com/) engine. The last game version was built using Unity version [**5.3.6f1**](https://unity3d.com/get-unity/download/archive). This means that most of the game code is stored in the ```AmazingWorld_Data\Managed\Assembly-CSharp.dll``` file. It is possible to decompile and compile it back using the following tools:

- [ILSpy](https://github.com/icsharpcode/ILSpy): decompiles almost all the code with no issues (use **C# 4.0 / VS 2010**).
- [dnSpy](https://github.com/dnSpy/dnSpy): use where ILSpy failed (only one file in our case).
- [Visual Studio](https://visualstudio.microsoft.com/): edit code and build dll file.

### Decompilation notes:

- ILSpy may decompile **DownloadAgent.cs** with errors.
- Check **LoadAssetFromWeb()** method if it has something like
  ```'((<LoadAssetFromWeb>c__Iterator3)(object)this).<>__Finally0();'```
- Use dnSpy to get this code parts right.

## How to repack

The following tools were made to perform repack and tests (**\tools** directory):

- Amazing-Repacker Visual Studio solution: There is already an empty Assembly-CSharp folder that you can use for source code.
- Build Assembly-CSharp project with Visual Studio.
- replace_dll.bat: a fast way to swap new library in the game folder.

### Repacking notes

You can use ```Debugger.LogWarning("");``` method to log stuff into the ```AmazingWorld_Data\output_log.txt``` file.

## Exploring asset caches

The game client used to download Unity asset bundles from the server. Luckily, we have a lot of cached assets straight from the Cache folder. You can find [zip archive here](https://github.com/dv1x3r/amazing-core/releases) under the very first release with 1914 files inside.

### But how do I know...

That is a good question. Some of the files are asset bundles, some are deprecated and do not work with the latest game version, there could be anything. I wanted to find Player_Base.unity3d asset, so I created a PowerShell script to extract each file as Unity asset bundle. Tools used:

- [UtinyRipper](https://github.com/mafaca/UtinyRipper): slightly modified version to close console window if failed.
- extract_cache.ps1: PowerShell script to iterate every file with UtinyRipper.

When script did its job, I got 673 folders with prefabs and other Unity stuff.

### Build your own caches

Using the Amazing-Bundler tool you can build your own asset bundles. There is an editor script which allows you to do that using right click on prefab.
