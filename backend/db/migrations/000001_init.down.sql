DROP TABLE inbox;
DROP TABLE outbox;

DROP TABLE bff.peers;

DROP TABLE directory.entries;

DROP TABLE conversations.message_files;
DROP TABLE conversations.files;
DROP TABLE conversations.messages;
DROP TABLE conversations.group_conversations;
DROP TABLE conversations.direct_conversations;
DROP TABLE conversations.memberships;
DROP TABLE conversations.conversations;

DROP TABLE identity.sessions;
DROP TABLE identity.bots;
DROP TABLE identity.users;
DROP TABLE identity.accounts;

DROP TYPE outbox_kind;
DROP TYPE outbox_status;
DROP TYPE account_type;

DROP SCHEMA bff;
DROP SCHEMA directory;
DROP SCHEMA identity;
DROP SCHEMA conversations;
