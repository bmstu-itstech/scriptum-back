ALTER TABLE
    scripts
ADD
    CONSTRAINT unique_script_path UNIQUE(path);