package hibpdownloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

const apiEndpoint = "https://api.pwnedpasswords.com/range/%s%s"

type hibpDownloader struct {
	ntlm         bool
	file         string
	n_workers    uint64
	fp           *os.File
	hex5         <-chan string
	data         chan []byte
	quit         chan bool
	progress_max uint64
	progress_bar *progressbar.ProgressBar
}

func Download(file string, parallelism uint64, overwrite bool, ntlm bool, showProgress bool) {
	var err error
	hd := &hibpDownloader{
		ntlm:      ntlm,
		file:      file,
		n_workers: parallelism,
		data:      make(chan []byte, 100),
		quit:      make(chan bool),
	}
	hd.progress_max, hd.hex5 = hex5generator()

	if _, err := os.Stat(hd.file); err == nil {
		if !overwrite {
			panic(fmt.Errorf("file `%s` already exist", hd.file))
		}
	}

	hd.fp, err = os.OpenFile(hd.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0440)
	if err != nil {
		panic(fmt.Errorf("failed open file `%s`: %w", hd.file, err))
	}
	defer hd.fp.Close()

	hashType := ""
	if hd.ntlm {
		hashType += "NTLM"
	} else {
		hashType += "SHA1"
	}

	fmt.Printf("Donwloading %s hashes with %d workers\n\n", hashType, hd.n_workers)

	if showProgress {
		desc := fmt.Sprintf("%s (%s)", path.Base(hd.file), hashType)
		hd.progress_bar = progressbar.NewOptions64(
			int64(hd.progress_max),
			progressbar.OptionSetTheme(progressbar.ThemeASCII),
			progressbar.OptionSetDescription(desc),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetItsString("pcs"),
			progressbar.OptionSetElapsedTime(true),
			progressbar.OptionShowCount(),
			progressbar.OptionShowElapsedTimeOnFinish(),
			progressbar.OptionShowIts(),
			progressbar.OptionThrottle(time.Second*5))
	}

	var ww sync.WaitGroup
	ww.Add(1)
	go func() {
		hd.writer(hd.data, hd.quit)
		ww.Done()
	}()

	var wg sync.WaitGroup
	for range hd.n_workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hd.downloader()
		}()
	}
	wg.Wait()
	time.Sleep(time.Second * 5)
	hd.quit <- true
	ww.Wait()
}

func (hd *hibpDownloader) downloader() {
	client := &http.Client{
		Timeout: time.Duration(time.Second * 5),
	}
	ntlmStr := "?mode=ntlm"
	if !hd.ntlm {
		ntlmStr = ""
	}
	for hex5 := range hd.hex5 {

		ok := false
		retry_count := 30
		for !ok {
			if retry_count > 0 {
				retry_count--
			}
			url := fmt.Sprintf(apiEndpoint, hex5, ntlmStr)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				panic(fmt.Errorf("failed to create http request for `%s`: %w", url, err))
			}
			resp, err := client.Do(req)
			if err != nil {
				if retry_count <= 0 {
					fmt.Printf("error fetching %s: %s\n", hex5, err)
				}
				continue
			}
			if resp.StatusCode != http.StatusOK {
				if retry_count <= 0 {
					fmt.Printf("error fetching %s: status %d\n", hex5, resp.StatusCode)
				}
				continue
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				if retry_count <= 0 {
					fmt.Printf("error reading %s: %s\n", hex5, err)
				}
				resp.Body.Close()
				continue
			}
			hd.data <- hd.applyHex5Prefix(hex5, data)
			ok = true
		}
	}
}

func (hd *hibpDownloader) applyHex5Prefix(hex5 string, data []byte) []byte {
	var b strings.Builder
	for _, s := range strings.Split(string(data), "\r\n") {
		fmt.Fprintf(&b, "%s%s\r\n", hex5, strings.ToUpper(s))
	}
	return []byte(b.String())
}

func (hd *hibpDownloader) writer(data chan []byte, quit chan bool) {
	for {
		select {
		case blob := <-data:
			_, err := hd.fp.Write(blob)
			if err != nil {
				panic(err)
			}

			hd.progress_bar.Add64(1)
		case <-quit:
			return
		}
	}
}
