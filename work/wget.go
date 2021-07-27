package work

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Wget 下载配置文件专用
func Wget(URL, fileName string) error {
	res, httperr := http.Get(URL)
	if httperr != nil {
		fmt.Println(httperr)
		return httperr
	}
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
