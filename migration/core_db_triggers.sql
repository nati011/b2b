
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

CREATE OR REPLACE FUNCTION create_last_modified_trigger_for_new_table()
RETURNS event_trigger AS $$
DECLARE
    obj record;
BEGIN
    FOR obj IN SELECT * FROM pg_event_trigger_ddl_commands()
    WHERE command_tag = 'CREATE TABLE'
    LOOP
        IF EXISTS (
            SELECT 1 
            FROM pg_inherits i
            JOIN pg_class parent ON i.inhparent = parent.oid
            JOIN pg_class child ON i.inhrelid = child.oid
            WHERE parent.relname = 'base'
            AND child.relname = obj.object_identity::regclass::text
        ) THEN
            EXECUTE format(
                'CREATE TRIGGER update_last_modified_trigger 
                BEFORE UPDATE ON %s 
                FOR EACH ROW 
                EXECUTE FUNCTION update_last_modified()',
                obj.object_identity
            );
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;


CREATE EVENT TRIGGER create_last_modified_trigger_event 
ON ddl_command_end
WHEN TAG IN ('CREATE TABLE')
EXECUTE FUNCTION create_last_modified_trigger_for_new_table();