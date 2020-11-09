package main

import (
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-gorp/gorp"
	"time"
	"fmt"
)
/*
type artistInfo struct {
	Name 	string
	ID 		string
	Url 	string
	IconUrl string
}*/
type ArtistInfo struct {
    Id			string 		`db:"ID"`
	Name 		string 		`db:"name"`
	Url 		string 		`db:"url"`
	IconUrl		string 		`db:iconUrl`
	Timestamp	time.Time	`db:"timestamp"`
}

type Relate struct {
	RelateId	int64 	`db:relateId`
	Name1		string 	`db:name1`
	Name2		string 	`db:name2`
}

func artistExists(dbmap *gorp.DbMap, id string) (time.Time, error) {
	obj, err := dbmap.Get(ArtistInfo{}, id)
	if err != nil {
		fmt.Println(err)
		return time.Time{}, err
	}
	if obj == nil {
		return time.Time{}, fmt.Errorf("no data")
	}
	return obj.(*ArtistInfo).Timestamp, nil
}

func insertArtist(dbmap *gorp.DbMap, id, name, url, iconUrl string) error {
	timestamp, err := artistExists(dbmap, id)
	if err != nil {
		err = dbmap.Insert(&ArtistInfo{Id:id, Name:name, Url:url, IconUrl:iconUrl, Timestamp:time.Now()})
		//fmt.Println("47")
	} else {
		if (time.Since(timestamp)).Hours() > 730.001 {
			_, err = dbmap.Update(&ArtistInfo{id, name, url, iconUrl, time.Now()})
		} else {
			err = nil
		}
	}
	//err := dbmap.Insert(&ArtistInfo{id, name, url, iconUrl, time.Now()})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func updateArtist(dbmap *gorp.DbMap, artist ArtistInfo) error {
	_, err := dbmap.Update(&artist)
	if err != nil {
		return err
	}
	return nil
}

func relationExists(id string) error {
	var relation []Relate
	_, err := dbmap.Select(&relation, "select * from Relate where name1 = ? limit 1", id)
	if err != nil {
		return err
	} else {
		if len(relation) > 0 {
			return nil
		} else {
			return fmt.Errorf("no data")
		}
	}
}

func insertRelation(dbmap *gorp.DbMap, id1, id2 string) error {
	err := dbmap.Insert(&Relate{0, id1, id2})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getRelation(dbmap *gorp.DbMap, id string) ([]ArtistInfo, error) {
	var artists []ArtistInfo
	cmd := "select Artist.ID, Artist.Name, Artist.url, Artist.iconUrl from Artist inner join Relate on Artist.Id = Relate.Name2 where Relate.Name1 = ?"
	_, err := dbmap.Select(&artists, cmd, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return artists, nil
}