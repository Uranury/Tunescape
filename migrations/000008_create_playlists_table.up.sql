CREATE TABLE playlists (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    spotify_playlist_id TEXT NOT NULL,
    name                TEXT NOT NULL,
    external_url        TEXT NOT NULL,
    embed_url           TEXT NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_playlists_user ON playlists(user_id);
