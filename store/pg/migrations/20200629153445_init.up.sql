CREATE TABLE transactions (
	id uuid NOT NULL,
	client_id uuid NOT NULL,
	status smallint NOT NULL,
	change money NOT NULL,
	created_at timestamp default current_timestamp,
	CONSTRAINT "pk_transaction_id" PRIMARY KEY (id)
);

