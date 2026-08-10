CREATE TABLE IF NOT EXISTS games (
                       id                  TEXT PRIMARY KEY,
                       title               TEXT NOT NULL,
                       company             TEXT,
                       year                INTEGER,
                       description         TEXT,

                       cover_path          TEXT,
                       background_type     TEXT CHECK (
                           background_type IN ('image', 'video')
                               OR background_type IS NULL
                           ),
                       background_path     TEXT,
                       bgm_path            TEXT,
                       bgm_enabled         INTEGER NOT NULL DEFAULT 0 CHECK (bgm_enabled IN (0, 1)),

                       executable_path     TEXT NOT NULL,
                       working_directory   TEXT,

                       favorite            INTEGER NOT NULL DEFAULT 0 CHECK (favorite IN (0, 1)),
                       progress            TEXT NOT NULL DEFAULT 'not_started'
                           CHECK (progress IN ('not_started', 'playing', 'completed')),

                       total_play_seconds  INTEGER NOT NULL DEFAULT 0,
                       last_played_at      INTEGER
);

CREATE TABLE IF NOT EXISTS (
                      id      TEXT PRIMARY KEY,
                      name    TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS (
                      game_id TEXT NOT NULL,
                      tag_id TEXT NOT NULL,
                      PRIMARY KEY (game_id, tag_id),
                      FOREIGN KEY (game_id) REFERENCES games(id),
                      FOREIGN KEY (tag_id) REFERENCES tags(id)
)

CREATE TABLE IF NOT EXISTS (
                      queue_type TEXT NOT NULL
                      CHECK (queue_type IN ('default', 'current')),

                      game_id    TEXT NOT NULL,
                      position   INTEGER NOT NULL,

                      PRIMARY KEY (queue_type, game_id),
                      UNIQUE (queue_type, position),
                      FOREIGN KEY (game_id) REFERENCES games(id)
);