package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/display/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
	}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterDisplayServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "magic", Value: int64(12345)},
	}
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	t := template.New("page")
	t, err := t.Parse(`</head>
	<body>
		<div id="container">	
			<div class="artwork"></div>
			<section id="main">
				<img class="art_image" src="https://img.discogs.com/VsImbPqwzP7fNEM_Ws_y7Lkbh7Q=/fit-in/600x598/filters:strip_icc():format(jpeg):mode_rgb():quality(90)/discogs-images/R-10604586-1500801461-4954.jpeg.jpg" width="500" height="500">
				<div class="text">
					<div class="artist">Bernard Estardy</div>
					<div class="album">Piano & Orgues</div>
				</div>		
			</section>		
		</div>
	</body>
	</html>`)
	s.Log(fmt.Sprintf("PARSED: %v", err))
	t.Execute(w, nil)
}

func main() {
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("display", false, true)
	if err != nil {
		return
	}

	http.HandleFunc("/", server.handler)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		server.Log(fmt.Sprintf("Unable to server: %v", err))
	}()

	fmt.Printf("%v", server.Serve())
}
