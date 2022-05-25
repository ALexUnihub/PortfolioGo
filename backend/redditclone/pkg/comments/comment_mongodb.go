package comments

import (
	"errors"
	"log"
	"math/rand"
	"myRedditClone/pkg/author"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CommentRepoStruct struct{}

func NewCommentRepoStruct() *CommentRepoStruct {
	return &CommentRepoStruct{}
}

func (c *CommentRepoStruct) GetAllComments(session *mgo.Session, postID string) ([]*Comment, error) {
	collection := session.DB("posts").C(postID)
	comm := []*Comment{}

	err := collection.Find(bson.M{}).All(&comm)
	if err != nil {
		log.Println("pckg comments, GetAllComments: ", err.Error())
		return nil, err
	}
	return comm, nil
}

func (c *CommentRepoStruct) CreateComment(session *mgo.Session, postID, data, userID, userLogin string) (*Comment, error) {
	collection := session.DB("posts").C(postID)

	rand.Seed(time.Now().UnixNano())
	randID := string(randomBytes(10))

	newComm := &Comment{
		IDBson: bson.NewObjectId(),
		Author: &author.Author{
			UserID:   userID,
			UserName: userLogin,
		},
		CommentBody:  data,
		CreationData: getFormatData(),
		CommentID:    randID,
	}

	err := collection.Insert(&newComm)
	if err != nil {
		log.Println("pckg comments, CreateComment, collection.Insert", err.Error())
		return nil, err
	}

	return newComm, nil
}

func (c *CommentRepoStruct) DeleteCommentFromRepo(session *mgo.Session, postID, commID, userID, userLogin string) error {
	collection := session.DB("posts").C(postID)

	commCheck := &Comment{}
	err := collection.Find(bson.M{"id": commID}).One(commCheck)
	if err != nil {
		log.Println("pckg comments, DeleteCommentFromRepo, collection.Find ", err.Error())
		return err
	}

	if commCheck.Author.UserID != userID || commCheck.Author.UserName != userLogin {
		log.Println("pckg comments, DeleteCommentFromRepo, wrong user ID or user login ")
		return errors.New("wrong user ID or user login")
	}

	err = collection.Remove(bson.M{"id": commID})
	if err != nil {
		log.Println("pckg comments, DeleteCommentFromRepo, collection.Remove ", err.Error())
		return err
	}

	return nil
}

func (c *CommentRepoStruct) AddStaticComment(session *mgo.Session, postID string) (*Comment, error) {
	collection := session.DB("posts").C(postID)
	checkComm := &Comment{}

	err := collection.Find(bson.M{"id": "23"}).One(checkComm)
	if err == nil {
		return nil, errors.New("static comm already exists")
	}

	newCommObj := &Comment{
		IDBson: bson.NewObjectId(),
		Author: &author.Author{
			UserID:   "1st user",
			UserName: "DeepThought",
		},
		CommentBody:  "Yes",
		CreationData: "2022-04-18T16:20:00",
		CommentID:    "23",
	}

	err = collection.Insert(&newCommObj)
	if err != nil {
		log.Println("pckg comments, AddStaticComment: ", err.Error())
		return nil, err
	}

	return newCommObj, nil
}

func getFormatData() string {
	date := time.Now()
	dateUMD := date.Format("2006-01-02")
	dateHMS := date.Format("15:04:05")
	finDate := dateUMD + "T" + dateHMS
	return finDate
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
