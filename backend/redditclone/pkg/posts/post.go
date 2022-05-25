package posts

import (
	"myRedditClone/pkg/author"
	"myRedditClone/pkg/comments"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

type Vote struct {
	UserID string `json:"user"`
	Vote   int32  `json:"vote"`
}

type Post struct {
	IDBson           bson.ObjectId       `json:"-" bson:"_id"`
	AuthorBson       string              `json:"-" bson:"authorBson"`
	Author           *author.Author      `json:"author" bson:"author"`
	Category         string              `json:"category" bson:"category"`
	Comments         []*comments.Comment `json:"comments" bson:"comments"`
	CreationData     string              `json:"created" bson:"created"`
	PostID           string              `json:"id" bson:"id"`
	Score            int32               `json:"score" bson:"score"`
	Type             string              `json:"type" bson:"type"`
	Data             string              `json:"data" bson:"data"`
	Title            string              `json:"title" bson:"title"`
	UpvotePersentage int32               `json:"upvotePercentage" bson:"upvotePercentage"`
	Views            uint32              `json:"views" bson:"views"`
	Votes            []*Vote             `json:"votes" bson:"votes"`
}

//go:generate mockgen -source=post.go -destination=repo_mock.go -package=posts PostsRepo
type PostsRepo interface {
	GetAllPosts() ([]*Post, error)
	GetPost(id string) (*Post, error)
	GetCategory(category string) ([]*Post, error)
	GetAllUserPosts(userName string) ([]*Post, error)

	CreatePost(r *http.Request, userID, userLogin string) (*Post, error)
	DeletePostFromRepo(r *http.Request) error

	AddCommentInPostRepository(postID, data, userID, userLogin string) (*Post, error)
	DeleteCommentInPostRepository(postID, commID, userID, userLogin string) (*Post, error)

	UpVotePost(postID, userID string) (*Post, error)
	DownVotePost(postID, userID string) (*Post, error)
	UnVotePost(postID, userID string) (*Post, error)
}
