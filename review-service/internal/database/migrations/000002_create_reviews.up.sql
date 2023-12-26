CREATE TABLE IF NOT EXISTS reviews(
    id BIGSERIAL PRIMARY KEY,
    offer_id BIGINT NOT NULL,
    reviewer_id BIGINT NOT NULL, 
    offer_status INT,
    description VARCHAR(200),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_offer FOREIGN KEY(offer_id) REFERENCES offers(id)

);