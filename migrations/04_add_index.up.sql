CREATE INDEX idx_parameters_value ON parameters(value);

CREATE INDEX idx_scripts_path ON scripts(path);

CREATE INDEX idx_scripts_owner_id ON scripts(owner_id);

CREATE INDEX idx_jobs_script_id ON jobs(script_id);

CREATE INDEX idx_job_params_parameter_id ON job_params(parameter_id);