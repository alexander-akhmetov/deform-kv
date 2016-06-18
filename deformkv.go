package deformkv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// http://deformio.github.io/docs/quickstart/
// deform project create -d '{"_id": "project_name", "name": "My key-value storage"}'
// http://deformio.github.io/docs/quickstart/#creating-a-token

type document struct {
	Id    string `json:"_id"`
	Value string `json:"value"`
}

type DeformError struct {
	Message string
	Key     string
}

func (err DeformError) Error() string {
	return fmt.Sprintf("Deform error. Message: '%s'. Key: '%s'", err.Message, err.Key)
}

type Deform struct {
	token      string
	project    string
	collection string
}

func NewClient(project string, collection string, token string) Deform {
	return Deform{project: project, token: token, collection: collection}
}

func (deform Deform) Get(key string) (string, error) {
	response, err := deform.getRequest(deform.getDocumentUrl(key))
	defer response.Body.Close()

	if err != nil {
		return "", DeformError{Message: err.Error(), Key: key}
	}

	responseString, err := deform.readResponseBody(response)
	if err != nil {
		return "", DeformError{Message: err.Error(), Key: key}
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return "", DeformError{Message: responseString, Key: key}
	}
	return deform.getValueFromResponseData(responseString, key)
}

func (deform Deform) getValueFromResponseData(responseString string, key string) (string, error) {
	var responseDocument document
	err := json.Unmarshal([]byte(responseString), &responseDocument)
	if err != nil {
		return "", DeformError{Message: err.Error(), Key: key}
	}

	return responseDocument.Value, nil
}

func (deform Deform) getRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	deform.setRequestHeaders(req)
	return client.Do(req)
}

func (deform Deform) readResponseBody(response *http.Response) (string, error) {
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(responseData), nil
}

func (deform Deform) Set(key string, value string) error {
	_, err := deform.unsafeRequest(
		"PUT",
		deform.getDocumentUrl(key),
		document{Id: key, Value: value},
	)
	return err
}

func (deform Deform) unsafeRequest(method string, url string, doc document) (*http.Response, error) {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	data := bytes.NewBuffer([]byte(jsonData))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	deform.setRequestHeaders(req)

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (deform Deform) setRequestHeaders(request *http.Request) {
	request.Header.Set("Authorization", fmt.Sprintf("Token %s", deform.token))
	request.Header.Set("Content-Type", "application/json")
}

func (deform Deform) getDocumentUrl(documentId string) string {
	return fmt.Sprintf("%sdocuments/%s/", deform.getCollectionUrl(), documentId)
}

func (deform Deform) getCollectionUrl() string {
	return fmt.Sprintf("%s/collections/%s/", deform.getApiEndpoint(), deform.collection)
}

func (deform Deform) getApiEndpoint() string {
	return fmt.Sprintf("https://%s.deform.io/api", deform.project)
}
