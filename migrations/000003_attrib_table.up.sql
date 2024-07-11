CREATE TABLE IF NOT EXISTS product_attrib (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	quantity NUMERIC(4,2),-- in kilo grams,
	price NUMERIC(8,2),
	status VARCHAR(1), -- 1 available,0 sold
	product_id UUID REFERENCES products(id) ON UPDATE CASCADE ON DELETE CASCADE,
	updated_at TIMESTAMP
);
