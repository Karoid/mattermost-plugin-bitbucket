package webhookpayload

// User는 Bitbucket 사용자 정보를 나타냅니다.
type User struct {
	Username string `json:"username"`
	Links    Links  `json:"links"`
}

// Links는 각종 링크 정보를 담고 있는 구조체입니다.
type Links struct {
	HTML Link `json:"html"`
}

// Link는 단일 링크 정보를 담고 있는 구조체입니다.
type Link struct {
	Href string `json:"href"`
}

// Push는 푸시 이벤트 정보를 담고 있는 구조체입니다.
type Push struct {
	Changes []Change `json:"changes"`
}

// Change는 변경 사항을 담고 있는 구조체입니다.
type Change struct {
	Forced  bool         `json:"forced"`
	New     ChangeBranch `json:"new"`
	Old     ChangeBranch `json:"old"`
	Links   Links        `json:"links"`
	Commits []Commit     `json:"commits"`
}

// ChangeBranch는 브랜치 정보를 담고 있는 구조체입니다.
type ChangeBranch struct {
	Name  string `json:"name"`
	Links Links  `json:"links"`
}

// Commit은 커밋 정보를 담고 있는 구조체입니다.
type Commit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Links   Links  `json:"links"`
	Author  Author `json:"author"`
}

// Author는 커밋 작성자 정보를 담고 있는 구조체입니다.
type Author struct {
	User User `json:"user"`
} 