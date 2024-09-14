-- +goose Up
-- +goose StatementBegin
ALTER TABLE SharedSecret
ADD CONSTRAINT unique_user_secret UNIQUE (shared_with_user, secret_id),
ADD CONSTRAINT unique_group_secret UNIQUE (shared_with_group, secret_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE SharedSecret
DROP CONSTRAINT IF EXISTS unique_user_secret,
DROP CONSTRAINT IF EXISTS unique_group_secret;
-- +goose StatementEnd
