package posts

import (
	"math/rand"
	"myRedditClone/pkg/author"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//go:generate mockgen -source=DBHelper.go -destination=mock_dbhelper.go -package=posts DBHelperRepo
type DBHelperRepo interface {
	FindAll(repo *ItemMemoryRepository) ([]*Post, error)
	FindPostByID(repo *ItemMemoryRepository, id string) (*Post, error)
	FindCategory(repo *ItemMemoryRepository, category string) ([]*Post, error)
	FindAllUserPosts(repo *ItemMemoryRepository, userName string) ([]*Post, error)

	InsertPost(repo *ItemMemoryRepository, post *Post, userID, userLogin string) error
	RemovePost(repo *ItemMemoryRepository, pathURL string) error
	UpdatePost(repo *ItemMemoryRepository, postID string, post *Post) error
}

type DBHelperStruct struct{}

func NewDBHelperRepo() *DBHelperStruct {
	return &DBHelperStruct{}
}

func (db *DBHelperStruct) FindAll(repo *ItemMemoryRepository) ([]*Post, error) {
	posts := []*Post{}
	err := repo.data.Find(bson.M{}).All(&posts)
	return posts, err
}

func (db *DBHelperStruct) FindPostByID(repo *ItemMemoryRepository, id string) (*Post, error) {
	post := &Post{}
	err := repo.data.Find(bson.M{"id": id}).One(&post)
	return post, err
}

func (db *DBHelperStruct) FindCategory(repo *ItemMemoryRepository, category string) ([]*Post, error) {
	categoryPosts := []*Post{}
	err := repo.data.Find(bson.M{"category": category}).All(&categoryPosts)
	return categoryPosts, err
}

func (db *DBHelperStruct) FindAllUserPosts(repo *ItemMemoryRepository, userName string) ([]*Post, error) {
	userPosts := []*Post{}
	err := repo.data.Find(bson.M{"authorBson": userName}).All(&userPosts)
	return userPosts, err
}

func (db *DBHelperStruct) InsertPost(repo *ItemMemoryRepository, post *Post, userID, userLogin string) error {
	post.IDBson = bson.NewObjectId()
	post.Author = &author.Author{
		UserID:   userID,
		UserName: userLogin,
	}
	post.Votes = []*Vote{
		{
			UserID: userID,
			Vote:   1,
		},
	}

	rand.Seed(time.Now().UnixNano())
	randID := string(randomBytes(10))

	post.PostID = randID

	err := repo.data.Insert(&post)
	return err
}

func (db *DBHelperStruct) RemovePost(repo *ItemMemoryRepository, pathURL string) error {
	err := repo.data.Remove(bson.M{"id": pathURL})
	return err
}

func (db *DBHelperStruct) UpdatePost(repo *ItemMemoryRepository, postID string, post *Post) error {
	err := repo.data.Update(bson.M{"id": postID}, &post)
	return err
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func randomBytes(len int) []byte {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(97, 122))
	}
	return bytes
}
