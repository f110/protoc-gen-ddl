// Generated by protoc-ddl.
// protoc-gen-dao: v0.1
package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"go.f110.dev/protoc-ddl/sample"
)

type User struct {
	conn *sql.DB
}

func NewUser(conn *sql.DB) *User {
	return &User{
		conn: conn,
	}
}

func (d *User) Select(ctx context.Context, id int32) (*sample.User, error) {
	row := d.conn.QueryRowContext(ctx, "SELECT * FROM `users` WHERE `id` = ?", id)

	v := &sample.User{}
	if err := row.Scan(&v.Id, &v.Age, &v.Name, &v.CreatedAt); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	v.ResetMark()
	return v, nil
}

func (d *User) ListAll(ctx context.Context) ([]*sample.User, error) {
	rows, err := d.conn.QueryContext(
		ctx,
		"SELECT * FROM user",
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	res := make([]*sample.User, 0)
	for rows.Next() {
		r := &sample.User{}
		if err := rows.Scan(&r.Id, &r.Age, &r.Name, &r.CreatedAt); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		r.ResetMark()
		res = append(res, r)
	}

	return res, nil
}

func (d *User) ListOverTwenty(ctx context.Context) ([]*sample.User, error) {
	rows, err := d.conn.QueryContext(
		ctx,
		"SELECT * FROM user WHERE age > 20",
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	res := make([]*sample.User, 0)
	for rows.Next() {
		r := &sample.User{}
		if err := rows.Scan(&r.Id, &r.Age, &r.Name, &r.CreatedAt); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		r.ResetMark()
		res = append(res, r)
	}

	return res, nil
}

func (d *User) Create(ctx context.Context, v *sample.User) (*sample.User, error) {
	res, err := d.conn.ExecContext(
		ctx,
		"INSERT INTO `task` (`age`, `name`, `created_at`) VALUES (?, ?, ?)", v.Age, v.Name, v.CreatedAt,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return nil, sql.ErrNoRows
	}

	v = v.Copy()
	insertedId, err := res.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	v.Id = int32(insertedId)

	v.ResetMark()
	return v, nil
}

func (d *User) Delete(ctx context.Context, id int32) error {
	res, err := d.conn.ExecContext(ctx, "DELETE FROM `users` WHERE `id` = ?", id)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *User) Update(ctx context.Context, v *sample.User) error {
	if !v.IsChanged() {
		return nil
	}

	changedColumn := v.ChangedColumn()
	cols := make([]string, len(changedColumn)+1)
	values := make([]interface{}, len(changedColumn)+1)
	for i := range changedColumn {
		cols[i] = "`" + changedColumn[i].Name + "` = ?"
		values[i] = changedColumn[i].Value
	}

	query := fmt.Sprintf("UPDATE `users` SET %s WHERE `id` = ?", strings.Join(cols, ", "))
	res, err := d.conn.ExecContext(
		ctx,
		query,
		append(values, v.Id)...,
	)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	v.ResetMark()
	return nil
}

type Blog struct {
	conn *sql.DB

	user *User
}

func NewBlog(conn *sql.DB) *Blog {
	return &Blog{
		conn: conn,
		user: NewUser(conn),
	}
}

func (d *Blog) Select(ctx context.Context, id int64) (*sample.Blog, error) {
	row := d.conn.QueryRowContext(ctx, "SELECT * FROM `blog` WHERE `id` = ?", id)

	v := &sample.Blog{}
	if err := row.Scan(&v.Id, &v.UserId, &v.Title, &v.Body, &v.CategoryId, &v.Attach, &v.CreatedAt, &v.UpdatedAt); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	{
		rel, err := d.user.Select(ctx, v.UserId)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		v.User = rel
	}

	v.ResetMark()
	return v, nil
}

func (d *Blog) ListByTitle(ctx context.Context, title string) ([]*sample.Blog, error) {
	rows, err := d.conn.QueryContext(
		ctx,
		"SELECT * FROM blog WHERE title = ?",
		title,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	res := make([]*sample.Blog, 0)
	for rows.Next() {
		r := &sample.Blog{}
		if err := rows.Scan(&r.Id, &r.UserId, &r.Title, &r.Body, &r.CategoryId, &r.Attach, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		r.ResetMark()
		res = append(res, r)
	}
	if len(res) > 0 {
		for _, v := range res {
			{
				rel, err := d.user.Select(ctx, v.UserId)
				if err != nil {
					return nil, xerrors.Errorf(": %w", err)
				}
				v.User = rel
			}
		}
	}

	return res, nil
}

func (d *Blog) ListByUserAndCategory(ctx context.Context, userId int32, categoryId int32) ([]*sample.Blog, error) {
	rows, err := d.conn.QueryContext(
		ctx,
		"SELECT * FROM blog WHERE user_id = ? AND category_id = ?",
		userId,
		categoryId,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	res := make([]*sample.Blog, 0)
	for rows.Next() {
		r := &sample.Blog{}
		if err := rows.Scan(&r.Id, &r.UserId, &r.Title, &r.Body, &r.CategoryId, &r.Attach, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		r.ResetMark()
		res = append(res, r)
	}
	if len(res) > 0 {
		for _, v := range res {
			{
				rel, err := d.user.Select(ctx, v.UserId)
				if err != nil {
					return nil, xerrors.Errorf(": %w", err)
				}
				v.User = rel
			}
		}
	}

	return res, nil
}

func (d *Blog) Create(ctx context.Context, v *sample.Blog) (*sample.Blog, error) {
	res, err := d.conn.ExecContext(
		ctx,
		"INSERT INTO `task` (`user_id`, `title`, `body`, `category_id`, `attach`, `created_at`) VALUES (?, ?, ?, ?, ?, ?)", v.UserId, v.Title, v.Body, v.CategoryId, v.Attach, time.Now(),
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return nil, sql.ErrNoRows
	}

	v = v.Copy()
	insertedId, err := res.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	v.Id = int64(insertedId)

	v.ResetMark()
	return v, nil
}

func (d *Blog) Delete(ctx context.Context, id int64) error {
	res, err := d.conn.ExecContext(ctx, "DELETE FROM `blog` WHERE `id` = ?", id)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *Blog) Update(ctx context.Context, v *sample.Blog) error {
	if !v.IsChanged() {
		return nil
	}

	changedColumn := v.ChangedColumn()
	cols := make([]string, len(changedColumn)+1)
	values := make([]interface{}, len(changedColumn)+1)
	for i := range changedColumn {
		cols[i] = "`" + changedColumn[i].Name + "` = ?"
		values[i] = changedColumn[i].Value
	}
	cols[len(cols)-1] = "`updated_at` = ?"
	values[len(values)-1] = time.Now()

	query := fmt.Sprintf("UPDATE `blog` SET %s WHERE `id` = ?", strings.Join(cols, ", "))
	res, err := d.conn.ExecContext(
		ctx,
		query,
		append(values, v.Id)...,
	)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	v.ResetMark()
	return nil
}

type CommentImage struct {
	conn *sql.DB

	comment *Comment
	like    *Like
}

func NewCommentImage(conn *sql.DB) *CommentImage {
	return &CommentImage{
		conn:    conn,
		comment: NewComment(conn),
		like:    NewLike(conn),
	}
}

func (d *CommentImage) Select(ctx context.Context, commentBlogId int64, commentUserId int32, likeId uint64) (*sample.CommentImage, error) {
	row := d.conn.QueryRowContext(ctx, "SELECT * FROM `comment_image` WHERE `comment_blog_id` = ? AND `comment_user_id` = ? AND `like_id` = ?", commentBlogId, commentUserId, likeId)

	v := &sample.CommentImage{}
	if err := row.Scan(&v.CommentBlogId, &v.CommentUserId, &v.LikeId); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	{
		rel, err := d.comment.Select(ctx, v.CommentBlogId, v.CommentUserId)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		v.Comment = rel
	}
	{
		rel, err := d.like.Select(ctx, v.LikeId)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		v.Like = rel
	}

	v.ResetMark()
	return v, nil
}

func (d *CommentImage) ListByLikeId(ctx context.Context, likeId uint64) ([]*sample.CommentImage, error) {
	rows, err := d.conn.QueryContext(
		ctx,
		"SELECT * FROM comment_image WHERE like_id = ?",
		likeId,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	res := make([]*sample.CommentImage, 0)
	for rows.Next() {
		r := &sample.CommentImage{}
		if err := rows.Scan(&r.CommentBlogId, &r.CommentUserId, &r.LikeId); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		r.ResetMark()
		res = append(res, r)
	}
	if len(res) > 0 {
		for _, v := range res {
			{
				rel, err := d.comment.Select(ctx, v.CommentBlogId, v.CommentUserId)
				if err != nil {
					return nil, xerrors.Errorf(": %w", err)
				}
				v.Comment = rel
			}
			{
				rel, err := d.like.Select(ctx, v.LikeId)
				if err != nil {
					return nil, xerrors.Errorf(": %w", err)
				}
				v.Like = rel
			}
		}
	}

	return res, nil
}

func (d *CommentImage) Create(ctx context.Context, v *sample.CommentImage) (*sample.CommentImage, error) {
	res, err := d.conn.ExecContext(
		ctx,
		"INSERT INTO `task` (`comment_blog_id`, `comment_user_id`, `like_id`) VALUES (?, ?, ?)", v.CommentBlogId, v.CommentUserId, v.LikeId,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return nil, sql.ErrNoRows
	}

	v = v.Copy()

	v.ResetMark()
	return v, nil
}

func (d *CommentImage) Delete(ctx context.Context, commentBlogId int64, commentUserId int32, likeId uint64) error {
	res, err := d.conn.ExecContext(ctx, "DELETE FROM `comment_image` WHERE `comment_blog_id` = ? AND `comment_user_id` = ? AND `like_id` = ?", commentBlogId, commentUserId, likeId)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *CommentImage) Update(ctx context.Context, v *sample.CommentImage) error {
	if !v.IsChanged() {
		return nil
	}

	changedColumn := v.ChangedColumn()
	cols := make([]string, len(changedColumn)+1)
	values := make([]interface{}, len(changedColumn)+1)
	for i := range changedColumn {
		cols[i] = "`" + changedColumn[i].Name + "` = ?"
		values[i] = changedColumn[i].Value
	}

	query := fmt.Sprintf("UPDATE `comment_image` SET %s WHERE `comment_blog_id` = ? AND `comment_user_id` = ? AND `like_id` = ?", strings.Join(cols, ", "))
	res, err := d.conn.ExecContext(
		ctx,
		query,
		append(values, v.CommentBlogId, v.CommentUserId, v.LikeId)...,
	)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	v.ResetMark()
	return nil
}

type Comment struct {
	conn *sql.DB

	user *User
	blog *Blog
}

func NewComment(conn *sql.DB) *Comment {
	return &Comment{
		conn: conn,
		user: NewUser(conn),
		blog: NewBlog(conn),
	}
}

func (d *Comment) Select(ctx context.Context, blogId int64, userId int32) (*sample.Comment, error) {
	row := d.conn.QueryRowContext(ctx, "SELECT * FROM `comment` WHERE `blog_id` = ? AND `user_id` = ?", blogId, userId)

	v := &sample.Comment{}
	if err := row.Scan(&v.BlogId, &v.UserId); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	{
		rel, err := d.blog.Select(ctx, v.BlogId)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		v.Blog = rel
	}
	{
		rel, err := d.user.Select(ctx, v.UserId)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		v.User = rel
	}

	v.ResetMark()
	return v, nil
}

func (d *Comment) Create(ctx context.Context, v *sample.Comment) (*sample.Comment, error) {
	res, err := d.conn.ExecContext(
		ctx,
		"INSERT INTO `task` (`blog_id`, `user_id`) VALUES (?, ?)", v.BlogId, v.UserId,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return nil, sql.ErrNoRows
	}

	v = v.Copy()

	v.ResetMark()
	return v, nil
}

func (d *Comment) Delete(ctx context.Context, blogId int64, userId int32) error {
	res, err := d.conn.ExecContext(ctx, "DELETE FROM `comment` WHERE `blog_id` = ? AND `user_id` = ?", blogId, userId)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *Comment) Update(ctx context.Context, v *sample.Comment) error {
	if !v.IsChanged() {
		return nil
	}

	changedColumn := v.ChangedColumn()
	cols := make([]string, len(changedColumn)+1)
	values := make([]interface{}, len(changedColumn)+1)
	for i := range changedColumn {
		cols[i] = "`" + changedColumn[i].Name + "` = ?"
		values[i] = changedColumn[i].Value
	}

	query := fmt.Sprintf("UPDATE `comment` SET %s WHERE `blog_id` = ? AND `user_id` = ?", strings.Join(cols, ", "))
	res, err := d.conn.ExecContext(
		ctx,
		query,
		append(values, v.BlogId, v.UserId)...,
	)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	v.ResetMark()
	return nil
}

type Reply struct {
	conn *sql.DB

	comment *Comment
}

func NewReply(conn *sql.DB) *Reply {
	return &Reply{
		conn:    conn,
		comment: NewComment(conn),
	}
}

func (d *Reply) Select(ctx context.Context, id int32) (*sample.Reply, error) {
	row := d.conn.QueryRowContext(ctx, "SELECT * FROM `reply` WHERE `id` = ?", id)

	v := &sample.Reply{}
	if err := row.Scan(&v.Id, &v.CommentBlogId, &v.CommentUserId, &v.Body); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	{
		if v.CommentBlogId != nil && v.CommentUserId != nil {
			rel, err := d.comment.Select(ctx, *v.CommentBlogId, *v.CommentUserId)
			if err != nil {
				return nil, xerrors.Errorf(": %w", err)
			}
			v.Comment = rel
		}
	}

	v.ResetMark()
	return v, nil
}

func (d *Reply) ListByBody(ctx context.Context, body string) ([]*sample.Reply, error) {
	rows, err := d.conn.QueryContext(
		ctx,
		"SELECT * FROM reply WHERE body = ?",
		body,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	res := make([]*sample.Reply, 0)
	for rows.Next() {
		r := &sample.Reply{}
		if err := rows.Scan(&r.Id, &r.CommentBlogId, &r.CommentUserId, &r.Body); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		r.ResetMark()
		res = append(res, r)
	}
	if len(res) > 0 {
		for _, v := range res {
			{
				if v.CommentBlogId != nil && v.CommentUserId != nil {
					rel, err := d.comment.Select(ctx, *v.CommentBlogId, *v.CommentUserId)
					if err != nil {
						return nil, xerrors.Errorf(": %w", err)
					}
					v.Comment = rel
				}
			}
		}
	}

	return res, nil
}

func (d *Reply) Create(ctx context.Context, v *sample.Reply) (*sample.Reply, error) {
	res, err := d.conn.ExecContext(
		ctx,
		"INSERT INTO `task` (`comment_blog_id`, `comment_user_id`, `body`) VALUES (?, ?, ?)", v.CommentBlogId, v.CommentUserId, v.Body,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return nil, sql.ErrNoRows
	}

	v = v.Copy()
	insertedId, err := res.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	v.Id = int32(insertedId)

	v.ResetMark()
	return v, nil
}

func (d *Reply) Delete(ctx context.Context, id int32) error {
	res, err := d.conn.ExecContext(ctx, "DELETE FROM `reply` WHERE `id` = ?", id)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *Reply) Update(ctx context.Context, v *sample.Reply) error {
	if !v.IsChanged() {
		return nil
	}

	changedColumn := v.ChangedColumn()
	cols := make([]string, len(changedColumn)+1)
	values := make([]interface{}, len(changedColumn)+1)
	for i := range changedColumn {
		cols[i] = "`" + changedColumn[i].Name + "` = ?"
		values[i] = changedColumn[i].Value
	}

	query := fmt.Sprintf("UPDATE `reply` SET %s WHERE `id` = ?", strings.Join(cols, ", "))
	res, err := d.conn.ExecContext(
		ctx,
		query,
		append(values, v.Id)...,
	)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	v.ResetMark()
	return nil
}

type Like struct {
	conn *sql.DB

	user *User
	blog *Blog
}

func NewLike(conn *sql.DB) *Like {
	return &Like{
		conn: conn,
		user: NewUser(conn),
		blog: NewBlog(conn),
	}
}

func (d *Like) Select(ctx context.Context, id uint64) (*sample.Like, error) {
	row := d.conn.QueryRowContext(ctx, "SELECT * FROM `like` WHERE `id` = ?", id)

	v := &sample.Like{}
	if err := row.Scan(&v.Id, &v.UserId, &v.BlogId); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	{
		rel, err := d.blog.Select(ctx, v.BlogId)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		v.Blog = rel
	}
	{
		rel, err := d.user.Select(ctx, v.UserId)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		v.User = rel
	}

	v.ResetMark()
	return v, nil
}

func (d *Like) Create(ctx context.Context, v *sample.Like) (*sample.Like, error) {
	res, err := d.conn.ExecContext(
		ctx,
		"INSERT INTO `task` (`user_id`, `blog_id`) VALUES (?, ?)", v.UserId, v.BlogId,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return nil, sql.ErrNoRows
	}

	v = v.Copy()
	insertedId, err := res.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	v.Id = uint64(insertedId)

	v.ResetMark()
	return v, nil
}

func (d *Like) Delete(ctx context.Context, id uint64) error {
	res, err := d.conn.ExecContext(ctx, "DELETE FROM `like` WHERE `id` = ?", id)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *Like) Update(ctx context.Context, v *sample.Like) error {
	if !v.IsChanged() {
		return nil
	}

	changedColumn := v.ChangedColumn()
	cols := make([]string, len(changedColumn)+1)
	values := make([]interface{}, len(changedColumn)+1)
	for i := range changedColumn {
		cols[i] = "`" + changedColumn[i].Name + "` = ?"
		values[i] = changedColumn[i].Value
	}

	query := fmt.Sprintf("UPDATE `like` SET %s WHERE `id` = ?", strings.Join(cols, ", "))
	res, err := d.conn.ExecContext(
		ctx,
		query,
		append(values, v.Id)...,
	)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	v.ResetMark()
	return nil
}

type PostImage struct {
	conn *sql.DB
}

func NewPostImage(conn *sql.DB) *PostImage {
	return &PostImage{
		conn: conn,
	}
}

func (d *PostImage) Select(ctx context.Context, id int32) (*sample.PostImage, error) {
	row := d.conn.QueryRowContext(ctx, "SELECT * FROM `post_image` WHERE `id` = ?", id)

	v := &sample.PostImage{}
	if err := row.Scan(&v.Id, &v.Url); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	v.ResetMark()
	return v, nil
}

func (d *PostImage) Create(ctx context.Context, v *sample.PostImage) (*sample.PostImage, error) {
	res, err := d.conn.ExecContext(
		ctx,
		"INSERT INTO `task` (`id`, `url`) VALUES (?, ?)", v.Id, v.Url,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return nil, sql.ErrNoRows
	}

	v = v.Copy()

	v.ResetMark()
	return v, nil
}

func (d *PostImage) Delete(ctx context.Context, id int32) error {
	res, err := d.conn.ExecContext(ctx, "DELETE FROM `post_image` WHERE `id` = ?", id)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *PostImage) Update(ctx context.Context, v *sample.PostImage) error {
	if !v.IsChanged() {
		return nil
	}

	changedColumn := v.ChangedColumn()
	cols := make([]string, len(changedColumn)+1)
	values := make([]interface{}, len(changedColumn)+1)
	for i := range changedColumn {
		cols[i] = "`" + changedColumn[i].Name + "` = ?"
		values[i] = changedColumn[i].Value
	}

	query := fmt.Sprintf("UPDATE `post_image` SET %s WHERE `id` = ?", strings.Join(cols, ", "))
	res, err := d.conn.ExecContext(
		ctx,
		query,
		append(values, v.Id)...,
	)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	v.ResetMark()
	return nil
}

type Task struct {
	conn *sql.DB

	postImage *PostImage
}

func NewTask(conn *sql.DB) *Task {
	return &Task{
		conn:      conn,
		postImage: NewPostImage(conn),
	}
}

func (d *Task) Select(ctx context.Context, id int32) (*sample.Task, error) {
	row := d.conn.QueryRowContext(ctx, "SELECT * FROM `task` WHERE `id` = ?", id)

	v := &sample.Task{}
	if err := row.Scan(&v.Id, &v.ImageId); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	{
		rel, err := d.postImage.Select(ctx, v.ImageId)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		v.Image = rel
	}

	v.ResetMark()
	return v, nil
}

func (d *Task) ListAll(ctx context.Context) ([]*sample.Task, error) {
	rows, err := d.conn.QueryContext(
		ctx,
		"SELECT * FROM task",
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	res := make([]*sample.Task, 0)
	for rows.Next() {
		r := &sample.Task{}
		if err := rows.Scan(&r.Id, &r.ImageId); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		r.ResetMark()
		res = append(res, r)
	}
	if len(res) > 0 {
		for _, v := range res {
			{
				rel, err := d.postImage.Select(ctx, v.ImageId)
				if err != nil {
					return nil, xerrors.Errorf(": %w", err)
				}
				v.Image = rel
			}
		}
	}

	return res, nil
}

func (d *Task) Create(ctx context.Context, v *sample.Task) (*sample.Task, error) {
	res, err := d.conn.ExecContext(
		ctx,
		"INSERT INTO `task` (`image_id`) VALUES (?)", v.ImageId,
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return nil, sql.ErrNoRows
	}

	v = v.Copy()
	insertedId, err := res.LastInsertId()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	v.Id = int32(insertedId)

	v.ResetMark()
	return v, nil
}

func (d *Task) Delete(ctx context.Context, id int32) error {
	res, err := d.conn.ExecContext(ctx, "DELETE FROM `task` WHERE `id` = ?", id)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *Task) Update(ctx context.Context, v *sample.Task) error {
	if !v.IsChanged() {
		return nil
	}

	changedColumn := v.ChangedColumn()
	cols := make([]string, len(changedColumn)+1)
	values := make([]interface{}, len(changedColumn)+1)
	for i := range changedColumn {
		cols[i] = "`" + changedColumn[i].Name + "` = ?"
		values[i] = changedColumn[i].Value
	}

	query := fmt.Sprintf("UPDATE `task` SET %s WHERE `id` = ?", strings.Join(cols, ", "))
	res, err := d.conn.ExecContext(
		ctx,
		query,
		append(values, v.Id)...,
	)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return xerrors.Errorf(": %w", err)
	} else if n == 0 {
		return sql.ErrNoRows
	}

	v.ResetMark()
	return nil
}
