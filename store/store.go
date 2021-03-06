package store

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"time"
)

// About entity type
// 0 : post
// 1 : comment

type EntityType int

const (
	Post EntityType = iota
	Comment
)

type Pager struct {
	SinceId int64
	Limit   int
}

type AdminStore interface {
	// counter - user/active/discussion
	GetTotalUser() (int, error)
	GetNewUser(day time.Time) (int, error)
	GetNewUserDaily(from, to time.Time) ([]*model.DailyCount, error)

	GetUserActive(day time.Time) (int, error)
	GetUserActiveDaily(from, to time.Time) ([]*model.DailyCount, error)

	GetTotalDiscussion() (int, error)
	GetNewDiscussion(day time.Time) (int, error)
	GetNewDiscussionDaily(from, to time.Time) ([]*model.DailyCount, error)

	// search
	SearchUser(params model.QueryParams, page, size int) ([]*model.User, error)
	SearchUserCount(params model.QueryParams) (int, error)

	SearchDiscussion(params model.QueryParams, page, size int) ([]*model.Discussion, error)
	SearchDiscussionCount(params model.QueryParams) (int, error)

	SearchPost(params model.QueryParams, page, size int) ([]*model.Post, error)
	SearchPostCount(params model.QueryParams) (int, error)

	SearchReport(params model.QueryParams, page, size int) ([]*model.Report, error)
	SearchReportCount(params model.QueryParams) (int, error)

	SearchSignUser(page, size int) ([]*model.User, error)
	SearchSignUserCount(page, size int) (int, error)
}

type UserStore interface {
	GetUser(id int64) (*model.User, error)
	GetUserIdList(ids []int64) ([]*model.User, error)
	GetUserByLogin(string) (*model.User, error)
	GetUserByPhone(string) (*model.User, error)
	CreateUser(usr *model.User) error
	UpdateUser(usr *model.User) error
	DeleteUser(id int64) error

	// bind
	GetUserByBind(string) (*model.User, error)
	GetUserByUnion(string) (*model.User, error)
	CreateBind(bind *model.UserBind) error
	DeleteBind(bindId string) error
	DeleteBindByUser(uid int64) error
	DeleteBindUser(uid int64) error
	CreateBindUser(user *model.User, bind *model.UserBind) error

	// user-exp-sign
	GetUserListSign(page int, size int) ([]*model.User, error)
	GetUserListExp(page int, size int) ([]*model.User, error)
	UpdateUserExp(uid int64, exp int) error
	UpdateUserSign(uid int64, count int) error
}

type DiscussionStore interface {
	GetDiscussionList(page, size int, filter string) ([]*model.Discussion, error)
	GetDiscussionListCommentCount(page int, size int) ([]*model.Discussion, error)
	GetDiscussionListUser(uid int64, page, size int) ([]*model.Discussion, error)
	GetDiscussionListByIds(ids []int64) ([]*model.Discussion, error)
	GetDiscussion(id int64) (*model.Discussion, error)
	CreateDiscussion(p *model.Discussion) error
	UpdateDiscussion(p *model.Discussion) error
	DeleteDiscussion(id int64) error
}

type PostStore interface {
	GetPostList(did int64, page, size int) ([]*model.Post, error)
	GetPostListUser(uid int64, page, size int) ([]*model.Post, error)
	GetPostListByIds(ids []int64) ([]*model.Post, error)
	GetPostCount(pid int64) (int, error)
	GetPost(id int64) (*model.Post, error)

	CreatePost(p *model.Post) error
	UpdatePost(p *model.Post) error
	DeletePost(id int64) error
}

type LikeStore interface {
	GetLikeList(pid int64, page, size int) ([]*model.Like, error)
	GetLikeListUser(uid int64, page, size int) ([]*model.Like, error)
	GetLikeCount(pid int64) (int, error)
	GetLike(pid, uid int64) (*model.Like, error)
	GetLikeId(id int64) (*model.Like, error)
	CreateLike(f *model.Like) error
	UpdateLike(f *model.Like) error
	DeleteLike(int64) error

	GetLikePostList(uid int64, pids []int64) ([]int64, error)
}

type TagStore interface {
	GetDiscussionsByTag(tag string, page, size int) ([]*model.Discussion, error)
	GetTagList() ([]*model.Tag, error)
	GetTagId(id int64) (*model.Tag, error)
	LinkTagDiscussion(pid int64, tag []string) error
	CreateTag(tag, summary string) (*model.Tag, error)
	DeleteTag(id int64) error
	UpdateTag(t *model.Tag) error
}

type MessageStore interface {
	// notification
	GetNotificationCount(to int64) (*model.MessageCount, error)
	GetNotificationById(id int64) (*model.Notification, error)
	GetNotificationListType(to int64, mt model.NotType, page, size int) ([]*model.Notification, error)
	CreateNotification(n *model.Notification) error
	SetNotificationReadId(uid, id int64) error
	SetNotificationReadType(to int64, mt model.NotType) error

	// chat
	GetChatMsgList(from, to int64, page, size int) ([]*model.Chat, error)
	GetChatMessageById(id int64) (*model.Chat, error)
	CreateChatMessage(m *model.Chat) error
	SetChatMsgAsRead(from, to int64) error
	GetChatUserList(to int64, page, size int) ([]*model.Chat, error)
}

type MetaStore interface {
	GetMetaData() (map[string]string, error)
	GetMetaValue(key string) (string, error)
	PutMetaData(map[string]interface{}) error
}

type MediaStore interface {
	CreateMedia(m *model.Media) error
	GetMediaListByPostIds(pids []int64) ([]*model.Media, error)
	GetMediaByPostId(pid int64) (*model.Media, error)
	DeleteMediaByPostId(pid int64) error
}

type ReportStore interface {
	CreateReport(rpt *model.Report) error
	GetReport(id int64) (*model.Report, error)
	UpdateReport(*model.Report) error
}

type Store interface {
	AdminStore

	// user
	UserStore

	// discussion
	DiscussionStore

	// posts
	PostStore

	// like
	LikeStore

	// tag
	TagStore

	// message
	MessageStore

	// meta-data
	MetaStore

	// images
	MediaStore

	// report
	ReportStore
}

// user
func GetUser(c echo.Context, id int64) (*model.User, error) {
	return FromContext(c).GetUser(id)
}

func GetUserLogin(c echo.Context, login string) (*model.User, error) {
	return FromContext(c).GetUserByLogin(login)
}

func GetUserPhone(c echo.Context, phone string) (*model.User, error)  {
	return FromContext(c).GetUserByPhone(phone)
}

func GetUserBind(c echo.Context, bind string) (*model.User, error) {
	return FromContext(c).GetUserByBind(bind)
}

func GetUserUnion(c echo.Context, bind string) (*model.User, error) {
	return FromContext(c).GetUserByUnion(bind)
}

func CreateUser(c echo.Context, usr *model.User) error {
	return FromContext(c).CreateUser(usr)
}

func CreateBind(c echo.Context, bind *model.UserBind) error {
	return FromContext(c).CreateBind(bind)
}

func DeleteBind(c echo.Context, bindId string) error {
	return FromContext(c).DeleteBind(bindId)
}

func UpdateUser(c echo.Context, usr *model.User) error {
	return FromContext(c).UpdateUser(usr)
}

func DeleteUser(c echo.Context, id int64) error {
	return FromContext(c).DeleteUser(id)
}

// discussion
func GetDiscussionList(c echo.Context, page, size int, filter string) ([]*model.Discussion, error) {
	return FromContext(c).GetDiscussionList(page, size, filter)
}

func GetDiscussion(c echo.Context, id int64) (*model.Discussion, error) {
	return FromContext(c).GetDiscussion(id)
}

func CreateDiscussion(c echo.Context, p *model.Discussion) error {
	return FromContext(c).CreateDiscussion(p)
}

func UpdateDiscussion(c echo.Context, p *model.Discussion) error {
	return FromContext(c).UpdateDiscussion(p)
}

func DeleteDiscussion(c echo.Context, id int64) error {
	return FromContext(c).DeleteDiscussion(id)
}

// posts
func GetPostList(c echo.Context, pid int64, page, size int) ([]*model.Post, error) {
	return FromContext(c).GetPostList(pid, page, size)
}

func GetPost(c echo.Context, id int64) (*model.Post, error) {
	return FromContext(c).GetPost(id)
}

func CreatePost(c echo.Context, cmt *model.Post) error {
	return FromContext(c).CreatePost(cmt)
}

func UpdatePost(c echo.Context, cmt *model.Post) error {
	return FromContext(c).UpdatePost(cmt)
}

func DeletePost(c echo.Context, id int64) error {
	return FromContext(c).DeletePost(id)
}

// like
func GetLikeList(c echo.Context, pid int64, page, size int) ([]*model.Like, error) {
	return FromContext(c).GetLikeList(pid, page, size)
}

func GetFavorCount(c echo.Context, pid int64) (int, error) {
	return FromContext(c).GetLikeCount(pid)
}

func CreateFavor(c echo.Context, f *model.Like) error {
	return FromContext(c).CreateLike(f)
}

func DeleteFavor(c echo.Context, id int64) error {
	return FromContext(c).DeleteLike(id)
}

// tags
func GetTagList(c echo.Context) ([]*model.Tag, error) {
	return FromContext(c).GetTagList()
}

func GetDiscussionByTag(c echo.Context, tag string, page, size int) ([]*model.Discussion, error) {
	return FromContext(c).GetDiscussionsByTag(tag, page, size)
}

func LinkTagPost(c echo.Context, pid int64, tags []string) error {
	return FromContext(c).LinkTagDiscussion(pid, tags)
}

// message
func GetNotificationCount(c echo.Context, to int64) (*model.MessageCount, error) {
	return FromContext(c).GetNotificationCount(to)
}

func GetNotificationListType(c echo.Context, to int64, mt model.NotType, page, size int) ([]*model.Notification, error) {
	return FromContext(c).GetNotificationListType(to, mt, page, size)
}

func SetNotificationReadId(c echo.Context, uid, id int64) error {
	return FromContext(c).SetNotificationReadId(uid, id)
}

func GetNotificationId(c echo.Context, id int64) (*model.Notification, error) {
	return FromContext(c).GetNotificationById(id)
}

func SetNotificationReadType(c echo.Context, to int64, mt model.NotType) error {
	return FromContext(c).SetNotificationReadType(to, mt)
}

func SetNotificationReadFromTo(c echo.Context, from, to int64) error {
	return FromContext(c).SetChatMsgAsRead(from, to)
}

func CreateNotification(c echo.Context, n *model.Notification) error {
	return FromContext(c).CreateNotification(n)
}

// meta
func GetMetaValue(c echo.Context, key string) (string, error) {
	return FromContext(c).GetMetaValue(key)
}

// counter

func GetDailyActiveUser(c echo.Context, from, to time.Time) ([]*model.DailyCount, error) {
	return FromContext(c).GetUserActiveDaily(from, to)
}

func GetDailyNewUser(c echo.Context, from, to time.Time) ([]*model.DailyCount, error) {
	return FromContext(c).GetNewUserDaily(from, to)
}

func GetDailyNewDiscussion(c echo.Context, from, to time.Time) ([]*model.DailyCount, error)  {
	return FromContext(c).GetNewDiscussionDaily(from, to)
}

func GetDailyNewComment(c echo.Context, from, to time.Time) ([]*model.DailyCount, error) {
	return nil, nil
}

func GetDailyNewFavor(c echo.Context, from, to time.Time) ([]*model.DailyCount, error)  {
	return nil, nil
}

// search
func SearchUser(c echo.Context, params model.QueryParams, page, size int) ([]*model.User, error) {
	return FromContext(c).SearchUser(params, page, size)
}

func SearchUserCount(c echo.Context, params model.QueryParams) (int, error) {
	return FromContext(c).SearchUserCount(params)
}

func SearchDiscussion(c echo.Context, params model.QueryParams, page, size int) ([]*model.Discussion, error) {
	return FromContext(c).SearchDiscussion(params, page, size)
}

func SearchDiscussionCount(c echo.Context, params model.QueryParams) (int, error) {
	return FromContext(c).SearchDiscussionCount(params)
}

func SearchPost(c echo.Context, params model.QueryParams, page, size int) ([]*model.Post, error) {
	return FromContext(c).SearchPost(params, page, size)
}

func SearchPostCount(c echo.Context, params model.QueryParams) (int, error) {
	return FromContext(c).SearchPostCount(params)
}

func SearchReport(c echo.Context, params model.QueryParams, page, size int) ([]*model.Report, error)  {
	return FromContext(c).SearchReport(params, page, size)
}

func SearchReportCount(c echo.Context, params model.QueryParams) (int, error) {
	return FromContext(c).SearchReportCount(params)
}
