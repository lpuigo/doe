package xlsextract

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
)

const (
	kizeoImgMaxSize int = 1600
	kizeoImgQuality int = 75
)

type PoleRecord struct {
	Date           string
	Hour           string
	SRO            string
	Ref            string
	Comment        string
	Images         map[string]string
	long           float64
	lat            float64
	ExtractPicture bool
}

func (pr *PoleRecord) String() string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "Pole %s %s (created on %s %s)\n", pr.SRO, pr.Ref, pr.Date, pr.Hour)
	fmt.Fprintf(&sb, "\tGPS : %+.8f, %+.8f\n", pr.lat, pr.long)
	for img, link := range pr.Images {
		fmt.Fprintf(&sb, "\timage %s: %s\n", img, link)
	}
	return sb.String()
}

func (pr *PoleRecord) GetImageLabels() []string {
	res := make([]string, len(pr.Images))
	i := 0
	for label, _ := range pr.Images {
		res[i] = label
		i++
	}
	sort.Strings(res)
	return res
}

func (pr *PoleRecord) GetImage(dir, label string) error {
	url, found := pr.Images[label]
	if !found {
		return fmt.Errorf("could not find image with label '%s'\n", label)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	imgtype := resp.Header.Get("Content-Type")
	_ = imgtype // TODO Check img type for "image/jpeg"

	imgName := fmt.Sprintf("%s %s.jpg", pr.GetSafeRef(), label)

	imgFullName := filepath.Join(dir, imgName)
	imgFile, err := os.Create(imgFullName)
	if err != nil {
		return fmt.Errorf("could not create image file '%s': %s\n", imgFullName, err.Error())
	}
	defer imgFile.Close()

	// process image
	imgBuff, err := processImage(resp.Body)
	if err != nil {
		return fmt.Errorf("could not process image file '%s': %s\n", imgFullName, err.Error())
	}

	// write image to file
	_, err = io.Copy(imgFile, imgBuff)
	if err != nil {
		return fmt.Errorf("could not write image file '%s': %s\n", imgFullName, err.Error())
	}
	return nil
}

// processImage ensure resize / sharpen and quality downgrade on given img body
func processImage(img io.Reader) (io.Reader, error) {
	im, err := imaging.Decode(img, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("image decode returned: %s", err.Error())
	}
	// Resize image if too large
	if im.Bounds().Size().X > im.Bounds().Size().Y && im.Bounds().Size().X > kizeoImgMaxSize {
		im = imaging.Resize(im, kizeoImgMaxSize, 0, imaging.CatmullRom)
	}
	if im.Bounds().Size().Y > im.Bounds().Size().X && im.Bounds().Size().Y > kizeoImgMaxSize {
		im = imaging.Resize(im, 0, kizeoImgMaxSize, imaging.CatmullRom)
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, im, &jpeg.Options{Quality: kizeoImgQuality})
	if err != nil {
		return nil, fmt.Errorf("image jpeg encode returned: %s", err.Error())
	}
	return bytes.NewReader(buf.Bytes()), nil
}

func (pr *PoleRecord) GetAllImages(dir string, parallel int) error {
	workers := make(chan struct{}, parallel)

	for _, imgLabel := range pr.GetImageLabels() {
		workers <- struct{}{}
		go func(dir, label string) {
			err := pr.GetImage(dir, label)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
			<-workers
		}(dir, imgLabel)
	}

	// wait for goroutine completion
	for i := 0; i < parallel; i++ {
		workers <- struct{}{}
	}
	return nil
}

func (pr *PoleRecord) WriteComment(dir string) error {
	if pr.Comment == "" {
		return nil
	}

	commentName := fmt.Sprintf("%s %s.txt", pr.GetSafeRef(), "Commentaire")
	commentFullName := filepath.Join(dir, commentName)
	commentFile, err := os.Create(commentFullName)
	if err != nil {
		return fmt.Errorf("could not create comment file '%s': %s\n", commentFullName, err.Error())
	}
	defer commentFile.Close()
	_, err = fmt.Fprint(commentFile, pr.Comment)
	if err != nil {
		return fmt.Errorf("could not write to comment file '%s': %s\n", commentFullName, err.Error())
	}
	return nil
}

func (pr *PoleRecord) GetSRORef() string {
	ref := ""
	if pr.SRO != "" {
		ref = strings.Trim(pr.SRO, " ") + " "
	}
	ref += strings.Trim(pr.Ref, " ")
	return ref
}

func (pr *PoleRecord) GetSafeSRO() string {
	return safeName(pr.SRO)
}

func (pr *PoleRecord) GetSafeRef() string {
	return safeName(pr.Ref)
}

func safeName(name string) string {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, ":", "")
	name = strings.ReplaceAll(name, "  ", " ")
	name = strings.Trim(name, " \t")
	return name
}
