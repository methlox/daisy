CREATE TABLE form {
    id SERIAL PRIMARY KEY
    title TEXT NOT NULL
    description TEXT
    created_at DATETIME
}

CREATE TABLE question {
    id SERIAL PRIMARY KEY
    form_id SERIAL
    question_text TEXT NOT NULL
    question_order INT
    created_at DATETIME
}

CREATE TABLE formResponse {
    id SERIAL PRIMARY KEY
    form_id SERIAL
    responded_at DATETIME
}

CREATE TABLE response {
    id SERIAL PRIMARY KEY
    response_id SERIAL
    question_id SERIAL
    answer TEXT NOT NULL
}