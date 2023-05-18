-- Table: public.ratings

-- DROP TABLE IF EXISTS public.ratings;

CREATE TABLE IF NOT EXISTS public.ratings
(
    id integer,
    rating integer NOT NULL,
    product_id integer NOT NULL,
    user_id integer,
    created_at time with time zone
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.ratings
    OWNER to postgres;