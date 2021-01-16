package myapp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexPathHandler(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	mux := NewHTTPHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, err := ioutil.ReadAll(res.Body) // res.Body 자체는 버퍼값이므로 이걸 ioutil.ReadAll이 읽어줘야 한다.
	assert.NoError(err)
	assert.Equal("Hello World", string(data))
}

func TestBarPathHandler_WithoutName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	mux := NewHTTPHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body) // res.Body 자체는 버퍼값이므로 이걸 ioutil.ReadAll이 읽어줘야 한다.
	assert.Equal("Hello Bar!", string(data))
}

func TestBarPathHandler_WithName(t *testing.T) {
	assert := assert.New(t)
	testName := "greatlaboratory"

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name="+testName, nil)

	mux := NewHTTPHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body) // res.Body 자체는 버퍼값이므로 이걸 ioutil.ReadAll이 읽어줘야 한다.
	assert.Equal("Hello "+testName+"!", string(data))
}

func TestFooHandler_WithoutJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil) // request body에 nil을 넣어서 테스트

	mux := NewHTTPHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code) // request body가 없으면 bad request인 400코드를 반환해야 한다.
}

func TestFooHandler_WithJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(`{
		"first_name": "MG", 
		"last_name": "Kim", 
		"email": "wowo0201@gmaill.com"
	}`))

	mux := NewHTTPHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusCreated, res.Code)

	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)
	assert.Nil(err)
	assert.Equal("MG", user.FirstName)
	assert.Equal("Kim", user.LastName)
	assert.Equal("wowo0201@gmaill.com", user.Email)
}

func TestUploadHandler(t *testing.T) {
	assert := assert.New(t)
	path := "C:/Users/GreatLaboratory/Desktop/lalaland.png"
	file, err := os.Open(path)
	assert.NoError(err)
	defer file.Close()

	os.RemoveAll("./uploads")

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path))
	assert.NoError(err)

	io.Copy(multi, file)
	writer.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf)
	req.Header.Set("content-type", writer.FormDataContentType())

	mux := NewHTTPHandler()
	mux.ServeHTTP(res, req)
	assert.Equal(http.StatusOK, res.Code)

	uploadFilePath := "./uploads/" + filepath.Base(path)
	_, err = os.Stat(uploadFilePath)
	assert.NoError(err)

	uploadFile, _ := os.Open(uploadFilePath)
	originFile, _ := os.Open(path)
	defer uploadFile.Close()
	defer originFile.Close()

	uploadData := []byte{}
	originData := []byte{}
	uploadFile.Read(uploadData)
	uploadFile.Read(originData)
	assert.Equal(originData, uploadData)
}

func TestGetAllUserInfo(t *testing.T) {
	assert := assert.New(t)

	testServer := httptest.NewServer(NewHTTPHandler())
	defer testServer.Close()

	// user 정보가 하나도 없을 경우
	res, err := http.Get(testServer.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusNotFound, res.StatusCode)
	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)
	assert.Equal(string(data), "No Users")

	// 첫번째 user instance 생성
	res, err = http.Post(testServer.URL+"/user", "application/json", strings.NewReader(`{
		"first_name": "MG", 
		"last_name": "Kim",  
		"email": "wowo0201@gmaill.com"
		}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	// 두번째 user instance 생성
	res, err = http.Post(testServer.URL+"/user", "application/json", strings.NewReader(`{
		"first_name": "YM", 
		"last_name": "Kim",  
		"email": "rhawl97@gmaill.com"
	}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	// user 정보가 2개 있을 경우
	res, err = http.Get(testServer.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	users := []*User{}
	err = json.NewDecoder(res.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))
}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	testServer := httptest.NewServer(NewHTTPHandler())
	defer testServer.Close()

	res, err := http.Get(testServer.URL + "/user/89")
	assert.NoError(err)
	assert.Equal(http.StatusNotFound, res.StatusCode)
	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)
	assert.Contains(string(data), "No User Id:89")

	res, err = http.Get(testServer.URL + "/user/56")
	assert.NoError(err)
	assert.Equal(http.StatusNotFound, res.StatusCode)
	data, err = ioutil.ReadAll(res.Body)
	assert.NoError(err)
	assert.Contains(string(data), "No User Id:56")
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	testServer := httptest.NewServer(NewHTTPHandler())
	defer testServer.Close()

	res, err := http.Post(testServer.URL+"/user", "application/json", strings.NewReader(`{
		"first_name": "MG", 
		"last_name": "Kim",  
		"email": "wowo0201@gmaill.com"
	}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	// 클라에서 post method로 보내는 request body에 있는 user json
	user := new(User)
	err = json.NewDecoder(res.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	id := user.ID
	resp, err := http.Get(testServer.URL + "/user/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// get method로 받아진 response body에 있는 user json
	user2 := new(User)
	err = json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID)
	assert.Equal(user.FirstName, user2.FirstName)
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	testServer := httptest.NewServer(NewHTTPHandler())
	defer testServer.Close()

	// 삭제할 user instance를 생성
	res, err := http.Post(testServer.URL+"/user", "application/json", strings.NewReader(`{
		"first_name": "MG", 
		"last_name": "Kim",  
		"email": "wowo0201@gmaill.com"
	}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	// id가 1인 user instance는 존재하므로 삭제 성공
	req, _ := http.NewRequest("DELETE", testServer.URL+"/user/1", nil)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "Deleted User Id:1")

	// id가 2인 user instance는 존재하지 않으므로 404
	req, _ = http.NewRequest("DELETE", testServer.URL+"/user/2", nil)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusNotFound, res.StatusCode)
	data, _ = ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "No User Id:2")
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)

	testServer := httptest.NewServer(NewHTTPHandler())
	defer testServer.Close()

	// 수정할 user instance를 생성
	res, err := http.Post(testServer.URL+"/user", "application/json", strings.NewReader(`{
		"first_name": "MG", 
		"last_name": "Kim",  
		"email": "wowo0201@gmaill.com"
	}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	// 생성된 user instance를 응답으로 받아서 user struct에 디코드해서 저장
	user := new(User)
	err = json.NewDecoder(res.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	// id가 1인 user instance는 존재하므로 수정 성공
	req, _ := http.NewRequest("PUT", testServer.URL+"/user", strings.NewReader(`{
		"id": 1,
		"email": "updated wowo0201@gmaill.com"
	}`))
	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	updatedUser := new(User)
	err = json.NewDecoder(res.Body).Decode(updatedUser)
	assert.NoError(err)
	assert.Equal(updatedUser.ID, user.ID)
	assert.Equal(updatedUser.FirstName, user.FirstName)
	assert.Equal(updatedUser.LastName, user.LastName)
	assert.Equal("updated wowo0201@gmaill.com", updatedUser.Email)

	// id가 2인 user instance는 존재하지 않으므로 수정 실패
	req, _ = http.NewRequest("PUT", testServer.URL+"/user", strings.NewReader(`{
		"id": 2,
		"first_name": "updated MG", 
		"last_name": "updated Kim",  
		"email": "updated wowo0201@gmaill.com"
	}`))
	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusNotFound, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "No User Id:2")
}
