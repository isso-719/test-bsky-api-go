package main

import (
	lexutil "github.com/bluesky-social/indigo/lex/util"

	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/util"
	"github.com/bluesky-social/indigo/xrpc"
	"time"
)

func main() {
	// クライアントを作成
	client := &xrpc.Client{
		Host: "https://bsky.social",
	}

	// 認証
	actor := "dummy.bsky.social" // ここに自分の ID を入れる
	pw := "dummy"                // ここに自分のパスワードを入れる
	input := &atproto.ServerCreateSession_Input{
		Identifier: actor,
		Password:   pw,
	}

	ctx := context.Background()
	output, err := atproto.ServerCreateSession(ctx, client, input) // ここで Bearer トークンが取得できる
	if err != nil {
		panic(err)
	}

	client.Auth = &xrpc.AuthInfo{
		AccessJwt:  output.AccessJwt,
		RefreshJwt: output.RefreshJwt,
		Handle:     output.Handle,
		Did:        output.Did,
	}

	// 投稿を作成
	inp := &atproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			&bsky.FeedPost{
				Text:      "Hello, World!",
				CreatedAt: time.Now().Format(util.ISO8601),
				Langs:     []string{"ja"},
			},
		},
	}

	res, err := atproto.RepoCreateRecord(ctx, client, inp)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Uri)

	// actor (自分) のフィードを新しい順に 10 件取得
	feed, err := bsky.FeedGetAuthorFeed(ctx, client, actor, "", "", 10)
	if err != nil {
		panic(err)
	}

	for _, f := range feed.Feed {
		fmt.Printf("%s\n", f.Post.Record.Val.(*bsky.FeedPost).Text)
	}
}
