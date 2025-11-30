CREATE SCHEMA box;
CREATE SCHEMA job;

CREATE TABLE IF NOT EXISTS box.boxes (
    id          VARCHAR(8)      PRIMARY KEY,
    owner_id    BIGINT          NOT NULL,
    archive_id  VARCHAR         NOT NULL,
    name        VARCHAR         NOT NULL,
    "desc"      VARCHAR                     DEFAULT NULL,
    vis         VISIBILITY_T    NOT NULL,
    created_at  TIMESTAMPTZ     NOT NULL    DEFAULT now(),
    deleted_at  TIMESTAMPTZ                 DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS box.input_fields (
    box_id   VARCHAR(8)         NOT NULL,
    index    INTEGER            NOT NULL,
    type     VALUE_TYPE_T       NOT NULL,
    name     VARCHAR            NOT NULL,
    "desc"   VARCHAR                        DEFAULT NULL,
    "unit"   VARCHAR                        DEFAULT NULL,

    PRIMARY KEY (box_id, index),

    FOREIGN KEY (box_id)
        REFERENCES box.boxes (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS box.output_fields (
    box_id   VARCHAR(8)         NOT NULL,
    index    INTEGER            NOT NULL,
    type     VALUE_TYPE_T       NOT NULL,
    name     VARCHAR            NOT NULL,
    "desc"   VARCHAR                        DEFAULT NULL,
    "unit"   VARCHAR                        DEFAULT NULL,

    PRIMARY KEY (box_id, index),

    FOREIGN KEY (box_id)
        REFERENCES box.boxes (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS job.jobs (
    id              VARCHAR(8)  PRIMARY KEY,
    box_id          VARCHAR(8)  NOT NULL,
    archive_id      VARCHAR     NOT NULL,
    owner_id        BIGINT      NOT NULL,
    state           JOB_STATE_T NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL    DEFAULT now(),

    started_at      TIMESTAMPTZ             DEFAULT NULL,
    result_code     INTEGER                 DEFAULT NULL,
    result_msg      VARCHAR                 DEFAULT NULL,
    finished_at     TIMESTAMPTZ             DEFAULT NULL,

    FOREIGN KEY (box_id)
        REFERENCES box.boxes (id)
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
