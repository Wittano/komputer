module github.com/wittano/komputer

go 1.22.2

toolchain go1.22.3

replace (
	github.com/wittano/komputer/bot => ./bot
	github.com/wittano/komputer/server => ./server
)

require (
	github.com/wittano/komputer/bot v0.0.0-20240610193926-c76ee77cc375
	github.com/wittano/komputer/server v0.0.0
	go.mongodb.org/mongo-driver v1.15.1
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/bwmarrin/dgvoice v0.0.0-20210225172318-caaac756e02e // indirect
	github.com/bwmarrin/discordgo v0.28.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240424034433-3c2c7870ae76 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240624140628-dc46fd24d27d // indirect
	layeh.com/gopus v0.0.0-20210501142526-1ee02d434e32 // indirect
)
