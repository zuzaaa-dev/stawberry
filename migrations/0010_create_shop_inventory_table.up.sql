CREATE TABLE shop_inventory (
    product_id INT NOT NULL,
    shop_id INT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (product_id, shop_id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (shop_id) REFERENCES shops(id)
);

CREATE INDEX idx_shop_inventory_product_id ON shop_inventory(product_id);
CREATE INDEX idx_shop_inventory_shop_id ON shop_inventory(shop_id);
