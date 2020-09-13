package main

import (
	"github.com/zmb3/spotify"
	"github.com/deckarep/golang-set"
	"fmt"
)

type artistInfo struct {
	Name string
	ID string
	Url string
	IconUrl string
}

type ID spotify.ID

func getArtistInfo(artist spotify.FullArtist) (string, string, string, string) {
	name := artist.SimpleArtist.Name
	ID := artist.SimpleArtist.ID.String()
	url := artist.SimpleArtist.ExternalURLs["spotify"]
	iconUrl := artist.Images[0].URL
	return name, ID, url, iconUrl
}

func getFollowingArtists(client spotify.Client) ([]artistInfo, []string) {
	lastId := ""
	var artists []artistInfo
	var artistsId []string
	for {
		following, err := client.CurrentUsersFollowedArtistsOpt(50, lastId)
		if err != nil {
			fmt.Println(err)
		}
		for _, following := range following.Artists {
			var name, ID, url, iconUrl string
			name, ID, url, iconUrl = getArtistInfo(following)
			lastId = ID
			//_ := Set(name, ID, conn)
			artists = append(artists, artistInfo{name, ID, url, iconUrl})
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
		recommendIds = append(recommendIds, track.Artists[0].ID.String())
	}
	return recommendIds
}

func getRecommendedArtists(client spotify.Client, recommendIds, artistsId []string) []artistInfo {
	var recommendedArtists []artistInfo
	// TODO: use redis or set to figure out if id is in following
	artistsSet := mapset.NewSet()
	for _, id := range artistsId {
		artistsSet.Add(id)
	}
	
	for _, id := range recommendIds {
		if artistsSet.Contains(id) {
			continue
		}
		fullArtists, err := client.GetRelatedArtists(spotify.ID(id))
		if err != nil {
			fmt.Println("GetRelatedArtists", err)
		}
		for _, fullArtist := range fullArtists {
			name, ID, url, iconUrl := getArtistInfo(fullArtist)
			if artistsSet.Contains(ID) {
				continue
			}
			artistsSet.Add(ID)
			recommendedArtists = append(recommendedArtists, artistInfo{name, ID, url, iconUrl})	
		}
	}
	return recommendedArtists
}