CREATE TABLE chats(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL,
    name TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- CREATE TABLE private_chats(
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     chat_id UUID REFERENCES chats ON DELETE CASCADE,
--     user1_id UUID NOT NULL,
--     user2_id UUID NOT NULL
-- );

CREATE TABLE chat_participants(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID REFERENCES chats ON DELETE CASCADE,
    participant_id UUID NOT NULL,
    role TEXT NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT unique_chat_participant UNIQUE(chat_id, participant_id)
);
CREATE INDEX chat_participants_chat_id_index ON chat_participants (chat_id);
CREATE INDEX chat_participants_participant_id_index ON chat_participants (participant_id);

CREATE TABLE messages(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL,
    chat_id UUID NOT NULL REFERENCES chats ON DELETE CASCADE,
    body TEXT NOT NULL,
    reply_to_message_id UUID REFERENCES messages,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE files(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hash BYTEA NOT NULL,
    name TEXT,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE messages_files(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID REFERENCES messages ON DELETE CASCADE,
    file_id UUID REFERENCES files ON DELETE CASCADE,
    name TEXT
);

CREATE TABLE events(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type TEXT NOT NULL,
    event_payload JSONB,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    status TEXT NOT NULL DEFAULT 'pending',
    claimed_at TIMESTAMP WITH TIME ZONE,
    attempts INT NOT NULL DEFAULT 0,
    error TEXT,
    next_retry_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE event_handler_executions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events ON DELETE CASCADE,
    handler_type TEXT NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    error TEXT,
    next_retry_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT event_handler_executions_unique UNIQUE(event_id, handler_type)
);