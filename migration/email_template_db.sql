CREATE TABLE IF NOT EXISTS public."templates"
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    html BYTEA,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

COMMENT ON TABLE public."templates" IS 'stores email templates.';

-- procedures----------------------------
CREATE OR REPLACE PROCEDURE writer(
   name int,
   html BYTEA
)
language plpgsql    
as $$
begin
   
   -- create teamplate
    INSERT INTO public."templates"
    Values(name, html)
    commit;
end;$$; 
----------------------------------------
CREATE OR REPLACE PROCEDURE reader(
   name int
)
language plpgsql    
as $$
begin
   
   -- read teamplate
    SELECT * FROM public."templates" as t
    WHERE t.name = name 
    commit;
end;$$; 
------------------------------------------