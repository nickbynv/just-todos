package todo

type List struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UsersLists struct {
	Id     int
	UserId int
	ListId int
}

type Item struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type ListsItems struct {
	Id     int
	ListId int
	ItemId int
}
