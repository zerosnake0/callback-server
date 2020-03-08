package client

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"callback-server/pkg/log"
	"callback-server/pkg/server"
)

func init() {
	log.InitLog()
}

func TestClient(t *testing.T) {
	must := require.New(t)

	port := 9000 + rand.Intn(1000)

	go func() {
		gin.SetMode(gin.ReleaseMode)
		server.Run(port)
	}()

	time.Sleep(time.Second)

	t.Run("call first", func(t *testing.T) {
		ctx := context.TODO()
		cli := NewClient("http://localhost:"+strconv.Itoa(port),
			WithHttpClient(http.DefaultClient))

		id := "call first"
		ct := "test/ct1"
		body := []byte("test/data1")

		err := cli.Call(ctx, id, &Data{ContentType: ct, Body: body})
		must.NoError(err)

		recv, err := cli.Wait(ctx, id)
		must.NoError(err)

		must.Equal(ct, recv.ContentType)
		must.Equal(body, recv.Body)
	})

	t.Run("wait first", func(t *testing.T) {
		ctx := context.TODO()
		cli := NewClient("http://localhost:"+strconv.Itoa(port),
			WithHttpClient(http.DefaultClient))

		id := "wait first"
		ct := "test/ct2"
		body := []byte("test/data2")
		done := make(chan struct{})

		go func() {
			recv, err := cli.Wait(ctx, id)
			must.NoError(err)

			must.Equal(ct, recv.ContentType)
			must.Equal(body, recv.Body)

			close(done)
		}()

		time.Sleep(time.Second)

		err := cli.Call(ctx, id, &Data{ContentType: ct, Body: body})
		must.NoError(err)

		<-done
	})
}
