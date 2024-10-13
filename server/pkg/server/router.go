package server

func (s *Server) registerRouter() error {
	s.engine.GET("/users", s.listUsers)
	s.engine.GET("/users/:id/posts", s.listPosts)
	return nil
}
