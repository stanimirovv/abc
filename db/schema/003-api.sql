create extension pgcrypto;

CREATE TABLE apis(
    id SERIAL PRIMARY KEY,
    api_key TEXT NOT NULL DEFAULT md5(now()::text || (random() * 10000 + 1)::text),
    key_get_param_name TEXT NOT NULL default 'api_key', 
    is_https BOOLEAN NOT NULL DEFAULT false,
    https_cert_file TEXT
);
