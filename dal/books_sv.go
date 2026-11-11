package dal

import (
	"context"

	"github.com/phuslu/log"
)

func GetAllBooks(ctx context.Context) ([]Book, error) {
	var books []Book
	err := DB.NewSelect().Model(&books).Scan(ctx)
	if err == nil {
		log.Warn().Msgf("Retrieved all books: #%d", len(books))
	}
	return books, err
}

func GetBookByID(ctx context.Context, id uint) (*Book, error) {
	var book Book
	err := DB.NewSelect().Model(&book).Where("id = ?", id).Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get book with id: %d", id)
		return nil, err
	}
	return &book, nil
}

func CreateBook(ctx context.Context, book *Book) error {
	_, err := DB.NewInsert().Model(book).Exec(ctx)
	if err == nil {
		log.Info().Msgf("Created book: %v", book)
	} else {
		log.Error().Err(err).Msgf("Failed to create book: %v", book)
	}
	return err
}

func UpdateBook(ctx context.Context, book *Book) error {
	_, err := DB.NewUpdate().Model(book).Where("id = ?", book.ID).Exec(ctx)
	if err == nil {
		log.Info().Msgf("Updated book: %v", book)
	} else {
		log.Error().Err(err).Msgf("Failed to update book: %v", book)
	}
	return err
}

func DeleteBook(ctx context.Context, id uint) error {
	_, err := DB.NewDelete().Model(&Book{}).Where("id = ?", id).Exec(ctx)
	if err == nil {
		log.Info().Msgf("Deleted book with id: %d", id)
	} else {
		log.Error().Err(err).Msgf("Failed to delete book with id: %d", id)
	}
	return err
}
