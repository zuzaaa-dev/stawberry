CREATE TABLE shop_point_inventory (
    shop_point_id INT NOT NULL,
    product_id INT NOT NULL,
    price INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (shop_point_id) REFERENCES shop_points(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX idx_shop_point_inventory_shop_point_id ON shop_point_inventory(shop_point_id);
CREATE INDEX idx_shop_point_inventory_product_id ON shop_point_inventory(product_id);
