package main

import (
        "encoding/json"
        "fmt"
        "net/http"
        "io/ioutil"
        "bytes"
        "strconv"
        )

type HTTPError struct {
  Code int
  Message string
}

func (e HTTPError) Error() string {
  return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

type Client struct {
  url string
  authToken string
}

func (c *Client) Set(key string, value interface {},  optional ...int64) error {
  var ex int64 = 0

  url := fmt.Sprintf("%s/set/%s", c.url, key)

  if len(optional) > 0 {
    ex = optional[0]
  }

  if ex != 0 {
    url += fmt.Sprintf("?ex=%s", strconv.FormatInt(ex, 10))
  }

  jsonValue, err := json.Marshal(value)

  if err != nil {return err}
  resp, err := c.makeRequest("POST", url, jsonValue)

  if err != nil { return  err}
  defer resp.Body.Close()

  if err := handleAuth(resp); err == nil {
    return handleUnprocessableEntity(resp)
  }

  return err
}


func (c *Client) Get(key string) (interface {}, bool) {
  var value interface {}

  url := fmt.Sprintf("%s/get/%s", c.url, key)

  resp, err := c.makeRequest("GET", url, nil)

  if err != nil { panic(err) }

  defer resp.Body.Close()

  if err := handleNotFound(resp); err != nil {
    return nil, false
  }

  content, err := ioutil.ReadAll(resp.Body)

  if (err != nil) { panic(err) }

  json.Unmarshal(content, &value)

  return value, true
}

func (c *Client) Update(key, value interface {}) error {
  url := fmt.Sprintf("%s/update/%s", c.url, key)

  jsonValue, err := json.Marshal(value)

  if err != nil {return err}
  resp, err := c.makeRequest("PUT", url, jsonValue)

  if err != nil { return  err}
  defer resp.Body.Close()

  if err := handleAuth(resp); err == nil {
    return handleUnprocessableEntity(resp)
  }
  return err
}

func (c *Client) Delete(key string) {
  url := fmt.Sprintf("%s/delete/%s", c.url, key)

  resp, err := c.makeRequest("DELETE", url, nil)

  if err != nil { panic(err) }

  resp.Body.Close()
}

func (c *Client) Keys(pattern string) {

}

func (c *Client) makeRequest(requestType string, url string, json []byte) (resp *http.Response, err error) {
  req, err := http.NewRequest(requestType, url, bytes.NewBuffer(json))
  req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.authToken))
  req.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  return client.Do(req)
}

func handleAuth(resp *http.Response) error {
  if (resp.StatusCode == 401) {
    return HTTPError{
      401,
      "Unauthorized",
    }
  }
  return nil
}

func handleUnprocessableEntity(resp *http.Response) error {
  if (resp.StatusCode == 422) {
    content, _ := ioutil.ReadAll(resp.Body)
    return HTTPError{
      422,
      string(content),
    }
  }
  return nil
}

func handleNotFound(resp *http.Response) error {
  if (resp.StatusCode == 404) {
    return HTTPError{
      422,
      "Not found",
    }
  }
  return nil
}
