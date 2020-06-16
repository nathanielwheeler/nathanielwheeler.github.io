package controllers

import (
	"fmt"
	"net/http"

	"nathanielwheeler.com/models"
	"nathanielwheeler.com/views"
)

// Posts will hold information about views and services
type Posts struct {
	New *views.View
	ps  models.PostsService
}

// NewPosts is a constructor for Posts struct
func NewPosts(ps models.PostsService) *Posts {
	return &Posts{
		New: views.NewView("app", "posts/new"),
		ps:  ps,
	}
}

// PostForm will hold information for creating a new post
type PostForm struct {
	Title string `schema:"title"`
}

// Create : POST /posts
func (p *Posts) Create(res http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form PostForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, vd)
		return
	}
	post := models.Post{
		Title: form.Title,
	}
	if err := p.ps.Create(&post); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, vd)
		return
	}
	fmt.Fprintln(res, post)
}

// Update : PUT /posts/:id
func (p *Posts) Update(res http.ResponseWriter, req *http.Request) {
	// TODO implement
}

// Delete : DELETE /posts/:id
func (p *Posts) Delete(res http.ResponseWriter, req *http.Request) {
	// TODO implement
}