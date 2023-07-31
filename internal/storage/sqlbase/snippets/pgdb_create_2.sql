-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler version: 1.0.2
-- PostgreSQL version: 15.0
-- Project Site: pgmodeler.io
-- Model Author: ---

-- Database creation must be performed outside a multi lined SQL file. 
-- These commands were put in this file only as a convenience.
-- 
-- object: store | type: DATABASE --
-- DROP DATABASE IF EXISTS store;
-- CREATE DATABASE store
-- 	ENCODING = 'UTF8'
-- 	LC_COLLATE = 'Russian_Russia.1251'
-- 	LC_CTYPE = 'Russian_Russia.1251'
-- 	TABLESPACE = pg_default
-- 	OWNER = postgres;
-- ddl-end --


-- object: metmon | type: SCHEMA --
-- DROP SCHEMA IF EXISTS metmon CASCADE;
CREATE SCHEMA metmon;
-- ddl-end --
ALTER SCHEMA metmon OWNER TO postgres;
-- ddl-end --

SET search_path TO pg_catalog,public,metmon;
-- ddl-end --

-- object: metmon.gauges | type: TABLE --
-- DROP TABLE IF EXISTS metmon.gauges CASCADE;
CREATE TABLE metmon.gauges (
	gauge_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START WITH 1 CACHE 1 ),
	"timestamp" timestamp with time zone NOT NULL DEFAULT now(),
	gvalue double precision NOT NULL,
	metric_id integer NOT NULL,
	CONSTRAINT gauges_pk PRIMARY KEY (gauge_id)
);
-- ddl-end --
ALTER TABLE metmon.gauges OWNER TO postgres;
-- ddl-end --

-- -- object: metmon.gauges_gauge_id_seq | type: SEQUENCE --
-- -- DROP SEQUENCE IF EXISTS metmon.gauges_gauge_id_seq CASCADE;
-- CREATE SEQUENCE metmon.gauges_gauge_id_seq
-- 	INCREMENT BY 1
-- 	MINVALUE 1
-- 	MAXVALUE 2147483647
-- 	START WITH 1
-- 	CACHE 1
-- 	NO CYCLE
-- 	OWNED BY NONE;
-- 
-- -- ddl-end --
-- ALTER SEQUENCE metmon.gauges_gauge_id_seq OWNER TO postgres;
-- -- ddl-end --
-- 
-- object: metmon.counters | type: TABLE --
-- DROP TABLE IF EXISTS metmon.counters CASCADE;
CREATE TABLE metmon.counters (
	counter_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START WITH 1 CACHE 1 ),
	"timestamp" timestamp with time zone NOT NULL DEFAULT now(),
	cvalue bigint NOT NULL,
	metric_id integer NOT NULL,
	CONSTRAINT counters_pk PRIMARY KEY (counter_id),
	CONSTRAINT counters_uq UNIQUE (metric_id)
);
-- ddl-end --
ALTER TABLE metmon.counters OWNER TO postgres;
-- ddl-end --

-- -- object: metmon.counters_counter_id_seq | type: SEQUENCE --
-- -- DROP SEQUENCE IF EXISTS metmon.counters_counter_id_seq CASCADE;
-- CREATE SEQUENCE metmon.counters_counter_id_seq
-- 	INCREMENT BY 1
-- 	MINVALUE 1
-- 	MAXVALUE 2147483647
-- 	START WITH 1
-- 	CACHE 1
-- 	NO CYCLE
-- 	OWNED BY NONE;
-- 
-- -- ddl-end --
-- ALTER SEQUENCE metmon.counters_counter_id_seq OWNER TO postgres;
-- -- ddl-end --
-- 
-- object: metmon.metrics | type: TABLE --
-- DROP TABLE IF EXISTS metmon.metrics CASCADE;
CREATE TABLE metmon.metrics (
	metric_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START WITH 1 CACHE 1 ),
	name character varying NOT NULL,
	type boolean NOT NULL,
	CONSTRAINT metrics_pk PRIMARY KEY (metric_id),
	CONSTRAINT uniq_name UNIQUE (name)
);
-- ddl-end --
COMMENT ON COLUMN metmon.metrics.type IS E'is gauge';
-- ddl-end --
ALTER TABLE metmon.metrics OWNER TO postgres;
-- ddl-end --

-- -- object: metmon.metrics_metric_id_seq | type: SEQUENCE --
-- -- DROP SEQUENCE IF EXISTS metmon.metrics_metric_id_seq CASCADE;
-- CREATE SEQUENCE metmon.metrics_metric_id_seq
-- 	INCREMENT BY 1
-- 	MINVALUE 1
-- 	MAXVALUE 2147483647
-- 	START WITH 1
-- 	CACHE 1
-- 	NO CYCLE
-- 	OWNED BY NONE;
-- 
-- -- ddl-end --
-- ALTER SEQUENCE metmon.metrics_metric_id_seq OWNER TO postgres;
-- -- ddl-end --
-- 
-- object: metmon.current_metrics | type: VIEW --
-- DROP VIEW IF EXISTS metmon.current_metrics CASCADE;
CREATE VIEW metmon.current_metrics
AS 

SELECT COALESCE(gm."timestamp", c."timestamp") AS "time",
    m.name,
    m.type,
    COALESCE(c.cvalue, 0) as "cvalue",
    COALESCE(gm.gvalue, 0) as "gvalue"
   FROM ((metmon.metrics m
     LEFT JOIN metmon.counters c ON ((m.metric_id = c.metric_id)))
     LEFT JOIN ( SELECT DISTINCT ON (g.metric_id) g."timestamp",
            g.metric_id,
            g.gvalue
           FROM metmon.gauges g
          ORDER BY g.metric_id, g."timestamp" DESC) gm ON ((m.metric_id = gm.metric_id)));
-- ddl-end --
ALTER VIEW metmon.current_metrics OWNER TO postgres;
-- ddl-end --

-- object: metrics_fk | type: CONSTRAINT --
-- ALTER TABLE metmon.gauges DROP CONSTRAINT IF EXISTS metrics_fk CASCADE;
ALTER TABLE metmon.gauges ADD CONSTRAINT metrics_fk FOREIGN KEY (metric_id)
REFERENCES metmon.metrics (metric_id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: metrics_fk | type: CONSTRAINT --
-- ALTER TABLE metmon.counters DROP CONSTRAINT IF EXISTS metrics_fk CASCADE;
ALTER TABLE metmon.counters ADD CONSTRAINT metrics_fk FOREIGN KEY (metric_id)
REFERENCES metmon.metrics (metric_id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --


