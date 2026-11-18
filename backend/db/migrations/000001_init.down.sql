DROP TABLE inbox;
DROP TABLE outbox;

DROP TABLE bff.peers;

DROP TABLE directory.user_profiles;
DROP TABLE directory.search_index;
DROP TABLE directory.handles;

DROP TABLE conversations.message_attachments;
DROP TABLE conversations.messages;
DROP TABLE conversations.group_conversations;
DROP TABLE conversations.direct_conversations;
DROP TABLE conversations.memberships;
DROP TABLE conversations.conversations;

DROP TABLE identity.sessions;
DROP TABLE identity.bots;
DROP TABLE identity.users;
DROP TABLE identity.accounts;

DROP TABLE media.uploads;
DROP TABLE media.assets;

DROP TYPE outbox_kind;
DROP TYPE outbox_status;
DROP TYPE account_type;
DROP TYPE asset_status;
DROP TYPE upload_status;

DROP SCHEMA media;
DROP SCHEMA bff;
DROP SCHEMA directory;
DROP SCHEMA identity;
DROP SCHEMA conversations;
