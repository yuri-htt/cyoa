package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Choose Your Own Adventure</title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
        <p>{{.}}</p>
        {{end}}
        <ul>
        {{range .Options}}
            <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
    </body>
</html>
`

// 構造体のHandlerはhttpパッケージで定義されており、ServeHTTP(ResponseWriter, *Rewust)で構成される
func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	// Slice strings ex) "/intro" => "intro"
	path = path[1:]

	//                   ["intro"]
	if chapter, ok := h.s[path]; ok {
		// Execute: 解析されたテンプレートを指定されたデータオブジェクトに適用し、出力をwrに書き込む
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

// io.Reader: バイト列を読むためのReadメソッドを提供する
func JsonStory(r io.Reader) (Story, error) {
	// デコード: 別ファイルを読み取ってjsonに変換し直す
	d := json.NewDecoder(r)
	var story Story
	// 構造体Chapterを要素とする配列の構造体のポインタ(アドレスのみ)を渡してデコード
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

// structタグ: 構造体のフィールドにタグをつけることができる
// type Demo struct {
// 	Name string `json:"name,omitempty"`
// 	Age  int    `json:"age,omitempty"`
// }
