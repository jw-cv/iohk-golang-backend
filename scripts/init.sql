-- Create the customer table
CREATE TABLE IF NOT EXISTS customers (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    number INT NOT NULL,
    gender VARCHAR(15) NOT NULL CHECK (gender IN ('Male', 'Female')),
    country VARCHAR(50) NOT NULL,
    dependants INT NOT NULL DEFAULT 0 CHECK (dependants >= 0),
    birth_date DATE NOT NULL CHECK (birth_date <= CURRENT_DATE)
);

-- Insert sample data into the customer table using TO_DATE to handle the MM/DD/YYYY date format
INSERT INTO customers (name, surname, number, gender, country, dependants, birth_date) VALUES
('Jack', 'Front', 123, 'Male', 'Latvia', 5, TO_DATE('10/3/1981', 'MM/DD/YYYY')),
('Jill', 'Human', 654, 'Female', 'Spain', 0, TO_DATE('6/2/1983', 'MM/DD/YYYY')),
('Robert', 'Pullman', 456, 'Male', 'Germany', 2, TO_DATE('5/4/1999', 'MM/DD/YYYY')),
('Chun Li', 'Suzuki', 987, 'Female', 'China', 1, TO_DATE('11/9/2001', 'MM/DD/YYYY')),
('Sarah', 'Van Que', 587, 'Female', 'Latvia', 4, TO_DATE('6/22/1989', 'MM/DD/YYYY'));
