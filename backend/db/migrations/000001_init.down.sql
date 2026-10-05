DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS outbox;

DROP TABLE IF EXISTS conversations.message_files;
DROP TABLE IF EXISTS conversations.files;
DROP TABLE IF EXISTS conversations.messages;
DROP TABLE IF EXISTS conversations.direct_conversations;
DROP TABLE IF EXISTS conversations.memberships;
DROP TABLE IF EXISTS conversations.conversations;

DROP TABLE IF EXISTS auth.sessions;
DROP TABLE IF EXISTS auth.bot_tokens;
DROP TABLE IF EXISTS auth.bots;
DROP TABLE IF EXISTS auth.users;

DROP SCHEMA IF EXISTS auth;
DROP SCHEMA IF EXISTS conversations;
