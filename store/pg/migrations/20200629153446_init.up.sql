CREATE TABLE clients (
	id uuid NOT NULL,
	bio text NOT NULL,
	balance money NOT NULL,
	created_at timestamp default current_timestamp,
	CONSTRAINT "pk_client_id" PRIMARY KEY (id)
);

