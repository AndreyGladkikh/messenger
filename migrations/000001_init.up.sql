CREATE TABLE chats(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL,
    name TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE TABLE chat_participants(
    chat_id UUID REFERENCES chats ON DELETE CASCADE,
    participant_id UUID NOT NULL,
    role TEXT NOT NULL,
    joined_at timestamp NOT NULL,
    PRIMARY KEY (chat_id, participant_id)
);

CREATE TABLE messages(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL,
    chat_id UUID NOT NULL REFERENCES chats ON DELETE CASCADE,
    body TEXT NOT NULL,
    response_to_message_id UUID REFERENCES messages,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE TABLE files(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hash BYTEA NOT NULL,
    name TEXT,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE messages_files(
    message_id UUID REFERENCES messages ON DELETE CASCADE,
    file_id UUID REFERENCES files ON DELETE CASCADE,
    name TEXT,
    PRIMARY KEY (message_id, file_id)
);
