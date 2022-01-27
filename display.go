package main

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/brotherlogic/filecopier/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	"github.com/brotherlogic/goserver/utils"
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
		curr:     int32(-1),
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
	Extra  string
}

func (s *Server) backgroundBuild() {
	go func() {
		ctx, cancel := utils.ManualContext("display-background", time.Minute*10)
		defer cancel()
		s.buildPage(ctx)
	}()
}

func (s *Server) buildPage(ctx context.Context) {
	conn, err := s.FDialServer(ctx, "recordgetter")
	if err == nil {
		defer conn.Close()
		client := pbrg.NewRecordGetterClient(conn)

		r, err := client.GetRecord(ctx, &pbrg.GetRecordRequest{})

		if status.Convert(err).Code() == codes.FailedPrecondition {
			r, err = client.GetRecord(ctx, &pbrg.GetRecordRequest{Type: pbrg.RequestType_DIGITAL})
		}

		if err == nil {
			if r.GetRecord().GetRelease().GetInstanceId() != s.curr {
				extra := ""
				if r.GetDisk() > 1 {
					extra = fmt.Sprintf("{Disk %v", r.GetDisk())
				}

				if r.GetRecord().GetMetadata().GetCategory() == rcpb.ReleaseMetadata_UNKNOWN {
					extra = "(Want)"
				}
				if r.GetRecord().GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_DIGITAL ||
					r.GetRecord().GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_CD {
					extra = "(Digital)"
				}
				err := s.handler(ctx, r.GetRecord().GetRelease().GetTitle(), r.GetRecord().GetRelease().GetArtists()[0].GetName(), r.GetRecord().GetRelease().GetImages()[0].GetUri(), extra)
				if err == nil {
					s.curr = r.GetRecord().GetRelease().GetInstanceId()
				} else {
					s.Log(fmt.Sprintf("Bad build: %v", err))
				}
			}
		}
	}
}

func (s *Server) handler(ctx context.Context, title, artist, image, extra string) error {
	t := template.New("page")
	t, err := t.Parse(`<html>
	<link rel="stylesheet" href="normalize.css">
	<link rel="stylesheet" href="style.css">
	<style>
		.artwork {
			background-image: url("image.jpeg");
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
					<div class="number">{{.Extra}}</div>
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
		Image:  image,
		Extra:  extra})
	buildStyle()
	buildCssNorm()

	err = exec.Command("curl", image, "-o", "/media/scratch/display/image-raw.jpeg").Run()
	if err != nil {
		return fmt.Errorf("Bad download: %v", err)
	}
	err2 := exec.Command("/usr/bin/convert", "/media/scratch/display/image-raw.jpeg", "-resize", "500x500", "/media/scratch/display/image.jpeg").Run()
	if err2 != nil {
		return fmt.Errorf("Bad convert of %v: %v", image, err2)
	}

	conn, err := s.FDialServer(ctx, "filecopier")
	if err != nil {
		return err
	}
	defer conn.Close()

	fc := fcpb.NewFileCopierServiceClient(conn)
	_, err = fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/display.html",
		OutputServer: "rdisplay",
		OutputFile:   "/home/simon/index.html",
	})
	if err != nil {
		return err
	}
	_, err = fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/style.css",
		OutputServer: "rdisplay",
		OutputFile:   "/home/simon/style.css",
	})
	if err != nil {
		return err
	}
	_, err = fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/normalize.css",
		OutputServer: "rdisplay",
		OutputFile:   "/home/simon/normalize.css",
	})
	if err != nil {
		return err
	}
	_, err = fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/image.jpeg",
		OutputServer: "rdisplay",
		OutputFile:   "/home/simon/image.jpeg",
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("display", false, true)
	if err != nil {
		return
	}

	err2 := exec.Command("sudo", "apt", "install", "imagemagick", "-y").Run()
	server.Log(fmt.Sprintf("INSTALLED: %v", err2))
	server.backgroundBuild()

	fmt.Printf("%v", server.Serve())
}
