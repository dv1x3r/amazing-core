# New player registration

| Action                          | Where | What happens                                                                       |
| ------------------------------- | ----- | ---------------------------------------------------------------------------------- |
| 1. Entering registration scene  |       | Sending GSFGetPublicItemCategoriesSvc                                              |
|                                 |       | Sending GSFGetSiteFrameSvc                                                         |
| 2. Requesting random zing name  |       | Sending GSFGetRandomNamesSvc (name_part_type = 'second_name')                      |
| 3. Confirming zing name         |       | Sending GSFValidateNameSvc, if filter_name is not empty then name is not valid     |
| 4. Building family name         |       | Sending GSFGetRandomNamesSvc (name_part_type = 'Family_1', 'Family_2', 'Family_3') |
| 5. Confirming family name       |       | Sending GSFSelectedPlayerNameSvc                                                   |
| 6. Submitting registration form |       | Sending GSFCheckUsernameSvc                                                        |
|                                 |       | Sending GSFRegisterPlayerSvc                                                       |
|                                 |       | Sending GSFRegisterAvatarForRegistrationSvc                                        |
