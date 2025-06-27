package main

import (
	"encoding/json"
	"fmt"
	"self_utilities/fileio"
	"self_utilities/http/get"
)

const keypath = "key"
const mediaServer = "http://10.2.0.6:8096"

func main() {
	keyParam, err := getKey("key")
	if err != nil {
		panic(err)
	}
	adminUid, err := getAdminUser(mediaServer + "/Users?" + keyParam)
	if err != nil {
		panic(err)
	}

	libraries, err := getLibraries(mediaServer + "/Users/" + adminUid + "/Items/?" + keyParam)
	if err != nil {
		panic(err)
	}

	libraryUrl := mediaServer + "/Users/" + adminUid + "/Items/"
	for _, v := range libraries {
		fmt.Println(v)
		itemUrl := libraryUrl + "?Recursive=true&ParentId=" + v.Id + "&fields=Path&" + keyParam
		fmt.Println(itemUrl)
	}
}

func getLibraries(url string) (l []libraryStruct, e error) {
	libraries, e := get.HttpGetURL(url)
	var items itemsList
	json.Unmarshal(libraries, &items)
	l = items.Items
	return
}

func getAdminUser(url string) (a string, e error) {
	data, e := get.HttpGetURL(url)
	if e != nil {
		return
	}

	var userList []jellyUser
	json.Unmarshal(data, &userList)
	for _, v := range userList {
		if v.Policy["IsAdministrator"].(bool) {
			a = v.Id
			break
		}
	}

	return
}

func getKey(s string) (k string, e error) {
	k, e = fileio.StrFromPath(s)
	k = "apikey=" + k
	return
}

type jellyUser struct {
	Name     string                 `json:"Name"`
	ServerId string                 `json:"ServerId"`
	Id       string                 `json:"Id"`
	Policy   map[string]interface{} `json:"Policy"`
}

type libraryStruct struct {
	Name           string `json:"Name"`
	Id             string `json:"Id"`
	Type           string `json:"Type"`
	CollectionType string `json:"CollectionType"`
	LocatationType string `json:"LocationType"`
}
type itemsList struct {
	Items            []libraryStruct `json:"Items"`
	TotalRecordCount int             `json:"TotalRecordCount"`
	StartIndex       int             `json:"StartIndex"`
}
