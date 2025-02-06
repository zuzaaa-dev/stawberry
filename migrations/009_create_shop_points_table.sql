-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_points (
    id SERIAL PRIMARY KEY,
    shop_id INT NOT NULL,
    address TEXT NOT NULL,
    phone VARCHAR(30) NOT NULL,
    FOREIGN KEY (shop_id) REFERENCES shops(id)
);

CREATE INDEX idx_shop_points_shop_id ON shop_points(shop_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shop_points;
-- +goose StatementEnd
