package wordpress_to_md

import (
	"fmt"
	"log"
	"time"

	"github.com/tmc/wxr"
)

// Post represents a blog post
type Post struct {
	Filename string
	Contents []byte
}

// RSSToFiles returns Posts from a wxr.RSS object
func RSSToFiles(rss *wxr.RSS) ([]Post, error) {

	result := []Post{}
	for _, item := range rss.Channel.Items {
		filename := time.Time(item.PubDate).Format("2006-01-02-") + item.Name + ".md"
		body := "---\n"

		body += "title: '" + item.Title + "'\n"
		body += "layout: " + item.Type + "\n"
		body += "type: " + item.Type + "\n"
		body += "author: " + item.Author + "\n"

		tags, categories := []string{}, []string{}
		for _, c := range item.Categories {
			if c.Domain == "post_tag" {
				tags = append(tags, c.Slug)
			} else {
				// assume the rest are categories
				categories = append(categories, c.Slug)
			}
		}
		body += fmt.Sprintf("categories: %s\n", categories)
		body += fmt.Sprintf("tags: %s\n", tags)

		body += "permalink: " + item.Link + "\n"
		body += "---\n"
		body += fmt.Sprintf("%+v\n", item.Content)
		longest := ""
		for _, c := range item.Content {
			if len(c) > len(longest) {
				longest = c
			}
		}
		body += longest
		log.Println(body)

		result = append(result, Post{
			filename,
			[]byte(body),
		})
	}

	return result, nil
}
