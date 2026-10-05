CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS conversations;

CREATE TABLE auth.users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    display_name TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE auth.bots(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES auth.users,
    display_name TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE auth.bot_tokens(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_id UUID NOT NULL REFERENCES auth.bots ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE
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

CREATE TABLE conversations.conversations(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE conversations.direct_conversations(
    conversation_id UUID NOT NULL REFERENCES conversations.conversations ON DELETE CASCADE,
    participant_id_low UUID NOT NULL,
    participant_id_high UUID NOT NULL,
    UNIQUE(participant_id_low, participant_id_high)
);

CREATE TABLE conversations.memberships(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations.conversations ON DELETE CASCADE,
    participant_id UUID NOT NULL,
    role TEXT NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    UNIQUE(conversation_id, participant_id)
);
CREATE INDEX conversation_memberships_conversation_id_index ON conversations.memberships (conversation_id);
CREATE INDEX conversation_memberships_participant_id_index ON conversations.memberships (participant_id);

CREATE TABLE conversations.messages(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL,
    conversation_id UUID NOT NULL REFERENCES conversations.conversations ON DELETE CASCADE,
    body TEXT NOT NULL,
    reply_to_message_id UUID REFERENCES conversations.messages,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE conversations.files(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hash BYTEA NOT NULL,
    name TEXT,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE conversations.message_files(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID REFERENCES conversations.messages ON DELETE CASCADE,
    file_id UUID REFERENCES conversations.files ON DELETE CASCADE,
    name TEXT
);

CREATE TABLE outbox(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
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
    event_id UUID NOT NULL,
    handler TEXT NOT NULL,
    executed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(event_id, handler)
);