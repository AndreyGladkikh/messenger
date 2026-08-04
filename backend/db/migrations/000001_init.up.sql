CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS messaging;

CREATE TABLE auth.users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE auth.sessions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users ON DELETE CASCADE,
    refresh_token_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    user_agent TEXT NOT NULL,
    ip INET NOT NULL
);
CREATE INDEX session_user_id_index on auth.sessions (user_id);
CREATE INDEX session_refresh_token_index on auth.sessions (refresh_token_hash);

CREATE TABLE messaging.chats(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kind TEXT NOT NULL,
    name TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE messaging.chat_participants(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID REFERENCES messaging.chats ON DELETE CASCADE,
    participant_id UUID NOT NULL,
    role TEXT NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    UNIQUE(chat_id, participant_id)
);
CREATE INDEX chat_participants_chat_id_index ON messaging.chat_participants (chat_id);
CREATE INDEX chat_participants_participant_id_index ON messaging.chat_participants (participant_id);

CREATE TABLE messaging.private_chats(
    chat_id UUID NOT NULL REFERENCES messaging.chats ON DELETE CASCADE,
    first_participant_id UUID NOT NULL,
    second_participant_id UUID NOT NULL,
    UNIQUE(first_participant_id, second_participant_id)
);

CREATE TABLE messaging.messages(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL,
    chat_id UUID NOT NULL REFERENCES messaging.chats ON DELETE CASCADE,
    body TEXT NOT NULL,
    reply_to_message_id UUID REFERENCES messaging.messages,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE messaging.files(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hash BYTEA NOT NULL,
    name TEXT,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE messaging.message_files(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID REFERENCES messaging.messages ON DELETE CASCADE,
    file_id UUID REFERENCES messaging.files ON DELETE CASCADE,
    name TEXT
);

CREATE TABLE outbox(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    event_payload JSONB,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    status TEXT NOT NULL DEFAULT 'pending',
    claimed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    attempts INT NOT NULL DEFAULT 0,
    errors TEXT[],
    next_retry_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX outbox_occurred_at_index ON outbox(occurred_at) WHERE status = 'pending';
CREATE INDEX outbox_next_retry_at_occurred_at_index ON outbox(next_retry_at, occurred_at) WHERE status = 'retry';

CREATE TABLE inbox(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id TEXT NOT NULL,
    handler TEXT NOT NULL,
    executed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(event_id, handler)
);