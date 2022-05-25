package posts

import (
	"encoding/json"
	"fmt"
	"myRedditClone/pkg/author"
	"myRedditClone/pkg/comments"
	"net/http/httptest"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"gopkg.in/mgo.v2/bson"
)

var (
	FoundedPosts = []*Post{
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
			Votes: []*Vote{
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
			Votes: []*Vote{
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

	newPostTextForm = &NewPostForm{
		Category: "programming",
		Text:     "Test text",
		Title:    "some title",
		Type:     "text",
	}
	newPostURLForm = &NewPostForm{
		Category: "programming",
		URL:      "/local",
		Title:    "some title",
		Type:     "link",
	}
	newPostText = &Post{
		// IDBson:     bson.NewObjectId(),
		AuthorBson: "testLogin",
		// Author: &author.Author{
		// 	UserID:   "testID",
		// 	UserName: "testLogin",
		// },
		Category:         "programming",
		Score:            1,
		Data:             "Test text",
		Title:            "some title",
		Type:             "text",
		UpvotePersentage: 100,
		Views:            0,
		// Votes: []*Vote{
		// 	{
		// 		UserID: "testID",
		// 		Vote:   1,
		// 	},
		// },
	}
	newPostURL = &Post{
		AuthorBson:       "testLogin",
		Category:         "programming",
		Score:            1,
		Data:             "/local",
		Title:            "some title",
		Type:             "link",
		UpvotePersentage: 100,
		Views:            0,
	}

	postWithComm = &Post{
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
		Votes: []*Vote{
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
	}
	newComm = &comments.Comment{
		IDBson: bson.NewObjectId(),
		Author: &author.Author{
			UserID:   "newCommAuthorID",
			UserName: "newCommAuthorLogin",
		},
		CommentBody: "newCommBody",
		CommentID:   "newCommID",
	}

	upVotedPost = &Post{
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
		Votes: []*Vote{
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
				UserID: "upVoteID",
				Vote:   1,
			},
		},
	}
)

// нужно мокать DB

func TestGetAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
	}

	// Founded users
	st.EXPECT().FindAll(handlerPosts).Return(FoundedPosts, nil)
	_, err := handlerPosts.GetAllPosts()

	if err != nil {
		t.Errorf("unexpected error")
		return
	}

	// err in founded users
	st.EXPECT().FindAll(handlerPosts).Return(nil, fmt.Errorf("DB err: FindAll"))
	_, err = handlerPosts.GetAllPosts()
	if err.Error() != "DB err: FindAll" {
		t.Errorf("unexpected error, expected: DB err: FindAll")
		return
	}
}

func TestGetPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
	}

	// founded post
	st.EXPECT().FindPostByID(handlerPosts, "1stPost").Return(FoundedPosts[0], nil)
	post, err := handlerPosts.GetPost("1stPost")
	if err != nil || post.PostID != "1stPost" {
		t.Errorf("unexpected error")
		return
	}

	// founded post err
	st.EXPECT().FindPostByID(handlerPosts, "1stPost").Return(nil, fmt.Errorf("DB err: FindPostByID"))
	_, err = handlerPosts.GetPost("1stPost")
	if err.Error() != "DB err: FindPostByID" {
		t.Errorf("unexpected error, expected: DB err: FindPostByID")
		return
	}
}

func TestGetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
	}

	// founded category
	st.EXPECT().FindCategory(handlerPosts, "programming").Return(FoundedPosts, nil)
	_, err := handlerPosts.GetCategory("programming")
	if err != nil {
		t.Errorf("unexpected error")
		return
	}

	// founded category err
	st.EXPECT().FindCategory(handlerPosts, "programming").Return(nil, fmt.Errorf("DB err: FindCategory"))
	_, err = handlerPosts.GetCategory("programming")
	if err.Error() != "DB err: FindCategory" {
		t.Errorf("unexpected error, expected: DB err: FindCategory")
		return
	}

	// founded category: 0 posts
	st.EXPECT().FindCategory(handlerPosts, "programming").Return([]*Post{}, nil)
	posts, err := handlerPosts.GetCategory("programming")
	if err != nil || len(posts) != 0 {
		t.Errorf("unexpected error")
		return
	}
}

func TestGetAllUserPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
	}

	// founded all user posts
	st.EXPECT().FindAllUserPosts(handlerPosts, "DeepThought").Return(FoundedPosts, nil)
	posts, err := handlerPosts.GetAllUserPosts("DeepThought")
	if err != nil || len(posts) != 2 {
		t.Errorf("unexpected error")
		return
	}

	// founded all user posts err
	st.EXPECT().FindAllUserPosts(handlerPosts, "DeepThought").Return(nil, fmt.Errorf("DB err: FindAllUserPosts"))
	_, err = handlerPosts.GetAllUserPosts("DeepThought")
	if err.Error() != "DB err: FindAllUserPosts" {
		t.Errorf("unexpected error, expected: DB err: FindAllUserPosts")
		return
	}

	// founded no posts
	st.EXPECT().FindAllUserPosts(handlerPosts, "DeepThought").Return([]*Post{}, nil)
	posts, err = handlerPosts.GetAllUserPosts("DeepThought")
	if err != nil || len(posts) != 0 {
		t.Errorf("unexpected error")
		return
	}
}

func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
	}

	// created post text
	byteValue, err := json.Marshal(newPostTextForm)
	if err != nil {
		t.Errorf("TestCreatePost, json.Marsha, err : %#v", err.Error())
		return
	}
	bodyReader := strings.NewReader(string(byteValue))

	req := httptest.NewRequest("POST", "/api/posts", bodyReader)
	req.Header.Set("Content-Type", "application/json")

	newPostText.CreationData = getFormatData()

	st.EXPECT().InsertPost(handlerPosts, newPostText, "testID", "testLogin").Return(nil)
	_, err = handlerPosts.CreatePost(req, "testID", "testLogin")

	if err != nil {
		t.Errorf("unexpected error")
		return
	}

	// created post text err
	bodyReader = strings.NewReader(string(byteValue))
	req = httptest.NewRequest("POST", "/api/posts", bodyReader)
	req.Header.Set("Content-Type", "application/json")

	st.EXPECT().InsertPost(handlerPosts, newPostText, "testID", "testLogin").Return(fmt.Errorf("DB err: InsertPost"))
	_, err = handlerPosts.CreatePost(req, "testID", "testLogin")

	if err.Error() != "DB err: InsertPost" {
		t.Errorf("unexpected error, expected: DB err: InsertPost, got %#v", err.Error())
		return
	}

	// created post url
	byteValue, err = json.Marshal(newPostURLForm)
	if err != nil {
		t.Errorf("TestCreatePost, json.Marsha, err : %#v", err.Error())
		return
	}
	bodyReader = strings.NewReader(string(byteValue))

	req = httptest.NewRequest("POST", "/api/posts", bodyReader)
	req.Header.Set("Content-Type", "application/json")

	newPostURL.CreationData = getFormatData()

	st.EXPECT().InsertPost(handlerPosts, newPostURL, "testID", "testLogin").Return(nil)
	_, err = handlerPosts.CreatePost(req, "testID", "testLogin")

	if err != nil {
		t.Errorf("unexpected error")
		return
	}

	// created post url err
	bodyReader = strings.NewReader(string(byteValue))
	req = httptest.NewRequest("POST", "/api/posts", bodyReader)
	req.Header.Set("Content-Type", "application/json")

	st.EXPECT().InsertPost(handlerPosts, newPostURL, "testID", "testLogin").Return(fmt.Errorf("DB err: InsertPost"))
	_, err = handlerPosts.CreatePost(req, "testID", "testLogin")

	if err.Error() != "DB err: InsertPost" {
		t.Errorf("unexpected error, expected: DB err: InsertPost, got %#v", err.Error())
		return
	}

	// bad header
	req = httptest.NewRequest("POST", "/api/posts", nil)

	_, err = handlerPosts.CreatePost(req, "testID", "testLogin")

	if err.Error() != "unknown payload" {
		t.Errorf("unexpected error, expected: unknown payload, got %#v", err.Error())
		return
	}

	// bad JSON
	bodyReader = strings.NewReader("")
	req = httptest.NewRequest("POST", "/api/posts", bodyReader)
	req.Header.Set("Content-Type", "application/json")

	_, err = handlerPosts.CreatePost(req, "testID", "testLogin")

	if err.Error() != "cant unpack payload" {
		t.Errorf("unexpected error, expected: cant unpack payload, got %#v", err.Error())
		return
	}
}

func TestDeletePostFromRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
	}

	// good req
	req := httptest.NewRequest("GET", "/api/post/1stPost", nil)

	st.EXPECT().RemovePost(handlerPosts, "1stPost").Return(nil)
	err := handlerPosts.DeletePostFromRepo(req)

	if err != nil {
		t.Errorf("unexpected error, got %#v", err.Error())
		return
	}

	// good req err
	st.EXPECT().RemovePost(handlerPosts, "1stPost").Return(fmt.Errorf("DB err"))
	err = handlerPosts.DeletePostFromRepo(req)

	if err.Error() != "DB err" {
		t.Errorf("unexpected error, expected: DB err, got %#v", err.Error())
		return
	}
}

func TestAddCommentInPostRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	cmRepo := comments.NewMockCommentsRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
		commRepo: cmRepo,
	}

	// good req
	comms := []*comments.Comment{}
	comms = append(comms, newComm)
	postWithComm.Comments = comms

	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(postWithComm, nil)
	cmRepo.EXPECT().
		CreateComment(handlerPosts.sessionDB, FoundedPosts[0].PostID, FoundedPosts[0].Data, newComm.Author.UserID, newComm.Author.UserName).
		Return(nil, nil)
	cmRepo.EXPECT().GetAllComments(handlerPosts.sessionDB, FoundedPosts[0].PostID).Return(comms, nil)
	st.EXPECT().UpdatePost(handlerPosts, postWithComm.PostID, postWithComm).Return(nil)

	_, err := handlerPosts.
		AddCommentInPostRepository(FoundedPosts[0].PostID, FoundedPosts[0].Data, newComm.Author.UserID, newComm.Author.UserName)
	if err != nil {
		t.Errorf("unexpected error, got %#v", err.Error())
		return
	}

	// FindPostByID err
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(nil, fmt.Errorf("FindPostByID err"))

	_, err = handlerPosts.
		AddCommentInPostRepository(FoundedPosts[0].PostID, FoundedPosts[0].Data, newComm.Author.UserID, newComm.Author.UserName)
	if err.Error() != "FindPostByID err" {
		t.Errorf("unexpected error, expected: FindPostByID err, got %#v", err.Error())
		return
	}

	// CreateComment err
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(postWithComm, nil)
	cmRepo.EXPECT().
		CreateComment(handlerPosts.sessionDB, FoundedPosts[0].PostID, FoundedPosts[0].Data, newComm.Author.UserID, newComm.Author.UserName).
		Return(nil, fmt.Errorf("CreateComment err"))

	_, err = handlerPosts.
		AddCommentInPostRepository(FoundedPosts[0].PostID, FoundedPosts[0].Data, newComm.Author.UserID, newComm.Author.UserName)
	if err.Error() != "CreateComment err" {
		t.Errorf("unexpected error, expected: CreateComment err, got %#v", err.Error())
		return
	}

	// GetAllComments err
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(postWithComm, nil)
	cmRepo.EXPECT().
		CreateComment(handlerPosts.sessionDB, FoundedPosts[0].PostID, FoundedPosts[0].Data, newComm.Author.UserID, newComm.Author.UserName).
		Return(nil, nil)
	cmRepo.EXPECT().GetAllComments(handlerPosts.sessionDB, FoundedPosts[0].PostID).
		Return(nil, fmt.Errorf("GetAllComments err"))

	_, err = handlerPosts.
		AddCommentInPostRepository(FoundedPosts[0].PostID, FoundedPosts[0].Data, newComm.Author.UserID, newComm.Author.UserName)
	if err.Error() != "GetAllComments err" {
		t.Errorf("unexpected error, expected: CreateComment err, got %#v", err.Error())
		return
	}

	// FindPostByID err
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(postWithComm, nil)
	cmRepo.EXPECT().
		CreateComment(handlerPosts.sessionDB, FoundedPosts[0].PostID, FoundedPosts[0].Data, newComm.Author.UserID, newComm.Author.UserName).
		Return(nil, nil)
	cmRepo.EXPECT().GetAllComments(handlerPosts.sessionDB, FoundedPosts[0].PostID).Return(comms, nil)
	st.EXPECT().UpdatePost(handlerPosts, postWithComm.PostID, postWithComm).Return(fmt.Errorf("FindPostByID err"))

	_, err = handlerPosts.
		AddCommentInPostRepository(FoundedPosts[0].PostID, FoundedPosts[0].Data, newComm.Author.UserID, newComm.Author.UserName)
	if err.Error() != "FindPostByID err" {
		t.Errorf("uunexpected error, expected: FindPostByID err, got %#v", err.Error())
		return
	}
}

func TestDeleteCommentInPostRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	cmRepo := comments.NewMockCommentsRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
		commRepo: cmRepo,
	}

	// good req
	comms := []*comments.Comment{}
	comms = append(comms, newComm)
	postWithComm.Comments = comms

	st.EXPECT().FindPostByID(handlerPosts, postWithComm.PostID).Return(postWithComm, nil)
	cmRepo.EXPECT().
		DeleteCommentFromRepo(handlerPosts.sessionDB, postWithComm.PostID, postWithComm.Comments[0].CommentID, newComm.Author.UserID, newComm.Author.UserName).
		Return(nil)
	cmRepo.EXPECT().GetAllComments(handlerPosts.sessionDB, postWithComm.PostID).Return([]*comments.Comment{}, nil)
	st.EXPECT().UpdatePost(handlerPosts, postWithComm.PostID, postWithComm).Return(nil)

	_, err := handlerPosts.
		DeleteCommentInPostRepository(postWithComm.PostID, postWithComm.Comments[0].CommentID, newComm.Author.UserID, newComm.Author.UserName)
	if err != nil {
		t.Errorf("unexpected error, got %#v", err.Error())
		return
	}

	// FindPostByID err
	postWithComm.Comments = comms
	st.EXPECT().FindPostByID(handlerPosts, postWithComm.PostID).Return(nil, fmt.Errorf("FindPostByID err"))

	_, err = handlerPosts.
		DeleteCommentInPostRepository(postWithComm.PostID, postWithComm.Comments[0].CommentID, newComm.Author.UserID, newComm.Author.UserName)
	if err.Error() != "FindPostByID err" {
		t.Errorf("unexpected error, expected: FindPostByID err, got %#v", err.Error())
		return
	}

	// DeleteCommentFromRepo err
	st.EXPECT().FindPostByID(handlerPosts, postWithComm.PostID).Return(postWithComm, nil)
	cmRepo.EXPECT().
		DeleteCommentFromRepo(handlerPosts.sessionDB, postWithComm.PostID, postWithComm.Comments[0].CommentID, newComm.Author.UserID, newComm.Author.UserName).
		Return(fmt.Errorf("DeleteCommentFromRepo err"))

	_, err = handlerPosts.
		DeleteCommentInPostRepository(postWithComm.PostID, postWithComm.Comments[0].CommentID, newComm.Author.UserID, newComm.Author.UserName)
	if err.Error() != "DeleteCommentFromRepo err" {
		t.Errorf("unexpected error, expected: DeleteCommentFromRepo err, got %#v", err.Error())
		return
	}

	// GetAllComments err
	st.EXPECT().FindPostByID(handlerPosts, postWithComm.PostID).Return(postWithComm, nil)
	cmRepo.EXPECT().
		DeleteCommentFromRepo(handlerPosts.sessionDB, postWithComm.PostID, postWithComm.Comments[0].CommentID, newComm.Author.UserID, newComm.Author.UserName).
		Return(nil)
	cmRepo.EXPECT().GetAllComments(handlerPosts.sessionDB, postWithComm.PostID).Return(nil, fmt.Errorf("GetAllComments err"))

	_, err = handlerPosts.
		DeleteCommentInPostRepository(postWithComm.PostID, postWithComm.Comments[0].CommentID, newComm.Author.UserID, newComm.Author.UserName)
	if err.Error() != "GetAllComments err" {
		t.Errorf("unexpected error, expected: GetAllComments err, got %#v", err.Error())
		return
	}

	// UpdatePost err
	postWithComm.Comments = comms
	st.EXPECT().FindPostByID(handlerPosts, postWithComm.PostID).Return(postWithComm, nil)
	cmRepo.EXPECT().
		DeleteCommentFromRepo(handlerPosts.sessionDB, postWithComm.PostID, postWithComm.Comments[0].CommentID, newComm.Author.UserID, newComm.Author.UserName).
		Return(nil)
	cmRepo.EXPECT().GetAllComments(handlerPosts.sessionDB, postWithComm.PostID).Return([]*comments.Comment{}, nil)
	st.EXPECT().UpdatePost(handlerPosts, postWithComm.PostID, postWithComm).Return(fmt.Errorf("UpdatePost err"))

	_, err = handlerPosts.
		DeleteCommentInPostRepository(postWithComm.PostID, postWithComm.Comments[0].CommentID, newComm.Author.UserID, newComm.Author.UserName)
	if err.Error() != "UpdatePost err" {
		t.Errorf("unexpected error, expected: UpdatePost err, got %#v", err.Error())
		return
	}
}

func TestUpVotePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
	}

	// good req
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(FoundedPosts[0], nil)
	newVote := &Vote{
		UserID: "upVoteID",
		Vote:   1,
	}
	FoundedPosts[0].Votes = append(FoundedPosts[0].Votes, newVote)
	FoundedPosts[0].Score = 5
	st.EXPECT().UpdatePost(handlerPosts, FoundedPosts[0].PostID, FoundedPosts[0]).Return(nil)

	_, err := handlerPosts.UpVotePost(FoundedPosts[0].PostID, "upVoteID")
	if err != nil {
		t.Errorf("unexpected error, got %#v", err.Error())
		return
	}

	// FindPostByID err
	FoundedPosts[0].Votes = FoundedPosts[0].Votes[:len(FoundedPosts[0].Votes)-1]
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(nil, fmt.Errorf("FindPostByID err"))

	_, err = handlerPosts.UpVotePost(FoundedPosts[0].PostID, "upVoteID")
	if err.Error() != "FindPostByID err" {
		t.Errorf("unexpected error, expected: FindPostByID err, got %#v", err.Error())
		return
	}

	// UpdatePost err
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(FoundedPosts[0], nil)
	newVote = &Vote{
		UserID: "upVoteID",
		Vote:   1,
	}
	FoundedPosts[0].Votes = append(FoundedPosts[0].Votes, newVote)
	FoundedPosts[0].Score = 5
	st.EXPECT().UpdatePost(handlerPosts, FoundedPosts[0].PostID, FoundedPosts[0]).Return(fmt.Errorf("UpdatePost err"))

	FoundedPosts[0].Votes = FoundedPosts[0].Votes[:len(FoundedPosts[0].Votes)-1]
	_, err = handlerPosts.UpVotePost(FoundedPosts[0].PostID, "upVoteID")
	if err.Error() != "UpdatePost err" {
		t.Errorf("unexpected error, expected: UpdatePost err, got %#v", err.Error())
		return
	}
}

func TestDownVotePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
	}

	// good req
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(FoundedPosts[0], nil)
	newVote := &Vote{
		UserID: "upVoteID",
		Vote:   -1,
	}
	FoundedPosts[0].Votes = append(FoundedPosts[0].Votes, newVote)
	FoundedPosts[0].Score = 5
	st.EXPECT().UpdatePost(handlerPosts, FoundedPosts[0].PostID, FoundedPosts[0]).Return(nil)

	_, err := handlerPosts.DownVotePost(FoundedPosts[0].PostID, "upVoteID")
	if err != nil {
		t.Errorf("unexpected error, got %#v", err.Error())
		return
	}

	// FindPostByID err
	FoundedPosts[0].Votes = FoundedPosts[0].Votes[:len(FoundedPosts[0].Votes)-1]
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(nil, fmt.Errorf("FindPostByID err"))

	_, err = handlerPosts.DownVotePost(FoundedPosts[0].PostID, "upVoteID")
	if err.Error() != "FindPostByID err" {
		t.Errorf("unexpected error, expected: FindPostByID err, got %#v", err.Error())
		return
	}

	// UpdatePost err
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(FoundedPosts[0], nil)
	newVote = &Vote{
		UserID: "upVoteID",
		Vote:   -1,
	}
	FoundedPosts[0].Votes = append(FoundedPosts[0].Votes, newVote)
	FoundedPosts[0].Score = 5
	st.EXPECT().UpdatePost(handlerPosts, FoundedPosts[0].PostID, FoundedPosts[0]).Return(fmt.Errorf("UpdatePost err"))

	FoundedPosts[0].Votes = FoundedPosts[0].Votes[:len(FoundedPosts[0].Votes)-1]
	_, err = handlerPosts.DownVotePost(FoundedPosts[0].PostID, "upVoteID")
	if err.Error() != "UpdatePost err" {
		t.Errorf("unexpected error, expected: UpdatePost err, got %#v", err.Error())
		return
	}
}

func TestUnVotePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockDBHelperRepo(ctrl)
	handlerPosts := &ItemMemoryRepository{
		dbHelper: st,
	}

	// good req
	newVote := &Vote{
		UserID: "upVoteID",
		Vote:   -1,
	}
	FoundedPosts[0].Votes = append(FoundedPosts[0].Votes, newVote)
	FoundedPosts[0].Score = 5

	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(FoundedPosts[0], nil)
	FoundedPosts[0].Votes = FoundedPosts[0].Votes[:len(FoundedPosts[0].Votes)-1]
	st.EXPECT().UpdatePost(handlerPosts, FoundedPosts[0].PostID, FoundedPosts[0]).Return(nil)

	_, err := handlerPosts.UnVotePost(FoundedPosts[0].PostID, "upVoteID")
	if err != nil {
		t.Errorf("unexpected error, got %#v", err.Error())
		return
	}

	// FindPostByID err
	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(nil, fmt.Errorf("FindPostByID err"))

	_, err = handlerPosts.UnVotePost(FoundedPosts[0].PostID, "upVoteID")
	if err.Error() != "FindPostByID err" {
		t.Errorf("unexpected error, expected: FindPostByID err, got %#v", err.Error())
		return
	}

	// UpdatePost err
	newVote = &Vote{
		UserID: "upVoteID",
		Vote:   -1,
	}
	FoundedPosts[0].Votes = append(FoundedPosts[0].Votes, newVote)
	FoundedPosts[0].Score = 5

	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(FoundedPosts[0], nil)
	st.EXPECT().UpdatePost(handlerPosts, FoundedPosts[0].PostID, FoundedPosts[0]).Return(fmt.Errorf("UpdatePost err"))

	_, err = handlerPosts.UnVotePost(FoundedPosts[0].PostID, "upVoteID")
	if err.Error() != "UpdatePost err" {
		t.Errorf("unexpected error, expected: UpdatePost err, got %#v", err.Error())
		return
	}

	// no such vote found err
	newVote = &Vote{
		UserID: "upVoteID",
		Vote:   -1,
	}
	FoundedPosts[0].Votes = append(FoundedPosts[0].Votes, newVote)
	FoundedPosts[0].Score = 5

	st.EXPECT().FindPostByID(handlerPosts, FoundedPosts[0].PostID).Return(FoundedPosts[0], nil)

	_, err = handlerPosts.UnVotePost(FoundedPosts[0].PostID, "wrongID")
	if err.Error() != "no such vote found" {
		t.Errorf("unexpected error, expected: no such vote found, got %#v", err.Error())
		return
	}
}
