package file

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	us "os/user"
	"quote/auth"
	"quote/config"
	"quote/user"
	"quote/util"

	"path/filepath"
	"strings"

	"github.com/ashtyn3/zi/api"
	zi "github.com/ashtyn3/zi/pkg"
)

type Ops struct {
	Group string
}

func Set(path string, endPath string, options Ops) {
	//
	// godotenv.Load("../.env")
	// url := os.Getenv("url")
	// pd := os.Getenv("pd")
	z, err := zi.Zi(auth.Auth().Url, auth.Auth().Pd)
	if err != nil {
		log.Fatalln(err)
	}
	if strings.Index(path, "/") != -1 {
		group := strings.Split(path, "/")
		f, _ := ioutil.ReadFile(strings.Join(group[1:], "/"))
		usr, err := us.Current()
		if err != nil {
			log.Fatal(err)
		}
		home := usr.HomeDir
		user := user.Get(group[0])
		if user.ID == "" {
			fmt.Println("The user with the name: " + group[0] + " does not exist.")
		}
		image := false
		if strings.HasSuffix(strings.Join(group[1:], "/"), ".png") == true || strings.HasSuffix(strings.Join(group[1:], "/"), ".jpg") == true || strings.HasSuffix(strings.Join(group[1:], "/"), ".jpeg") == true || len(f) > 2000 {
			f = []byte(base64.StdEncoding.EncodeToString(f))
			image = true
		}
		if user.ID == config.Get("name") {
			p, _ := filepath.Abs(strings.Join(group[1:], "/"))
			if endPath != "" {
				data := File{Content: f, Name: endPath, ID: user.ID, Group: user.PrvTok, Image: image}
				if options.Group != "" {
					data.Group = options.Group
				}
				// item, _ := json.Marshal(data)
				fID, _ := util.RanString(6)
				if data.Image == true {
					d := util.ChunkString(string(data.Content), 1000)
					for i, c := range d {
						fmt.Printf("\r\033[K%d/%d", i+1, len(d))
						data = File{Content: []byte(c), Name: endPath, ID: user.ID, Group: user.PrvTok, Image: image}
						n, _ := json.Marshal(data)
						// z.Set(api.Pair{Key: user.ID + "/" + fID, Value: string(n)})
						z.Del(user.ID + "/" + fID)
						z.Dump(api.Pair{Key: user.ID + "/" + fID, Value: string(n)}, fID+".zi")
						// time.Sleep(1 * time.Second)
					}
					bFile, _ := json.Marshal(File{Name: endPath, Group: user.PrvTok, Image: image})
					z.Set(api.Pair{Key: user.ID + "/" + fID + "/pointer", Value: string(bFile)})
				} else {
					data.Content = []byte(base64.StdEncoding.EncodeToString(f))
					item, _ := json.Marshal(data)
					z.Set(api.Pair{Key: user.ID + "/" + fID, Value: string(item)})
				}
				fmt.Printf("\r\033[K")
				fmt.Printf("\033[F")
				fmt.Println("\n" + data.Name)
			} else {

				finalName := strings.Replace(p, home+"/", "", 1)
				f = []byte(base64.StdEncoding.EncodeToString(f))

				data := File{Content: f, Name: finalName, ID: user.ID, Group: user.PrvTok, Image: image}
				if options.Group != "" {
					data.Group = options.Group
				}
				item, _ := json.Marshal(data)
				fID, _ := util.RanString(6)
				if data.Image == true {
					d := util.ChunkString(string(data.Content), 2000)
					for i, c := range d {
						fmt.Printf("\r\033[K%d/%d", i+1, len(d))
						data = File{Content: []byte(c), Name: finalName, ID: user.ID, Group: user.PrvTok, Image: image}
						n, _ := json.Marshal(data)
						// z.Set(api.Pair{Key: user.ID + "/" + fID, Value: string(n)})
						z.Del(user.ID + "/" + fID)
						z.Dump(api.Pair{Key: user.ID + "/" + fID, Value: string(n)}, fID+".zi")
						// time.Sleep(1 * time.Second)
					}
					bFile, _ := json.Marshal(File{Name: finalName, Group: user.PrvTok, Image: image})
					z.Set(api.Pair{Key: user.ID + "/" + fID + "/pointer", Value: string(bFile)})
				} else {
					z.Set(api.Pair{Key: user.ID + "/" + fID, Value: string(item)})
				}
				fmt.Printf("\r\033[K")
				fmt.Printf("\033[F")
				fmt.Println("\n" + data.Name)
			}
		} else {
			fmt.Println("Cannot set file for user without write access.")
		}

	}
}
