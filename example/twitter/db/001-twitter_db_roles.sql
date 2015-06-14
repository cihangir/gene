
-- Drop role
DROP ROLE IF EXISTS "twitter_db_role";
-- Create role
CREATE ROLE "twitter_db_role";
-- Create shadow user for future extensibility
DROP USER IF EXISTS "twitter_db_roleapplication";
CREATE USER "twitter_db_roleapplication" PASSWORD 'twitter_db_roleapplication';
GRANT "twitter_db_role" TO "twitter_db_roleapplication";