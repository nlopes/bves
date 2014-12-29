package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	TempFilename = "bves"
	Name         = "bves"
	Version      = "0.1.0"
	Author       = "Norberto Lopes <nlopes.ml+bves __at__ gmail.com>"

	VagrantBoxESUrl = "http://www.vagrantbox.es"
	CacheDuration   = 24 // hours
)

var (
	app     = kingpin.New(Name, "A cmdline interface to vagrantbox.es.")
	debug   = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	timeout = app.Flag("timeout", "Timeout for download.").Default("10s").Short('t').Duration()

	list = app.Command("list", "List boxes.")

	show    = app.Command("show", "Show box information.")
	show_id = show.Arg("id", "ID for base box.").Required().Int()

	url    = app.Command("url", "Show box URL")
	url_id = url.Arg("id", "ID for base box.").Required().Int()

	clear_cache = app.Command("clear-cache", "Clear local cache.")
)

type VBNodes struct {
	Nodes []VBNode `json:"nodes"`
}

func download(url string) (body []byte, err error) {
	client := http.Client{
		Timeout: *timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not download from %s: %s", url, err))
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read response from %s: %s", url, err))
	}
	return body, nil
}

func download_vms_details() (vms_details []VBNode, err error) {
	body, err := download(VagrantBoxESUrl)
	if err != nil {
		return nil, err
	}
	vms_details, err = parse_vbes_html(body)
	if err != nil {
		return nil, err
	}
	return vms_details, nil
}

func save_file(file_path string) {
	vms, err := download_vms_details()
	if err != nil {
		log.Fatalln(err)
	}
	data, err := json.Marshal(VBNodes{Nodes: vms})
	if err != nil {
		log.Fatalln(err)
	}
	jsonfile, err := os.Create(file_path)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonfile.Close()
	jsonfile.Write(data)
}

func get_temp_filepath() (temp_file_fullpath string) {
	return fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, TempFilename)
}

func get_cached_file() (fh *os.File) {
	temp_file_fullpath := get_temp_filepath()
	file_info, err := os.Stat(temp_file_fullpath)
	re_download := false
	if os.IsNotExist(err) {
		log.Printf("Creating local cache.")
		re_download = true
		// If file older than 1 day, re-download file
	} else if time.Now().Sub(file_info.ModTime()).Hours() > CacheDuration {
		if *debug {
			log.Printf("Local cache is older than %d hours.", CacheDuration)
		}
		re_download = true
	}
	if re_download {
		if *debug {
			log.Printf("Re-downloading list from %s.\n", VagrantBoxESUrl)
		}
		save_file(temp_file_fullpath)
	}
	fh, err = os.Open(temp_file_fullpath)
	if err != nil {
		log.Fatalln(err)
	}
	return fh
}

func wipe_cache() {
	if *debug {
		log.Printf("Removing %s \n", get_temp_filepath())
	}
	err := os.Remove(get_temp_filepath())
	if err != nil && !os.IsNotExist(err) {
		log.Fatalln(err)
	}
}

func validate_id(number_boxes, box_id int) {
	if box_id > number_boxes-1 || box_id < 0 {
		log.Fatalf("Invalid box id `%d`\n", box_id)
	}
}

func main() {
	app.Version(Version)
	arg := kingpin.MustParse(app.Parse(os.Args[1:]))
	if arg == clear_cache.FullCommand() {
		wipe_cache()
		return
	}
	fh := get_cached_file()
	defer fh.Close()
	json_decoder := json.NewDecoder(fh)
	vms := VBNodes{}
	if err := json_decoder.Decode(&vms); err != nil {
		log.Fatalln(err)
	}
	switch arg {
	case url.FullCommand():
		validate_id(len(vms.Nodes), *url_id)
		fmt.Println(vms.Nodes[*url_id-1].URL)
	case show.FullCommand():
		validate_id(len(vms.Nodes), *show_id)
		node := vms.Nodes[*show_id-1]
		fmt.Printf("Name:     %s\n", strings.Split(node.Name, "\n")[0])
		fmt.Printf("Details:  %s\n", strings.Join(strings.Split(node.Name, "\n")[0:], "\n"))
		fmt.Printf("Provider: %s\n", node.Provider)
		fmt.Printf("URL:      %s\n", node.URL)
		fmt.Printf("Size:     %.1f MB\n", node.Size)
	case list.FullCommand():
		for _, node := range vms.Nodes {
			fmt.Printf("%5d %s\n", node.Id, strings.Split(node.Name, "\n")[0])
		}
	case "":
		app.Usage(os.Stdout)
	}
}
