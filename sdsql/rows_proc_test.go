package sdsql

type testUser struct {
	Id    string
	Name  string
	Posts []*testPost
}

type testPost struct {
	Id      string
	Title   string
	Uid     string
	Content string
	User    *testUser
}

type testData struct {
	Users []*testUser
	Posts []*testPost
}

func testNewData() *testData {
	return &testData{
		Users: []*testUser{
			{Id: "uid1", Name: "Alice"},
			{Id: "uid2", Name: "Bob"},
		},
		Posts: []*testPost{
			{Id: "postId1", Title: "Title1", Uid: "uid1", Content: "Content1"},
			{Id: "postId2", Title: "Title2", Uid: "uid1", Content: "Content2"},
			{Id: "postId3", Title: "Title3", Uid: "uid2", Content: "Content3"},
		},
	}
}
