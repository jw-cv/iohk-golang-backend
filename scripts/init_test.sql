CREATE TABLE IF NOT EXISTS customers_test (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    number INTEGER NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('Male', 'Female')),
    country VARCHAR(255),
    dependants INTEGER CHECK (dependants >= 0),
    birth_date DATE CHECK (birth_date <= CURRENT_DATE)
);

INSERT INTO customers_test (name, surname, number, gender, country, dependants, birth_date)
VALUES
    ('John', 'Doe', 123, 'Male', 'USA', 2, '1990-01-01'),
    ('Jane', 'Smith', 456, 'Female', 'Canada', 1, '1985-05-15');