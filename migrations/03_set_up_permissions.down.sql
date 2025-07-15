REVOKE USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public FROM app_user;

REVOKE USAGE ON SCHEMA public FROM app_user;

REVOKE SELECT ON script_fields FROM app_user;

REVOKE SELECT, INSERT, UPDATE, DELETE ON
    fields,
    scripts,
    parameters,
    jobs,
    job_params
FROM app_user;
