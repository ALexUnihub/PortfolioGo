package posts

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"myRedditClone/pkg/author"
	"myRedditClone/pkg/comments"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ItemMemoryRepository struct {
	data      *mgo.Collection
	sessionDB *mgo.Session

	dbHelper DBHelperRepo
	commRepo comments.CommentsRepo
}

type NewPostForm struct {
	Category string `json:"category"`
	Text     string `json:"text"`
	URL      string `json:"url"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

func NewMemoryRepo(sess *mgo.Session) *ItemMemoryRepository {
	collection := sess.DB("posts").C("postRepo")

	AddStaticPosts(collection, sess)

	return &ItemMemoryRepository{
		data:      collection,
		sessionDB: sess,

		dbHelper: NewDBHelperRepo(),
		commRepo: comments.NewCommentRepoStruct(),
	}
}

func (repo *ItemMemoryRepository) GetAllPosts() ([]*Post, error) {
	// posts := []*Post{}
	// err := repo.data.Find(bson.M{}).All(&posts)

	posts, err := repo.dbHelper.FindAll(repo)
	if err != nil {
		log.Println("pckg posts, GetAllPosts : ", err.Error())
		return nil, err
	}

	return posts, nil
}

func (repo *ItemMemoryRepository) GetPost(id string) (*Post, error) {
	// post := &Post{}
	// err := repo.data.Find(bson.M{"id": id}).One(&post)

	post, err := repo.dbHelper.FindPostByID(repo, id)
	if err != nil {
		log.Println("pckg posts, GetPost, repo.data.Find ", err.Error())
		return nil, err
	}

	return post, nil
}

func (repo *ItemMemoryRepository) GetCategory(category string) ([]*Post, error) {
	// categoryPosts := []*Post{}
	// err := repo.data.Find(bson.M{"category": category}).All(&categoryPosts)

	categoryPosts, err := repo.dbHelper.FindCategory(repo, category)
	if err != nil {
		log.Println("pckg posts, GetCategory, repo.data.Find ", err.Error())
		return nil, err
	}

	if len(categoryPosts) == 0 {
		noPost := make([]*Post, 0)
		return noPost, nil
	}

	return categoryPosts, nil
}

func (repo *ItemMemoryRepository) GetAllUserPosts(userName string) ([]*Post, error) {
	// userPosts := []*Post{}
	// err := repo.data.Find(bson.M{"authorBson": userName}).All(&userPosts)

	userPosts, err := repo.dbHelper.FindAllUserPosts(repo, userName)
	if err != nil {
		log.Println("pckg posts, GetAllUserPosts, repo.data.Find ", err.Error())
		return nil, err
	}

	if len(userPosts) == 0 {
		noPost := make([]*Post, 0)
		return noPost, nil
	}

	return userPosts, nil
}

func (repo *ItemMemoryRepository) CreatePost(r *http.Request, userID, userLogin string) (*Post, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, errors.New(`unknown payload`)
	}

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	fd := &NewPostForm{}
	err := json.Unmarshal(body, fd)
	if err != nil {
		return nil, errors.New(`cant unpack payload`)
	}

	// rand.Seed(time.Now().UnixNano())
	// randID := string(randomBytes(10))

	newPost := &Post{
		// IDBson:     bson.NewObjectId(),
		AuthorBson: userLogin,
		// Author:     &author.Author{
		// UserID:   userID,
		// UserName: userLogin,
		// },
		Category:     fd.Category,
		CreationData: getFormatData(),
		// PostID:           randID,
		Score:            1,
		Title:            fd.Title,
		Type:             fd.Type,
		UpvotePersentage: 100,
		// Votes:            []*Vote{
		// {
		// 	UserID: userID,
		// 	Vote:   1,
		// },
		// },
	}

	if fd.Text != "" {
		newPost.Data = fd.Text
		// err = repo.data.Insert(&newPost)
		err = repo.dbHelper.InsertPost(repo, newPost, userID, userLogin)
		if err != nil {
			log.Println("pckg posts, CreatePost, text, repo.data.Insert", err.Error())
			return nil, err
		}

		return newPost, nil
	}

	newPost.Data = fd.URL
	// err = repo.data.Insert(&newPost)
	err = repo.dbHelper.InsertPost(repo, newPost, userID, userLogin)
	if err != nil {
		log.Println("pckg posts, CreatePost, url,  repo.data.Insert", err.Error())
		return nil, err
	}

	return newPost, nil
}

func (repo *ItemMemoryRepository) DeletePostFromRepo(r *http.Request) error {
	pathURL := r.URL.Path[len("/api/post/"):]

	// err := repo.data.Remove(bson.M{"id": pathURL})
	err := repo.dbHelper.RemovePost(repo, pathURL)
	if err != nil {
		log.Println("pckg posts, DeletePostFromRepo, repo.data.Remove", err.Error())
		return err
	}

	return nil
}

func (repo *ItemMemoryRepository) AddCommentInPostRepository(postID, data, userID, userLogin string) (*Post, error) {
	// post := &Post{}
	// err := repo.data.Find(bson.M{"id": postID}).One(&post)

	post, err := repo.dbHelper.FindPostByID(repo, postID)
	if err != nil {
		log.Println("pckg posts, AddCommentInPostRepository, repo.data.Find", err.Error())
		return nil, err
	}

	_, err = repo.commRepo.CreateComment(repo.sessionDB, postID, data, userID, userLogin)
	if err != nil {
		log.Println("pckg posts, AddCommentInPostRepository, CreateComment", err.Error())
		return nil, err
	}

	post.Comments, err = repo.commRepo.GetAllComments(repo.sessionDB, post.PostID)
	if err != nil {
		log.Println("pckg posts, AddCommentInPostRepository, GetAllComments", err.Error())
		return nil, err
	}

	// err = repo.data.Update(bson.M{"id": postID}, &post)
	err = repo.dbHelper.UpdatePost(repo, postID, post)
	if err != nil {
		log.Println("pckg posts, AddCommentInPostRepository, repo.data.Update", err.Error())
		return nil, err
	}

	return post, nil
}

func (repo *ItemMemoryRepository) DeleteCommentInPostRepository(postID, commID, userID, userLogin string) (*Post, error) {
	// post := &Post{}
	// err := repo.data.Find(bson.M{"id": postID}).One(&post)

	post, err := repo.dbHelper.FindPostByID(repo, postID)
	if err != nil {
		log.Println("pckg posts, DeleteCommentInPostRepository, repo.data.Find", err.Error())
		return nil, err
	}

	err = repo.commRepo.DeleteCommentFromRepo(repo.sessionDB, postID, commID, userID, userLogin)
	if err != nil {
		log.Println("pckg posts, DeleteCommentInPostRepository, comments.DeleteCommentFromRepo", err.Error())
		return nil, err
	}

	post.Comments, err = repo.commRepo.GetAllComments(repo.sessionDB, postID)
	if err != nil {
		log.Println("pckg posts, DeleteCommentInPostRepository, comments.GetAllComments", err.Error())
		return nil, err
	}

	err = repo.dbHelper.UpdatePost(repo, postID, post)
	if err != nil {
		log.Println("pckg posts, DeleteCommentInPostRepository, UpdatePost", err.Error())
		return nil, err
	}

	return post, err
}

func (repo *ItemMemoryRepository) UpVotePost(postID, userID string) (*Post, error) {
	// post := &Post{}
	// err := repo.data.Find(bson.M{"id": postID}).One(&post)

	post, err := repo.dbHelper.FindPostByID(repo, postID)
	if err != nil {
		log.Println("pckg posts, UpVotePost, repo.data.Find", err.Error())
		return nil, err
	}

	_ = changeVote(post, userID, 1)
	post.Score = setPostScore(post)
	post.UpvotePersentage = int32(setUpvotePersentage(post))

	// err = repo.data.Update(bson.M{"id": postID}, &post)
	err = repo.dbHelper.UpdatePost(repo, postID, post)
	if err != nil {
		log.Println("pckg posts, UpVotePost, repo.data.Update", err.Error())
		return nil, err
	}

	return post, nil
}

func (repo *ItemMemoryRepository) DownVotePost(postID, userID string) (*Post, error) {
	// post := &Post{}
	// err := repo.data.Find(bson.M{"id": postID}).One(&post)

	post, err := repo.dbHelper.FindPostByID(repo, postID)
	if err != nil {
		log.Println("pckg posts, DownVotePost, repo.data.Find", err.Error())
		return nil, err
	}

	_ = changeVote(post, userID, -1)
	post.Score = setPostScore(post)
	post.UpvotePersentage = int32(setUpvotePersentage(post))

	// err = repo.data.Update(bson.M{"id": postID}, &post)
	err = repo.dbHelper.UpdatePost(repo, postID, post)
	if err != nil {
		log.Println("pckg posts, DownVotePost, repo.data.Update", err.Error())
		return nil, err
	}

	return post, nil
}

func (repo *ItemMemoryRepository) UnVotePost(postID, userID string) (*Post, error) {
	// post := &Post{}
	// err := repo.data.Find(bson.M{"id": postID}).One(&post)

	post, err := repo.dbHelper.FindPostByID(repo, postID)
	if err != nil {
		log.Println("pckg posts, UnVotePost, repo.data.Find", err.Error())
		return nil, err
	}

	_, err = removeVote(post, userID)
	if err != nil {
		log.Println("pckg posts, UnVotePost, removeVote", err.Error())
		return nil, err
	}

	post.Score = setPostScore(post)
	post.UpvotePersentage = int32(setUpvotePersentage(post))

	// err = repo.data.Update(bson.M{"id": postID}, &post)
	err = repo.dbHelper.UpdatePost(repo, postID, post)
	if err != nil {
		log.Println("pckg posts, UnVotePost, repo.data.Update", err.Error())
		return nil, err
	}

	return post, nil
}

func removeVote(post *Post, userID string) (*Vote, error) {
	for idx, item := range post.Votes {
		if item.UserID == userID {
			copy(post.Votes[idx:], post.Votes[idx+1:])
			post.Votes[len(post.Votes)-1] = &Vote{}
			post.Votes = post.Votes[:len(post.Votes)-1]
			return item, nil
		}
	}

	return nil, errors.New("no such vote found")
}

func changeVote(post *Post, userID string, vote int32) *Vote {
	for _, item := range post.Votes {
		if item.UserID == userID {
			item.Vote = vote
			return item
		}
	}

	newVote := &Vote{
		UserID: userID,
		Vote:   vote,
	}
	post.Votes = append(post.Votes, newVote)
	return newVote
}

func setUpvotePersentage(post *Post) int {
	votesPersentage := 0
	for _, item := range post.Votes {
		if item.Vote == 1 {
			votesPersentage++
		}
	}

	if len(post.Votes) == 0 {
		return 0
	}

	votesPersentage = votesPersentage * 100 / len(post.Votes)
	return votesPersentage
}

func setPostScore(post *Post) int32 {
	score := 0
	for _, item := range post.Votes {
		score += int(item.Vote)
	}

	return int32(score)
}

func getFormatData() string {
	date := time.Now()
	dateUMD := date.Format("2006-01-02")
	dateHMS := date.Format("15:04:05")
	finDate := dateUMD + "T" + dateHMS
	return finDate
}

func AddStaticPosts(collection *mgo.Collection, sess *mgo.Session) {
	if n, _ := collection.Count(); n == 0 {
		collection.Insert(&Post{
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
		})

		collection.Insert(&Post{
			IDBson:     bson.NewObjectId(),
			AuthorBson: "Joe",
			Author: &author.Author{
				UserID:   "3",
				UserName: "Joe",
			},
			Category:         "programming",
			CreationData:     "2007-01-01T04:20:00",
			PostID:           "2ndPost",
			Score:            3,
			Data:             "Some test",
			Title:            "2nd post",
			Type:             "text",
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
		})

		collection.Insert(&Post{
			IDBson:     bson.NewObjectId(),
			AuthorBson: "DeepThought",
			Author: &author.Author{
				UserID:   "1st user",
				UserName: "DeepThought",
			},
			Category:         "news",
			CreationData:     "2007-01-01T04:20:00",
			PostID:           "3rdPost",
			Score:            2,
			Type:             "text",
			Data:             "nothing here",
			Title:            "check this",
			UpvotePersentage: 25,
			Views:            7,
			Votes: []*Vote{
				{
					UserID: "1stUser",
					Vote:   1,
				},
				{
					UserID: "2ndUser",
					Vote:   -1,
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
		})
	}

	fstPost := &Post{}
	err := collection.Find(bson.M{"id": "1stPost"}).One(&fstPost)
	if err != nil {
		log.Println("pckg posts, AddStaticPosts, collection.Find ", err.Error())
		return
	}

	tempRepo := comments.NewCommentRepoStruct()
	_, err = tempRepo.AddStaticComment(sess, fstPost.PostID)
	if err != nil && err.Error() != "static comm already exists" {
		log.Println("pckg posts, AddStaticPosts, AddStaticComment", err.Error())
		return
	}

	fstPost.Comments, err = tempRepo.GetAllComments(sess, fstPost.PostID)
	if err != nil {
		log.Println("pckg posts, AddStaticPosts, GetAllComments", err.Error())
		return
	}

	err = collection.Update(bson.M{"id": "1stPost"}, &fstPost)
	if err != nil {
		log.Println("pckg posts, AddStaticPosts, collection.Update", err.Error())
		return
	}
}
