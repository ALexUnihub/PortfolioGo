package comments

import (
	"myRedditClone/pkg/author"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Comment struct {
	IDBson       bson.ObjectId  `json:"-" bson:"_id"`
	Author       *author.Author `json:"author" bson:"author"`
	CommentBody  string         `json:"body" bson:"body"`
	CreationData string         `json:"created" bson:"created"`
	CommentID    string         `json:"id" bson:"id"`
}

//go:generate mockgen -source=comment.go -destination=mock_comments.go -package=comments CommentsRepo
type CommentsRepo interface {
	GetAllComments(session *mgo.Session, postID string) ([]*Comment, error)
	CreateComment(session *mgo.Session, postID, data, userID, userLogin string) (*Comment, error)
	DeleteCommentFromRepo(session *mgo.Session, postID, commID, userID, userLogin string) error

	AddStaticComment(session *mgo.Session, postID string) (*Comment, error)
}
