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


ALTER TABLE public.applications OWNER TO chechyotka;

--
-- Name: car_images; Type: TABLE; Schema: public; Owner: chechyotka
--

CREATE TABLE public.car_images (
    id integer NOT NULL,
    url text,
    car_id character varying
);


ALTER TABLE public.car_images OWNER TO chechyotka;

--
-- Name: car_images_id_seq; Type: SEQUENCE; Schema: public; Owner: chechyotka
--

CREATE SEQUENCE public.car_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.car_images_id_seq OWNER TO chechyotka;

--
-- Name: car_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chechyotka
--

ALTER SEQUENCE public.car_images_id_seq OWNED BY public.car_images.id;


--
-- Name: cars; Type: TABLE; Schema: public; Owner: chechyotka
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


ALTER TABLE public.cars OWNER TO chechyotka;

--
-- Name: cars_ratings; Type: TABLE; Schema: public; Owner: chechyotka
--

CREATE TABLE public.cars_ratings (
    rate double precision,
    comment text,
    car_id character varying,
    rate_by_user uuid,
    CONSTRAINT cars_ratings_rate_check CHECK (((rate <= (5.0)::double precision) AND (rate >= (1.0)::double precision)))
);


ALTER TABLE public.cars_ratings OWNER TO chechyotka;

--
-- Name: cars_users; Type: TABLE; Schema: public; Owner: chechyotka
--

CREATE TABLE public.cars_users (
    car_id character varying(8),
    user_id uuid
);


ALTER TABLE public.cars_users OWNER TO chechyotka;

--
-- Name: reservations; Type: TABLE; Schema: public; Owner: chechyotka
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


ALTER TABLE public.reservations OWNER TO chechyotka;

--
-- Name: reservations_id_seq; Type: SEQUENCE; Schema: public; Owner: chechyotka
--

CREATE SEQUENCE public.reservations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reservations_id_seq OWNER TO chechyotka;

--
-- Name: reservations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chechyotka
--

ALTER SEQUENCE public.reservations_id_seq OWNED BY public.reservations.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: chechyotka
--

CREATE TABLE public.roles (
    role character varying(5) DEFAULT 'USER'::character varying,
    user_id uuid,
    CONSTRAINT roles_role_check CHECK (((role)::text = ANY ((ARRAY['USER'::character varying, 'ADMIN'::character varying])::text[])))
);


ALTER TABLE public.roles OWNER TO chechyotka;

--
-- Name: users; Type: TABLE; Schema: public; Owner: chechyotka
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


ALTER TABLE public.users OWNER TO chechyotka;

--
-- Name: users_ratings; Type: TABLE; Schema: public; Owner: chechyotka
--

CREATE TABLE public.users_ratings (
    rate double precision,
    comment text,
    user_id uuid,
    rate_by_user uuid,
    CONSTRAINT ratings_rate_check CHECK (((rate <= (5.0)::double precision) AND (rate >= (1.0)::double precision)))
);


ALTER TABLE public.users_ratings OWNER TO chechyotka;

--
-- Name: car_images id; Type: DEFAULT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.car_images ALTER COLUMN id SET DEFAULT nextval('public.car_images_id_seq'::regclass);


--
-- Name: reservations id; Type: DEFAULT; Schema: public; Owner: chechyotka
--

ALTER TABLE ONLY public.reservations ALTER COLUMN id SET DEFAULT nextval('public.reservations_id_seq'::regclass);


--
-- Data for Name: applications; Type: TABLE DATA; Schema: public; Owner: chechyotka
--

COPY public.applications (user_id, filename, is_visited) FROM stdin;
e6f8015d-de60-4179-a24c-5da3e243c28a	7c174963-8024-4985-ac22-e7b5bfd5e45a	f
fcadd4ce-ea3f-4e2c-acfa-5e589edc9065	89cbd6b1-9326-420c-91d1-6357ffcce340	t
d011af53-5036-4512-a5c8-add6647ab62a	589d4a13-63e4-4313-8fe2-7428213b3f41	t
a86167b2-3ee7-416d-94bd-b17553a4f133	c10920c2-536c-4356-a79c-3a0f31cebdef	t
\.


--
-- Data for Name: car_images; Type: TABLE DATA; Schema: public; Owner: chechyotka
--

COPY public.car_images (id, url, car_id) FROM stdin;
13	6175e26d-0bd5-4e61-929f-242e1087c71c	1919QQ-3
14	2792b4f1-c5c2-42a0-bee7-227c9914d2f5	1919QQ-1
15	c61e7f9d-a60d-40c6-bb42-ab223714f571	1919QQ-1
16	1f1fa0e4-fc56-49d7-953e-a993c138c59e	1919QQ-1
17	8590b5d7-f26c-48ce-8294-85ce6a552eca	1119QQ-1
18	c1fe9bdf-bb3a-493c-9772-019dc612cd9a	1119QQ-1
19	589f98e2-17ff-4e9c-8f89-c581bd28d5cf	1119QQ-1
24	f25a50cb-f725-4d43-abe3-33ec222803ad	1219QQ-1
25	d07cd69e-e90f-45d3-96de-725c9f4b0bb1	1219QQ-1
26	779560dd-844d-49b7-891a-4e21c4a1e2f1	1219QQ-1
27	8984d977-664f-4505-9920-c0593925c4df	9999OO-7
28	e978743a-3168-4818-9bd7-0392feba4648	9999OO-1
29	3bfa3dd2-5766-4436-9587-30b51a07c0d8	9999OO-1
30	a66a080d-b26f-48fc-b840-797700f29e3c	9999OO-1
31	0d411534-1209-465c-957e-bf851016194d	9999OO-1
32	54edb5e1-6964-4081-9db0-7c2ff06958ab	9999OO-1
33	4b0e9efc-d9f6-4c3e-b506-c7e0d2a80188	9999OO-1
34	63ffcd13-1a58-410d-acde-6e60eb8e9aad	9999OO-1
35	b76008e9-f841-4a19-8ed0-ccc31b4f0b51	9999OO-1
36	695b2176-4f99-425d-9c12-bc337b1110cf	9999OO-1
37	36f82ab6-fa30-418a-b911-076dc28eb710	9999OO-1
38	a0924ba9-89ea-4636-bbac-c90cbc2333df	9999OO-2
39	1afd8dbc-6f07-40df-a768-350c2ee31fde	6351FR-1
40	b79928d1-65e2-464e-8ee5-4827859b3a18	1234YY-4
41	9fff9114-bcd1-46ae-9fa1-d7d3a420f0ad	1223PQ-1
\.


--
-- Data for Name: cars; Type: TABLE DATA; Schema: public; Owner: chechyotka
--

COPY public.cars (id, brand, model, year, price_per_day, is_available, rating) FROM stdin;
1111AA-3	BMW	X6	2023	300.10	t	0.00
1111XX-7	BMW	X99	1233	234.00	t	0.00
1111XX-4	BMW	X99	1233	234.00	t	0.00
1111XX-1	BMW	X99	1233	0.00	t	0.00
1234QW-1	Ferrari	QWERTY	2023	111111.00	t	0.00
1241QW-1	Ferrari	QWERTY	2023	111111.00	t	0.00
9999QQ-7	LAMBA	LOLOLO	2023	0.00	t	0.00
9919QQ-7	LAMBA	LOLOLO	2023	0.00	t	0.00
1919QQ-7	LAMBA	LOLOLO	2023	0.00	t	0.00
1919QQ-3	LAMBA	LOLOLO	2023	0.00	t	0.00
1919QQ-1	LAMBA	LOLOLO	0	122.10	t	0.00
1119QQ-1	LAMBA	LOLOLO	2012	122.10	t	0.00
1219QQ-1	LA	L	2012	122.10	t	0.00
9999OO-7	TEST	test	2023	123.00	t	0.00
9999OO-1	TEST1	test1	2023	1232.00	t	0.00
6351FR-7	airbnb	Frida	0	0.00	t	0.00
6351FR-1	airbnb	Frida	0	15541.00	f	3.67
9999OO-2	TEST1	test1	2023	1232.00	f	0.00
1234YY-7	Audi	Q7	2023	123.20	t	0.00
1234YY-6	Audi	Q7	2023	123.20	t	0.00
1234YY-5	Audi	Q7	2023	123.20	t	0.00
1234YY-4	Audi	Q7	2023	123.20	t	0.00
1223PQ-1	test	TEST	2012	21312.20	f	4.00
\.


--
-- Data for Name: cars_ratings; Type: TABLE DATA; Schema: public; Owner: chechyotka
--

COPY public.cars_ratings (rate, comment, car_id, rate_by_user) FROM stdin;
4	cool car	6351FR-1	e6f8015d-de60-4179-a24c-5da3e243c28a
4	cool car	6351FR-1	e6f8015d-de60-4179-a24c-5da3e243c28a
3		6351FR-1	a86167b2-3ee7-416d-94bd-b17553a4f133
4	Great Car, it was awesome	1223PQ-1	d011af53-5036-4512-a5c8-add6647ab62a
\.


--
-- Data for Name: cars_users; Type: TABLE DATA; Schema: public; Owner: chechyotka
--

COPY public.cars_users (car_id, user_id) FROM stdin;
1219QQ-1	fcadd4ce-ea3f-4e2c-acfa-5e589edc9065
9999OO-7	fcadd4ce-ea3f-4e2c-acfa-5e589edc9065
9999OO-1	e6f8015d-de60-4179-a24c-5da3e243c28a
9999OO-2	e6f8015d-de60-4179-a24c-5da3e243c28a
6351FR-7	e6f8015d-de60-4179-a24c-5da3e243c28a
6351FR-1	e6f8015d-de60-4179-a24c-5da3e243c28a
1234YY-7	d011af53-5036-4512-a5c8-add6647ab62a
1234YY-6	d011af53-5036-4512-a5c8-add6647ab62a
1234YY-5	d011af53-5036-4512-a5c8-add6647ab62a
1234YY-4	d011af53-5036-4512-a5c8-add6647ab62a
1223PQ-1	d011af53-5036-4512-a5c8-add6647ab62a
\.


--
-- Data for Name: reservations; Type: TABLE DATA; Schema: public; Owner: chechyotka
--

COPY public.reservations (id, customer_id, seller_id, car_id, start_date, end_date, total_price) FROM stdin;
1	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-7	2023-04-26 00:00:00	2023-04-28 00:00:00	0.00
2	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-7	2023-04-26 00:00:00	2023-04-28 00:00:00	0.00
3	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-7	2023-04-26 00:00:00	2023-04-28 00:00:00	0.00
4	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-1	2023-04-26 00:00:00	2023-04-27 00:00:00	15541.00
5	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-1	2023-04-26 00:00:00	2023-04-27 00:00:00	15541.00
6	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-1	2023-04-26 00:00:00	2023-04-27 00:00:00	15541.00
7	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-1	2023-04-27 00:00:00	2023-04-28 00:00:00	15541.00
8	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-1	2023-04-27 00:00:00	2023-04-28 00:00:00	15541.00
9	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-1	2023-04-27 00:00:00	2023-04-28 00:00:00	15541.00
10	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-1	2023-04-28 00:00:00	2023-04-29 00:00:00	15541.00
11	a86167b2-3ee7-416d-94bd-b17553a4f133	e6f8015d-de60-4179-a24c-5da3e243c28a	6351FR-1	2023-04-24 00:00:00	2023-04-26 00:00:00	31082.00
38	6037703d-9649-4ba8-ba6f-ab19a5633f4b	e6f8015d-de60-4179-a24c-5da3e243c28a	9999OO-2	2020-12-09 00:00:00	2020-12-19 00:00:00	12320.00
39	6037703d-9649-4ba8-ba6f-ab19a5633f4b	e6f8015d-de60-4179-a24c-5da3e243c28a	9999OO-2	2020-12-09 15:00:00	2020-12-19 16:00:00	12371.33
40	6037703d-9649-4ba8-ba6f-ab19a5633f4b	e6f8015d-de60-4179-a24c-5da3e243c28a	9999OO-2	2020-12-09 16:00:00	2020-12-19 17:00:00	12371.33
41	6037703d-9649-4ba8-ba6f-ab19a5633f4b	e6f8015d-de60-4179-a24c-5da3e243c28a	9999OO-2	2020-12-09 17:00:00	2020-12-19 18:00:00	12371.00
42	6037703d-9649-4ba8-ba6f-ab19a5633f4b	e6f8015d-de60-4179-a24c-5da3e243c28a	9999OO-2	2022-12-10 18:00:00	2022-12-19 18:00:00	11088.00
43	6037703d-9649-4ba8-ba6f-ab19a5633f4b	e6f8015d-de60-4179-a24c-5da3e243c28a	9999OO-2	2020-12-10 18:00:00	2020-12-19 18:00:00	11088.00
44	6037703d-9649-4ba8-ba6f-ab19a5633f4b	e6f8015d-de60-4179-a24c-5da3e243c28a	9999OO-2	2020-12-10 18:00:00	2020-12-19 18:00:00	11088.00
46	d011af53-5036-4512-a5c8-add6647ab62a	e6f8015d-de60-4179-a24c-5da3e243c28a	9999OO-2	2023-12-10 10:00:00	2023-12-31 19:00:00	26334.00
47	d011af53-5036-4512-a5c8-add6647ab62a	e6f8015d-de60-4179-a24c-5da3e243c28a	9999OO-2	2023-12-10 10:00:00	2023-12-31 19:00:00	26334.00
48	d011af53-5036-4512-a5c8-add6647ab62a	d011af53-5036-4512-a5c8-add6647ab62a	1223PQ-1	2023-05-12 23:00:00	2023-05-13 12:00:00	11544.00
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: chechyotka
--

COPY public.roles (role, user_id) FROM stdin;
ADMIN	a86167b2-3ee7-416d-94bd-b17553a4f133
USER	6037703d-9649-4ba8-ba6f-ab19a5633f4b
USER	e2e406e0-9a57-4cfd-b0f5-301cd7e04156
USER	d011af53-5036-4512-a5c8-add6647ab62a
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: chechyotka
--

COPY public.users (id, email, password, full_name, telephone_number, is_premium, city, rating, is_verified) FROM stdin;
c5cc18be-5483-45bd-8beb-f02ba4174d52	senya@gmail.com	$2a$10$mY4m..K1dd6SHZyWNaQY9uJFYCJhvCfs8rljgIpNelr9w4RWSrtT2	Roman	+375444535067	f	\N	0.00	f
084708a5-b456-4667-8947-7719504bc2f1	senya2004alishevich@gmail.com	$2a$10$rtrkvAO/5u5o7VUiavTBz.qNDivI4rcuz1dmI7XN8DZ5qk6.pW8G.	Roman	+375444535867	f	\N	0.00	f
e6f8015d-de60-4179-a24c-5da3e243c28a	cecetkin99@gmail.com	$2a$10$vi1B8sJhWRkg4Fb9O9OX6e.Rf6cSv9MRynjwN1u/f9/CNPe.RODvS	Roman	+375444125867	f	\N	1.50	t
6037703d-9649-4ba8-ba6f-ab19a5633f4b	test123@gmail.com	$2a$10$83bcogq.AH3TQ4RACMdQQ.SUbHsn9g8z1gRcAWXXM/RcJKvRgxAEm	qweqw2eqweqwe	+375440908009	f	\N	0.00	t
fcadd4ce-ea3f-4e2c-acfa-5e589edc9065	test@gmail.com	$2a$10$B13UW/GIp.TIYXAk7qn0P.6G3XUzaN3B0cOpd1488GCJFpetuHn0q	Test	+375447534067	f	\N	5.00	t
e2e406e0-9a57-4cfd-b0f5-301cd7e04156	romac@gmail.com	$2a$10$ZMkfuAvynfL1Q2z7fd2k6eBe3zY3wV9pkTkoW1IhOF2/a7qOE/k6G	string	+375449934067	f	\N	0.00	f
d011af53-5036-4512-a5c8-add6647ab62a	romac2@gmail.com	$2a$10$b8LKxoCIzzZ4ZHFrZ7hTMu24ij0obsNSPG9gPKnPIrS7KQE/krN82	qwerty	+375448934067	f	\N	2.00	t
a86167b2-3ee7-416d-94bd-b17553a4f133	romanchechyotkin@gmail.com	$2a$10$aF1V54AfPHporMxXZtpuSeClHvBHOcXjq5ucVQvfHfHHYa9pTrz6m	qwrfqaf	+375444534067	f	\N	0.00	t
\.


--
-- Data for Name: users_ratings; Type: TABLE DATA; Schema: public; Owner: chechyotka
--

COPY public.users_ratings (rate, comment, user_id, rate_by_user) FROM stdin;
1	loh	e6f8015d-de60-4179-a24c-5da3e243c28a	a86167b2-3ee7-416d-94bd-b17553a4f133
2		e6f8015d-de60-4179-a24c-5da3e243c28a	fcadd4ce-ea3f-4e2c-acfa-5e589edc9065
5	great person	fcadd4ce-ea3f-4e2c-acfa-5e589edc9065	e6f8015d-de60-4179-a24c-5da3e243c28a
2	Stupid guy	d011af53-5036-4512-a5c8-add6647ab62a	a86167b2-3ee7-416d-94bd-b17553a4f133
\.


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

