-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    question_id INT NOT NULL,
    user_id UUID NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (question_id) REFERENCES questions(id) ON UPDATE CASCADE ON DELETE CASCADE
);
INSERT INTO answers (question_id, user_id, text)
VALUES
    -- Вопрос 1: What is the capital of France?
    (1, uuid_generate_v4(), 'Paris'),
    (1, uuid_generate_v4(), 'London'),
    (1, uuid_generate_v4(), 'Berlin'),
    (1, uuid_generate_v4(), 'Madrid'),
    (1, uuid_generate_v4(), 'Rome'),

    -- Вопрос 2: Who wrote the novel "War and Peace"?
    (2, uuid_generate_v4(), 'Leo Tolstoy'),
    (2, uuid_generate_v4(), 'Fyodor Dostoevsky'),
    (2, uuid_generate_v4(), 'Alexander Pushkin'),
    (2, uuid_generate_v4(), 'Anton Chekhov'),
    (2, uuid_generate_v4(), 'Ivan Turgenev'),

    -- Вопрос 3: Which planet is known as the Red Planet?
    (3, uuid_generate_v4(), 'Mars'),
    (3, uuid_generate_v4(), 'Venus'),
    (3, uuid_generate_v4(), 'Earth'),
    (3, uuid_generate_v4(), 'Jupiter'),
    (3, uuid_generate_v4(), 'Saturn'),

    -- Вопрос 4: How many continents are there in the world?
    (4, uuid_generate_v4(), 'Seven'),
    (4, uuid_generate_v4(), 'Six'),
    (4, uuid_generate_v4(), 'Eight'),
    (4, uuid_generate_v4(), 'Five'),
    (4, uuid_generate_v4(), 'Four'),

    -- Вопрос 5: What is the largest ocean on Earth?
    (5, uuid_generate_v4(), 'Pacific Ocean'),
    (5, uuid_generate_v4(), 'Atlantic Ocean'),
    (5, uuid_generate_v4(), 'Indian Ocean'),
    (5, uuid_generate_v4(), 'Arctic Ocean'),
    (5, uuid_generate_v4(), 'Southern Ocean'),

    -- Вопрос 6: Who painted the Mona Lisa?
    (6, uuid_generate_v4(), 'Leonardo da Vinci'),
    (6, uuid_generate_v4(), 'Vincent van Gogh'),
    (6, uuid_generate_v4(), 'Claude Monet'),
    (6, uuid_generate_v4(), 'Salvador Dalí'),
    (6, uuid_generate_v4(), 'Pablo Picasso'),

    -- Вопрос 7: What year did World War II end?
    (7, uuid_generate_v4(), '1945'),
    (7, uuid_generate_v4(), '1946'),
    (7, uuid_generate_v4(), '1944'),
    (7, uuid_generate_v4(), '1947'),
    (7, uuid_generate_v4(), '1943'),

    -- Вопрос 8: What is the chemical symbol for gold?
    (8, uuid_generate_v4(), 'Au'),
    (8, uuid_generate_v4(), 'Ag'),
    (8, uuid_generate_v4(), 'Fe'),
    (8, uuid_generate_v4(), 'Cu'),
    (8, uuid_generate_v4(), 'Zn'),

    -- Вопрос 9: What is the tallest mountain in the world?
    (9, uuid_generate_v4(), 'Mount Everest'),
    (9, uuid_generate_v4(), 'K2'),
    (9, uuid_generate_v4(), 'Kilimanjaro'),
    (9, uuid_generate_v4(), 'Denali'),
    (9, uuid_generate_v4(), 'Aconcagua'),

    -- Вопрос 10: What is the boiling point of water at sea level?
    (10, uuid_generate_v4(), '100°C'),
    (10, uuid_generate_v4(), '90°C'),
    (10, uuid_generate_v4(), '110°C'),
    (10, uuid_generate_v4(), '80°C'),
    (10, uuid_generate_v4(), '120°C')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_answers_question_id;
DROP TABLE IF EXISTS answers;
-- +goose StatementEnd
