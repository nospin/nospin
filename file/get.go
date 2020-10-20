package file

import (
	"encoding/json"
	"fmt"
	"log"
	"nospin/config"
	"nospin/user"
	us "os/user"
	"path/filepath"
	"strings"

	"github.com/vitecoin/zi/api"
	zi "github.com/vitecoin/zi/pkg"
)

type File struct {
	ID      string
	Name    string
	Content []byte
	Group   string
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
func Get(id string) File {
	z, err := zi.Zi("https://b3b8cd3fd294.ngrok.io/")
	if err != nil {
		log.Fatalln(err)
	}
	if strings.Index(id, "/") != -1 {
		group := strings.Split(id, "/")
		dir := []api.Pair{}
		g := z.GetAll()
		for _, v := range g {
			if strings.Contains(v.Key, "/") == true {
				data := z.Get(config.Get("name"))
				var z user.User
				json.Unmarshal([]byte(data.Value), &z)
				var file File
				json.Unmarshal([]byte(v.Value), &file)
				tok := z.PrvTok
				group := strings.Split(file.Group, ",")
				if _, ok := find(group, tok); ok == true || file.ID == z.ID {
					dir = append(dir, v)
				}
			}
		}
		usr, err := us.Current()
		if err != nil {
			log.Fatal(err)
		}

		home := usr.HomeDir
		p, _ := filepath.Abs(strings.Join(group[1:], "/"))
		p = strings.Replace(p, home+"/", "", 1)
		for _, f := range dir {
			var file File
			json.Unmarshal([]byte(f.Value), &file)
			fmt.Println(group[1:])
			if p == file.Name || file.Name == strings.Join(group[1:], "/") {
				return file
			}
		}
	}
	return File{}
}