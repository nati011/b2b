
CREATE OR REPLACE FUNCTION update_last_modified()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_modified = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_last_modified_triggers()
RETURNS void AS $$
DECLARE
    table_record RECORD;
BEGIN
    FOR table_record IN 
        SELECT tablename 
        FROM pg_tables 
        WHERE schemaname = 'public' 
        AND tablename != 'base'
    LOOP
        IF EXISTS (
            SELECT 1 
            FROM pg_inherits i
            JOIN pg_class parent ON i.inhparent = parent.oid
            JOIN pg_class child ON i.inhrelid = child.oid
            WHERE parent.relname = 'base'
            AND child.relname = table_record.tablename
        ) THEN
            EXECUTE format(
                'CREATE TRIGGER update_last_modified_trigger 
                BEFORE UPDATE ON public.%I 
                FOR EACH ROW 
                EXECUTE FUNCTION update_last_modified()',
                table_record.tablename
            );
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT create_last_modified_triggers();
