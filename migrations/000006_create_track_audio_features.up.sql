CREATE TABLE IF NOT EXISTS track_audio_features (
  track_id          UUID    NOT NULL PRIMARY KEY REFERENCES tracks(id) ON DELETE CASCADE,
  danceability      FLOAT   NOT NULL,
  valence           FLOAT   NOT NULL,
  energy            FLOAT   NOT NULL,
  acousticness      FLOAT   NOT NULL,
  instrumentalness  FLOAT   NOT NULL,
  liveness          FLOAT   NOT NULL,
  speechiness       FLOAT   NOT NULL,
  tempo             FLOAT   NOT NULL,
  loudness          FLOAT   NOT NULL
);