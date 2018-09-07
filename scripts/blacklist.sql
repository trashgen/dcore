CREATE TABLE IF NOT EXISTS blacklist (
	id bigserial NOT NULL,
	ip varchar(45) NOT NULL,
	created timestamp default now()
)
WITH (OIDS=false);

CREATE UNIQUE INDEX IF NOT EXISTS blacklist_ip_idx ON blacklist (ip);