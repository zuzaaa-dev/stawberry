-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_inventories (
    product_id INT NOT NULL,
    shop_id INT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (product_id, shop_id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (shop_id) REFERENCES shops(id)
);

CREATE INDEX idx_shop_inventories_product_id ON shop_inventories(product_id);
CREATE INDEX idx_shop_inventories_shop_id ON shop_inventories(shop_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shop_inventories;
-- +goose StatementEnd
