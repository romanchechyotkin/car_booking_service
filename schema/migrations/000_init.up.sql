--
-- PostgreSQL database dump
--

-- Dumped from database version 15.1
-- Dumped by pg_dump version 15.1

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: applications; Type: TABLE; Schema: public; Owner: chechyotka
--

CREATE TABLE public.applications (
    user_id uuid,
    filename text,
    is_visited boolean DEFAULT false
);



--
-- Name: car_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.car_images (
    id integer NOT NULL,
    url text,
    car_id character varying
);



--
-- Name: car_images_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.car_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- Name: car_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.car_images_id_seq OWNED BY public.car_images.id;


--
-- Name: cars; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cars (
    id character varying(8) NOT NULL,
    brand text,
    model text,
    year numeric(4,0),
    price_per_day numeric(10,2),
    is_available boolean DEFAULT true,
    rating numeric(3,2) DEFAULT 0.0
);



--
-- Name: cars_ratings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cars_ratings (
    rate double precision,
    comment text,
    car_id character varying,
    rate_by_user uuid,
    CONSTRAINT cars_ratings_rate_check CHECK (((rate <= (5.0)::double precision) AND (rate >= (1.0)::double precision)))
);



--
-- Name: cars_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cars_users (
    car_id character varying(8),
    user_id uuid
);



--
-- Name: reservations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reservations (
    id integer NOT NULL,
    customer_id uuid,
    seller_id uuid,
    car_id character varying(8),
    start_date timestamp without time zone,
    end_date timestamp without time zone,
    total_price numeric(10,2)
);



--
-- Name: reservations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.reservations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- Name: reservations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reservations_id_seq OWNED BY public.reservations.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    role character varying(5) DEFAULT 'USER'::character varying,
    user_id uuid,
    CONSTRAINT roles_role_check CHECK (((role)::text = ANY ((ARRAY['USER'::character varying, 'ADMIN'::character varying])::text[])))
);



--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    full_name text NOT NULL,
    telephone_number text NOT NULL,
    is_premium boolean DEFAULT false,
    city text,
    rating numeric(3,2) DEFAULT 0,
    is_verified boolean DEFAULT false
);



--
-- Name: users_ratings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_ratings (
    rate double precision,
    comment text,
    user_id uuid,
    rate_by_user uuid,
    CONSTRAINT ratings_rate_check CHECK (((rate <= (5.0)::double precision) AND (rate >= (1.0)::double precision)))
);



--
-- Name: car_images id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.car_images ALTER COLUMN id SET DEFAULT nextval('public.car_images_id_seq'::regclass);


--
-- Name: reservations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations ALTER COLUMN id SET DEFAULT nextval('public.reservations_id_seq'::regclass);


--
-- Name: car_images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chechyotka
--

SELECT pg_catalog.setval('public.car_images_id_seq', 41, true);


--
-- Name: reservations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chechyotka
--

SELECT pg_catalog.setval('public.reservations_id_seq', 48, true);


--
-- Name: applications applications_user_id_key; Type: CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.applications
    ADD CONSTRAINT applications_user_id_key UNIQUE (user_id);


--
-- Name: car_images car_images_pkey; Type: CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.car_images
    ADD CONSTRAINT car_images_pkey PRIMARY KEY (id);


--
-- Name: cars cars_pkey; Type: CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.cars
    ADD CONSTRAINT cars_pkey PRIMARY KEY (id);


--
-- Name: reservations reservations_pkey; Type: CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_telephone_number_key; Type: CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_telephone_number_key UNIQUE (telephone_number);


--
-- Name: applications applications_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.applications
    ADD CONSTRAINT applications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: car_images car_images_car_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.car_images
    ADD CONSTRAINT car_images_car_id_fkey FOREIGN KEY (car_id) REFERENCES public.cars(id);


--
-- Name: cars_ratings cars_ratings_car_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.cars_ratings
    ADD CONSTRAINT cars_ratings_car_id_fkey FOREIGN KEY (car_id) REFERENCES public.cars(id);


--
-- Name: cars_ratings cars_ratings_rate_by_user_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.cars_ratings
    ADD CONSTRAINT cars_ratings_rate_by_user_fkey FOREIGN KEY (rate_by_user) REFERENCES public.users(id);


--
-- Name: cars_users cars_users_car_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.cars_users
    ADD CONSTRAINT cars_users_car_id_fkey FOREIGN KEY (car_id) REFERENCES public.cars(id);


--
-- Name: cars_users cars_users_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.cars_users
    ADD CONSTRAINT cars_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: users_ratings ratings_rate_by_user_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.users_ratings
    ADD CONSTRAINT ratings_rate_by_user_fkey FOREIGN KEY (rate_by_user) REFERENCES public.users(id);


--
-- Name: users_ratings ratings_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.users_ratings
    ADD CONSTRAINT ratings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: reservations reservations_car_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_car_id_fkey FOREIGN KEY (car_id) REFERENCES public.cars(id);


--
-- Name: reservations reservations_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.users(id);


--
-- Name: reservations reservations_seller_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES public.users(id);


--
-- Name: roles roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

