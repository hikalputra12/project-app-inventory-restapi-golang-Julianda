--
-- PostgreSQL database dump
--

\restrict PR5Lv13QB9IQkXtlsm5WCjNX2Jv5anIYPxKG4fOVnWnPW8b7ytGLSvgfyFpahXb

-- Dumped from database version 16.11
-- Dumped by pg_dump version 16.11

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
-- Name: category_inventory; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.category_inventory (
    category_inventory_id integer NOT NULL,
    rack_inventory_id integer,
    name character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


ALTER TABLE public.category_inventory OWNER TO postgres;

--
-- Name: category_inventory_category_inventory_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.category_inventory_category_inventory_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.category_inventory_category_inventory_id_seq OWNER TO postgres;

--
-- Name: category_inventory_category_inventory_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.category_inventory_category_inventory_id_seq OWNED BY public.category_inventory.category_inventory_id;


--
-- Name: inventories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.inventories (
    inventory_id integer NOT NULL,
    category_inventory_id integer,
    name character varying(255),
    price integer,
    stock integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


ALTER TABLE public.inventories OWNER TO postgres;

--
-- Name: inventories_inventory_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.inventories_inventory_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.inventories_inventory_id_seq OWNER TO postgres;

--
-- Name: inventories_inventory_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.inventories_inventory_id_seq OWNED BY public.inventories.inventory_id;


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.permissions (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    description text
);


ALTER TABLE public.permissions OWNER TO postgres;

--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.permissions_id_seq OWNER TO postgres;

--
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- Name: rack_inventory; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rack_inventory (
    rack_inventory_id integer NOT NULL,
    name character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    warehouse_inventory_id integer
);


ALTER TABLE public.rack_inventory OWNER TO postgres;

--
-- Name: rack_inventory_rack_inventory_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.rack_inventory_rack_inventory_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.rack_inventory_rack_inventory_id_seq OWNER TO postgres;

--
-- Name: rack_inventory_rack_inventory_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rack_inventory_rack_inventory_id_seq OWNED BY public.rack_inventory.rack_inventory_id;


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role_permissions (
    role_id integer NOT NULL,
    permission_id integer NOT NULL
);


ALTER TABLE public.role_permissions OWNER TO postgres;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    description text
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.roles_id_seq OWNER TO postgres;

--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: sales_item; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sales_item (
    sales_item_id integer NOT NULL,
    user_id integer,
    inventory_id integer,
    quantity integer,
    price integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


ALTER TABLE public.sales_item OWNER TO postgres;

--
-- Name: sales_item_sales_item_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sales_item_sales_item_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sales_item_sales_item_id_seq OWNER TO postgres;

--
-- Name: sales_item_sales_item_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sales_item_sales_item_id_seq OWNED BY public.sales_item.sales_item_id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    session_id uuid NOT NULL,
    user_id integer,
    expired_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    last_active timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: user_permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_permissions (
    user_id integer NOT NULL,
    permission_id integer NOT NULL,
    effect character varying(10),
    CONSTRAINT user_permissions_effect_check CHECK (((effect)::text = ANY ((ARRAY['allow'::character varying, 'deny'::character varying])::text[])))
);


ALTER TABLE public.user_permissions OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    user_id integer NOT NULL,
    name character varying(255),
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    role_id integer
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_user_id_seq OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_user_id_seq OWNED BY public.users.user_id;


--
-- Name: warehouse_inventory; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.warehouse_inventory (
    warehouse_inventory_id integer NOT NULL,
    name character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    location character varying(100)
);


ALTER TABLE public.warehouse_inventory OWNER TO postgres;

--
-- Name: warehouse_inventory_warehouse_inventory_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.warehouse_inventory_warehouse_inventory_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.warehouse_inventory_warehouse_inventory_id_seq OWNER TO postgres;

--
-- Name: warehouse_inventory_warehouse_inventory_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.warehouse_inventory_warehouse_inventory_id_seq OWNED BY public.warehouse_inventory.warehouse_inventory_id;


--
-- Name: category_inventory category_inventory_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_inventory ALTER COLUMN category_inventory_id SET DEFAULT nextval('public.category_inventory_category_inventory_id_seq'::regclass);


--
-- Name: inventories inventory_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventories ALTER COLUMN inventory_id SET DEFAULT nextval('public.inventories_inventory_id_seq'::regclass);


--
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: rack_inventory rack_inventory_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rack_inventory ALTER COLUMN rack_inventory_id SET DEFAULT nextval('public.rack_inventory_rack_inventory_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: sales_item sales_item_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales_item ALTER COLUMN sales_item_id SET DEFAULT nextval('public.sales_item_sales_item_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_user_id_seq'::regclass);


--
-- Name: warehouse_inventory warehouse_inventory_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouse_inventory ALTER COLUMN warehouse_inventory_id SET DEFAULT nextval('public.warehouse_inventory_warehouse_inventory_id_seq'::regclass);


--
-- Data for Name: category_inventory; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.category_inventory (category_inventory_id, rack_inventory_id, name, created_at, updated_at, deleted_at) FROM stdin;
2	1	Smartphone	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N
3	2	Kabel & Charger	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N
4	3	Stationary	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N
5	3	peralatan bayi	2026-01-05 04:43:10.377529	2026-01-05 04:43:10.377529	\N
1	3	Laptop & PC	2026-01-02 01:48:20.022833	2026-01-11 21:37:45.740816	\N
\.


--
-- Data for Name: inventories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.inventories (inventory_id, category_inventory_id, name, price, stock, created_at, updated_at, deleted_at) FROM stdin;
2	1	Asus ROG Zephyrus	25000000	5	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N
3	2	iPhone 15	18000000	20	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N
5	3	Kabel USB-C Baseus	50000	100	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N
1	1	MacBook Pro M2	20000000	4	2026-01-02 01:48:20.022833	2026-01-05 22:46:19.949162	\N
8	1	xiaomi2pro	110000	10	2026-01-11 20:31:03.084821	2026-01-11 20:31:03.084821	\N
9	1	xiaomi2pro	110000	10	2026-01-11 21:13:45.422267	2026-01-11 21:13:45.422267	\N
4	2	Samsung S24 Ultra	19000000	120	2026-01-02 01:48:20.022833	2026-01-11 21:13:56.754061	\N
\.


--
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.permissions (id, code, description) FROM stdin;
1	inventory:view	Melihat daftar barang
2	inventory:create	Menambah barang baru
3	inventory:edit	Mengedit barang
4	inventory:delete	Menghapus barang
5	user:view	Melihat daftar user
6	user:manage	Mengelola user (tambah/hapus)
7	category:view	Melihat category barang
8	category:manage	Mengelola category(membuat,edit,melihat dan delete)
9	rack:view	Melihat rack barang
10	rack:manage	mengedit,membuat dan menghapus rack
11	warehouse:view	Melihat daftar gudang barang
12	warehouse:manage	Membuat,mengedit dan menghapus daftar gudang barang
13	transaction:manage	Mengelola transaksi
14	report:view	melihat report penjualan
15	stock:view	check stock minimal 5
\.


--
-- Data for Name: rack_inventory; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rack_inventory (rack_inventory_id, name, created_at, updated_at, deleted_at, warehouse_inventory_id) FROM stdin;
2	Rak B - Aksesoris	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N	1
3	Rak C - Umum	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N	2
1	Rak A - Elektronik	2026-01-02 01:48:20.022833	2026-01-11 22:52:39.228685	\N	2
\.


--
-- Data for Name: role_permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.role_permissions (role_id, permission_id) FROM stdin;
1	1
1	2
1	3
1	4
1	5
1	6
2	1
2	2
2	3
2	4
2	5
3	1
2	6
1	8
2	8
3	7
1	7
2	7
1	9
2	9
3	9
1	10
2	10
1	11
2	11
3	11
1	12
2	12
1	13
2	13
3	13
1	14
2	14
1	15
2	15
3	15
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id, name, description) FROM stdin;
1	super admin	\N
2	admin	\N
3	staf	\N
\.


--
-- Data for Name: sales_item; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sales_item (sales_item_id, user_id, inventory_id, quantity, price, created_at, updated_at, deleted_at) FROM stdin;
1	1	1	1	20000000	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N
2	1	3	1	18000000	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N
3	2	5	1	50000	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (session_id, user_id, expired_at, revoked_at, created_at, last_active) FROM stdin;
e47f9e6e-c45b-438c-9499-aadc2ced87f6	2	2026-01-12 03:48:25.472264+08	2026-01-11 03:49:47.909929+08	2026-01-11 00:44:57.855125	2026-01-11 03:48:25.472716
aeb416c0-aa34-4050-b836-86cf3857d851	2	2026-01-12 16:48:53.653578+08	2026-01-11 16:49:07.334506+08	2026-01-11 03:50:14.228475	2026-01-11 16:48:53.65417
0b3a9696-9633-4974-bdf6-eafb3bb75d45	2	2026-01-12 17:56:53.225189+08	2026-01-11 18:06:15.057744+08	2026-01-11 17:52:59.045144	2026-01-11 17:56:53.225941
07ae9122-40e5-4a7a-84ae-d01831eb3a4d	9	2026-01-12 18:06:51.250938+08	2026-01-11 18:07:12.340074+08	2026-01-11 18:06:51.250938	2026-01-11 18:06:51.250938
2613a98f-6d4f-46fe-bd70-d7bbf573d846	2	2026-01-12 18:08:48.532499+08	\N	2026-01-11 18:08:48.532499	2026-01-11 18:08:48.532499
b22d29e5-38cd-415b-a408-14df7404b1e3	2	2026-01-12 23:53:01.271085+08	\N	2026-01-11 18:13:45.250554	2026-01-11 23:53:01.271156
\.


--
-- Data for Name: user_permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_permissions (user_id, permission_id, effect) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (user_id, name, email, password_hash, created_at, updated_at, deleted_at, role_id) FROM stdin;
9	cahya	cahya123@gmail.com	$2a$14$fk1rhXM6Z3joh.bq2zdTMOM3GK8Bue19AD1clCTGAYkhHIvbisoN6	2026-01-11 17:56:54.726332	2026-01-11 17:56:54.726332	\N	3
8	hendry	hendry@gmail.com	$2a$14$M1IgomCUpWNCe08xpAUA8Ov4WGG0D3bDIjVIo2lYlD.XrHPqUTdL6	2026-01-11 17:53:32.102882	2026-01-11 17:53:32.102882	\N	3
1	Julianda Putra	julianda@admin.com	$2a$14$Y.uIs6s/LQ3d.uxZRLlhtOhVexO.WrHdWitInLVbWEO2s2vaFVqca	2026-01-02 01:48:20.022833	2026-01-11 20:17:48.840147	\N	1
2	Budi Santoso	budi@gudang.com	$2a$14$vzqNsJBhd0Tk5v3MLpg0oelT63jPxu1LltMw9jmZMprLQuyvPCOHa	2026-01-02 01:48:20.022833	2026-01-11 20:18:26.398357	\N	2
3	Siti Aminah	siti@kasir.com	$2a$14$BPBeFnb1NGra9rlszLPtFeIzwY3crzft/9yb2Uuy73wO7mxM.CWLG	2026-01-02 01:48:20.022833	2026-01-11 20:18:45.397801	\N	3
\.


--
-- Data for Name: warehouse_inventory; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.warehouse_inventory (warehouse_inventory_id, name, created_at, updated_at, deleted_at, location) FROM stdin;
1	Gudang Utama Makassar	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N	Makassar, kecamatan panakukang
2	Gudang Cabang Jakarta	2026-01-02 01:48:20.022833	2026-01-02 01:48:20.022833	\N	Jakarta, Bundaran HI
4	Gudang Bau-Bau	2026-01-11 23:36:25.832464	2026-01-11 23:39:46.862765	\N	bukit wolio indah kecamatan wolio, Kota Bau-Bau
\.


--
-- Name: category_inventory_category_inventory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.category_inventory_category_inventory_id_seq', 15, true);


--
-- Name: inventories_inventory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.inventories_inventory_id_seq', 9, true);


--
-- Name: permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.permissions_id_seq', 6, true);


--
-- Name: rack_inventory_rack_inventory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rack_inventory_rack_inventory_id_seq', 4, true);


--
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 3, true);


--
-- Name: sales_item_sales_item_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sales_item_sales_item_id_seq', 4, true);


--
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_user_id_seq', 9, true);


--
-- Name: warehouse_inventory_warehouse_inventory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.warehouse_inventory_warehouse_inventory_id_seq', 4, true);


--
-- Name: category_inventory category_inventory_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_inventory
    ADD CONSTRAINT category_inventory_pkey PRIMARY KEY (category_inventory_id);


--
-- Name: inventories inventories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventories
    ADD CONSTRAINT inventories_pkey PRIMARY KEY (inventory_id);


--
-- Name: permissions permissions_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_code_key UNIQUE (code);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: rack_inventory rack_inventory_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rack_inventory
    ADD CONSTRAINT rack_inventory_pkey PRIMARY KEY (rack_inventory_id);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (role_id, permission_id);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: sales_item sales_item_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales_item
    ADD CONSTRAINT sales_item_pkey PRIMARY KEY (sales_item_id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (session_id);


--
-- Name: user_permissions user_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_permissions
    ADD CONSTRAINT user_permissions_pkey PRIMARY KEY (user_id, permission_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: warehouse_inventory warehouse_inventory_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouse_inventory
    ADD CONSTRAINT warehouse_inventory_pkey PRIMARY KEY (warehouse_inventory_id);


--
-- Name: category_inventory category_inventory_rack_inventory_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_inventory
    ADD CONSTRAINT category_inventory_rack_inventory_id_fkey FOREIGN KEY (rack_inventory_id) REFERENCES public.rack_inventory(rack_inventory_id);


--
-- Name: sessions fk_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- Name: inventories inventories_category_inventory_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventories
    ADD CONSTRAINT inventories_category_inventory_id_fkey FOREIGN KEY (category_inventory_id) REFERENCES public.category_inventory(category_inventory_id);


--
-- Name: role_permissions role_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: role_permissions role_permissions_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: sales_item sales_item_inventory_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales_item
    ADD CONSTRAINT sales_item_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES public.inventories(inventory_id);


--
-- Name: user_permissions user_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_permissions
    ADD CONSTRAINT user_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: user_permissions user_permissions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_permissions
    ADD CONSTRAINT user_permissions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- Name: users users_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- PostgreSQL database dump complete
--

\unrestrict PR5Lv13QB9IQkXtlsm5WCjNX2Jv5anIYPxKG4fOVnWnPW8b7ytGLSvgfyFpahXb

