-- Table: public.products

-- DROP TABLE IF EXISTS public.products;

CREATE TABLE IF NOT EXISTS public.products
(
    id integer NOT NULL DEFAULT 'nextval('products_id_seq'::regclass)',
    name text COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default",
    price numeric,
    rating numeric,
    comments text[] COLLATE pg_catalog."default",
    CONSTRAINT products_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.products
    OWNER to postgres;