    -- alter table urls, if column old exist - rename, else add is_deleted flag
    DO $$
    BEGIN   
        IF EXISTS(
            SELECT *
            FROM information_schema.columns
            WHERE table_name='urls' and column_name='__is_deleted'
        )
        THEN
            ALTER TABLE urls
            RENAME COLUMN "__is_deleted" TO "is_deleted";
        ELSE
            ALTER TABLE urls
            ADD COLUMN is_deleted boolean;
        END IF;
END $$;