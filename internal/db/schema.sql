-- ルーム状態のENUM（重複作成は無視）
DO $$ BEGIN
  CREATE TYPE room_status AS ENUM ('waiting','started','ended');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- users：フロント生成のUUIDで識別
CREATE TABLE IF NOT EXISTS users (
  id         BIGSERIAL PRIMARY KEY,
  name       TEXT NOT NULL,
  uuid       UUID NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- rooms：参加コードで入室、進行状態を保持
CREATE TABLE IF NOT EXISTS rooms (
  id             BIGSERIAL PRIMARY KEY,
  code           TEXT    NOT NULL UNIQUE,                 -- 例: 6～8桁英数
  host_user_id   BIGINT  NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  status         room_status NOT NULL DEFAULT 'waiting',  -- waiting|started|ended
  current_round  INTEGER NOT NULL DEFAULT 0,
  created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
  started_at     TIMESTAMP,
  ended_at       TIMESTAMP
);

-- players：誰がどのルームに参加してるか。ラウンド中クリア状態のみ持つ
CREATE TABLE IF NOT EXISTS players (
  id           BIGSERIAL PRIMARY KEY,
  room_id      BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  has_cleared  BOOLEAN NOT NULL DEFAULT FALSE,            -- ラウンドごとにリセット運用
  joined_at    TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE (room_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_players_room ON players(room_id);

-- history：1ラウンド・1人の提出結果（スコア＆写真）
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
  result_json  JSONB,                 -- 任意: 判定詳細(ヒストグラム等)

  -- 画像（MVPはDB保存、将来は外部ストレージに移行可）
  photo_bytes  BYTEA,
  photo_mime   TEXT,
  photo_url    TEXT,

  created_at   TIMESTAMP NOT NULL DEFAULT NOW(),

  -- 1人1ラウンド1件（差し替えたいならON CONFLICTでUPDATE）
  UNIQUE (room_id, round_no, user_id)
);
CREATE INDEX IF NOT EXISTS idx_history_room_round ON history(room_id, round_no);
CREATE INDEX IF NOT EXISTS idx_history_user       ON history(user_id);