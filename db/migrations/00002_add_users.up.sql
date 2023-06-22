BEGIN TRANSACTION ;

    /* 
        add table users
        no need index because creating pk - add b-tree index auto
    */
	CREATE TABLE IF NOT EXISTS users
	(
		id bigint NOT NULL,
		CONSTRAINT users_pkey PRIMARY KEY (id)
	);

    ALTER TABLE urls
    ADD COLUMN user_id bigint
    REFERENCES users(id);

COMMIT ;

--     -- alter table urls, if column old exist - rename, else add fk for users by user.id
--     DO $$
--     BEGIN   
--         IF EXISTS(
--             SELECT *
--             FROM information_schema.columns
--             WHERE table_name='urls' and column_name='__user_id'
--         )
--         THEN
--             ALTER TABLE urls
--             RENAME COLUMN "__user_id" TO "user_id";
--         ELSE
--             ALTER TABLE urls
--             ADD COLUMN user_id bigserial
--             REFERENCES users(id);
--         END IF;
-- END $$;