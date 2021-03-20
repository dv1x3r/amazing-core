from enum import Enum


class ServiceClass(Enum):
    USER_SERVER = 18
    SYNC_SERVER = 19
    LOCATION = 20
    CLIENT = -1


class UserMessageTypes(Enum):
    MF_CRISP = 10
    MF_SYNCHRONOUS_API = 11
    MF_AGENT = 12
    MAX_FAMILY = 12
    GET_AVATARS = 1
    CLAIM_GIFT = 2
    REGISTER_PLAYER = 3
    GET_PLAYER_MAZES = 4
    START_MAZE_EDIT = 5
    GET_PLAYER_MAZE = 6
    GET_COMMUNITY_MAZE = 7
    END_MAZE_EDIT = 8
    DELETE_MAZE = 9
    GET_DYNAMIC_SURPRISE = 10
    GET_ZONE_MAZES = 11
    GET_COMMUNITY_MAZES = 12
    GET_QUEST_MAZES = 13
    PUBLISH_MAZE = 14
    LOGIN = 15
    GET_PLAYER_NPCS = 16
    CHECK_USERNAME = 17
    REGISTER_FCSS_CODE = 18
    LOGOUT = 19
    GET_INVENTORY_OBJECTS = 20
    GET_BUILD_OBJECTS = 21
    GET_COLLECTION_OBJECTS = 22
    UPDATE_AVATAR_NAME = 23
    SEARCH_VILLAGES = 24
    MOVE_VILLAGE = 25
    RELOGIN = 26
    INIT_LOCATION = 27
    GET_LANG_LOCALE = 28
    ENTER_BUILDING = 29
    BUY_ITEM = 30
    GET_SHOPPING_ITEMS = 31
    LIST_SHOPPING_CATEGORIES = 32
    GET_FRIEND_LIST = 33
    REMOVE_FRIEND = 34
    MANAGE_FRIEND_REQUEST = 35
    GET_BLOCKED_PLAYERS = 36
    MANAGE_BLOCK_PLAYER = 37
    ADD_FRIEND = 38
    UNPUBLISH_MAZE = 39
    GET_RANDOM_WORLD_NAME = 40
    GET_ONLINE_STATUS = 41
    MANAGE_TEST_PLAYER = 42
    UPDATE_QUEST_ITEM = 43
    PLACE_MAZE_ITEM = 44
    REMOVE_MAZE_ITEM = 45
    GET_PLAYER_MAZE_THUMBNAILS = 46
    GET_PLAYER_MAZE_THUMBNAIL = 47
    UPDATE_PLAYER_MAZE_THUMBNAIL = 48
    GET_COMMUNITY_MAZE_THUMBNAILS = 49
    ACCEPT_QUEST = 50
    COMPLETE_QUEST = 51
    ADD_QUEST_ITEM = 52
    END_GAME = 53
    GET_PLAYER_DYNAMIC_SURPRISES = 54
    CLAIM_DYNAMIC_SURPRISE = 55
    GET_NOTIFICATIONS = 56
    GET_NOTIFICATION_BY_PLAYER_NOTIFICATION_ID = 57
    GET_NOTIFICATION_BY_PLAYER_ID = 58
    UPDATE_NOTIFICATION = 59
    GET_PLAYER_MAZE_RATING = 60
    GET_PLAYER_MAZE_RATINGS = 61
    GET_SYSTEM_MAZE_RATING = 62
    RATE_PLAYER_MAZE = 63
    GET_NOTIFICATION_OPTIONS = 64
    UPDATE_NOTIFICATION_OPTION = 65
    CLEAR_NOTIFICATIONS = 66
    GET_HOME_MAZE = 67
    UPDATE_HOME_MAZE = 68
    GET_NOTIFICATION_CATEGORIES = 69
    GET_NOTIFICATION_OPTION_BY_CATEGORY = 70
    ACKNOWLEDGE_NOTIFICATION = 71
    GET_NOTIFICATION_OPTION_BY_PLAYER_NOTIFICATION_OPTION_ID = 72
    GET_NOTIFICATION_CATEGORY = 73
    CLEAR_NOTIFICATIONS_BY_PLAYER_ID = 74
    SET_NOTIFICATION_OPTIONS = 75
    GET_CMS_COLLECTION_ITEMS = 76
    ADD_COLLECTION_ITEM = 77
    DROP_COLLECTION_ITEM = 78
    RESET_COLLECTION = 79
    APPROVE_MAZE_PUBLISHING = 80
    UPDATE_FRIEND_COMMENT = 81
    GET_ACTIVE_FRIEND_LIST = 82
    GET_FRIEND_REQUESTS = 83
    GET_FRIEND_BY_FRIEND_PLAYER_ID = 84
    FIND_FRIEND_BY_PLAYER_FRIEND_ID = 85
    GET_QUEST_TYPES = 86
    DELETE_FRIEND_BY_PLAYER_FRIEND_ID = 87
    SEND_PRIVATE_CHAT_GROUP_INVITE = 88
    SEND_MESSAGE = 89
    GET_CMS_ITEMCATEGORY_IDS = 90
    GET_CMS_ITEMCATEGORY_BY_ID = 91
    GET_PLAYER_MAZE_PLAY = 92
    CONSUME_INVENTORY_ITEM = 93
    MOVE_INVENTORY_ITEM = 94
    SWAP_INVENTORY_ITEMS = 95
    REDEEM_FEATURE_CODE = 96
    GET_MAZE_ITEMS = 97
    GET_YARD_ITEMS = 98
    PLACE_YARD_ITEM = 99
    REMOVE_YARD_ITEM = 100
    GET_SYSTEM_MAZE_PLAY = 101
    GET_AVATAR_ITEMS = 102
    GET_OUTFIT_ITEMS = 103
    START_PRIVATE_CHAT_GROUP = 104
    ACCEPT_PRIVATE_CHAT_GROUP_INVITE = 105
    LEAVE_PRIVATE_GROUP = 106
    GET_CHAT_CHANNEL_TYPES = 107
    START_MAZE_PLAY = 108
    END_MAZE_PLAY = 109
    DRESS_AVATAR = 110
    GET_OUTFITS = 111
    START_QUEST = 112
    HEARTBEAT = 113
    GET_FRIEND_ACCOUNTS_BY_BIRTHDAY = 114
    GET_PRIVATE_CHAT_GROUP_MEMBERS = 115
    REJOIN_PRIVATE_CHAT_GROUP = 116
    ADD_PLAYER_INFRACTION = 117
    GET_PLAYER_CHAT_HISTORY = 118
    GET_QUEST_ITEMS = 119
    START_GAME = 120
    GET_MAZE_SETS = 121
    DECLINE_PRIVATE_CHAT_GROUP_INVITE = 122
    GET_MAZE_SET = 123
    FIND_PRIVATE_CHAT_GROUP_MEMBER = 124
    FIND_PRIVATE_CHAT_GROUP_ID = 125
    COMPLETE_DAILY_ACTIVITY = 126
    GET_QUEST_BY_ID = 127
    GET_WEEKLY_AWARD_STATUS = 128
    GET_PLAYER_WEEKLY_AWARD_HISTORY = 129
    GET_WEEKLY_AWARD_STATUS_BY_PLAYER_WEEKLY_AWARD_ID = 130
    GET_SLOTS = 131
    GET_GAMES = 132
    GET_QUESTS = 133
    GET_DAILY_AWARD_STATUS = 134
    UPDATE_PLAYER_ACTIVE_AVATAR = 135
    GET_SITE_FRAME = 136
    GET_CMS_ITEMCATEGORIES = 137
    NEXT_CAPTCHA = 138
    CHECK_CAPTCHA = 139
    GET_PLAYER_GIFTS = 140
    UPDATE_HOME_THEME = 141
    GET_HOME_THEMES = 142
    GET_PLAYER_HOME_THEMES = 143
    REJECT_QUEST = 144
    GET_PLAYER_RECEIVED_GIFTS = 145
    MANAGE_GIFT_REQUEST = 146
    GET_ANNOUNCEMENTS = 147
    GET_ANNOUNCEMENT = 148
    GET_ACHIEVEMENT_BY_ID = 149
    GET_ACHIEVEMENTS = 150
    GET_ADOPTION_OBJECTS_BY_AVATAR = 151
    GET_OBJECT_IDS_BY_ADOPTION_NUMBER = 152
    FIND_GIFT_BY_PLAYER_GIFT_ID = 153
    GET_ZONES = 154
    LIST_STORE_CATEGORIES = 155
    LIST_STORE_INVENTORY = 156
    LIST_STORES = 157
    PURCHASE_ITEMS = 158
    COMPLETE_TUTORIAL = 159
    UPDATE_LANGLOCALE = 160
    UPDATE_PLAYER_NAME = 161
    GET_PLAYER_ACHIEVEMENTS = 162
    GET_SHARDS = 163
    GET_CRAFTABLE_ITEMS = 164
    GET_CRAFTABLE_ITEM_BY_ID = 165
    PAUSE_NPC = 166
    RESUME_NPC = 167
    GET_PLAYER_QUESTS = 168
    NPC_INTERACTION = 169
    CRAFT_ITEM_BY_CRAFTABLE_ITEM_ID = 170
    GET_NPC_RELATIONSHIPS = 171
    CRAFT_ITEM_BY_ITEMS = 172
    GET_NPCS = 173
    DRESS_AVATAR_ITEMS = 174
    LOCK_HOME = 175
    UNDRESS_AVATAR = 176
    GET_MAZE_PIECES_BY_PLAYER_MAZE_ID = 177
    GET_STORE_THEMES = 178
    GET_ASSETS_BY_OIDS = 179
    GET_CURRENCIES = 180
    ADD_OUTFIT_ITEMS = 181
    ADD_OUTFIT = 182
    REMOVE_OUTFIT_ITEMS = 183
    REMOVE_OUTFIT = 184
    REPLACE_OUTFIT_ITEMS = 185
    SET_CURRENT_OUTFIT = 186
    GET_FRIEND_AVATARS = 187
    GET_PLAYER_QUESTS_BY_QUEST_IDS = 188
    GET_PLAYER_GAMES_BY_ZONE = 189
    GET_HOME_INVITATIONS = 190
    MANAGE_HOME_INVITATIONS = 191
    FIND_PLAYER_BY_NICKNAME = 192
    GET_PLAYER_AVATARS_BY_BIRTHDAY = 193
    SEND_QUEST_INVITE = 194
    GET_HOSTED_QUESTS = 195
    GET_INVITED_QUESTS = 196
    GET_PLAYERS_IN_QUEST = 197
    REMOVE_QUEST_INVITE = 198
    ACCEPT_QUEST_INVITE = 199
    DECLINE_QUEST_INVITE = 200
    GET_PLAYER_RECEIVED_QUEST_INVITE = 201
    GET_PLAYER_QUEST_INVITE = 202
    CLEANUP_PLAYER_MAZES = 204
    GET_PLAYER_STATS = 205
    LIST_STORE_INVENTORY_ITEMS = 206
    GET_OTHER_PLAYER_DETAILS = 207
    GET_TEST_PLAYERS = 208
    GET_STATS_TYPE = 209
    TEST_INVENTORY_UPDATES = 210
    TEST_INVENTORY_REMOVALS = 211
    SET_PLAYER_GIFT_STATUS = 212
    GET_REQUIRED_EXPERIENCE = 213
    GET_QUEST_BY_PARENT_ID = 214
    UPDATE_ONLINE_STATUS = 215
    CHANGE_STORE_ITEM_STOCK = 216
    LOGOUT_SESSION = 217
    ABANDON_QUEST = 218
    GET_AWARD_SETS = 219
    UPDATE_ARM = 220
    UPDATE_RULED_OBJECT = 221
    UPDATE_RULE = 222
    ADD_AWARD = 223
    UPDATE_AWARD = 224
    UPDATE_STATE_OBJECT = 225
    GET_STATE_OBJECTS = 226
    GET_AWARD_BY_ID = 227
    GET_RULE_BY_ID = 228
    GET_RULED_OBJECT_BY_HIERARCHY = 229
    ADD_STATE_OBJECT = 230
    ADD_RULE = 231
    ADD_ARM = 232
    GET_ARM_BY_AWARD = 233
    GET_ARM_BY_RULE = 234
    GET_ARM_BY_STATE_OBJECT = 235
    GET_AWARD_BY_ARM = 236
    GET_RULE_BY_ARM = 237
    GET_STATE_OBJECT_BY_ARM = 238
    OBJECT_INFO_ROWS = 239
    ASSIGN_RULED_SET_TO_RULED_OBJECT = 240
    ADD_RULED_OBJECT = 241
    GET_RULE_SETS = 242
    SAVE_OBJECT_ATTRIBUTES = 243
    GET_OBJECT_TYPES = 244
    GET_CMS_NOTIFICATIONS = 245
    JOIN_QUEST = 246
    GET_ALL_OPERATORS = 247
    GET_ONLINE_STATUSES = 248
    CSTOOL_GET_FULL_ACCOUNT_INFORMATION = 249
    REMOVE_FRIEND_CSTOOL = 250
    UPDATE_YARD = 251
    UPDATE_CHAT_AVAILABILITY = 252
    GET_FRIEND_LIST_CSTOOL = 253
    CREATE_OR_UPDATE_CURRENCY_CSTOOL = 254
    GET_ALL_RULES = 255
    GET_ALL_RULED_OBJECT = 256
    GET_ALL_AWARDS = 257
    DELETE_RULE = 258
    DELETE_RULED_OBJECT = 259
    DELETE_AWARD = 260
    GET_FULL_RULE_BY_ID = 261
    GET_RULED_OBJECT_BY_ID = 262
    GET_CURRENCY_CSTOOL = 263
    GET_ENERGY_LEVEL_CSTOOL = 264
    GET_USER_LEVEL_CSTOOL = 265
    GET_USER_XP_CSTOOL = 266
    GET_USER_NPC_RELATIONSHIP_LEVEL = 267
    CREATE_OR_UPDATE_ENERGY_LEVEL_CSTOOL = 268
    CREATE_OR_UPDATE_USER_LEVEL_CSTOOL = 269
    CREATE_OR_UPDATE_USER_XP_CSTOOL = 270
    CREATE_OR_UPDATE_USER_NPC_RELATIONSHIP_LEVEL = 271
    GET_FULL_RULED_OBJECT_BY_ID = 272
    GET_FULL_AWARD_BY_ID = 273
    DELETE_ASAM = 274
    DELETE_RSRM = 275
    GET_FORMULA_BY_TYPE = 276
    ADD_BUG_REPORT = 277
    FIND_PLAYER_ACCOUNT_CSTOOL = 278
    FIND_NPCS_CSTOOL = 279
    FIND_AVATAR_BY_SKU_CSTOOL = 280
    REGISTER_AVATAR_CSTOOL = 281
    GET_ALL_CATEGORIES = 282
    UNREGISTER_TEST_PLAYERS = 283
    GET_ALL_ATTRIBUTES = 284
    UPDATE_STORE_ITEM_NOTIFY = 285
    GET_ALL_ATTRIBUTES_BY_CATEGORY = 286
    ADD_OBJECT_ATTRIBUTES = 287
    UPDATE_OBJECT_ATTRIBUTES = 288
    GET_ARM_BY_ID = 289
    DELETE_RULE_FROM_RSRM = 290
    DELETE_AWARD_FROM_ASAM = 291
    GET_ALL_ARM = 292
    ADD_RULES_TO_RSRM = 293
    ADD_AWARDS_TO_ASAM = 294
    GET_NPC_WITH_MOST_RELATIONSHIP = 295
    STATE_OBJECT_ARM = 296
    GET_STORE_ITEMS = 297
    MANAGE_CRISP_ACTION = 298
    GET_CRISP_ACTIONS = 299
    SEARCH_ITEMS = 300
    ADD_ITEMS_TO_PLAYER = 301
    ADD_TO_PLAYER_STAT = 302
    RESEND_REG_CONF_KEY = 303
    ACKNOWLEDGE_CONF_KEY = 304
    GET_CMS_MISSIONS = 305
    GET_RULE_INFO = 306
    GET_PLAYER_QUESTS_CSTOOL = 307
    START_TASK_CSTOOL = 308
    GET_OTHER_PLAYER_NPC_RELATIONSHIPS = 309
    GET_BUY_BACK_STORE_ITEMS = 310
    UPDATE_CHAT_AVAILABILITY_FOR_PLAYER_CSTOOL = 311
    GET_CHAT_AVAILABILITY_CSTOOL = 312
    GET_STORE_ITEMS_FOR_ITEMS = 313
    PURCHASE_QUEST = 314
    GET_CMS_EVENTS = 315
    ADD_EVENT = 316
    REQUEST_EVENT = 317
    GET_PENDING_EVENTS = 318
    ACCEPT_EVENT = 319
    REJECT_EVENT = 320
    SAVE_PLAYER_SETTINGS = 321
    LOAD_PLAYER_SETTINGS = 322
    GET_RULED_OBJECT_BY_OBJECT_ID = 323
    SUGGEST_FRIENDS = 324
    SEARCH_USERS = 325
    CSTOOL_UNLOCK_PLAYER_ACCOUNT = 326
    GET_TIERS = 327
    GET_DEBUG_QUESTS = 328
    GET_NPCS_WITH_QUEST_OFFER = 329
    CREATE_PLAYER_MISSION = 330
    GET_PLAYER_MISSION_DETAIL = 331
    GET_PLAYER_MISSIONS = 332
    DELETE_ITEM_FROM_PLAYER = 333
    DEBUG_CREATE_QUEST = 334
    REGISTER_AVATAR_FOR_REGISTRATION = 335
    UPDATE_AVATAR_NAME_FOR_REGISTRATION = 336
    ENTER_STORE = 337
    EXIT_STORE = 338
    GET_ALL_QUESTS = 339
    GET_ITEMS_BY_CATEGORY = 340
    GET_ITEM_CATEGORIES_CSTOOL = 341
    GET_ALL_ZONES = 342
    GET_ALL_NPCS = 343
    GET_ALL_STORE_ITEMS = 344
    UPDATE_PLAYER_EMAIL = 345
    ADD_USER_STORE_ITEM = 346
    ALLOCATE_PLAYER_VILLAGE = 347
    SEARCH_USER_STORE_ITEMS = 348
    UPDATE_CRISP_AVAILABILITY_CS_TOOL = 349
    GET_ALL_VILLAGES = 350
    UPDATE_USER_STORE_ITEM = 351
    REMOVE_USER_STORE_ITEM = 352
    GET_CRISP_AVAILABILITY_CS_TOOL = 353
    GET_PLAYER_CRISP_DATA = 354
    UPDATE_CHAT_BLOCKED_BY_PARENT_CS_TOOL = 355
    UPDATE_CHAT_BLOCKED_BY_PARENT = 356
    GET_PLAYER_WORLDNAME_BY_GASID = 357
    GET_SCS_INVALID_CODE_STATUS = 358
    GET_CRISP_ACTIONS_FOR_NONLOGGED_INSESSION = 359
    GET_PLAYER_CONTAINER_BY_PLAYER = 360
    GET_PLAYER_CONTAINER_TYPE = 361
    GET_WEB_CONTENT = 362
    LIST_ITEM_BY_PLAYER = 363
    LIST_ITEM_BY_CONTAINER = 364
    GET_WEB_CONTENT_BY_PTAG = 365
    GET_EULA = 366
    GET_PLAYER_ACCOUNT_BY_GASID = 367
    GET_PLAYER_ID_BY_GASID = 368
    LIST_VILLAGE_PLOTS = 369
    SEND_VILLAGE_INVITE = 370
    ACCEPT_VILLAGE_INVITE = 371
    REJECT_VILLAGE_INVITE = 372
    REPORT_ABUSE = 373
    GET_ABUSE_REPORTS = 374
    GET_ABUSE_REPORTS_BY_PLAYER = 375
    UPDATE_ABUSE_REPORT = 376
    PURCHASE_USER_STORE_ITEMS = 377
    DECORATE_USER_STORE = 378
    REDEEM_USER_STORE_SALES = 379
    GET_RULE_INSTANCE_INFO = 380
    ACCEPT_EULA = 381
    ADD_PLAYER_LIKE = 382
    REMOVE_PLAYER_LIKE = 383
    LIST_TOP_LIKES = 384
    UNDECORATE_USER_STORE = 385
    USER_STORE_INFO = 386
    ADD_AWARD_TO_AWARD_QUEUE = 387
    REGISTRATION_RECORD = 388
    GET_OBJECT_FROM_SKU = 389
    GET_INFRACTION_TYPES = 390
    GET_ALL_CMS_MISSIONS = 391
    CHECK_EMAIL_AVAILABILITY = 392
    ADD_BOOKMARK = 393
    REMOVE_BOOKMARK = 394
    LIST_BOOKMARKS = 395
    LIST_LIMITS = 396
    LIST_VILLAGE_USERS = 397
    GET_DAILY_AWARDS = 398
    LIST_PLAYER_LIKES = 399
    VILLAGE_INFO = 400
    DECORATED_USER_STORE_ITEMS = 401
    LIST_BOOKMARK_OBJECT_TYPES = 402
    LINK_FACEBOOK_ACCOUNT = 403
    UNLINK_FACEBOOK_ACCOUNT = 404
    LIST_FRIEND_FACEBOOK_IDS = 405
    LIST_PLAYER_ACCOUNTS_BY_FACEBOOK_IDS = 406
    GET_ALL_USER_STORE_ITEMS = 407
    ADD_QUEST_ITEM_NOTIFY = 408
    SET_PLAYER_FINDABLE = 409
    GET_PLAYER_CHAT_RECEIVED_HISTORY = 410
    GET_ALL_STORE_THEMES = 411
    SELL_ITEM = 412
    GET_ALL_BUILDINGS = 413
    GET_ALL_ITEMCATEGORIES = 414
    PLANT_PLAYER_ITEM = 415
    GET_FEATURED_ITEMS = 416
    GET_NPC_RELATIONSHIP_LEVELS = 417
    SAVE_GAME_STATE = 418
    LIST_SAVED_GAMES = 419
    LOAD_GAME = 420
    SEND_GAME = 421
    ACCEPT_GAME = 422
    REJECT_GAME = 423
    CMS_GET_GAME_STATE = 427
    GET_ALL_FEATURED_ITEMS = 428
    MANAGE_NPC_FRIEND_REQUEST = 429
    PLACE_VILLAGE_ITEM = 430
    GET_CMS_VILLAGE_ROLES = 431
    SEND_EMAIL_MESSAGE = 432
    LIST_INBOX_MESSAGES = 433
    GET_INBOX_MESSAGE = 434
    MARK_INBOX_MESSAGE = 435
    GET_RULE_COUNT = 436
    LIST_SENT_MESSAGES = 437
    GET_VILLAGE_ROLES = 438
    ASSIGN_VILLAGE_ROLE = 439
    GET_CMS_VILLAGE_TEMPLATES_LOCKED = 440
    CONFIRM_FINDABLE = 441
    CREATE_PRIVATE_VILLAGE = 442
    GET_RULE_INFO_LIST = 443
    GET_RULE_INSTANCE_INFO_LIST = 444
    GET_OTHER_PLAYER_DETAILS_LIST = 445
    ATTACH_ITEMS = 446
    DETACH_ITEMS = 447
    ADD_AVATAR_CSTOOL = 448
    KICK_OUT = 449
    HARVEST_PLAYER_ITEM = 450
    GET_NPC_GIFTS = 451
    REMOVE_VILLAGE_ITEM = 452
    LIST_VILLAGE_ITEMS = 453
    GET_CMS_VILLAGE_TEMPLATES_UNLOCKED = 454
    LIST_FINDABLE_PLAYERS = 455
    GET_ALL_CRAFTABLE_ITEMS = 456
    LIST_CMS_MESSAGES = 457
    GET_CMS_MESSAGE = 458
    GET_FRIENDSHIP_REQUEST_COUNTS = 459
    GET_NPC_INTERACTIONS = 460
    GET_CRAFTABLE_TYPES = 461
    VILLAGE_INVITE_STATUS_NOTIFY = 462
    FIND_NPCS_BY_NAME = 463
    CHANGE_PASSWORD = 464
    FORGOT_PASSWORD = 465
    SET_PASSWORD = 466
    ACCEPT_EVENT_NOTIFY = 467
    LIST_QUEST_SPAWN_IDS = 468
    GET_RULE_INSTANCE_INFO_FOR_UI = 469
    GET_RULE_INFO_LIST_FOR_UI = 470
    SEND_TEST_MESSAGE = 471
    SAVE_AVATAR_IMAGE = 472
    GET_AVATAR_IMAGES = 473
    GET_OBJECTS_LOCK_INFO = 474
    GET_RULE_TEMPLATE_ID_FOR_UI = 475
    CRAFT_ITEMS_BY_CRAFTABLE_ITEM_IDS = 476
    GET_ALL_GAMES = 477
    UPDATE_GENDER_AND_LOCATION = 478
    GET_GEOGRAPHICAL_LOCATIONS = 479
    GET_CANNED_MESSAGE_CATEGORIES = 480
    PRE_FILTER_NAME_CHECK_AVAILABILITY = 481
    CRAFT_BRACELET = 482
    GET_NPC_ITEMS = 483
    ITEM_INTERACTION = 484
    CREATE_QUEST_GAME = 485
    ENTER_MAZE = 486
    REMOVE_AND_DETACH_MAZE_ITEM = 487
    ATTACH_AND_PLACE_ITEM = 488
    START_CRAFTING_BRACELET = 489
    GET_QUEST_FROM_PARENT = 490
    CREATE_QUEST = 491
    GET_PLAYER_ATTACHMENT_ITEMS_BY_SENDER_ID = 492
    SET_FRIEND_ORDER = 493
    GET_RULE_INSTANCE_DATA = 494
    GET_MEMBERSHIP_SUBSCRIPTION_BY_ID = 495
    CLAIM_MEMBERSHIP_SUBSCRIPTION = 496
    GET_USER_SUBSCRIPTION_INFO = 497
    PROCESS_CRISP_MESSAGE = 498
    LIST_SENDING_PLAYER_ATTACHMENT_ITEMS = 499
    OID_TO_DBID = 500
    DBID_TO_OID = 501
    FILTER_BAD_WORD = 502
    DOCK_ITEM = 503
    GET_MESSAGES_COUNT = 504
    GET_VIEWS = 505
    DISCOVER_ONLINE_USER = 506
    GET_UNLOCK_INFO_PEER = 507
    DELETE_ZING_FROM_PLAYER = 508
    SEND_NOTIFICATION_CS_TOOL = 509
    FRIEND_STATUS_NOTIFY = 510
    GET_RULE_COUNT_LIST = 511
    SELL_PLAYER_ITEMS = 512
    SEND_MESSAGE_CS_TOOL = 513
    GIFT_STATUS_NOTIFY = 514
    LEVEL_STATUS_NOTIFY = 515
    PLAYER_NOTIFICATION_NOTIFY = 516
    GET_SYSTEM_NOTIFICATIONS = 517
    GET_SYSTEM_MESSAGES = 518
    ATTACH_SINGLE_ITEM = 519
    MANAGE_CRISP_NOTIFY = 520
    ADD_PLAYER_AWARD_NOTIFY = 521
    GET_STATEFUL_INSTANCE = 522
    MARK_ANNOUNCEMENT_READ = 523
    REORDER_FRIENDS_ORDINAL = 524
    GET_ESTORE_TRANS_INFO_FOR_CSTOOL = 525
    GET_ESTORE_POINTS_FOR_CSTOOL = 526
    GET_PLAYER_EXTERNAL_SITE_MAP = 527
    SET_PLAYER_EXTERNAL_SITE_MAP = 528
    UPDATE_TOKEN_PLAYER_EXTERNAL_SITE_MAP = 529
    GET_AUTHENTICATION_PLAYER_EXTERNAL_SITE_MAP = 530
    SEND_EMAIL_ATTACHMENT = 531
    SET_BRACELETS_ARM_ORDER = 532
    SEND_NOTIFICATION_TEMPLATE = 533
    CHECK_ONLINE_STATUS_TEMPLATE = 534
    GET_OTHER_PLAYER_DETAILS_TEMPLATE = 535
    CROSS_USER_SERVER_PLAYER_TEMPLATE = 536
    CROSS_USER_SERVER_PLAYER_TEMPLATE_NOTIFY = 537
    GET_PLAYER_INFO_ON_OTHER_SHARD_TEMPLATE = 538
    MULTI_DAR_TEMPLATE = 539
    MULTI_DAR_COMPLEX_TEMPLATE = 540
    GET_RANDOM_NAMES = 541
    SELECT_PLAYER_NAME = 542
    VALIDATE_NAME = 543
    QUEST_EVENT_IN_PROGRESS = 544
    GET_QUEST_ALL_FROM_PARENT = 546
    GET_PLAYER_VOTED_LIST = 547
    CREATE_PLAYER_VOTED = 548
    VOTE_ON_PLAYER_VOTED = 549
    GET_PLAYER_VOTED = 550
    GET_COMPOSED_ITEM = 551
    SAVE_COMPOSED_ITEM = 552
    COMPLETE_PLAYER_VOTED = 553
    WITHDRAW_INSTANCE_FOR_VOTED = 554
    GET_LEADER_BOARD_INFO_FOR_PLAYER_VOTED = 555
    GET_FRIENDS_PLAYER_VOTEDS = 556
    MANAGE_SYNCHRONIZED_OBJECTS = 557
    PURCHASE_UNLOCK = 558
    VALIDATE_AND_REDEEM_MOBILE_PRODUCT_PURCHASE = 559
    GET_PLAYER_VOTED_DATA = 560
    UPDATE_ACCOUNT_REFERRAL = 561
    REMOVE_PLAYER_EXTERNAL_SITE_MAPPING = 562
    GET_INVITED_PLAYER_QUEST = 563
    ACCEPT_QUEST_TRANSACTION = 564
    FINALIZE_QUEST_TRANSACTION = 565
    GET_CLIENT_VERSION_INFO = 566
    PURCHASE_WALLET_ITEM = 567
    GET_PLAYER_ONLINE_STATUS = 568
    GET_ITEM_BY_ID = 569
    GET_PUBLIC_ASSETS_BY_OIDS = 570
    GET_PUBLIC_ITEMS_BY_OIDS = 571
    GET_PUBLIC_ITEM_CATEGORIES = 572
    GET_PARENT_BUILDING_ID = 573
    CREATE_RECIPE = 574
    GET_NPCS_BY_CHILD_HIERARCHIES = 575
    RECYCLE_RECIPE = 576
    ENHANCE_RECIPE = 577
    MAX_TYPE = 577


class SyncMessageTypes(Enum):
    ADD_OBJECT = 1
    MOVE_OBJECT = 2
    REMOVE_OBJECT = 3
    CHANGE_OBJECT = 4
    VILLAGE_HANDOFF_QUERY = 5
    SERVER_CHANGE_OBJECT = 6
    MOVE_PLAYER = 7
    REMOVE_PLAYER = 8
    CHAT = 9
    START_EVENT = 10
    STOP_EVENT = 11
    EMOTE = 12
    BIND_USER_NOTIFY = 13
    PAUSE_NPC = 14
    RESUME_NPC = 15
    UPDATE_NPC_SCRIPT = 16
    NOTIFICATION = 17
    CLIENT_TEST = 18
    ADD_ITEM = 19
    EVICT = 20
    ENTER_LOC = 21
    EXIT_LOC = 22
    ECHO = 23
    GET_VILLAGE = 24
    REFRESH = 25
    BIND = 26
    BIND_QUERY = 27
    BIND_VILLAGE_NOTIFY = 28
    VILLAGE_HANDOFF = 29
    FIND_SERVER = 30
    UPDATE_FILTER = 31
    SEND_NOTIFY = 32
    LOGIN = 33
    UPDATE_USER_INFO = 34
    UPDATE_USER_FRIENDS = 35
    UPDATE_USER_GROUPS = 36
    MANAGE_USER_FRIENDS = 37
    MANAGE_USER_GROUPS = 38
    LIST_SERVERS = 39
    LIST_VILLAGES = 40
    LIST_USERS = 41
    MOVE_VILLAGE = 42
    START_NPCS = 43
    STOP_NPCS = 44
    USER_SESSION_HANDOFF = 45
    EMOTE_SVC = 46
    ACTION = 47
    LOGOUT = 48
    CLOSE_ZONE = 49
    TERMINATE_USER_SESSION = 50
    GET_PLAYER_COUNT = 51
    UPDATE_LOCATION = 52
    HEARTBEAT_NOTIFY = 54
    RELOGIN = 55
    MAX_TYPE = 55


class ClientMessageTypes(Enum):
    ADD_OBJECT = 1
    MOVE_OBJECT = 2
    REMOVE_OBJECT = 3
    CHANGE_OBJECT = 4
    SERVER_CHANGE_OBJECT = 5
    MOVE_PLAYER = 6
    REMOVE_PLAYER = 7
    CHAT = 8
    START_EVENT = 9
    STOP_EVENT = 10
    EMOTE = 11
    CHANGE_WEIGHT = 12
    PAUSE = 13
    RESUME = 14
    NOTIFICATION = 16
    CHANGE_SERVER = 17
    SEND_NOTIFY = 18
    INT_LIST = 19
    ADD_PLAYER = 20
    UPDATE_NPCS = 21
    STOP_NPC = 22
    MINIMAP = 23
    POS_RECAP = 24
    ACTION = 25
    EVICT = 26
    ONLINE_STATUS = 27
    CHANGE_OBJECT_STATE = 28
    MAX_TYPE = 28


class ResultCode(Enum):
    INCOMPLETE = -1
    OK = 0
    APP = 5
    APP_DB = 6
    ERR = 10
    QUEUE = 11
    DB = 20
    DB_QUEUE = 21
    DB_NO_RETRY = 22
    NO_MEM = 30
    COMM = 40
    CONN_FAILED = 41
    DISCONNECT = 42
    SHUTDOWN = 43
    IO = 44
    TIMEOUT = 45
    BUSY = 46
    COMM_INIT = 47
    CANCEL = 48
    WOULD_BLOCK = 49
    PROTOCOL_VER = 50
    SERIALIZE = 51
    PENDING_IO = 52
    ASYNC = 53
    CONN = 54
    CHNL_CLOSED = 55
    CONN_EXHST = 56
    NO_DEST = 57
    CHRONO = 58
    NOT_READY = 59


class AppCode(Enum):
    ILG = -1
    OK = 0
    CONTINUE = 1
    ERR = 10
    TYPE = 11
    INTERLOCK = 12
    NO_INTERLOCK = 13
    INPUT = 14
    DUP_REQUEST = 15
    SVC_VER = 16
    MULT_ERRS = 17
    NOT_FOUND = 18
    NOT_IMPLEMENTED = 19
    MEMORY = 20
    NOT_READY = 21
    STATE = 22
    PERM = 40
    XPERM = 41
    AUTH = 42
    SESSION = 43
    NO_SPACE = 50
    FILE_PATH = 51
    FILE_NAME = 52
    FILE = 53
    FILE_IO = 54
    FILE_MODE = 55
    FILE_DAMAGED = 56
    FILE_ACCESS = 57
    FILE_NF = 58
    FILE_RENAME = 59
    FILE_MOVE = 60
    FILE_COPY = 61
    DB = 70
    DUP_KEY = 71
    DBRC_UNKNOWN = 72
    INVALID_USER = 100
    INVALID_NAME = 101
    INVALID_RELATIONSHIP = 102
    INVALID_PRIMARY_KEY = 103
    PIN_MISMATCH = 104
    SITE_CREATION_FAILED = 105
    INVALID_ITEM = 106
    NOT_SELLABLE_ITEM = 107
    INVALID_PLACEMENT = 109
    INSUFFICIENT_FUND = 110
    INVALID_INVENTORY_SEARCH = 111
    INVALID_FEATURE_CODE = 112
    INVENTORY_ITEM_NOT_EXIST = 113
    NOT_LOG_IN = 114
    DUPLICATE_NICKNAME = 115
    PLAYER_CREATION_FAILED = 116
    PLAYER_DELETION_FAILED = 117
    INVALID_AUTHENTICATION = 118
    BLANK_CREDENTIALS = 119
    INVALID_SITE_INFO = 120
    INVALID_INVENTORY_ORDINAL = 121
    CODE_STATUS_UPDATE_FAILED = 122
    SECRET_CODE_VERIFICATION_FAILED = 123
    INVALID_OBJECT_ESTORE_SKU = 124
    CREATE_PLAYER_ACCOUNT_FAILED = 125
    INVALID_PIN = 126
    INVALID_ASSET = 127
    BACKPACK_IS_FULL = 128
    NO_ACTIVE_VILLAGER = 129
    ITEM_NOT_IN_CONTAINER_OR_BACKPACK = 130
    ITEMS_IN_SAME_CONTAINER_OR_BACKPACK = 131
    INVENTORY_TYPE_IS_NULL = 132
    INVENTORY_TYPES_NOT_SAME = 133
    SWAPPED_ITEMS_ARE_SAME = 134
    ITEM_NOT_OWNED_BY_SESSION_PLAYER = 135
    INVALID_SITE_CONTENT = 136
    STARTER_TOWN_ID_NOT_CONFIGURED = 137
    PRESENTABLE_SLOTS_USED_UP = 138
    INVALID_GAME_ID = 139
    INVALID_TOKEN = 140
    TOO_MANY_POINTS = 141
    NO_CHECKPOINT = 142
    INAVLID_ACTION = 143
    NAME_CANNOT_BE_EMPTY = 160
    NOT_ACCEPT_PRESENTABLE_ITEM = 161
    PLAYER_HAS_NO_HOME_VILLAGE = 162
    ONLY_ONE_SEARCH_CRITERION_ALLOWED = 163
    ONE_SEARCH_CRITERION_MUST_BE_PROVIDED = 164
    INVALID_ITEM_RELATION_NAME = 165
    UNPAID_ACCOUNT = 166
    ALREADY_QUEUED = 167
    INVALID_COOLDOWN = 168
    INVALID_AUTH = 169
    INSUFFICIENT_PERMISSION = 170
    ALREADY_HAVE_A_TOURIST = 171
    INVALID_ITEM_OR_BUILDING_PLACEMENT = 172
    NOT_ENOUGH_SPACE = 173
    CRAFTING_FAIL = 174
    CHILD_ALREADY_PLACED = 175
    SLOT_ALREADY_USED = 176
    SLOT_NOT_ON_PARENT = 177
    INCOMPATIBLE_SLOT = 178
    PARENT_NOT_CRAFTABLE = 179
    INVALID_OUTFIT_NO = 180
    INSUFFICIENT_FUNDS = 181
    INAPPROPRIATE_LANGUAGE = 182
    INVALID_VILLAGE_NAME = 183
    PENDING_VILLAGE_EXISTS = 184
    CREATE_VILLAGE_FAILED = 185
    PLAYER_DOES_NOT_OWN_INVENTORY = 200
    NOT_EMOTE_INVENTORY = 201
    INVALID_EMOTE = 202
    INVALID_PLAYER_EMOTE_ID = 203
    INVALID_CODE = 210
    SCS_BLOCKED = 225
    INTERNAL_ERROR = 500
    TEST_ERROR_CODE = 9999
