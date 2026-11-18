-- name: UpsertDirectoryHandle :exec
INSERT INTO
  directory.handles (target_id, target_type, handle)
VALUES
  ($1, $2, $3)
ON CONFLICT (target_id, target_type) DO UPDATE
SET
  handle = excluded.handle
;

-- name: DeleteDirectoryHandle :exec
DELETE FROM directory.handles
WHERE
  target_id = $1
  AND target_type = $2
;

-- name: UpsertUserProfile :exec
INSERT INTO
  directory.user_profiles (user_id, display_name, avatar_id, bio, birthday)
VALUES
  ($1, $2, $3, $4, $5)
ON CONFLICT (user_id) DO UPDATE
SET
  display_name = excluded.display_name,
  avatar_id = excluded.avatar_id,
  bio = excluded.bio,
  birthday = excluded.birthday
;
