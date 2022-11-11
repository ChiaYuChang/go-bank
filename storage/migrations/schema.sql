--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5
-- Dumped by pg_dump version 14.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: tstatus; Type: TYPE; Schema: public; Owner: onlyfortest
--

CREATE TYPE public.tstatus AS ENUM (
    'created',
    'success',
    'failure'
);


ALTER TYPE public.tstatus OWNER TO onlyfortest;

--
-- Name: trigger_set_timestamp(); Type: FUNCTION; Schema: public; Owner: onlyfortest
--

CREATE FUNCTION public.trigger_set_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
    END;
$$;


ALTER FUNCTION public.trigger_set_timestamp() OWNER TO onlyfortest;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: onlyfortest
--

CREATE TABLE public.accounts (
    id bigint NOT NULL,
    owner character varying(30) NOT NULL,
    balance numeric DEFAULT 0 NOT NULL,
    currency integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.accounts OWNER TO onlyfortest;

--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: onlyfortest
--

CREATE SEQUENCE public.accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.accounts_id_seq OWNER TO onlyfortest;

--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: onlyfortest
--

ALTER SEQUENCE public.accounts_id_seq OWNED BY public.accounts.id;


--
-- Name: currencies; Type: TABLE; Schema: public; Owner: onlyfortest
--

CREATE TABLE public.currencies (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    abbr character varying(5) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.currencies OWNER TO onlyfortest;

--
-- Name: currencies_id_seq; Type: SEQUENCE; Schema: public; Owner: onlyfortest
--

CREATE SEQUENCE public.currencies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.currencies_id_seq OWNER TO onlyfortest;

--
-- Name: currencies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: onlyfortest
--

ALTER SEQUENCE public.currencies_id_seq OWNED BY public.currencies.id;


--
-- Name: entries; Type: TABLE; Schema: public; Owner: onlyfortest
--

CREATE TABLE public.entries (
    id bigint NOT NULL,
    account_id bigint NOT NULL,
    amount numeric DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.entries OWNER TO onlyfortest;

--
-- Name: COLUMN entries.amount; Type: COMMENT; Schema: public; Owner: onlyfortest
--

COMMENT ON COLUMN public.entries.amount IS 'could be any number';


--
-- Name: entries_id_seq; Type: SEQUENCE; Schema: public; Owner: onlyfortest
--

CREATE SEQUENCE public.entries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.entries_id_seq OWNER TO onlyfortest;

--
-- Name: entries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: onlyfortest
--

ALTER SEQUENCE public.entries_id_seq OWNED BY public.entries.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: onlyfortest
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO onlyfortest;

--
-- Name: transfers; Type: TABLE; Schema: public; Owner: onlyfortest
--

CREATE TABLE public.transfers (
    id bigint NOT NULL,
    src_id bigint NOT NULL,
    dst_id bigint NOT NULL,
    amount numeric DEFAULT 0 NOT NULL,
    status public.tstatus DEFAULT 'created'::public.tstatus NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.transfers OWNER TO onlyfortest;

--
-- Name: COLUMN transfers.amount; Type: COMMENT; Schema: public; Owner: onlyfortest
--

COMMENT ON COLUMN public.transfers.amount IS 'must gte 0';


--
-- Name: transfers_id_seq; Type: SEQUENCE; Schema: public; Owner: onlyfortest
--

CREATE SEQUENCE public.transfers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.transfers_id_seq OWNER TO onlyfortest;

--
-- Name: transfers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: onlyfortest
--

ALTER SEQUENCE public.transfers_id_seq OWNED BY public.transfers.id;


--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.accounts ALTER COLUMN id SET DEFAULT nextval('public.accounts_id_seq'::regclass);


--
-- Name: currencies id; Type: DEFAULT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.currencies ALTER COLUMN id SET DEFAULT nextval('public.currencies_id_seq'::regclass);


--
-- Name: entries id; Type: DEFAULT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.entries ALTER COLUMN id SET DEFAULT nextval('public.entries_id_seq'::regclass);


--
-- Name: transfers id; Type: DEFAULT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.transfers ALTER COLUMN id SET DEFAULT nextval('public.transfers_id_seq'::regclass);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: currencies currencies_pkey; Type: CONSTRAINT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.currencies
    ADD CONSTRAINT currencies_pkey PRIMARY KEY (id);


--
-- Name: entries entries_pkey; Type: CONSTRAINT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.entries
    ADD CONSTRAINT entries_pkey PRIMARY KEY (id);


--
-- Name: schema_migration schema_migration_pkey; Type: CONSTRAINT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.schema_migration
    ADD CONSTRAINT schema_migration_pkey PRIMARY KEY (version);


--
-- Name: transfers transfers_pkey; Type: CONSTRAINT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.transfers
    ADD CONSTRAINT transfers_pkey PRIMARY KEY (id);


--
-- Name: accounts_owner_idx; Type: INDEX; Schema: public; Owner: onlyfortest
--

CREATE INDEX accounts_owner_idx ON public.accounts USING btree (owner);


--
-- Name: currencies_abbr_idx; Type: INDEX; Schema: public; Owner: onlyfortest
--

CREATE INDEX currencies_abbr_idx ON public.currencies USING btree (abbr);


--
-- Name: entries_account_id_idx; Type: INDEX; Schema: public; Owner: onlyfortest
--

CREATE INDEX entries_account_id_idx ON public.entries USING btree (account_id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: onlyfortest
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: transfers_dst_id_idx; Type: INDEX; Schema: public; Owner: onlyfortest
--

CREATE INDEX transfers_dst_id_idx ON public.transfers USING btree (dst_id);


--
-- Name: transfers_src_id_dst_id_idx; Type: INDEX; Schema: public; Owner: onlyfortest
--

CREATE INDEX transfers_src_id_dst_id_idx ON public.transfers USING btree (src_id, dst_id);


--
-- Name: transfers_src_id_idx; Type: INDEX; Schema: public; Owner: onlyfortest
--

CREATE INDEX transfers_src_id_idx ON public.transfers USING btree (src_id);


--
-- Name: accounts set_timestamp; Type: TRIGGER; Schema: public; Owner: onlyfortest
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON public.accounts FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: currencies set_timestamp; Type: TRIGGER; Schema: public; Owner: onlyfortest
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON public.currencies FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: accounts accounts_currency_fkey; Type: FK CONSTRAINT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_currency_fkey FOREIGN KEY (currency) REFERENCES public.currencies(id);


--
-- Name: entries entries_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.entries
    ADD CONSTRAINT entries_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.accounts(id);


--
-- Name: transfers transfers_dst_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.transfers
    ADD CONSTRAINT transfers_dst_id_fkey FOREIGN KEY (dst_id) REFERENCES public.accounts(id);


--
-- Name: transfers transfers_src_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: onlyfortest
--

ALTER TABLE ONLY public.transfers
    ADD CONSTRAINT transfers_src_id_fkey FOREIGN KEY (src_id) REFERENCES public.accounts(id);


--
-- PostgreSQL database dump complete
--

