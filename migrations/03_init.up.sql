DO $$ BEGIN
    CREATE TYPE VISIBILITY AS ENUM ('public', 'private');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
    CREATE TYPE FIELD_TYPE AS ENUM ('integer', 'real', 'complex');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
    CREATE TYPE PARAM_TYPE AS ENUM ('in', 'out');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

CREATE TABLE fields (
    field_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CHECK (LENGTH(name) <= 100),
    description TEXT NOT NULL CHECK (LENGTH(description) <= 500),
    unit TEXT NOT NULL CHECK (LENGTH(unit) <= 20),
    field_type FIELD_TYPE NOT NULL,
    param PARAM_TYPE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE scripts (
    script_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CHECK (LENGTH(name) <= 100),        
    description TEXT CHECK (LENGTH(description) <= 500),   А
    path TEXT NOT NULL CHECK (LENGTH(path) <= 200),
    visibility VISIBILITY NOT NULL,
    owner_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE script_fields (
    script_id BIGINT NOT NULL,
    field_id BIGINT NOT NULL,
    PRIMARY KEY (script_id, field_id),
    CONSTRAINT fk_script_fields_script FOREIGN KEY (script_id) REFERENCES scripts(script_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_script_fields_field FOREIGN KEY (field_id) REFERENCES fields(field_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE parameters (
    parameter_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    field_id BIGINT NOT NULL,
    value TEXT NOT NULL CHECK (LENGTH(value) <= 100),
    CONSTRAINT fk_parameters_field FOREIGN KEY (field_id) REFERENCES fields(field_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE jobs (
    job_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status_code INT,
    error_message TEXT,
    script_id BIGINT NOT NULL,
    CONSTRAINT fk_jobs_script FOREIGN KEY (script_id) REFERENCES scripts(script_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE job_params (
    job_id BIGINT NOT NULL,
    parameter_id BIGINT NOT NULL,
    PRIMARY KEY (job_id, parameter_id),
    CONSTRAINT fk_job_params_job FOREIGN KEY (job_id) REFERENCES jobs(job_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_job_params_parameter FOREIGN KEY (parameter_id) REFERENCES parameters(parameter_id) ON DELETE CASCADE ON UPDATE CASCADE
);
