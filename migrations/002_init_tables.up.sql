CREATE SCHEMA blueprint;
CREATE SCHEMA job;

CREATE TABLE IF NOT EXISTS blueprint.blueprints (
    id          VARCHAR(8)      PRIMARY KEY,
    owner_id    VARCHAR(8)      NOT NULL,
    archive_id  VARCHAR         NOT NULL,
    name        VARCHAR         NOT NULL,
    "desc"      VARCHAR                     DEFAULT NULL,
    vis         VISIBILITY_T    NOT NULL,
    created_at  TIMESTAMPTZ     NOT NULL    DEFAULT now(),
    deleted_at  TIMESTAMPTZ                 DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS blueprint.input_fields (
    blueprint_id VARCHAR(8)     NOT NULL,
    index        INTEGER        NOT NULL,
    type         VALUE_TYPE_T   NOT NULL,
    name         VARCHAR        NOT NULL,
    "desc"       VARCHAR                    DEFAULT NULL,
    "unit"       VARCHAR                    DEFAULT NULL,

    PRIMARY KEY (blueprint_id, index),

    FOREIGN KEY (blueprint_id)
        REFERENCES blueprint.blueprints (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS blueprint.output_fields (
    blueprint_id VARCHAR(8)     NOT NULL,
    index        INTEGER        NOT NULL,
    type         VALUE_TYPE_T   NOT NULL,
    name         VARCHAR        NOT NULL,
    "desc"       VARCHAR                    DEFAULT NULL,
    "unit"       VARCHAR                    DEFAULT NULL,

    PRIMARY KEY (blueprint_id, index),

    FOREIGN KEY (blueprint_id)
        REFERENCES blueprint.blueprints (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS job.jobs (
    id              VARCHAR(8)  PRIMARY KEY,
    blueprint_id    VARCHAR(8)  NOT NULL,
    archive_id      VARCHAR(8)  NOT NULL,
    owner_id        VARCHAR(8)  NOT NULL,
    state           JOB_STATE_T NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL    DEFAULT now(),

    started_at      TIMESTAMPTZ             DEFAULT NULL,
    result_code     INTEGER                 DEFAULT NULL,
    result_msg      VARCHAR                 DEFAULT NULL,
    finished_at     TIMESTAMPTZ             DEFAULT NULL,

    FOREIGN KEY (blueprint_id)
        REFERENCES blueprint.blueprints (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS job.input_values (
    job_id  VARCHAR(8)          NOT NULL,
    index   INTEGER             NOT NULL,
    type    VALUE_TYPE_T        NOT NULL,
    value   VARCHAR             NOT NULL,

    PRIMARY KEY (job_id, index),

    FOREIGN KEY (job_id)
        REFERENCES job.jobs (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS job.output_values (
    job_id  VARCHAR(8)          NOT NULL,
    index   INTEGER             NOT NULL,
    type    VALUE_TYPE_T        NOT NULL,
    value   VARCHAR             NOT NULL,

    PRIMARY KEY (job_id, index),

    FOREIGN KEY (job_id)
        REFERENCES job.jobs (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS job.output_fields (
    job_id   VARCHAR(8)         NOT NULL,
    index    INTEGER            NOT NULL,
    type     VALUE_TYPE_T       NOT NULL,
    name     VARCHAR            NOT NULL,
    "desc"   VARCHAR                        DEFAULT NULL,
    "unit"   VARCHAR                        DEFAULT NULL,

    PRIMARY KEY (job_id, index),

    FOREIGN KEY (job_id)
        REFERENCES job.jobs (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users (
    id          VARCHAR(8)  NOT NULL,
    email       VARCHAR     NOT NULL    UNIQUE,
    name        VARCHAR     NOT NULL,
    role        ROLE_T      NOT NULL,
    passhash    VARCHAR     NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL    DEFAULT now(),
    deleted_at  TIMESTAMPTZ             DEFAULT NULL
);
