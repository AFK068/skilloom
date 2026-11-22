package role

//go:generate go run github.com/abice/go-enum@v0.9.5 --file=$GOFILE --noprefix

// ENUM(user, assistant)
type Role string
