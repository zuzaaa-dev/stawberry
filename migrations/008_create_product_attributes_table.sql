-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_attributes (
    product_id INT NOT NULL,
    attributes JSONB NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX idx_product_attributes_product_id ON product_attributes(product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_attributes;
-- +goose StatementEnd
