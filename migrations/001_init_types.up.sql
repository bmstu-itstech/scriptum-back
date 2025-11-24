DO $$ BEGIN
    CREATE TYPE VISIBILITY_T
    AS ENUM (
        'public',
        'private'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE JOB_STATE_T
    AS ENUM (
        'pending',
        'running',
        'finished'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE VALUE_TYPE_T
    AS ENUM (
        'integer',
        'real',
        'string'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
