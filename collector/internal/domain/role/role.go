package role

//go:generate go tool go-enum --file=$GOFILE --noprefix

// ENUM(user, assistant)
type Role string
