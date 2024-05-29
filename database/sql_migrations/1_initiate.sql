-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE patient (
    id BIGINT NOT NULL,
    name VARCHAR(256),
    password VARCHAR(256),
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE doctor (
    id BIGINT NOT NULL,
    name VARCHAR(256),
    password VARCHAR(256),
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE consultation (
    id BIGINT NOT NULL,
    meeting_link VARCHAR(256),
    patient_id BIGINT NOT NULL,
    doctor_id BIGINT NOT NULL,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE prescription (
    id BIGINT NOT NULL,
    content VARCHAR(256),
    patient_id BIGINT NOT NULL,
    doctor_id BIGINT NOT NULL,
    consultation_id BIGINT NOT NULL,
    created_at timestamp,
    updated_at timestamp
);

-- +migrate StatementEnd