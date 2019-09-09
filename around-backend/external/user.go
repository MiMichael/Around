package external

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	consts "../constant"
	"gopkg.in/olivere/elastic.v6"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int64  `json:"age"`
	Gender   string `json:"gender"`
}

func CheckUser(username, password string) error {
	client, err := elastic.NewClient(elastic.SetURL(consts.ES_URL), elastic.SetSniff(false))
	if err != nil {
		return err
	}

	query := elastic.NewTermQuery("username", username)
	searchResult, err := client.Search().
		Index(consts.USER_INDEX).
		Query(query).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return err
	}
	var utyp User
	for _, item := range searchResult.Each(reflect.TypeOf(utyp)) {
		if u, ok := item.(User); ok {
			if username == u.Username && password == u.Password {
				fmt.Printf("Login as %s\n", username)
				return nil
			}
		}
	}
	return errors.New("wrong username or password")

}

func AddUser(user User) error {
	client, err := elastic.NewClient(elastic.SetURL(consts.ES_URL), elastic.SetSniff(false))
	if err != nil {
		return err
	}

	query := elastic.NewTermQuery("username", user.Username)

	searchResult, err := client.Search().
		Index(consts.USER_INDEX).
		Query(query).
		Pretty(true).
		Do(context.Background())

	if err != nil {
		return err
	}

	if searchResult.TotalHits() > 0 {
		return errors.New("User already exists")
	}

	_, err = client.Index().
		Index(consts.USER_INDEX).
		Type(consts.USER_TYPE).
		Id(user.Username).
		BodyJson(user).
		Refresh("wait_for").
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("User is added: %s\n", user.Username)
	return nil

}
