package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/events"
	"github.com/kawaapp/kawaqing/spam"

	"net/http"
	"strconv"
	"database/sql"
)

// comment/reply


// user - api

func GetPostList(c echo.Context) error {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	page, size := getPageSize(c)

	posts, err := store.GetPostList(c, int64(pid), page, size)
	if err != nil {
		return err
	}

	p := makePayload(0, posts)
	if len(posts) == size {
		p.HasMore = true
	}

	if includes(c, "user") {
		attachUserToPost(c, posts, p)
	}
	if includes(c, "like") {
		attachLikeToPost(c, posts, p)
	}
	return c.JSON(200, p)
}

func GetPost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	post, err := store.GetPost(c, int64(id))
	if err != nil {
		return c.NoContent(404)
	}

	p := makePayload(0, post)
	if includes(c, "user") {
		attachUserToPost(c, []*model.Post{ post }, p)
	}
	if includes(c, "like") {
		attachLikeToPost(c, []*model.Post{ post }, p)
	}
	return c.JSON(200, p)
}

func CreatePost(c echo.Context) error {
	post := new(model.Post)
	if err := c.Bind(post); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// ensure author is valid
	user := session.User(c)
	if user.Login != "Admin" {
		if user.Silenced() {
			return c.String(http.StatusForbidden, "!user silenced")
		}
		post.AuthorID = user.ID
	}

	// spam test
	if checker := spam.FromContext(c); checker != nil {
		if ok, _ := checker.Validate(post.Content); !ok {
			return c.String(400, "4001:sensitive content!")
		}
	}

	// insert db
	if err := store.CreatePost(c, post); err != nil {
		return err
	}

	// publish comment message
	events.Dispatch(ePostCreated, c, post)

	return c.JSON(200, post)
}

func UpdatePost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	// ensure author is valid
	user := session.User(c)
	if user.Login != "Admin" {
		return c.NoContent(http.StatusForbidden)
	}

	cmt := new(model.Post)
	if err := c.Bind(cmt); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	cmt.ID = int64(id)
	if err := store.UpdatePost(c, cmt); err != nil {
		return err
	}
	return c.JSON(200, cmt)
}

func DeletePost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	cmt, err := store.GetPost(c, int64(id))
	if err == sql.ErrNoRows {
		return c.NoContent(404)
	} else if err != nil {
		return err
	}
	// ensure author is valid
	if user := session.User(c); user.Login != "Admin" && user.ID != cmt.AuthorID {
		return c.NoContent(http.StatusForbidden)
	}
	if err := store.DeletePost(c, int64(id)); err != nil {
		return err
	}
	events.Dispatch(ePostDeleted, c, cmt)
	return c.NoContent(200)
}

// uid(0) is the current user
func GetPostByUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	var user *model.User
	if id == 0 {
		user = session.User(c)
	} else {
		user, err = store.GetUser(c, int64(id))
	}

	if err == sql.ErrNoRows {
		return c.String(404, "user not found")
	} else if err != nil {
		return err
	}

	page, size := getPageSize(c)
	posts, err := store.FromContext(c).GetPostListUser(user.ID, page, size)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	p := makePayload(0, posts)
	if includes(c, "user") {
		attachUserToPost(c, posts, p)
	}
	if includes(c, "like") {
		attachLikeToPost(c, posts, p)
	}
	return c.JSON(200, p)
}

func attachUserToPost(c echo.Context, posts []*model.Post, p payload) {
	Ids := make([]int64, len(posts) * 2)
	kv := make(map[int64]*model.User)

	for i, v := range posts {
		Ids[i * 2 + 0] = v.AuthorID
		Ids[i * 2 + 1] = v.ReplyID
	}

	from, _ := store.FromContext(c).GetUserIdList(Ids)
	if from != nil {
		for _, v := range from {
			kv[v.ID] = v
		}
	}
	p.Entities["users"] = kv
}

func attachLikeToPost(c echo.Context, posts []*model.Post, p payload) {
	ids := make([]int64, len(posts))
	kv := make(map[int64]int64)

	// favored state
	for i, v := range posts {
		ids[i] = v.ID
	}
	var uid int64
	if user := session.User(c); user != nil {
		uid = user.ID
	}

	favors, _ := store.FromContext(c).GetLikePostList(uid, ids)
	if len(favors) > 0 {
		for _, v := range favors {
			kv[v] = v
		}
	}
	p.Entities["likes"] = kv
}

// events
func postOnLikeChanged(c echo.Context, v interface{}, getCount func(f *model.Like, num int) int) error {
	favor, ok := v.(*model.Like)
	if !ok {
		return typeError("Favor")
	}

	cmt, err := store.FromContext(c).GetPost(favor.PostId)
	if err != nil {
		return err
	}
	cmt.LikeCount = getCount(favor, cmt.LikeCount)
	err = store.UpdatePost(c, cmt)
	return err
}

func postOnLikeCreated(c echo.Context, v interface{}) error {
	return postOnLikeChanged(c, v, func(f *model.Like, num int) int {
		return num + 1
	})
}

func postOnLikeUpdated(c echo.Context, v interface{}) error {
	return postOnLikeChanged(c, v, func(f *model.Like, num int) int {
		if f.Status == 0 {
			return Max(num - 1, 0)
		} else {
			return num + 1
		}
	})
}

func init() {
	events.Subscribe(eLikeCreated, postOnLikeCreated)
	events.Subscribe(eLikeUpdated, postOnLikeUpdated)
}