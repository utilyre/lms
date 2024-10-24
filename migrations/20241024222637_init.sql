-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY,

    "name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(300) NOT NULL UNIQUE,
    "password" BYTEA NOT NULL,
    "role" VARCHAR(50)
);

CREATE TABLE "books" (
    "id" SERIAL PRIMARY KEY,

    "title" VARCHAR(100) NOT NULL,
    "author" VARCHAR(100),
    "isbn" CHAR(13) NOT NULL,
    "availibility_status" VARCHAR(50)
);

CREATE TABLE "loans" (
    "id" SERIAL PRIMARY KEY,

    "user_id" INTEGER NOT NULL REFERENCES "users" ON DELETE CASCADE,
    "book_id" INTEGER NOT NULL UNIQUE REFERENCES "books" ON DELETE CASCADE,

    "loan_date" TIMESTAMP NOT NULL,
    "due_date" TIMESTAMP NOT NULL,
    "return_date" TIMESTAMP
);

CREATE TABLE "reservations" (
    "id" SERIAL PRIMARY KEY,

    "user_id" INTEGER NOT NULL REFERENCES "users" ON DELETE CASCADE,
    "book_id" INTEGER NOT NULL UNIQUE REFERENCES "books" ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "reservations";
DROP TABLE "loans";
DROP TABLE "books";
DROP TABLE "users";
-- +goose StatementEnd
