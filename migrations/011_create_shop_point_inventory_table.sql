-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_point_inventories (
    shop_point_id INT NOT NULL,
    product_id INT NOT NULL,
    price INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (shop_point_id) REFERENCES shop_points(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX idx_shop_point_inventories_shop_point_id ON shop_point_inventories(shop_point_id);
CREATE INDEX idx_shop_point_inventories_product_id ON shop_point_inventories(product_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shop_point_inventories;
-- +goose StatementEnd
