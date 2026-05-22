CREATE TABLE friend_requests (
    id          BIGSERIAL PRIMARY KEY,
    sender_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status      TEXT NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'accepted', 'rejected')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (sender_id, receiver_id)
);

CREATE INDEX idx_friend_requests_receiver ON friend_requests(receiver_id, status);
CREATE INDEX idx_friend_requests_sender ON friend_requests(sender_id, status);

CREATE TABLE friends (
    id         BIGSERIAL PRIMARY KEY,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, friend_id)
);

CREATE INDEX idx_friends_user ON friends(user_id);
