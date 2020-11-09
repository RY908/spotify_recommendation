package main

import (
	"github.com/zmb3/spotify"
	"github.com/deckarep/golang-set"
	"fmt"
	"time"
)
/*
type ArtistInfo struct {
	Name 	string
	ID 		string
	Url 	string
	IconUrl string
}*/

type ID spotify.ID

func getArtistInfo(artist interface{}) (string, string, string, string) {
	var name, ID, url, iconUrl string
	switch t := artist.(type) {
	case spotify.FullArtist:
		name = t.SimpleArtist.Name
		ID = t.SimpleArtist.ID.String()
		url = t.SimpleArtist.ExternalURLs["spotify"]
		iconUrl = t.Images[0].URL
	case ArtistInfo:
		name = t.Name
		ID = t.Id
		url = t.Url
		iconUrl = t.IconUrl
	}
	return name, ID, url, iconUrl
}

func getArtistInfoFromTrack(client spotify.Client, id spotify.ID) (string, string, string, string) {
	artist, err := client.GetArtist(spotify.ID(id))
	if err != nil {
		fmt.Println("Couldn't get info")
		return "", "", "", ""
	}
	/*
	name, ID, url, iconUrl := getArtistInfo(&artist)
	fmt.Println(name, ID, url, iconUrl)
	*/
	name := artist.SimpleArtist.Name
	ID := artist.SimpleArtist.ID.String()
	url := artist.SimpleArtist.ExternalURLs["spotify"]
	iconUrl := artist.Images[0].URL
	return name, ID, url, iconUrl 
}

func getFollowingArtists(client spotify.Client) ([]ArtistInfo, []string) {
	lastId := ""
	var artists []ArtistInfo
	var artistsId []string
	for {
		following, err := client.CurrentUsersFollowedArtistsOpt(50, lastId)
		if err != nil {
			fmt.Println(err)
		}
		for _, following := range following.Artists {
			var name, ID, url, iconUrl string
			name, ID, url, iconUrl = getArtistInfo(following)
			err = insertArtist(dbmap, ID, name, url, iconUrl)
			//err = relationExists(dbmap, ID)
			lastId = ID
			//_ := Set(name, ID, conn)
			artists = append(artists, ArtistInfo{Id:ID, Name:name, Url:url, IconUrl:iconUrl, Timestamp:time.Now()})
			artistsId = append(artistsId, ID)
		}

		if len(following.Artists) < 50 {
			break
		}
	}
	return artists, artistsId
}

func makeSeeds(ids []string) spotify.Seeds {
	var seeds spotify.Seeds

	for _, id := range ids {
		seeds.Artists = append(seeds.Artists, spotify.ID(id))
	}

	return seeds
}

func getRecommendationId(client spotify.Client, ids []string) []string {
	seeds := makeSeeds(ids)
	attribute := spotify.NewTrackAttributes()
	recommendations, err := client.GetRecommendations(seeds, attribute, nil)
	if err != nil {
		fmt.Println(err)
	}
	tracks := recommendations.Tracks
	var recommendIds []string
	for _, track := range tracks {
		name, ID, url, iconUrl := getArtistInfoFromTrack(client, track.Artists[0].ID)
		_ = insertArtist(dbmap, ID, name, url, iconUrl)
		recommendIds = append(recommendIds, track.Artists[0].ID.String())
	}
	return recommendIds
}

func getRecommendedArtists(client spotify.Client, recommendIds, artistsId []string) []ArtistInfo {
	var recommendedArtists []ArtistInfo
	//var fullArtists interface{}
	// TODO: use redis or set to figure out if id is in following
	artistsSet := mapset.NewSet()
	for _, id := range artistsId {
		artistsSet.Add(id)
	}
	for _, id := range recommendIds {
		if artistsSet.Contains(id) {
			continue
		}
		err := relationExists(id)
		if err != nil {
			fullArtists, _ := client.GetRelatedArtists(spotify.ID(id))
			for _, fullArtist := range fullArtists {
				name, ID, url, iconUrl := getArtistInfo(fullArtist)
				insertArtist(dbmap, ID, name, url, iconUrl)
				insertRelation(dbmap, id, ID)
				if artistsSet.Contains(ID) {
					continue
				}
				artistsSet.Add(ID)
				recommendedArtists = append(recommendedArtists, ArtistInfo{Id:ID, Name:name, Url:url, IconUrl:iconUrl, Timestamp:time.Now()})	
			}
			//fullArtists = gfullArtists
		} else {
			fullArtists, _ := getRelation(dbmap, id)
			for _, fullArtist := range fullArtists {
				name, ID, url, iconUrl := getArtistInfo(fullArtist)
				if artistsSet.Contains(ID) {
					continue
				}
				artistsSet.Add(ID)
				recommendedArtists = append(recommendedArtists, ArtistInfo{Id:ID, Name:name, Url:url, IconUrl:iconUrl, Timestamp:time.Now()})	
			}
		}
	}
	return recommendedArtists
}