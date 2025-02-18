-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INT,
    description TEXT,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE INDEX idx_products_category_id ON products(category_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
