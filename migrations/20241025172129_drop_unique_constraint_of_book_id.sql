-- +goose Up
-- +goose StatementBegin
ALTER TABLE "loans" DROP CONSTRAINT "loans_book_id_key";
ALTER TABLE "reservations" DROP CONSTRAINT "reservations_book_id_key";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "reservations" ADD CONSTRAINT "reservations_book_id_key" UNIQUE ("book_id");
ALTER TABLE "loans" ADD CONSTRAINT "loans_book_id_key" UNIQUE ("book_id");
-- +goose StatementEnd
