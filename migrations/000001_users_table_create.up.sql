CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	first_name VARCHAR(255),
	last_name VARCHAR(255),
	email VARCHAR(255) UNIQUE,
	phone VARCHAR(255) UNIQUE,
	password VARCHAR(255),
	role VARCHAR(255),
	username VARCHAR(255) UNIQUE
);
