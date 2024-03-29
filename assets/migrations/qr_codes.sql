CREATE TABLE qr_codes (
    qrcode_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    data TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);