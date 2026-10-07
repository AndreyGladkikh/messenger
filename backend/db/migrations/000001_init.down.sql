DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS outbox;

DROP TABLE IF EXISTS conversations.message_files;
DROP TABLE IF EXISTS conversations.files;
DROP TABLE IF EXISTS conversations.messages;
DROP TABLE IF EXISTS conversations.direct_conversations;
DROP TABLE IF EXISTS conversations.memberships;
DROP TABLE IF EXISTS conversations.conversations;

DROP TABLE IF EXISTS identity.sessions;
DROP TABLE IF EXISTS identity.bot_tokens;
DROP TABLE IF EXISTS identity.bots;
DROP TABLE IF EXISTS identity.users;
DROP TABLE IF EXISTS identity.accounts;

DROP TYPE account_type;

DROP SCHEMA IF EXISTS identity;
DROP SCHEMA IF EXISTS conversations;
