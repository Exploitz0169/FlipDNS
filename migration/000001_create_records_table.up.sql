BEGIN;

CREATE TYPE record_type as ENUM(
    'A',
    'AAAA'
);

COMMIT;

CREATE TABLE IF NOT EXISTS record (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain_name VARCHAR(255) NOT NULL,
    record_data VARCHAR(255) NOT NULL,
    record_type record_type NOT NULL,
    ttl INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()  
);

CREATE INDEX records_domain_name_idx ON record (domain_name);
