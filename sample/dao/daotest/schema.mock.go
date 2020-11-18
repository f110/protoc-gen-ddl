// Generated by protoc-ddl.
// protoc-gen-dao-mock: v0.1
package daotest

import (
	"context"

	"go.f110.dev/protoc-ddl/mock"

	"go.f110.dev/protoc-ddl/sample"
	"go.f110.dev/protoc-ddl/sample/dao"
)

type User struct {
	*mock.Mock
}

func NewUser() *User {
	return &User{Mock: mock.New()}
}

func (d *User) Select(ctx context.Context, id int32) (*sample.User, error) {
	v, err := d.Call("Select", map[string]interface{}{"id": id})
	return v.(*sample.User), err
}

func (d *User) RegisterSelect(id int32, value *sample.User) {
	d.Register("Select", map[string]interface{}{"id": id}, value, nil)
}

func (d *User) ListAll(ctx context.Context, opt ...dao.ListOption) ([]*sample.User, error) {
	v, err := d.Call("ListAll", map[string]interface{}{})
	return v.([]*sample.User), err
}

func (d *User) RegisterListAll(value []*sample.User, err error) {
	d.Register("ListAll", map[string]interface{}{}, value, err)
}

func (d *User) ListOverTwenty(ctx context.Context, opt ...dao.ListOption) ([]*sample.User, error) {
	v, err := d.Call("ListOverTwenty", map[string]interface{}{})
	return v.([]*sample.User), err
}

func (d *User) RegisterListOverTwenty(value []*sample.User, err error) {
	d.Register("ListOverTwenty", map[string]interface{}{}, value, err)
}

func (d *User) Create(ctx context.Context, user *sample.User, opt ...dao.ExecOption) (*sample.User, error) {
	_, _ = d.Call("Create", map[string]interface{}{"user": user})
	return user, nil
}

func (d *User) Delete(ctx context.Context, id int32, opt ...dao.ExecOption) error {
	_, _ = d.Call("Delete", map[string]interface{}{"id": id})
	return nil
}

func (d *User) Update(ctx context.Context, user *sample.User, opt ...dao.ExecOption) error {
	_, _ = d.Call("Update", map[string]interface{}{"user": user})
	return nil
}

type Blog struct {
	*mock.Mock
}

func NewBlog() *Blog {
	return &Blog{Mock: mock.New()}
}

func (d *Blog) Select(ctx context.Context, id int64) (*sample.Blog, error) {
	v, err := d.Call("Select", map[string]interface{}{"id": id})
	return v.(*sample.Blog), err
}

func (d *Blog) RegisterSelect(id int64, value *sample.Blog) {
	d.Register("Select", map[string]interface{}{"id": id}, value, nil)
}

func (d *Blog) ListByTitle(ctx context.Context, title string, opt ...dao.ListOption) ([]*sample.Blog, error) {
	v, err := d.Call("ListByTitle", map[string]interface{}{"title": title})
	return v.([]*sample.Blog), err
}

func (d *Blog) RegisterListByTitle(title string, value []*sample.Blog, err error) {
	d.Register("ListByTitle", map[string]interface{}{"title": title}, value, err)
}

func (d *Blog) ListByUserAndCategory(ctx context.Context, userId int32, categoryId int32, opt ...dao.ListOption) ([]*sample.Blog, error) {
	v, err := d.Call("ListByUserAndCategory", map[string]interface{}{"userId": userId, "categoryId": categoryId})
	return v.([]*sample.Blog), err
}

func (d *Blog) RegisterListByUserAndCategory(userId int32, categoryId int32, value []*sample.Blog, err error) {
	d.Register("ListByUserAndCategory", map[string]interface{}{"userId": userId, "categoryId": categoryId}, value, err)
}

func (d *Blog) SelectByUserAndTitle(ctx context.Context, userId int32, title string) (*sample.Blog, error) {
	v, err := d.Call("SelectByUserAndTitle", map[string]interface{}{"userId": userId})
	return v.(*sample.Blog), err
}

func (d *Blog) RegisterSelectByUserAndTitle(userId int32, value *sample.Blog, err error) {
	d.Register("SelectByUserAndTitle", map[string]interface{}{"userId": userId}, value, err)
}

func (d *Blog) Create(ctx context.Context, blog *sample.Blog, opt ...dao.ExecOption) (*sample.Blog, error) {
	_, _ = d.Call("Create", map[string]interface{}{"blog": blog})
	return blog, nil
}

func (d *Blog) Delete(ctx context.Context, id int64, opt ...dao.ExecOption) error {
	_, _ = d.Call("Delete", map[string]interface{}{"id": id})
	return nil
}

func (d *Blog) Update(ctx context.Context, blog *sample.Blog, opt ...dao.ExecOption) error {
	_, _ = d.Call("Update", map[string]interface{}{"blog": blog})
	return nil
}

type CommentImage struct {
	*mock.Mock
}

func NewCommentImage() *CommentImage {
	return &CommentImage{Mock: mock.New()}
}

func (d *CommentImage) Select(ctx context.Context, commentBlogId int64, commentUserId int32, likeId uint64) (*sample.CommentImage, error) {
	v, err := d.Call("Select", map[string]interface{}{"commentBlogId": commentBlogId, "commentUserId": commentUserId, "likeId": likeId})
	return v.(*sample.CommentImage), err
}

func (d *CommentImage) RegisterSelect(commentBlogId int64, commentUserId int32, likeId uint64, value *sample.CommentImage) {
	d.Register("Select", map[string]interface{}{"commentBlogId": commentBlogId, "commentUserId": commentUserId, "likeId": likeId}, value, nil)
}

func (d *CommentImage) ListByLikeId(ctx context.Context, likeId uint64, opt ...dao.ListOption) ([]*sample.CommentImage, error) {
	v, err := d.Call("ListByLikeId", map[string]interface{}{"likeId": likeId})
	return v.([]*sample.CommentImage), err
}

func (d *CommentImage) RegisterListByLikeId(likeId uint64, value []*sample.CommentImage, err error) {
	d.Register("ListByLikeId", map[string]interface{}{"likeId": likeId}, value, err)
}

func (d *CommentImage) Create(ctx context.Context, commentImage *sample.CommentImage, opt ...dao.ExecOption) (*sample.CommentImage, error) {
	_, _ = d.Call("Create", map[string]interface{}{"commentImage": commentImage})
	return commentImage, nil
}

func (d *CommentImage) Delete(ctx context.Context, commentBlogId int64, commentUserId int32, likeId uint64, opt ...dao.ExecOption) error {
	_, _ = d.Call("Delete", map[string]interface{}{"commentBlogId": commentBlogId, "commentUserId": commentUserId, "likeId": likeId})
	return nil
}

func (d *CommentImage) Update(ctx context.Context, commentImage *sample.CommentImage, opt ...dao.ExecOption) error {
	_, _ = d.Call("Update", map[string]interface{}{"commentImage": commentImage})
	return nil
}

type Comment struct {
	*mock.Mock
}

func NewComment() *Comment {
	return &Comment{Mock: mock.New()}
}

func (d *Comment) Select(ctx context.Context, blogId int64, userId int32) (*sample.Comment, error) {
	v, err := d.Call("Select", map[string]interface{}{"blogId": blogId, "userId": userId})
	return v.(*sample.Comment), err
}

func (d *Comment) RegisterSelect(blogId int64, userId int32, value *sample.Comment) {
	d.Register("Select", map[string]interface{}{"blogId": blogId, "userId": userId}, value, nil)
}

func (d *Comment) SelectByUser(ctx context.Context, userId int32) (*sample.Comment, error) {
	v, err := d.Call("SelectByUser", map[string]interface{}{})
	return v.(*sample.Comment), err
}

func (d *Comment) RegisterSelectByUser(value *sample.Comment, err error) {
	d.Register("SelectByUser", map[string]interface{}{}, value, err)
}

func (d *Comment) Create(ctx context.Context, comment *sample.Comment, opt ...dao.ExecOption) (*sample.Comment, error) {
	_, _ = d.Call("Create", map[string]interface{}{"comment": comment})
	return comment, nil
}

func (d *Comment) Delete(ctx context.Context, blogId int64, userId int32, opt ...dao.ExecOption) error {
	_, _ = d.Call("Delete", map[string]interface{}{"blogId": blogId, "userId": userId})
	return nil
}

func (d *Comment) Update(ctx context.Context, comment *sample.Comment, opt ...dao.ExecOption) error {
	_, _ = d.Call("Update", map[string]interface{}{"comment": comment})
	return nil
}

type Reply struct {
	*mock.Mock
}

func NewReply() *Reply {
	return &Reply{Mock: mock.New()}
}

func (d *Reply) Select(ctx context.Context, id int32) (*sample.Reply, error) {
	v, err := d.Call("Select", map[string]interface{}{"id": id})
	return v.(*sample.Reply), err
}

func (d *Reply) RegisterSelect(id int32, value *sample.Reply) {
	d.Register("Select", map[string]interface{}{"id": id}, value, nil)
}

func (d *Reply) ListByBody(ctx context.Context, body string, opt ...dao.ListOption) ([]*sample.Reply, error) {
	v, err := d.Call("ListByBody", map[string]interface{}{"body": body})
	return v.([]*sample.Reply), err
}

func (d *Reply) RegisterListByBody(body string, value []*sample.Reply, err error) {
	d.Register("ListByBody", map[string]interface{}{"body": body}, value, err)
}

func (d *Reply) Create(ctx context.Context, reply *sample.Reply, opt ...dao.ExecOption) (*sample.Reply, error) {
	_, _ = d.Call("Create", map[string]interface{}{"reply": reply})
	return reply, nil
}

func (d *Reply) Delete(ctx context.Context, id int32, opt ...dao.ExecOption) error {
	_, _ = d.Call("Delete", map[string]interface{}{"id": id})
	return nil
}

func (d *Reply) Update(ctx context.Context, reply *sample.Reply, opt ...dao.ExecOption) error {
	_, _ = d.Call("Update", map[string]interface{}{"reply": reply})
	return nil
}

type Like struct {
	*mock.Mock
}

func NewLike() *Like {
	return &Like{Mock: mock.New()}
}

func (d *Like) Select(ctx context.Context, id uint64) (*sample.Like, error) {
	v, err := d.Call("Select", map[string]interface{}{"id": id})
	return v.(*sample.Like), err
}

func (d *Like) RegisterSelect(id uint64, value *sample.Like) {
	d.Register("Select", map[string]interface{}{"id": id}, value, nil)
}

func (d *Like) Create(ctx context.Context, like *sample.Like, opt ...dao.ExecOption) (*sample.Like, error) {
	_, _ = d.Call("Create", map[string]interface{}{"like": like})
	return like, nil
}

func (d *Like) Delete(ctx context.Context, id uint64, opt ...dao.ExecOption) error {
	_, _ = d.Call("Delete", map[string]interface{}{"id": id})
	return nil
}

func (d *Like) Update(ctx context.Context, like *sample.Like, opt ...dao.ExecOption) error {
	_, _ = d.Call("Update", map[string]interface{}{"like": like})
	return nil
}

type PostImage struct {
	*mock.Mock
}

func NewPostImage() *PostImage {
	return &PostImage{Mock: mock.New()}
}

func (d *PostImage) Select(ctx context.Context, id int32) (*sample.PostImage, error) {
	v, err := d.Call("Select", map[string]interface{}{"id": id})
	return v.(*sample.PostImage), err
}

func (d *PostImage) RegisterSelect(id int32, value *sample.PostImage) {
	d.Register("Select", map[string]interface{}{"id": id}, value, nil)
}

func (d *PostImage) Create(ctx context.Context, postImage *sample.PostImage, opt ...dao.ExecOption) (*sample.PostImage, error) {
	_, _ = d.Call("Create", map[string]interface{}{"postImage": postImage})
	return postImage, nil
}

func (d *PostImage) Delete(ctx context.Context, id int32, opt ...dao.ExecOption) error {
	_, _ = d.Call("Delete", map[string]interface{}{"id": id})
	return nil
}

func (d *PostImage) Update(ctx context.Context, postImage *sample.PostImage, opt ...dao.ExecOption) error {
	_, _ = d.Call("Update", map[string]interface{}{"postImage": postImage})
	return nil
}

type Task struct {
	*mock.Mock
}

func NewTask() *Task {
	return &Task{Mock: mock.New()}
}

func (d *Task) Select(ctx context.Context, id int32) (*sample.Task, error) {
	v, err := d.Call("Select", map[string]interface{}{"id": id})
	return v.(*sample.Task), err
}

func (d *Task) RegisterSelect(id int32, value *sample.Task) {
	d.Register("Select", map[string]interface{}{"id": id}, value, nil)
}

func (d *Task) ListAll(ctx context.Context, opt ...dao.ListOption) ([]*sample.Task, error) {
	v, err := d.Call("ListAll", map[string]interface{}{})
	return v.([]*sample.Task), err
}

func (d *Task) RegisterListAll(value []*sample.Task, err error) {
	d.Register("ListAll", map[string]interface{}{}, value, err)
}

func (d *Task) Create(ctx context.Context, task *sample.Task, opt ...dao.ExecOption) (*sample.Task, error) {
	_, _ = d.Call("Create", map[string]interface{}{"task": task})
	return task, nil
}

func (d *Task) Delete(ctx context.Context, id int32, opt ...dao.ExecOption) error {
	_, _ = d.Call("Delete", map[string]interface{}{"id": id})
	return nil
}

func (d *Task) Update(ctx context.Context, task *sample.Task, opt ...dao.ExecOption) error {
	_, _ = d.Call("Update", map[string]interface{}{"task": task})
	return nil
}
