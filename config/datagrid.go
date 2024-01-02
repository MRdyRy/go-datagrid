package config

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
author : @Rudy Ryanto
this class provide connection for Redhat Datagrid

*/

// configuration field
type Configuration struct {
	Protocol string
	Host     string
	Port     string
	User     string
	Password string
}

type DatagridClient struct {
	cacheConfig Configuration
}

// object
type Cache interface{}

const (
	ERROR_GENERAL  = "Request Failed with status code %d"
	UNKNOWN_ERROR  = "Error caused : "
	ERROR_BASE_URL = "Error to call : "
	AUTH           = `Digest username="ryan", realm="default", nonce="AAAACgAAFZvTnKoMD9ZSk/V4rDneqcuGNewJJPQ4sZRhjhQ5gM3wVYLHiXI=", uri="/rest/v2/caches/b", algorithm=MD5, response="c5bd946f378df7db2e4da4494ec1bd9d", opaque="00000000000000000000000000000000", qop=auth, nc=00000004, cnonce="01182048b5a909be"`
)

// var (
// 	body []byte
// )

func NewDatagridClient(protocol, host, port, user, password string) DatagridClient {
	return DatagridClient{GetCacheConfig(protocol, host, port, user, password)}
}

func GetCacheConfig(protocol, host, port, user, password string) Configuration {
	cacheConfig := Configuration{
		Protocol: protocol,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}
	return cacheConfig
}

// private func, provide http client to call datagrid server
func buildHttpClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	return client
}

// private func to provide baseurl
func (dg *DatagridClient) baseUrl() string {
	return dg.cacheConfig.Protocol + "://" + dg.cacheConfig.Host + ":" + dg.cacheConfig.Port + "/rest/v2"
}

// get all cache keys
func (dg *DatagridClient) GetAllKeysFromCache(cacheName string) ([]string, error) {
	log.Println("get all key run")
	req, err := http.NewRequest("GET", dg.baseUrl()+"/caches/"+cacheName+"?action=entries&content-negotiation=true&metadata=true&limit=100", nil)
	if err != nil {
		return nil, err
	}

	client := buildHttpClient()
	req.Header.Set("Accept", "*/*")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", AUTH)

	log.Println(req)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var keys []string

	log.Println(res)
	err = json.NewDecoder(res.Body).Decode(&keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// check key if exist, return map, error
// param slice string
func (dg *DatagridClient) CheckExistKey(cacheName string, keys ...string) (map[string]bool, error) {
	log.Println("check exist run")
	data := make(map[string]bool)
	for _, key := range keys {
		log.Println(key)
		req, err := http.NewRequest("HEAD", dg.baseUrl()+"/caches/"+cacheName+"/"+key, nil)
		if err != nil {
			fmt.Println(ERROR_BASE_URL + dg.baseUrl())
		}
		client := buildHttpClient()
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(UNKNOWN_ERROR + err.Error())
		}

		data[key] = res.StatusCode == http.StatusOK
	}
	return data, nil
}

// get data from datagrid
// param cachename, key
// return any, error
func (dg *DatagridClient) GetDataFromCache(cacheName, key string) (Cache, error) {
	req, err := http.NewRequest("GET", dg.baseUrl()+"/caches/"+cacheName+"/"+key, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	client := buildHttpClient()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return res.Body, nil
}

func (dg *DatagridClient) AddToCache(cacheName, key string, value Cache) error {
	// body, _ = json.Marshal(GenerateTemplate())
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	exists, err := dg.CheckExistKey(cacheName)
	if err != nil {
		return err
	}

	client := buildHttpClient()
	log.Println(exists[key])
	if exists[key] {
		//do update data

		req, err := http.NewRequest("PUT", dg.baseUrl()+"/caches/"+cacheName+"/"+key, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusNoContent {
			return fmt.Errorf(fmt.Sprintf(ERROR_GENERAL, res.StatusCode))
		}

	} else {
		// create
		req, err := http.NewRequest("POST", dg.baseUrl()+"/caches/"+cacheName+"/"+key, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusNoContent {
			return fmt.Errorf(fmt.Sprintf(ERROR_GENERAL, res.StatusCode))
		}

	}

	return nil
}

// delete cache by key from datagrid
func (dg *DatagridClient) DeleteFromDG(cacheName, key string) error {
	req, err := http.NewRequest("DELETE", dg.baseUrl()+"/caches/"+cacheName+"/"+key, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	client := buildHttpClient()
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
