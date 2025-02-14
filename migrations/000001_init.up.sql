CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
BEGIN
    CREATE USER admin WITH PASSWORD 'password' LOGIN;
    EXCEPTION WHEN DUPLICATE_OBJECT THEN RAISE NOTICE 'not creating role admin -- it already exists';
END
$$;

DO $$
BEGIN 
    CREATE USER readonly WITH PASSWORD 'password' LOGIN;
    EXCEPTION WHEN DUPLICATE_OBJECT THEN RAISE NOTICE 'not creating role readonly -- it already exists';
END
$$;