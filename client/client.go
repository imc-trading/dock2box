package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Client structure.
type Client struct {
	URL       string
	Host      HostResource
	Interface InterfaceResource
	Image     ImageResource
	Tag       TagResource
	Site      SiteResource
	Tenant    TenantResource
	Subnet    SubnetResource
	Debug     bool
}

// New client.
func New(url string) *Client {
	c := Client{
		URL: url,
	}
	c.Host.Client = &c
	c.Interface.Client = &c
	c.Image.Client = &c
	c.Tag.Client = &c
	c.Site.Client = &c
	c.Tenant.Client = &c
	c.Subnet.Client = &c
	return &c
}

// SetDebug enable debug.
func (c *Client) SetDebug() {
	c.Debug = true
}

// Info log
func (c *Client) Info(msg string) {
	if c.Debug {
		log.Print(msg)
	}
}

// Infof log
func (c *Client) Infof(fmt string, args ...interface{}) {
	if c.Debug {
		log.Printf(fmt, args...)
	}
}

// Fatal log and exit
func (c *Client) Fatal(msg string) {
	log.Fatal(msg)
}

// Fatalf log and exit
func (c *Client) Fatalf(fmt string, args ...interface{}) {
	log.Fatalf(fmt, args...)
}

// All resources.
func (c *Client) All(endp string) ([]byte, error) {
	url := c.URL + endp
	c.Infof("url: %s", url)

	resp, err := http.Get(url + "?envelope=false")
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	cont, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return cont, nil
}

// Get resource.
func (c *Client) Get(endp string, id string) ([]byte, error) {
	url := c.URL + endp + "/" + id
	c.Infof("url: %s", url)

	resp, err := http.Get(url + "?envelope=false")
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("Get %s: failed with status code %d", url, resp.StatusCode)
	}

	defer resp.Body.Close()
	cont, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return cont, nil
}

// Query for resources.
func (c *Client) Query(endp string, cond map[string]string) ([]byte, error) {
	url := c.URL + endp + "?envelope=false&"
	for k, v := range cond {
		url = url + k + "=" + v + "&"
	}
	url = strings.TrimRight(url, "&")
	c.Infof("url: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("Get %s: failed with status code %d", url, resp.StatusCode)
	}

	defer resp.Body.Close()
	cont, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return cont, nil
}

// Create resource.
func (c *Client) Create(endp string, s interface{}) ([]byte, error) {
	url := c.URL + endp
	c.Infof("header: application/json, method: POST, url: %s", url)

	b, _ := json.MarshalIndent(&s, "", "  ")
	fmt.Printf("Payload:\n%s\n", string(b))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	fmt.Println("Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body:\n%s\n", string(body))

	return body, nil
}

// Update resource.
func (c *Client) Update(endp string, id string, s interface{}) ([]byte, error) {
	url := c.URL + endp + "/" + id
	c.Infof("header: application/json, method: PUT, url: %s", url)

	b, _ := json.MarshalIndent(&s, "", "  ")
	fmt.Printf("Payload:\n%s\n", string(b))
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	fmt.Println("Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body:\n%s\n", string(body))

	return body, nil
}

// Delete resource.
func (c *Client) Delete(endp string, id string) ([]byte, error) {
	url := c.URL + endp + "/" + id
	c.Infof("url: %s", url)

	req, err := http.NewRequest("DELETE", url+"?envelope=false", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("Delete %s: failed with status code %d", url, resp.StatusCode)
	}

	defer resp.Body.Close()
	cont, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return cont, nil
}
