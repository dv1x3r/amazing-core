# Registration Messages

When the player clicks "I'm new!" in the main menu, the intro sequence begins.

Player chooses a name, a family name, and account credentials.

| # | Message                                                                  | Invoked at                                                    |
| - | ------------------------------------------------------------------------ | ------------------------------------------------------------- |
| 1 | [`GetPublicItemCategories`](./get-public-item-categories.md)             | `USPopulateIntroInventoryGrid.cs:46`                          |
| 2 | [`GetSiteFrame`](../system/get-site-frame.md)                            | `IntroManager.cs:156`                                         |
| 3 | [`GetRandomNames`](./get-random-names.md)                                | `IntroServiceCallManager.cs:74,88`                            |
| 4 | [`ValidateName`](./validate-name.md)                                     | `USValidateName.cs:36`, `USRegisterPlayer.cs:151`             |
| 5 | [`SelectPlayerName`](./select-player-name.md)                            | `USWorldName.cs:349`                                          |
| 6 | [`CheckUsername`](./check-username.md)                                   | `IntroServiceCallManager.cs:111`, `RegistrationManager.cs:52` |
| 7 | [`RegisterPlayer`](./register-player.md)                                 | `IntroServiceCallManager.cs:135`, `USRegisterPlayer.cs:125`   |
| 8 | [`RegisterAvatarForRegistration`](./register-avatar-for-registration.md) | `IntroServiceCallManager.cs:322`                              |
