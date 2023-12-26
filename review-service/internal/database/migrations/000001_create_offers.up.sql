CREATE TABLE IF NOT EXISTS offers(
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    price BIGINT,
    price_deal BOOLEAN NOT NULL,
    isbn VARCHAR(50) NOT NULL,
    publisher VARCHAR(100),
    edition INTEGER,
    description VARCHAR(512),
    status SMALLINT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS offer_images(
    id BIGSERIAL PRIMARY KEY,
    offer_id INT NOT NULL,
    image_url VARCHAR(256) NOT NULL,
    CONSTRAINT fk_offer
        FOREIGN KEY(offer_id)
            REFERENCES offers(id)
);

CREATE INDEX idx_offer_images_offer_id ON offer_images(offer_id);
