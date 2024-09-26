CREATE TYPE gender_type AS ENUM ('Male', 'Female');

CREATE TABLE IF NOT EXISTS customer (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    number INT NOT NULL,
    gender gender_type NOT NULL,
    country VARCHAR(50) NOT NULL,
    dependants INT NOT NULL DEFAULT 0 CHECK (dependants >= 0),
    birth_date DATE NOT NULL
);

INSERT INTO customer (name, surname, number, gender, country, dependants, birth_date) VALUES
('Jack', 'Front', 123, 'Male', 'Latvia', 5, '1981-10-03'),
('Jill', 'Human', 654, 'Female', 'Spain', 0, '1983-06-02'),
('Robert', 'Pullman', 456, 'Male', 'Germany', 2, '1999-05-04'),
('Chun Li', 'Suzuki', 987, 'Female', 'China', 1, '2001-11-09'),
('Sarah', 'Van Que', 587, 'Female', 'Latvia', 4, '1989-06-22');
