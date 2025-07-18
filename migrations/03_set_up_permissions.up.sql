GRANT SELECT ON script_fields TO app_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON
    fields,
    scripts,
    parameters,
    jobs,
    job_params
TO app_user;

GRANT USAGE ON SCHEMA public TO app_user;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;
