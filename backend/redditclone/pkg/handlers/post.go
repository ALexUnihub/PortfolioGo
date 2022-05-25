package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"myRedditClone/pkg/posts"
	"myRedditClone/pkg/session"
	"net/http"
)

type PostsHandler struct {
	PostsRepo posts.PostsRepo
	Sessions  *session.SessionManager
}

type DeleteMsg struct {
	Msg string `json:"message"`
}

type CommentMsg struct {
	Comment string `json:"comment"`
}

// for empty comment
type CommentErrors struct {
	CommErrors []*CommentErr `json:"errors"`
}

type CommentErr struct {
	Location string `json:"location"`
	MsgComm  string `json:"msg"`
	Param    string `json:"param"`
}

// /for empty comment

func (h *PostsHandler) List(w http.ResponseWriter, r *http.Request) {
	elems, err := h.PostsRepo.GetAllPosts()
	if err != nil {
		http.Error(w, `pckg handlers, List err`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, len(elems), elems)
	if err != nil {
		log.Printf("List err:%s", err.Error())
	}
}

func (h *PostsHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Path
	postID = postID[len("/api/post/"):]

	post, err := h.PostsRepo.GetPost(postID)
	if err != nil {
		http.Error(w, `pckg handlers, ShowPost no such post`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, 1, post)
	if err != nil {
		log.Printf("ShowPost err:%s", err.Error())
	}
}

func (h *PostsHandler) ShowCategory(w http.ResponseWriter, r *http.Request) {
	postCategory := r.URL.Path
	postCategory = postCategory[len("/api/posts/"):]

	categoryPosts, err := h.PostsRepo.GetCategory(postCategory)
	if err != nil {
		http.Error(w, `pckg handlers, ShowCategory no posts in category`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, len(categoryPosts), categoryPosts)
	if err != nil {
		log.Printf("ShowCategory err:%s", err.Error())
	}
}

func (h *PostsHandler) ShowAllUserPosts(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Path
	userName = userName[len("/api/user/"):]

	userPosts, err := h.PostsRepo.GetAllUserPosts(userName)
	if err != nil {
		http.Error(w, `pckg handlers, ShowAllUserPosts, user has no posts`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, len(userPosts), userPosts)
	if err != nil {
		log.Printf("ShowAllUserPosts err:%s", err.Error())
	}
}

func (h *PostsHandler) AddNewPost(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `pckg handlers, AddNewPost SessionFromContext err`, http.StatusInternalServerError)
		return
	}
	newPost, err := h.PostsRepo.CreatePost(r, sess.UserID, sess.UserLogin)
	if err != nil {
		http.Error(w, `pckg handlers, AddNewPost CreatePost err`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, 1, newPost)
	if err != nil {
		log.Printf("AddNewPost err:%s", err.Error())
	}
}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	err := h.PostsRepo.DeletePostFromRepo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := &DeleteMsg{
		Msg: "success",
	}

	err = sendJSON(w, 0, msg)
	if err != nil {
		log.Printf("DeletePost err:%s", err.Error())
	}
}

func (h *PostsHandler) AddNewComment(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `pckg handlers, CreateComment SessionFromContext err`, http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, `unknown payload`, http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	comm := &CommentMsg{}
	err = json.Unmarshal(body, comm)
	if err != nil {
		http.Error(w, "pkg handlers cant unpack payload in AddNewComment", http.StatusBadRequest)
		return
	}

	// if comment body is empty
	if comm.Comment == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		sendCommentErr(w)
		return
	}

	postID := r.URL.Path[len("/api/post/"):]

	updatedPost, err := h.PostsRepo.AddCommentInPostRepository(postID, comm.Comment, sess.UserID, sess.UserLogin)
	if err != nil {
		http.Error(w, `pckg handlers, CreateComment err`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, 1, updatedPost)
	if err != nil {
		log.Printf("AddNewComment err:%s", err.Error())
	}
}

// // for empty comment
func sendCommentErr(w http.ResponseWriter) {
	commErrors := &CommentErrors{}

	commErr := &CommentErr{
		Location: "body",
		MsgComm:  "is required",
		Param:    "comment",
	}

	commErrors.CommErrors = append(commErrors.CommErrors, commErr)
	err := sendJSON(w, 0, commErrors)
	if err != nil {
		log.Printf("AddNewComment: sendCommentErr err:%s", err.Error())
	}
}

// // /for empty comment

func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `pckg handlers, DeleteComment SessionFromContext err`, http.StatusInternalServerError)
		return
	}

	postID, commID, err := getPostCommentID(r.URL.Path)
	if err != nil {
		http.Error(w, `pckg handlers, DeleteComment getPostCommentID err`, http.StatusInternalServerError)
		return
	}

	updatedPost, err := h.PostsRepo.DeleteCommentInPostRepository(postID, commID, sess.UserID, sess.UserLogin)
	if err != nil {
		http.Error(w, `pckg handlers, DeleteComment DeleteCommentInPostRepository err`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, 1, updatedPost)
	if err != nil {
		log.Printf("DeleteComment err:%s", err.Error())
	}
}

func (h *PostsHandler) UpVote(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `pckg handlers, UpVote SessionFromContext err`, http.StatusInternalServerError)
		return
	}

	postID, _, err := getPostCommentID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedPost, err := h.PostsRepo.UpVotePost(postID, sess.UserID)
	if err != nil {
		http.Error(w, `pckg handlers, UpVote UpVotePost err`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, 1, updatedPost)
	if err != nil {
		log.Printf("UpVote err:%s", err.Error())
	}
}

func (h *PostsHandler) DownVote(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `pckg handlers, DownVote SessionFromContext err`, http.StatusInternalServerError)
		return
	}

	postID, _, err := getPostCommentID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedPost, err := h.PostsRepo.DownVotePost(postID, sess.UserID)
	if err != nil {
		http.Error(w, `pckg handlers, DownVote DownVotePost err`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, 1, updatedPost)
	if err != nil {
		log.Printf("DownVote err:%s", err.Error())
	}
}

func (h *PostsHandler) UnVote(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `pckg handlers, UnVote SessionFromContext err`, http.StatusInternalServerError)
		return
	}

	postID, _, err := getPostCommentID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedPost, err := h.PostsRepo.UnVotePost(postID, sess.UserID)
	if err != nil {
		http.Error(w, `pckg handlers, UnVote UnVotePost err`, http.StatusInternalServerError)
		return
	}

	err = sendJSON(w, 1, updatedPost)
	if err != nil {
		log.Printf("UnVote err:%s", err.Error())
	}
}

func sendJSON(w http.ResponseWriter, dataLen int, data interface{}) error {
	byteValue, err := json.Marshal(data)
	if err != nil {
		return errors.New("sendJSON, json.Marshal err")
	}

	if dataLen > 0 {
		byteValue = validatePostsType(byteValue, dataLen)
	}

	_, err = w.Write(byteValue)
	if err != nil {
		return errors.New("sendJSON, w.Write err")
	}

	return nil
}

func validatePostsType(data []byte, dataLen int) []byte {
	bytesIdxShift := 0
	for idx := 0; idx < dataLen; idx++ {
		bytesIdxShift = bytes.Index(data[bytesIdxShift:], []byte(`type`)) + len([]byte(`"type":`)) + bytesIdxShift
		switch data[bytesIdxShift] {
		case 't':
			data = bytes.Replace(data, []byte("data"), []byte("text"), 1)
		case 'l':
			data = bytes.Replace(data, []byte("data"), []byte("url"), 1)
		}
	}

	return data
}

func getPostCommentID(link string) (string, string, error) {
	link = link[len("/api/post/"):]

	linkByte := []byte(link)
	slashIdx := bytes.Index(linkByte, []byte(`/`))

	if slashIdx == -1 {
		return "", "", errors.New(`incorrect link`)
	}

	// первым возвращается ID поста, потом комментария
	return string(linkByte[:slashIdx]), string(linkByte[slashIdx+1:]), nil
}
