DO $ $ BEGIN IF NOT EXISTS (
   SELECT
   FROM
      pg_catalog.pg_roles
   WHERE
      rolname = 'app_user'
) THEN CREATE ROLE app_user WITH LOGIN PASSWORD 'your_secure_password';

END IF;

END $ $;