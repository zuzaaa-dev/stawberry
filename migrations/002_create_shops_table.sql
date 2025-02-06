-- +goose Up
-- +goose StatementBegin
CREATE TABLE shops (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_shops_user_id ON shops(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shops;
-- +goose StatementEnd
