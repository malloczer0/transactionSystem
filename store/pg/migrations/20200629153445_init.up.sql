CREATE TABLE transactions (
	id uuid NOT NULL,
	client_id uuid NOT NULL,
	status smallint NOT NULL default 0,
	change float NOT NULL,
	created_at timestamp default current_timestamp,
	CONSTRAINT "pk_transaction_id" PRIMARY KEY (id)
);

