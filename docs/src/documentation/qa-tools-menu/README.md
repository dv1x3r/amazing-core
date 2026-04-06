# QA Tools Menu

The game contains a hidden QA tools menu that can be enabled through `ServerConfig.xml`.

## Enabling it

Add a `Version` attribute to the config:

```
Version = '1.8'
```

After launching the game with this flag enabled, press `F3` to open the QA tools menu.

`F2` also should open a related QA layout, but it probably requires an additional UI file.

## How It Works

The `Version` attribute is obfuscated, and is not treated as a normal version number.

The client decodes it as an encoded options flag in `GameSettings.cs:748`.

1. Parse the value as a float.
2. Multiply it by 10.
3. Round it to an integer.
4. If the result is divisible by 9, accept it as a valid options code.
5. Divide it by 10 and interpret the result as a bitmask.

## Bitmask Values

The decoded bitmask controls three flags:

```
1 -> debuggingEnabled
2 -> skipIntro
4 -> showVersion
```

In practice, valid values are:

```
1.8 -> debug enabled
2.7 -> skip intro
3.6 -> debug enabled + skip intro
4.5 -> show version
5.4 -> debug enabled + show version
6.3 -> skip intro + show version
7.2 -> debug enabled + skip intro + show version
```
