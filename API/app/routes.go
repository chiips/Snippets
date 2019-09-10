package app

//Routes initiates our Server's routes
func (s *Server) Routes() {

	//Sample user routes
	//authenticateJWT middleware on routes that require authorization
	s.Router.GET("/api/search", s.searchUsers())
	s.Router.POST("/api/signup", s.signup())
	s.Router.PUT("/api/profilephoto/:userid", s.authenticateJWT(s.editProfilePhoto()))
	s.Router.DELETE("/api/profile/:userid", s.authenticateJWT(s.deleteUser()))

	//Sample post routes
	//authenticateJWT middleware on routes that require authorization
	s.Router.GET("/api/posts", s.allPosts())
	s.Router.POST("/api/post", s.authenticateJWT(s.submitPost()))
	s.Router.PUT("/api/post", s.authenticateJWT(s.editPost()))
	s.Router.DELETE("/api/post/:postid", s.authenticateJWT(s.deletePost()))
}
