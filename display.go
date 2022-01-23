package main

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	fcpb "github.com/brotherlogic/filecopier/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pbrg "github.com/brotherlogic/recordgetter/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	curr int32
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
	rcpb.RegisterClientUpdateServiceServer(server, s)
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

type temp struct {
	Title  string
	Artist string
	Image  string
}

func (s *Server) buildPage(ctx context.Context) {
	conn, err := s.FDialServer(ctx, "recordgetter")
	if err == nil {
		defer conn.Close()
		client := pbrg.NewRecordGetterClient(conn)

		r, err := client.GetRecord(ctx, &pbrg.GetRecordRequest{Refresh: true})
		if err == nil {
			if r.GetRecord().GetRelease().GetInstanceId() != s.curr {
				s.handler(ctx, r.GetRecord().GetRelease().GetTitle(), r.GetRecord().GetRelease().GetArtists()[0].GetName(), r.GetRecord().GetRelease().GetImages()[0].GetUri())
				s.curr = r.GetRecord().GetRelease().GetInstanceId()
			}
		}
	}
}

func (s *Server) handler(ctx context.Context, title, artist, image string) {
	t := template.New("page")
	t, err := t.Parse(`<html>
	<link rel="stylesheet" href="normalize.css">
	<link rel="stylesheet" href="style.css">
	<style>
		.artwork {
			background-image: url("index.jpeg");
		}
	</style>
	<meta http-equiv="refresh" content="60">
      <body>
		<div id="container">	
			<div class="artwork"></div>
			<section id="main">
				<img class="art_image" src="image.jpeg" width="500" height="500">
				<div class="text">
					<div class="artist">{{.Artist}}</div>
					<div class="album">{{.Title}}</div>
				</div>		
			</section>		
		</div>
	</body>
	</html>`)
	s.Log(fmt.Sprintf("PARSED: %v", err))

	os.MkdirAll("/media/scratch/display/", 0777)
	os.Create("/media/scratch/display/display.html")
	f, _ := os.OpenFile("/media/scratch/display/display.html", os.O_WRONLY, 0777)
	defer f.Close()

	t.Execute(f, &temp{
		Title:  title,
		Artist: artist,
		Image:  image})
	buildStyle()
	buildCssNorm()

	err = exec.Command("curl", image, "-o", "/media/scratch/display/image-raw.jpeg").Run()
	err2 := exec.Command("/usr/bin/convert", "/media/scratch/display/image-raw.jpeg", "-resize", "500x500", "/media/scratch/display/image.jpeg")

	s.Log(fmt.Sprintf("Built everything: %v and %v", err, err2))

	conn, _ := s.FDialServer(ctx, "filecopier")
	defer conn.Close()

	fc := fcpb.NewFileCopierServiceClient(conn)
	fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/display.html",
		OutputServer: "rdisplay",
		OutputFile:   "/home/simon/index.html",
	})
	fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/style.css",
		OutputServer: "rdisplay",
		OutputFile:   "/home/simon/style.css",
	})
	fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/normalize.css",
		OutputServer: "rdisplay",
		OutputFile:   "/home/simon/normalize.css",
	})
	fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/image.jpeg",
		OutputServer: "rdisplay",
		OutputFile:   "/home/simon/image.jpeg",
	})
}

func main() {
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("display", false, true)
	if err != nil {
		return
	}

	exec.Command("sudo", "apt", "install", "imagemagick", "-y").Run()

	fmt.Printf("%v", server.Serve())
}
