-- ===== Extensions =====
-- UUID生成などに使う拡張
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ===== Enum (状態) =====
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'room_status') THEN
    CREATE TYPE room_status AS ENUM ('waiting', 'playing', 'finished');
  END IF;
END$$;

-- ===== users =====
CREATE TABLE IF NOT EXISTS users (
  id            BIGSERIAL PRIMARY KEY,
  uuid          UUID NOT NULL DEFAULT gen_random_uuid(),
  display_name  TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ===== rooms =====
CREATE TABLE IF NOT EXISTS rooms (
  id              BIGSERIAL PRIMARY KEY,
  host_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  room_name       TEXT NOT NULL,
  join_code       VARCHAR(8) NOT NULL UNIQUE,
  status          room_status NOT NULL DEFAULT 'waiting',
  time_limit_sec  INTEGER NOT NULL DEFAULT 30,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  started_at      TIMESTAMPTZ,
  ended_at        TIMESTAMPTZ
);

-- ===== rounds =====
CREATE TABLE IF NOT EXISTS rounds (
  id                BIGSERIAL PRIMARY KEY,
  room_id           BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  round_number      INTEGER NOT NULL,
  -- 表示用HEXと、計算用HSVを両方保存
  target_color_hex  CHAR(7) NOT NULL CHECK (target_color_hex ~* '^#[0-9A-F]{6}$'),
  target_h          INTEGER NOT NULL CHECK (target_h BETWEEN 0 AND 359),
  target_s          SMALLINT NOT NULL CHECK (target_s BETWEEN 0 AND 100),
  target_v          SMALLINT NOT NULL CHECK (target_v BETWEEN 0 AND 100),
  started_at        TIMESTAMPTZ,
  ended_at          TIMESTAMPTZ,
  UNIQUE (room_id, round_number)
);
CREATE INDEX IF NOT EXISTS idx_rounds_room_id ON rounds(room_id);

-- ===== players (room membership) =====
CREATE TABLE IF NOT EXISTS players (
  user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  room_id    BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  is_ready   BOOLEAN NOT NULL DEFAULT false,
  joined_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, room_id)
);
CREATE INDEX IF NOT EXISTS idx_players_room_id ON players(room_id);

-- ===== player_rounds (scores per round) =====
CREATE TABLE IF NOT EXISTS player_rounds (
  user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  round_id    BIGINT NOT NULL REFERENCES rounds(id) ON DELETE CASCADE,
  score       INTEGER NOT NULL DEFAULT 0,
  has_cleared BOOLEAN NOT NULL DEFAULT false,
  analyzed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, round_id)
);
CREATE INDEX IF NOT EXISTS idx_player_rounds_round_id ON player_rounds(round_id);