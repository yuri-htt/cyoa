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
		<section class="page">
			<h1>{{.Title}}</h1>
			{{range .Paragraphs}}
			<p>{{.}}</p>
			{{end}}
			<ul>
			{{range .Options}}
				<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
			{{end}}
			</ul>
		</section>
		<style>
		body {
			font-family: helvetica, arial;
		  }
		  h1 {
			text-align:center;
			position:relative;
		  }
		  .page {
			width: 80%;
			max-width: 500px;
			margin: auto;
			margin-top: 40px;
			margin-bottom: 40px;
			padding: 80px;
			background: #FFFCF6;
			border: 1px solid #eee;
			box-shadow: 0 10px 6px -6px #777;
		  }
		  ul {
			border-top: 1px dotted #ccc;
			padding: 10px 0 0 0;
			-webkit-padding-start: 0;
		  }
		  li {
			padding-top: 10px;
		  }
		  a,
		  a:visited {
			text-decoration: none;
			color: #6295b5;
		  }
		  a:active,
		  a:hover {
			color: #7792a2;
		  }
		  p {
			text-indent: 1em;
		  }
		</style>
    </body>
</html>
`

// 構造体のHandlerはhttpパッケージで定義されており、ServeHTTP(ResponseWriter, *Rewust)で構成される
// type Handler:
func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

// type handlerFunc
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
