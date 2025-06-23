BEGIN;


ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role user_role DEFAULT 'user';


UPDATE users
SET role = 'user'
WHERE role IS NULL;

COMMIT;
