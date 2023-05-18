-- Table: public.commenting

-- DROP TABLE IF EXISTS public.commenting;

CREATE TABLE IF NOT EXISTS public.commenting
(
    product_id integer,
    comment character varying COLLATE pg_catalog."default",
    id integer NOT NULL DEFAULT 'nextval('commenting_id_seq'::regclass)',
    CONSTRAINT commenting_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.commenting
    OWNER to postgres;