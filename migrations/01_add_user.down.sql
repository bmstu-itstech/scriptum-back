DO $ $ BEGIN IF EXISTS (
    SELECT
    FROM
        pg_catalog.pg_roles
    WHERE
        rolname = 'app_user'
) THEN DROP ROLE app_user;

END IF;

END $ $;