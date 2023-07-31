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
-- CREATE DATABASE store;
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
	gauge_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ,
	"timestamp" timestamptz NOT NULL DEFAULT NOW(),
	gvalue double precision NOT NULL,
	metric_id integer,
	CONSTRAINT gauges_pk PRIMARY KEY (gauge_id)
);
-- ddl-end --
ALTER TABLE metmon.gauges OWNER TO postgres;
-- ddl-end --

-- object: metmon.counters | type: TABLE --
-- DROP TABLE IF EXISTS metmon.counters CASCADE;
CREATE TABLE metmon.counters (
	counter_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ,
	"timestamp" timestamptz NOT NULL DEFAULT NOW(),
	cvalue bigint NOT NULL,
	metric_id integer,
	CONSTRAINT counters_pk PRIMARY KEY (counter_id)
);
-- ddl-end --
ALTER TABLE metmon.counters OWNER TO postgres;
-- ddl-end --

-- object: metmon.metrics | type: TABLE --
-- DROP TABLE IF EXISTS metmon.metrics CASCADE;
CREATE TABLE metmon.metrics (
	metric_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ,
	name varchar NOT NULL,
	type bool NOT NULL,
	CONSTRAINT metrics_pk PRIMARY KEY (metric_id),
	CONSTRAINT uniq_name UNIQUE (name)
);
-- ddl-end --
COMMENT ON COLUMN metmon.metrics.type IS E'is gauge';
-- ddl-end --
ALTER TABLE metmon.metrics OWNER TO postgres;
-- ddl-end --

-- object: metrics_fk | type: CONSTRAINT --
-- ALTER TABLE metmon.counters DROP CONSTRAINT IF EXISTS metrics_fk CASCADE;
ALTER TABLE metmon.counters ADD CONSTRAINT metrics_fk FOREIGN KEY (metric_id)
REFERENCES metmon.metrics (metric_id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: counters_uq | type: CONSTRAINT --
-- ALTER TABLE metmon.counters DROP CONSTRAINT IF EXISTS counters_uq CASCADE;
ALTER TABLE metmon.counters ADD CONSTRAINT counters_uq UNIQUE (metric_id);
-- ddl-end --

-- object: metrics_fk | type: CONSTRAINT --
-- ALTER TABLE metmon.gauges DROP CONSTRAINT IF EXISTS metrics_fk CASCADE;
ALTER TABLE metmon.gauges ADD CONSTRAINT metrics_fk FOREIGN KEY (metric_id)
REFERENCES metmon.metrics (metric_id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --


