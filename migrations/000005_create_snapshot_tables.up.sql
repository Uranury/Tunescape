CREATE TABLE IF NOT EXISTS tracks (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  spotify_id  TEXT NOT NULL UNIQUE,
  name        TEXT NOT NULL,
  popularity  INT  NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS snapshots (
  id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_snapshots_user_id ON snapshots(user_id);

CREATE TABLE IF NOT EXISTS snapshot_tracks (
  snapshot_id UUID NOT NULL REFERENCES snapshots(id) ON DELETE CASCADE,
  track_id    UUID NOT NULL REFERENCES tracks(id)    ON DELETE CASCADE,
  position    INT  NOT NULL,
  PRIMARY KEY (snapshot_id, track_id)
);