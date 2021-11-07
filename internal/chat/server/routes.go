package server

func (s *Server) routes() {
	s.router.HandleFunc("/user", s.postOnly(s.handleUserCreate()))
	s.router.HandleFunc("/user/login", s.postOnly(s.handleUserLogin()))
}