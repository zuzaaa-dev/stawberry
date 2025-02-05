CREATE TABLE image_keys (
    image_key VARCHAR(50) UNIQUE NOT NULL,
    product_id INT NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX idx_image_keys_product_id ON image_keys(product_id);