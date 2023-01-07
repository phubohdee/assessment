-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS news_articles_id_seq;

-- Table Definition
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
	title TEXT,
	amount FLOAT,
	note TEXT,
	tags TEXT[]
);
