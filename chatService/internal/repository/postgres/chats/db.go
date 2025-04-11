package chats

import (
	"context"

	"github.com/Masterminds/squirrel"
	"microservices/chatService/internal/client/db"
	"microservices/chatService/internal/model"
	"microservices/chatService/internal/repository"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type repo struct {
	DB db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{DB: db}
}

func (r *repo) CreateChat(ctx context.Context, chat *model.Chat) (int64, error) {
	// Создаем чат
	query := psql.
		Insert("chats").
		Columns("created_at").
		Values("NOW()").
		Suffix("RETURNING id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: sqlStr,
	}

	var chatID int64
	err = r.DB.QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	// Связываем чат с пользователями
	if len(chat.Usernames) > 0 {
		query2 := psql.
			Insert("chat_users").
			Columns("chat_id", "username")

		// Добавляем все пользователи в VALUES
		for _, username := range chat.Usernames {
			query2 = query2.
				Values(chatID, username)
		}

		// Генерируем SQL и аргументы
		sqlStr, args, err = query2.ToSql()
		if err != nil {
			return 0, err
		}

		queryExec := db.Query{
			Name:     "user_repository.CreateChatUsers",
			QueryRaw: sqlStr,
		}

		// Выполняем вставку пользователей
		_, err = r.DB.ExecContext(ctx, queryExec, args...)
		if err != nil {
			return 0, err
		}
	}

	return chatID, nil
}

// DeleteChat удаляет чат по ID
func (r *repo) DeleteChat(ctx context.Context, chatID int64) error {
	// Удаляем чат
	query := psql.
		Delete("chats").
		Where(squirrel.Eq{"id": chatID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: sqlStr,
	}

	_, err = r.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

// SendMessage сохраняет сообщение в чат
func (r *repo) SendMessage(ctx context.Context, req *model.Message) error {
	// Сохраняем сообщение
	query := psql.
		Insert("messages").
		Columns("from_user", "text", "timestamp").
		Values(req.FromUser, req.Text, req.Timestamp).
		Suffix("RETURNING id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.SendMessage",
		QueryRaw: sqlStr,
	}

	_, err = r.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
