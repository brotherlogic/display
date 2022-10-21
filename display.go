package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

var (
	activity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "display_activity",
		Help: "All the auditioned scores",
	}, []string{"message"})
)

// Server main server type
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

func convertArtist(artist string) string {
	if strings.HasSuffix(artist, ")") {
		return artist[:strings.LastIndex(artist, "(")]
	}
	return artist
}

func (s *Server) buildPage(ctx context.Context) {
	conn, err := s.FDialServer(ctx, "recordgetter")

	if err != nil {
		activity.With(prometheus.Labels{"message": fmt.Sprintf("DIAL: %v", err)}).Inc()
	}

	if err == nil {
		defer conn.Close()
		client := pbrg.NewRecordGetterClient(conn)

		var r *pbrg.GetRecordResponse
		r, err = client.GetRecord(ctx, &pbrg.GetRecordRequest{Refresh: true})

		if err != nil {
			activity.With(prometheus.Labels{"message": fmt.Sprintf("GET_RECORD: %v", err)}).Inc()
		}

		if status.Convert(err).Code() == codes.FailedPrecondition {
			r, err = client.GetRecord(ctx, &pbrg.GetRecordRequest{Type: pbrg.RequestType_DIGITAL})
		}

		if err == nil {
			if r.GetRecord().GetRelease().GetInstanceId() != s.curr {
				extra := ""
				if r.GetRecord().GetRelease().GetFormatQuantity() > 1 {
					extra = fmt.Sprintf("{Disk %v}", r.GetDisk())
				}

				if r.GetRecord().GetMetadata().GetCategory() == rcpb.ReleaseMetadata_STAGED_TO_SELL {
					extra += " {SALE}"
				}

				if r.GetRecord().GetMetadata().GetCategory() == rcpb.ReleaseMetadata_UNKNOWN {
					extra += " (Want)"
				}
				if r.GetRecord().GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_DIGITAL ||
					r.GetRecord().GetMetadata().GetFiledUnder() == rcpb.ReleaseMetadata_FILE_CD {
					extra += " (Digital)"
				}

				extra = strings.TrimSpace(extra)

				artist := "Unknown"
				if len(r.GetRecord().GetRelease().GetArtists()) > 0 {
					artist = r.GetRecord().GetRelease().GetArtists()[0].GetName()
				}

				url := "https://secure.gravatar.com/avatar/d44e93769ea7b6bada5578bb0f48f76f?s=300&r=pg&d=mm"
				if len(r.GetRecord().GetRelease().GetImages()) > 0 {
					url = r.GetRecord().GetRelease().GetImages()[0].GetUri()
				}
				err := s.handler(ctx, r.GetRecord().GetRelease().GetTitle(), artist, url, extra)
				if err == nil {
					s.curr = r.GetRecord().GetRelease().GetInstanceId()
				} else {
					s.CtxLog(ctx, fmt.Sprintf("Bad build: %v", err))
				}
			} else {
				s.CtxLog(ctx, fmt.Sprintf("Skipping logging because %v == %v", r.GetRecord().GetRelease().GetInstanceId(), s.curr))
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
				<center>
				<img class="art_image" src="image.jpeg" width=300" height="300">
				<div class="text">
					<div class="artist">{{.Artist}}</div>
					<div class="album">{{.Title}}</div>
					<div class="number">{{.Extra}}</div>
				</div>		
				</center>
			</section>		
		</div>
	</body>
	</html>`)
	if err != nil {
		s.CtxLog(ctx, fmt.Sprintf("PARSED: %v", err))
	}

	os.MkdirAll("/media/scratch/display/", 0777)
	os.Create("/media/scratch/display/display.html")
	f, _ := os.OpenFile("/media/scratch/display/display.html", os.O_WRONLY, 0777)
	defer f.Close()

	t.Execute(f, &temp{
		Title:  title,
		Artist: strings.TrimSpace(convertArtist(artist)),
		Image:  image,
		Extra:  extra})
	buildStyle()
	buildCssNorm()

	err = exec.Command("curl", image, "-o", "/media/scratch/display/image-raw.jpeg").Run()
	if err != nil {
		activity.With(prometheus.Labels{"message": fmt.Sprintf("DOWNLOAD: %v", err)}).Inc()
		return fmt.Errorf("Bad download: %v", err)
	}
	output, err2 := exec.Command("/usr/bin/convert", "/media/scratch/display/image-raw.jpeg", "-resize", "300x300", "/media/scratch/display/image.jpeg").CombinedOutput()
	if err2 != nil {
		activity.With(prometheus.Labels{"message": fmt.Sprintf("CONVERT: %v", err2)}).Inc()
		return fmt.Errorf("Bad convert of %v: %v -> %v", image, err2, string(output))
	}

	conn, err := s.FDialServer(ctx, "filecopier")
	if err != nil {
		activity.With(prometheus.Labels{"message": fmt.Sprintf("DIAL_COPIER: %v", err)}).Inc()
		return err
	}
	defer conn.Close()

	fc := fcpb.NewFileCopierServiceClient(conn)
	_, err = fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/display.html",
		OutputServer: "mdisplay",
		OutputFile:   "/home/simon/index.html",
		Override:     true,
	})
	if err != nil {
		activity.With(prometheus.Labels{"message": fmt.Sprintf("COPY_HTML: %v", err)}).Inc()
		return err
	}
	_, err = fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/style.css",
		OutputServer: "mdisplay",
		OutputFile:   "/home/simon/style.css",
		Override:     true,
	})
	if err != nil {
		activity.With(prometheus.Labels{"message": fmt.Sprintf("COPY_CSS: %v", err)}).Inc()
		return err
	}
	_, err = fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/normalize.css",
		OutputServer: "mdisplay",
		OutputFile:   "/home/simon/normalize.css",
		Override:     true,
	})
	if err != nil {
		activity.With(prometheus.Labels{"message": fmt.Sprintf("COPY_NORM: %v", err)}).Inc()
		return err
	}
	_, err = fc.Copy(ctx, &fcpb.CopyRequest{
		InputServer:  s.Registry.Identifier,
		InputFile:    "/media/scratch/display/image.jpeg",
		OutputServer: "mdisplay",
		OutputFile:   "/home/simon/image.jpeg",
		Override:     true,
	})
	if err != nil {
		activity.With(prometheus.Labels{"message": fmt.Sprintf("COPY_IMAGE: %v", err)}).Inc()
		return err
	}

	activity.With(prometheus.Labels{"message": "SUCCESS"}).Inc()
	return nil
}

func main() {
	server := Init()
	server.PrepServer("display")
	server.Register = server

	err := server.RegisterServerV2(false)
	if err != nil {
		return
	}

	out, err2 := exec.Command("sudo", "apt", "install", "imagemagick", "-y").Output()
	if err2 != nil {
		log.Fatalf("Unable to install imagemgick: %v -> %v", err2, string(out))
	}
	server.backgroundBuild()

	fmt.Printf("%v", server.Serve())
}
