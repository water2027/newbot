package model

// Post [...]
type Post struct {
    PostID        int       `json:"PostID"`
    UserID        int       `json:"UserID"`
    UserName      string    `json:"UserName"`
    UserScore     int       `json:"UserScore"`
    UserTelephone string    `json:"UserTelephone"`
    UserAvatar    string    `json:"UserAvatar"`
    UserIdentity  string    `json:"UserIdentity"`
    Title         string    `json:"Title"`
    Content       string    `json:"Content"`
    Like          int       `json:"Like"`
    Comment       int       `json:"Comment"`
    Browse        int       `json:"Browse"`
    Heat          int       `json:"Heat"`
    PostTime      string    `json:"PostTime"`
    IsSaved       bool      `json:"IsSaved"`
    IsLiked       bool      `json:"IsLiked"`
    Photos        string    `json:"Photos"`
    Tag           string    `json:"Tag"`
}
