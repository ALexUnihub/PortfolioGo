package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"myRedditClone/pkg/author"
	"myRedditClone/pkg/comments"
	"myRedditClone/pkg/posts"
	"myRedditClone/pkg/session"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"gopkg.in/mgo.v2/bson"
)

var (
	ResultPosts = []*posts.Post{
		{
			IDBson:     bson.NewObjectId(),
			AuthorBson: "DeepThought",
			Author: &author.Author{
				UserID:   "1st user",
				UserName: "DeepThought",
			},
			Category:         "programming",
			CreationData:     "2007-01-01T04:20:00",
			PostID:           "1stPost",
			Score:            4,
			Data:             "42 ?",
			Title:            "The answer to the Ultimate Question of Life, the Universe, and Everything",
			Type:             "text",
			UpvotePersentage: 100,
			Views:            1500,
			Votes: []*posts.Vote{
				{
					UserID: "1stUser",
					Vote:   1,
				},
				{
					UserID: "2ndUser",
					Vote:   1,
				},
				{
					UserID: "3rdUser",
					Vote:   1,
				},
				{
					UserID: "4thUser",
					Vote:   1,
				},
			},
		},
		{
			IDBson:     bson.NewObjectId(),
			AuthorBson: "DeepThought",
			Author: &author.Author{
				UserID:   "1st user",
				UserName: "DeepThought",
			},
			Category:         "programming",
			CreationData:     "2007-01-01T04:20:00",
			PostID:           "2ndPost",
			Score:            3,
			Data:             "/localhost",
			Title:            "2nd post",
			Type:             "link",
			UpvotePersentage: 100,
			Views:            250,
			Votes: []*posts.Vote{
				{
					UserID: "1stUser",
					Vote:   1,
				},
				{
					UserID: "2ndUser",
					Vote:   1,
				},
				{
					UserID: "3rdUser",
					Vote:   1,
				},
			},
		},
	}

	UpVotedPost = &posts.Post{
		IDBson:     bson.NewObjectId(),
		AuthorBson: "DeepThought",
		Author: &author.Author{
			UserID:   "1st user",
			UserName: "DeepThought",
		},
		Category:         "programming",
		CreationData:     "2007-01-01T04:20:00",
		PostID:           "1stPost",
		Score:            4,
		Data:             "42 ?",
		Title:            "The answer to the Ultimate Question of Life, the Universe, and Everything",
		Type:             "text",
		UpvotePersentage: 100,
		Views:            1500,
		Votes: []*posts.Vote{
			{
				UserID: "1stUser",
				Vote:   1,
			},
			{
				UserID: "2ndUser",
				Vote:   1,
			},
			{
				UserID: "3rdUser",
				Vote:   1,
			},
			{
				UserID: "4thUser",
				Vote:   1,
			},
			{
				UserID: "132",
				Vote:   1,
			},
		},
	}

	DownVotedPost = &posts.Post{
		IDBson:     bson.NewObjectId(),
		AuthorBson: "DeepThought",
		Author: &author.Author{
			UserID:   "1st user",
			UserName: "DeepThought",
		},
		Category:         "programming",
		CreationData:     "2007-01-01T04:20:00",
		PostID:           "1stPost",
		Score:            4,
		Data:             "42 ?",
		Title:            "The answer to the Ultimate Question of Life, the Universe, and Everything",
		Type:             "text",
		UpvotePersentage: 100,
		Views:            1500,
		Votes: []*posts.Vote{
			{
				UserID: "1stUser",
				Vote:   1,
			},
			{
				UserID: "2ndUser",
				Vote:   1,
			},
			{
				UserID: "3rdUser",
				Vote:   1,
			},
			{
				UserID: "4thUser",
				Vote:   1,
			},
			{
				UserID: "132",
				Vote:   -1,
			},
		},
	}

	NewPost = &posts.Post{
		IDBson:     bson.NewObjectId(),
		AuthorBson: "newUs",
		Author: &author.Author{
			UserID:   "132",
			UserName: "newUs",
		},
		Category:         "news",
		CreationData:     "2020-01-01T04:20:00",
		PostID:           "newPost",
		Score:            3,
		Data:             "some new info",
		Title:            "new post",
		Type:             "text",
		UpvotePersentage: 100,
		Views:            250,
		Votes: []*posts.Vote{
			{
				UserID: "1stUser",
				Vote:   1,
			},
			{
				UserID: "2ndUser",
				Vote:   1,
			},
			{
				UserID: "3rdUser",
				Vote:   1,
			},
		},
	}

	SessionKey    = "sessionKey"
	GlobalSession = session.NewSession(NewPost.Author.UserID, NewPost.Author.UserName)
)

func TestPostHandlerList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	req := httptest.NewRequest("GET", "/api/posts/", nil)
	w := httptest.NewRecorder()

	st.EXPECT().GetAllPosts().Return(ResultPosts, nil)

	pstsHandler.List(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`1stPost`)) || !bytes.Contains(body, []byte(`2ndPost`)) {
		t.Errorf("unexpected error")
		return
	}

	// bad DB answer
	req = httptest.NewRequest("GET", "/api/posts/", nil)
	w = httptest.NewRecorder()
	st.EXPECT().GetAllPosts().Return(nil, fmt.Errorf(`pckg handlers, List err`))

	pstsHandler.List(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, List err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, List err")
		return
	}
}

func TestPostHandlerShowPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	req := httptest.NewRequest("GET", "/api/posts/1stPost", nil)
	w := httptest.NewRecorder()

	st.EXPECT().GetPost("/1stPost").Return(ResultPosts[0], nil)

	pstsHandler.ShowPost(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`1stPost`)) {
		t.Errorf("unexpected error")
		return
	}

	// bad item ID
	req = httptest.NewRequest("GET", "/api/post/1stPost", nil)
	w = httptest.NewRecorder()

	st.EXPECT().GetPost("1stPost").Return(nil, fmt.Errorf("bad item id"))

	pstsHandler.ShowPost(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, ShowPost no such post`)) {
		t.Errorf("unexpected error, expected: pckg handlers, ShowPost no such post")
		return
	}
}

func TestPostHandlerShowCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	req := httptest.NewRequest("GET", "/api/posts/programming", nil)
	w := httptest.NewRecorder()

	st.EXPECT().GetCategory("programming").Return(ResultPosts, nil)

	pstsHandler.ShowCategory(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`1stPost`)) || !bytes.Contains(body, []byte(`2ndPost`)) {
		t.Errorf("unexpected error")
		return
	}

	// err in get category
	req = httptest.NewRequest("GET", "/api/posts/programming", nil)
	w = httptest.NewRecorder()

	st.EXPECT().GetCategory("programming").Return(nil, fmt.Errorf("err in DB"))

	pstsHandler.ShowCategory(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, ShowCategory no posts in category`)) {
		t.Errorf("unexpected error, expected: pckg handlers, ShowCategory no posts in category")
		return
	}
}

func TestPostHandlerShowAllUserPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	req := httptest.NewRequest("GET", "/api/user/DeepThought", nil)
	w := httptest.NewRecorder()

	st.EXPECT().GetAllUserPosts("DeepThought").Return(ResultPosts, nil)

	pstsHandler.ShowAllUserPosts(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`1stPost`)) {
		t.Errorf("unexpected error")
		return
	}

	// err in get category
	req = httptest.NewRequest("GET", "/api/user/DeepThought", nil)
	w = httptest.NewRecorder()

	st.EXPECT().GetAllUserPosts("DeepThought").Return(nil, fmt.Errorf("err in DB"))

	pstsHandler.ShowAllUserPosts(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, ShowAllUserPosts, user has no posts`)) {
		t.Errorf("unexpected error, expected: pckg handlers, ShowAllUserPosts, user has no posts")
		return
	}
}

func TestPostHandlerAddNewPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	byteValue, err := json.Marshal(NewPost)
	if err != nil {
		t.Errorf("TestPostHandlerAddNewPost, marshal err %#v", err.Error())
		return
	}
	bodyReader := strings.NewReader(string(byteValue))

	req := httptest.NewRequest("POST", "/api/posts", bodyReader)
	ctx := session.ContextWithSession(req.Context(), GlobalSession)

	w := httptest.NewRecorder()

	st.EXPECT().CreatePost(req.WithContext(ctx), NewPost.Author.UserID, NewPost.Author.UserName).Return(NewPost, nil)

	pstsHandler.AddNewPost(w, req.WithContext(ctx))

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`newPost`)) {
		t.Errorf("unexpected error")
		return
	}

	// err in get category
	w = httptest.NewRecorder()

	st.EXPECT().CreatePost(req.WithContext(ctx), NewPost.Author.UserID, NewPost.Author.UserName).Return(nil, fmt.Errorf("err in DB"))

	pstsHandler.AddNewPost(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, AddNewPost CreatePost err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, AddNewPost CreatePost err")
		return
	}

	// session from ctx err
	w = httptest.NewRecorder()

	pstsHandler.AddNewPost(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, AddNewPost SessionFromContext err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, AddNewPost SessionFromContext err")
		return
	}
}

func TestPostHandlerDeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	req := httptest.NewRequest("DELETE", "/api/post/1stPost", nil)
	w := httptest.NewRecorder()

	st.EXPECT().DeletePostFromRepo(req).Return(nil)

	pstsHandler.DeletePost(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`success`)) {
		t.Errorf("unexpected error")
		return
	}

	// err in DeletePostFromRepo
	w = httptest.NewRecorder()

	st.EXPECT().DeletePostFromRepo(req).Return(fmt.Errorf("err in DB"))

	pstsHandler.DeletePost(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`err in DB`)) {
		t.Errorf("unexpected error, expected: err in DB")
		return
	}
}

func TestPostHandlerAddNewComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	newPostWithComm := NewPost
	comm := &comments.Comment{
		IDBson: bson.NewObjectId(),
		Author: &author.Author{
			UserID:   "132",
			UserName: "newUs",
		},
		CommentBody:  "newPostWithCommBody",
		CreationData: "2018-01-01T04:20:00",
		CommentID:    "commID",
	}
	newPostWithComm.Comments = append(newPostWithComm.Comments, comm)

	bodyReader := strings.NewReader(string(`{"comment": "newPostWithCommBody"}`))

	req := httptest.NewRequest("POST", "/api/post/newPost", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	ctx := session.ContextWithSession(req.Context(), GlobalSession)

	w := httptest.NewRecorder()

	st.EXPECT().
		AddCommentInPostRepository(NewPost.PostID, comm.CommentBody, NewPost.Author.UserID, NewPost.Author.UserName).
		Return(newPostWithComm, nil)

	pstsHandler.AddNewComment(w, req.WithContext(ctx))

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`newPostWithCommBody`)) {
		t.Errorf("unexpected error")
		return
	}

	// err: session
	req = httptest.NewRequest("POST", "/api/post/newPost", nil)
	w = httptest.NewRecorder()

	pstsHandler.AddNewComment(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, CreateComment SessionFromContext err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, CreateComment SessionFromContext err")
		return
	}

	// err: bad header
	req.Header.Set("test", "test")
	w = httptest.NewRecorder()

	pstsHandler.AddNewComment(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`unknown payload`)) {
		t.Errorf("unexpected error, expected: unknown payload")
		return
	}

	// err: bad json for message comm
	bodyReader = strings.NewReader(string(``))
	req = httptest.NewRequest("POST", "/api/post/newPost", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	pstsHandler.AddNewComment(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pkg handlers cant unpack payload in AddNewComment`)) {
		t.Errorf("unexpected error, expected: pkg handlers cant unpack payload in AddNewComment")
		return
	}

	// err: empty comm
	bodyReader = strings.NewReader(string(`{"comment": ""}`))
	req = httptest.NewRequest("POST", "/api/post/newPost", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	pstsHandler.AddNewComment(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`is required`)) {
		t.Errorf("unexpected error")
		return
	}

	// err: DB
	bodyReader = strings.NewReader(string(`{"comment": "newPostWithCommBody"}`))
	req = httptest.NewRequest("POST", "/api/post/newPost", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	st.EXPECT().
		AddCommentInPostRepository(NewPost.PostID, comm.CommentBody, NewPost.Author.UserID, NewPost.Author.UserName).
		Return(nil, fmt.Errorf("DB err"))

	pstsHandler.AddNewComment(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, CreateComment err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, CreateComment err")
		return
	}
}

func TestPostHandlerDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	newPostWithComm := NewPost
	comm := &comments.Comment{
		IDBson: bson.NewObjectId(),
		Author: &author.Author{
			UserID:   "132",
			UserName: "newUs",
		},
		CommentBody:  "newPostWithCommBody",
		CreationData: "2018-01-01T04:20:00",
		CommentID:    "commID",
	}
	newPostWithComm.Comments = append(newPostWithComm.Comments, comm)

	// good req
	req := httptest.NewRequest("DELETE", "/api/post/newPost/commID", nil)
	ctx := session.ContextWithSession(req.Context(), GlobalSession)
	w := httptest.NewRecorder()

	st.EXPECT().
		DeleteCommentInPostRepository(newPostWithComm.PostID, comm.CommentID, comm.Author.UserID, comm.Author.UserName).
		Return(NewPost, nil)

	newPostWithComm.Comments = nil
	pstsHandler.DeleteComment(w, req.WithContext(ctx))

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if bytes.Contains(body, []byte(`newPostWithCommBody`)) {
		t.Errorf("unexpected error")
		return
	}

	// err : session
	req = httptest.NewRequest("DELETE", "/api/post/newPost/commID", nil)
	w = httptest.NewRecorder()

	pstsHandler.DeleteComment(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, DeleteComment SessionFromContext err`)) {
		t.Errorf("unexpected error: expected: pckg handlers, DeleteComment SessionFromContext err")
		return
	}

	// err : getPostCommentID
	req = httptest.NewRequest("DELETE", "/api/post/newPostcommID", nil)
	ctx = session.ContextWithSession(req.Context(), GlobalSession)
	w = httptest.NewRecorder()

	pstsHandler.DeleteComment(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, DeleteComment getPostCommentID err`)) {
		t.Errorf("unexpected error: expected: pckg handlers, DeleteComment getPostCommentID err")
		return
	}

	// err : DB
	req = httptest.NewRequest("DELETE", "/api/post/newPost/commID", nil)
	ctx = session.ContextWithSession(req.Context(), GlobalSession)
	w = httptest.NewRecorder()

	st.EXPECT().
		DeleteCommentInPostRepository(newPostWithComm.PostID, comm.CommentID, comm.Author.UserID, comm.Author.UserName).
		Return(nil, fmt.Errorf("DB err"))

	pstsHandler.DeleteComment(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, DeleteComment DeleteCommentInPostRepository err`)) {
		t.Errorf("unexpected error: expected: pckg handlers, DeleteComment DeleteCommentInPostRepository err")
		return
	}
}

func TestPostHandlerUpVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	req := httptest.NewRequest("GET", "/api/post/1stPost/upvote", nil)
	ctx := session.ContextWithSession(req.Context(), GlobalSession)
	w := httptest.NewRecorder()

	st.EXPECT().
		UpVotePost(ResultPosts[0].PostID, NewPost.Author.UserID).
		Return(UpVotedPost, nil)

	pstsHandler.UpVote(w, req.WithContext(ctx))

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`132`)) || !bytes.Contains(body, []byte(`"vote":1`)) {
		t.Errorf("unexpected error")
		return
	}

	// err : session
	w = httptest.NewRecorder()

	pstsHandler.UpVote(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, UpVote SessionFromContext err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, UpVote SessionFromContext err")
		return
	}

	// err : getPostCommentID
	req = httptest.NewRequest("GET", "/api/post/1stPostupvote", nil)
	ctx = session.ContextWithSession(req.Context(), GlobalSession)
	w = httptest.NewRecorder()

	pstsHandler.UpVote(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`incorrect link`)) {
		t.Errorf("unexpected error, expected: incorrect link")
		return
	}

	// err : DB
	req = httptest.NewRequest("GET", "/api/post/1stPost/upvote", nil)
	ctx = session.ContextWithSession(req.Context(), GlobalSession)
	w = httptest.NewRecorder()

	st.EXPECT().
		UpVotePost(ResultPosts[0].PostID, NewPost.Author.UserID).
		Return(nil, fmt.Errorf("DB err"))

	pstsHandler.UpVote(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, UpVote UpVotePost err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, UpVote UpVotePost err")
		return
	}
}

func TestPostHandlerDownVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	req := httptest.NewRequest("GET", "/api/post/1stPost/downvote", nil)
	ctx := session.ContextWithSession(req.Context(), GlobalSession)
	w := httptest.NewRecorder()

	st.EXPECT().
		DownVotePost(UpVotedPost.PostID, NewPost.Author.UserID).
		Return(DownVotedPost, nil)

	pstsHandler.DownVote(w, req.WithContext(ctx))

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`132`)) || !bytes.Contains(body, []byte(`"vote":-1`)) {
		t.Errorf("unexpected error")
		return
	}

	// err : session
	w = httptest.NewRecorder()

	pstsHandler.DownVote(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, DownVote SessionFromContext err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, DownVote SessionFromContext err")
		return
	}

	// err : getPostCommentID
	req = httptest.NewRequest("GET", "/api/post/1stPostdownvote", nil)
	ctx = session.ContextWithSession(req.Context(), GlobalSession)
	w = httptest.NewRecorder()

	pstsHandler.DownVote(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`incorrect link`)) {
		t.Errorf("unexpected error, expected: incorrect link")
		return
	}

	// err : DB
	req = httptest.NewRequest("GET", "/api/post/1stPost/downvote", nil)
	ctx = session.ContextWithSession(req.Context(), GlobalSession)
	w = httptest.NewRecorder()

	st.EXPECT().
		DownVotePost(UpVotedPost.PostID, NewPost.Author.UserID).
		Return(nil, fmt.Errorf("DB err"))

	pstsHandler.DownVote(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, DownVote DownVotePost err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, DownVote DownVotePost err")
		return
	}
}

func TestPostHandlerUnVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := posts.NewMockPostsRepo(ctrl)
	pstsHandler := &PostsHandler{
		PostsRepo: st,
	}

	// good req
	req := httptest.NewRequest("GET", "/api/post/1stPost/unvote", nil)
	ctx := session.ContextWithSession(req.Context(), GlobalSession)
	w := httptest.NewRecorder()

	st.EXPECT().
		UnVotePost(UpVotedPost.PostID, NewPost.Author.UserID).
		Return(ResultPosts[0], nil)

	pstsHandler.UnVote(w, req.WithContext(ctx))

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if bytes.Contains(body, []byte(`132`)) {
		t.Errorf("unexpected error")
		return
	}

	// err : session
	w = httptest.NewRecorder()

	pstsHandler.UnVote(w, req)

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, UnVote SessionFromContext err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, UnVote SessionFromContext err")
		return
	}

	// err : getPostCommentID
	req = httptest.NewRequest("GET", "/api/post/1stPostunvote", nil)
	ctx = session.ContextWithSession(req.Context(), GlobalSession)
	w = httptest.NewRecorder()

	pstsHandler.UnVote(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`incorrect link`)) {
		t.Errorf("unexpected error, expected: incorrect link")
		return
	}

	// err : DB
	req = httptest.NewRequest("GET", "/api/post/1stPost/unvote", nil)
	ctx = session.ContextWithSession(req.Context(), GlobalSession)
	w = httptest.NewRecorder()

	st.EXPECT().
		UnVotePost(UpVotedPost.PostID, NewPost.Author.UserID).
		Return(nil, fmt.Errorf("DB err"))

	pstsHandler.UnVote(w, req.WithContext(ctx))

	resp = w.Result()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err")
		return
	}

	if !bytes.Contains(body, []byte(`pckg handlers, UnVote UnVotePost err`)) {
		t.Errorf("unexpected error, expected: pckg handlers, UnVote UnVotePost err")
		return
	}
}
