package cyoa

import (
	"encoding/json"
	"fmt"
	"io"
)

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
