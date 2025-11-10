CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(128) NOT NULL,
    group_name VARCHAR(32) NOT NULL,
    variant INT NOT NULL,
    email VARCHAR(128),
    github VARCHAR(64)
);

INSERT INTO students (fullname, group_name, variant, email, github)
VALUES ('Кужир Владислав Витальевич', 'AC-64', 26, 'your_email@edu.bstu.ru', 'XD-cods')
ON CONFLICT DO NOTHING;
