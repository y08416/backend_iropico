-- 依存順に落とす
DROP INDEX IF EXISTS idx_history_room_round_user;
DROP INDEX IF EXISTS idx_history_user;
DROP INDEX IF EXISTS idx_history_room_round;
DROP INDEX IF EXISTS idx_players_room;
DROP INDEX IF EXISTS uq_rooms_code;

DROP TABLE IF EXISTS history;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS room_status;