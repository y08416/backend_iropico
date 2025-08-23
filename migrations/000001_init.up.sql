-- ルーム状態のENUM（存在したらスキップ）
DO $$ BEGIN
  CREATE TYPE room_status AS ENUM ('waiting','started','ended');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- users
CREATE TABLE IF NOT EXISTS users (
  id         BIGSERIAL PRIMARY KEY,
  name       TEXT NOT NULL,
  uuid       TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- rooms
CREATE TABLE IF NOT EXISTS rooms (
  id             BIGSERIAL PRIMARY KEY,
  code           TEXT NOT NULL UNIQUE,
  host_user_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  status         room_status NOT NULL DEFAULT 'waiting',
  current_round  INTEGER NOT NULL DEFAULT 0,
  created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
  started_at     TIMESTAMP,
  ended_at       TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_rooms_code ON rooms(code);

-- players
CREATE TABLE IF NOT EXISTS players (
  id           BIGSERIAL PRIMARY KEY,
  room_id      BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  has_cleared  BOOLEAN NOT NULL DEFAULT FALSE,
  joined_at    TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE (room_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_players_room ON players(room_id);

-- history（提出）
CREATE TABLE IF NOT EXISTS history (
  id           BIGSERIAL PRIMARY KEY,
  room_id      BIGINT  NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  round_no     INTEGER NOT NULL,
  user_id      BIGINT  NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  -- 出題色（HSV）
  target_h     INTEGER NOT NULL CHECK (target_h >= 0 AND target_h <= 360),
  target_s     REAL    NOT NULL CHECK (target_s >= 0 AND target_s <= 1),
  target_v     REAL    NOT NULL CHECK (target_v >= 0 AND target_v <= 1),

  -- 結果
  score        REAL    NOT NULL CHECK (score >= 0 AND score <= 1),
  result_json  JSONB,

  -- 画像は保存しない運用だが、跡を残すなら以下を利用可
  photo_mime   TEXT,
  photo_url    TEXT,

  created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE (room_id, round_no, user_id)
);
CREATE INDEX IF NOT EXISTS idx_history_room_round       ON history(room_id, round_no);
CREATE INDEX IF NOT EXISTS idx_history_user             ON history(user_id);
CREATE INDEX IF NOT EXISTS idx_history_room_round_user  ON history(room_id, round_no, user_id);