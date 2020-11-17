// Generated by protoc-ddl.
// protoc-gen-entity: v0.1
package sample

import (
	"bytes"
	"sync"
	"time"

	"go.f110.dev/protoc-ddl"
)

var _ = time.Time{}
var _ = bytes.Buffer{}

type Column struct {
	Name  string
	Value interface{}
}

type User struct {
	Id        int32
	Age       int32
	Name      string
	Title     string
	CreatedAt time.Time

	mu   sync.Mutex
	mark *User
}

func (e *User) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *User) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Age != e.mark.Age ||
		e.Name != e.mark.Name ||
		e.Title != e.mark.Title ||
		!e.CreatedAt.Equal(e.mark.CreatedAt)
}

func (e *User) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.Age != e.mark.Age {
		res = append(res, ddl.Column{Name: "age", Value: e.Age})
	}
	if e.Name != e.mark.Name {
		res = append(res, ddl.Column{Name: "name", Value: e.Name})
	}
	if e.Title != e.mark.Title {
		res = append(res, ddl.Column{Name: "title", Value: e.Title})
	}
	if !e.CreatedAt.Equal(e.mark.CreatedAt) {
		res = append(res, ddl.Column{Name: "created_at", Value: e.CreatedAt})
	}

	return res
}

func (e *User) Copy() *User {
	n := &User{
		Id:        e.Id,
		Age:       e.Age,
		Name:      e.Name,
		Title:     e.Title,
		CreatedAt: e.CreatedAt,
	}

	return n
}

type Blog struct {
	Id         int64
	UserId     int32
	Title      string
	Body       string
	CategoryId *int32
	Attach     []byte
	EditorId   int32
	Sign       []byte
	CreatedAt  time.Time
	UpdatedAt  *time.Time

	User   *User
	Editor *User

	mu   sync.Mutex
	mark *Blog
}

func (e *Blog) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *Blog) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.UserId != e.mark.UserId ||
		e.Title != e.mark.Title ||
		e.Body != e.mark.Body ||
		((e.CategoryId != nil && (e.mark.CategoryId == nil || *e.CategoryId != *e.mark.CategoryId)) || e.CategoryId == nil && e.mark.CategoryId != nil) ||
		!bytes.Equal(e.Attach, e.mark.Attach) ||
		e.EditorId != e.mark.EditorId ||
		!bytes.Equal(e.Sign, e.mark.Sign) ||
		!e.CreatedAt.Equal(e.mark.CreatedAt) ||
		((e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil))
}

func (e *Blog) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.UserId != e.mark.UserId {
		res = append(res, ddl.Column{Name: "user_id", Value: e.UserId})
	}
	if e.Title != e.mark.Title {
		res = append(res, ddl.Column{Name: "title", Value: e.Title})
	}
	if e.Body != e.mark.Body {
		res = append(res, ddl.Column{Name: "body", Value: e.Body})
	}
	if (e.CategoryId != nil && (e.mark.CategoryId == nil || *e.CategoryId != *e.mark.CategoryId)) || (e.CategoryId == nil && e.mark.CategoryId != nil) {
		if e.CategoryId != nil {
			res = append(res, ddl.Column{Name: "category_id", Value: *e.CategoryId})
		} else {
			res = append(res, ddl.Column{Name: "category_id", Value: nil})
		}
	}
	if !bytes.Equal(e.Attach, e.mark.Attach) {
		res = append(res, ddl.Column{Name: "attach", Value: e.Attach})
	}
	if e.EditorId != e.mark.EditorId {
		res = append(res, ddl.Column{Name: "editor_id", Value: e.EditorId})
	}
	if !bytes.Equal(e.Sign, e.mark.Sign) {
		res = append(res, ddl.Column{Name: "sign", Value: e.Sign})
	}
	if !e.CreatedAt.Equal(e.mark.CreatedAt) {
		res = append(res, ddl.Column{Name: "created_at", Value: e.CreatedAt})
	}
	if (e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil) {
		if e.UpdatedAt != nil {
			res = append(res, ddl.Column{Name: "updated_at", Value: *e.UpdatedAt})
		} else {
			res = append(res, ddl.Column{Name: "updated_at", Value: nil})
		}
	}

	return res
}

func (e *Blog) Copy() *Blog {
	n := &Blog{
		Id:        e.Id,
		UserId:    e.UserId,
		Title:     e.Title,
		Body:      e.Body,
		Attach:    e.Attach,
		EditorId:  e.EditorId,
		Sign:      e.Sign,
		CreatedAt: e.CreatedAt,
	}
	if e.CategoryId != nil {
		v := *e.CategoryId
		n.CategoryId = &v
	}
	if e.UpdatedAt != nil {
		v := *e.UpdatedAt
		n.UpdatedAt = &v
	}

	if e.Editor != nil {
		n.Editor = e.Editor.Copy()
	}
	if e.User != nil {
		n.User = e.User.Copy()
	}

	return n
}

type CommentImage struct {
	CommentBlogId int64
	CommentUserId int32
	LikeId        uint64

	Comment *Comment
	Like    *Like

	mu   sync.Mutex
	mark *CommentImage
}

func (e *CommentImage) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *CommentImage) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return false
}

func (e *CommentImage) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)

	return res
}

func (e *CommentImage) Copy() *CommentImage {
	n := &CommentImage{
		CommentBlogId: e.CommentBlogId,
		CommentUserId: e.CommentUserId,
		LikeId:        e.LikeId,
	}

	if e.Comment != nil {
		n.Comment = e.Comment.Copy()
	}
	if e.Like != nil {
		n.Like = e.Like.Copy()
	}

	return n
}

type Comment struct {
	BlogId int64
	UserId int32

	Blog *Blog
	User *User

	mu   sync.Mutex
	mark *Comment
}

func (e *Comment) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *Comment) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return false
}

func (e *Comment) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)

	return res
}

func (e *Comment) Copy() *Comment {
	n := &Comment{
		BlogId: e.BlogId,
		UserId: e.UserId,
	}

	if e.Blog != nil {
		n.Blog = e.Blog.Copy()
	}
	if e.User != nil {
		n.User = e.User.Copy()
	}

	return n
}

type Reply struct {
	Id            int32
	CommentBlogId *int64
	CommentUserId *int32
	Body          string

	Comment *Comment

	mu   sync.Mutex
	mark *Reply
}

func (e *Reply) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *Reply) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return ((e.CommentBlogId != nil && (e.mark.CommentBlogId == nil || *e.CommentBlogId != *e.mark.CommentBlogId)) || e.CommentBlogId == nil && e.mark.CommentBlogId != nil) ||
		((e.CommentUserId != nil && (e.mark.CommentUserId == nil || *e.CommentUserId != *e.mark.CommentUserId)) || e.CommentUserId == nil && e.mark.CommentUserId != nil) ||
		e.Body != e.mark.Body
}

func (e *Reply) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if (e.CommentBlogId != nil && (e.mark.CommentBlogId == nil || *e.CommentBlogId != *e.mark.CommentBlogId)) || (e.CommentBlogId == nil && e.mark.CommentBlogId != nil) {
		if e.CommentBlogId != nil {
			res = append(res, ddl.Column{Name: "comment_blog_id", Value: *e.CommentBlogId})
		} else {
			res = append(res, ddl.Column{Name: "comment_blog_id", Value: nil})
		}
	}
	if (e.CommentUserId != nil && (e.mark.CommentUserId == nil || *e.CommentUserId != *e.mark.CommentUserId)) || (e.CommentUserId == nil && e.mark.CommentUserId != nil) {
		if e.CommentUserId != nil {
			res = append(res, ddl.Column{Name: "comment_user_id", Value: *e.CommentUserId})
		} else {
			res = append(res, ddl.Column{Name: "comment_user_id", Value: nil})
		}
	}
	if e.Body != e.mark.Body {
		res = append(res, ddl.Column{Name: "body", Value: e.Body})
	}

	return res
}

func (e *Reply) Copy() *Reply {
	n := &Reply{
		Id:   e.Id,
		Body: e.Body,
	}
	if e.CommentBlogId != nil {
		v := *e.CommentBlogId
		n.CommentBlogId = &v
	}
	if e.CommentUserId != nil {
		v := *e.CommentUserId
		n.CommentUserId = &v
	}

	if e.Comment != nil {
		n.Comment = e.Comment.Copy()
	}

	return n
}

type Like struct {
	Id     uint64
	UserId int32
	BlogId int64

	User *User
	Blog *Blog

	mu   sync.Mutex
	mark *Like
}

func (e *Like) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *Like) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.UserId != e.mark.UserId ||
		e.BlogId != e.mark.BlogId
}

func (e *Like) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.UserId != e.mark.UserId {
		res = append(res, ddl.Column{Name: "user_id", Value: e.UserId})
	}
	if e.BlogId != e.mark.BlogId {
		res = append(res, ddl.Column{Name: "blog_id", Value: e.BlogId})
	}

	return res
}

func (e *Like) Copy() *Like {
	n := &Like{
		Id:     e.Id,
		UserId: e.UserId,
		BlogId: e.BlogId,
	}

	if e.Blog != nil {
		n.Blog = e.Blog.Copy()
	}
	if e.User != nil {
		n.User = e.User.Copy()
	}

	return n
}

type PostImage struct {
	Id  int32
	Url string

	mu   sync.Mutex
	mark *PostImage
}

func (e *PostImage) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *PostImage) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Url != e.mark.Url
}

func (e *PostImage) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.Url != e.mark.Url {
		res = append(res, ddl.Column{Name: "url", Value: e.Url})
	}

	return res
}

func (e *PostImage) Copy() *PostImage {
	n := &PostImage{
		Id:  e.Id,
		Url: e.Url,
	}

	return n
}

type Task struct {
	Id      int32
	ImageId int32

	Image *PostImage

	mu   sync.Mutex
	mark *Task
}

func (e *Task) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *Task) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.ImageId != e.mark.ImageId
}

func (e *Task) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.ImageId != e.mark.ImageId {
		res = append(res, ddl.Column{Name: "image_id", Value: e.ImageId})
	}

	return res
}

func (e *Task) Copy() *Task {
	n := &Task{
		Id:      e.Id,
		ImageId: e.ImageId,
	}

	if e.Image != nil {
		n.Image = e.Image.Copy()
	}

	return n
}
