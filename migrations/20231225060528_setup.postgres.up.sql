CREATE TABLE IF NOT EXISTS public.users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.contacts (
    id SERIAL PRIMARY KEY,
    user_id  INT NOT NULL,
    nama VARCHAR(255) NOT NULL,
    nomor_telepon VARCHAR(25),
    email VARCHAR(255),
    alamat VARCHAR(255),
    CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
