package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/axgle/mahonia"
	"golang.org/x/net/proxy"
)

const (
	path          = "/Users/dracarysX/Pictures/每日一wen"
	userAgent     = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
	socketAddress = "127.0.0.1:1080"
)

type image struct {
	title  string
	imgURL []string
}

type imageInfo struct {
	name    string
	imgByte []byte
	err     error
}

var imgChan = make(chan imageInfo)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func convertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func request(client *http.Client, url string) ([]byte, error) {
	log.Printf("request url: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	res, err := client.Do(req)
	if err != nil {
		log.Printf("request failure, url: %s, err: %v\n", url, err)
		return nil, err
	}
	body := res.Body
	defer body.Close()
	bodyByte, _ := ioutil.ReadAll(body)
	return bodyByte, nil
}

func parseHTML(body []byte) image {
	data := convertToString(string(body), "gbk", "utf-8")
	tableRegexp := regexp.MustCompile(`<br><br><br><table style="border:1px solid #D4EFF7;width:98%">(.*?)</table>`)
	table := tableRegexp.FindAllString(data, -1)[0]
	titleRegexp := regexp.MustCompile(`<span class="f24">(.*?)<br></span>`)
	title := titleRegexp.FindAllStringSubmatch(table, -1)[0][1]
	log.Printf("%v\n", title)
	imgRegexp := regexp.MustCompile(`<img src='(.*?)' onclick`)
	imgURL := make([]string, 0)
	for _, img := range imgRegexp.FindAllStringSubmatch(table, -1) {
		imgURL = append(imgURL, img[1])
	}
	return image{title, imgURL}
}

func getImgName(url string) (string, bool) {
	s := strings.Split(url, "/")
	if len(s) == 0 {
		return "", false
	}
	var name string
	// re := regexp.MustCompile(`\/[0-9a-zA-Z_-]+\.(jpg|png|jpeg)`)
	// name := re.FindAllStringSubmatch(url, -1)[0][0]
	name = s[len(s)-1]
	if len(strings.Split(name, ".")) == 1 {
		name = name + ".jpg"
	}
	// log.Println("image name: ", name)
	return name, true
}

func saveImage(info imageInfo, path string) {
	if info.err != nil {
		return
	}
	imgName := path + "/" + info.name
	exists, err := pathExists(imgName)
	if err != nil {
		log.Println("check img path failure: ", err)
		return
	}
	if exists {
		return
	}
	var fh *os.File
	fh, _ = os.Create(imgName)
	defer fh.Close()
	log.Println("saving image: ", imgName)
	fh.Write(info.imgByte)
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("please input start url.")
	}
	startURL := os.Args[1]
	dialer, err := proxy.SOCKS5("tcp", socketAddress, nil, proxy.Direct)
	if err != nil {
		log.Fatalln("can't connect to the proxy: ", err)
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	body, err := request(httpClient, startURL)
	if err != nil {
		log.Fatalln("request start url failure, err: ", err)
	}
	image := parseHTML(body)
	folderPath := path + "/" + image.title
	ok, err := pathExists(folderPath)
	if err != nil {
		log.Fatalln("check folderpath failure: ", err)
	}
	if !ok {
		err := os.Mkdir(folderPath, 0777)
		if err != nil {
			log.Fatalln("create folderpath failure: ", err)
		}
	}
	for _, url := range image.imgURL {
		name, ok := getImgName(url)
		if ok {
			go func(n string, u string) {
				imgByte, err := request(httpClient, u)
				if err != nil {
					log.Printf("download image failure, url: %s, err: %v\n", u, err)
				}
				imgChan <- imageInfo{n, imgByte, err}
			}(name, url)
		}
	}
	for i := 0; i < len(image.imgURL); i++ {
		select {
		case info := <-imgChan:
			go saveImage(info, folderPath)
		}
	}
	log.Printf("download image completes, counts: %d\n", len(image.imgURL))
}
