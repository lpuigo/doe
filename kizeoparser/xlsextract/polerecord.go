package xlsextract

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PoleRecord struct {
	Date   string
	Hour   string
	SRO    string
	Ref    string
	Images map[string]string
	long   float64
	lat    float64
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

func (pr PoleRecord) GetImage(dir, label string) error {
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

	imgName := fmt.Sprintf("%s %s.jpg", pr.Ref, label)

	imgFullName := filepath.Join(dir, imgName)
	imgFile, err := os.Create(imgFullName)
	if err != nil {
		return fmt.Errorf("could not create image file '%s': %s\n", imgFullName, err.Error())
	}
	defer imgFile.Close()
	_, err = io.Copy(imgFile, resp.Body)
	if err != nil {
		return fmt.Errorf("could not write image file '%s': %s\n", imgFullName, err.Error())
	}
	return nil
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
