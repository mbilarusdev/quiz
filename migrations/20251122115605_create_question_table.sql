-- +goose Up
-- +goose StatementBegin
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
INSERT INTO questions (text) VALUES
('What is the capital of France?'),
('Who wrote the novel "War and Peace"?'),
('Which planet is known as the Red Planet?'),
('How many continents are there in the world?'),
('What is the largest ocean on Earth?'),
('Who painted the Mona Lisa?'),
('What year did World War II end?'),
('What is the chemical symbol for gold?'),
('What is the tallest mountain in the world?'),
('What is the boiling point of water at sea level?');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS questions;
-- +goose StatementEnd
