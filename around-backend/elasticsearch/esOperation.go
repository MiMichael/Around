package elasticsearch

import (
	"context"
	"fmt"
	"reflect"

	consts "../constant"
	"../post"
	"gopkg.in/olivere/elastic.v6"
)

func CreateIndexIfNotExist() {
	client, err := elastic.NewClient(elastic.SetURL(consts.ES_URL), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	exists, err := client.IndexExists(consts.POST_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		mapping := `{
			"mappings":{
				"post":{
					"properties":{
						"location":{
							"type":"geo_point"
						}
					}
				}
			}
		}`
		_, err := client.CreateIndex(consts.POST_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

}
func SaveToES(post *post.Post, id string) error {
	//set sniff 调用变成自己名字的函数
	client, err := elastic.NewClient(elastic.SetURL(consts.ES_URL), elastic.SetSniff(false))
	if err != nil {
		return err
	}
	_, err = client.Index().
		Index(consts.POST_INDEX).
		Type(consts.POST_TYPE).
		Id(id).
		BodyJson(post).
		Refresh("wait_for").
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Post is saved to index: %s\n", post.Message)
	return nil

}
func ReadFromES(lat, lon float64, ran string) ([]post.Post, error) {
	client, err := elastic.NewClient(elastic.SetURL(consts.ES_URL), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	query := elastic.NewGeoDistanceQuery("location")
	query = query.Distance(ran).Lat(lat).Lon(lon)
	searchResult, err := client.Search().
		Index(consts.POST_INDEX).
		Query(query).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	var ptyp post.Post
	var posts []post.Post
	for _, item := range searchResult.Each(reflect.TypeOf(ptyp)) {
		if p, ok := item.(post.Post); ok {
			posts = append(posts, p)
			// fmt.Printf("Post by %s: %s\n", p.User, p.Message)
		}
	}
	return posts, nil
}
