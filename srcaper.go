package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/irfanfajru/rssagg/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticket := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticket.C {
		feeds, err := db.GetNextFeedToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds:", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("error marking as fetched", err)
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("cannot parse time ", item.PubDate, err)
		}
		log.Println("Found post", item.Title, "on feed", feed.Name)
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Description: item.Description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("cannot create post with error ", err)
		}

	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
