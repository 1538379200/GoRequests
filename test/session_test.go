package test

import (
	"fmt"
	"github.com/1538379200/GoRequests/session"
	"testing"
)

type TestSuite struct {
	session.Session
	BaseUrl string
}

func (ts *TestSuite) TestPost(t *testing.T) {
	data := map[string]interface{}{
		"app_id":  "cd86552a-4e63-44a0-8528-95c7248aba38",
		"app_sec": "cd86552a-4e63-44a0-8528-95c7248aba38",
	}
	accessToken := ts.Post(ts.BaseUrl+"user/login/app", data).Find("data.access_token").Str
	//accessToken := ts.Find(res, "name.access_token").Str
	ts.AddHeader("Authorization", "AppToken "+accessToken)
	t.Log(accessToken)
}

func (ts *TestSuite) TestPostSearch(t *testing.T) {
	data := map[string]interface{}{
		"search":        "customer",
		"field_filters": []string{},
		"page_size":     10,
		"page_num":      0,
		"order_fields":  []string{"ctime"},
		"order_method":  "ORDER_DESC",
		"template_id":   "",
	}
	res := ts.Post(ts.BaseUrl+"sys/user-manager/user/search", data).JsonFormat()
	fmt.Println(res)
	t.Log(res)
}

func (ts *TestSuite) TestGet(t *testing.T) {
	//res, _ := ts.Get("https://www.baidu.com", nil)
	res := ts.Get(ts.BaseUrl+"order/management/preset?prod_type=OrderProdType_3_GfIP", nil).Json()
	t.Log(res)
}

func TestRunner(t *testing.T) {
	s := session.New()
	s.AddHeader("Content-Type", "application/json")
	ts := TestSuite{
		*s,
		"http://zcloud.skynetcloud.com/api/",
	}
	ts.TestPost(t)
	ts.TestPostSearch(t)
}
