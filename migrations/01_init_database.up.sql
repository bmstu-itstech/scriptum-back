DO $ $ BEGIN IF NOT EXISTS (
   SELECT
   FROM
      pg_database
   WHERE
      datname = 'dev'
) THEN CREATE DATABASE dev;

END IF;

END $ $;