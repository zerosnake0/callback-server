package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	port  int
	cbMap sync.Map
)

type side int

const (
	sideCaller side = 0
	sideWaiter side = 1
)

type response struct {
	err  error
	ct   string
	data []byte
}

type cbItem struct {
	id   string
	side side
	ch   chan response
}

func newCbItem(id string, side side) *cbItem {
	return &cbItem{
		id:   id,
		side: side,
		ch:   make(chan response, 1),
	}
}

func init() {
	log.Logger = zerolog.New(zerolog.NewConsoleWriter()).With().
		Timestamp().Caller().Logger()
}

type errDuplicatedID string

func (e errDuplicatedID) Error() string {
	return fmt.Sprintf("%q already exists", string(e))
}

func run(port int) {
	engine := gin.New()

	engine.GET("/clear", func(c *gin.Context) {
		cbMap = sync.Map{}
		c.String(http.StatusOK, "OK")
	})

	engine.POST("/cb/*id", func(c *gin.Context) {
		id := c.Param("id")
		logger := log.With().Str("side", "caller").Str("id", id).Logger()

		item := newCbItem(id, sideCaller)
		if o, loaded := cbMap.LoadOrStore(id, item); loaded {
			item = o.(*cbItem)
			if item.side == sideCaller {
				logger.Error().Msg("already created by another caller")
				c.AbortWithError(http.StatusBadRequest, errDuplicatedID(id))
				return
			}
			logger.Debug().Msg("already created")
			cbMap.Delete(id) // delete the item because waiter must have it already
		}

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error().Err(err).Msg("unable to read body")
			// DO NOT RETURN
		}

		ctx := c.Request.Context()
		select {
		case <-ctx.Done():
			logger.Error().Err(ctx.Err()).Msg("context done")
		case item.ch <- response{
			err:  err,
			ct:   c.ContentType(),
			data: data,
		}:
			c.Status(http.StatusOK)
		}
	})

	engine.GET("/cb/*id", func(c *gin.Context) {
		id := c.Param("id")
		logger := log.With().Str("side", "waiter").Str("id", id).Logger()

		item := newCbItem(id, sideWaiter)
		if o, loaded := cbMap.LoadOrStore(id, item); loaded {
			item = o.(*cbItem)
			if item.side == sideWaiter {
				logger.Error().Msg("already created by another waiter")
				c.AbortWithError(http.StatusBadRequest, errDuplicatedID(id))
				return
			}
			logger.Debug().Msg("already created by caller")
			cbMap.Delete(id) // delete the item because caller must have it already
		}

		ctx := c.Request.Context()
		select {
		case <-ctx.Done():
			log.Error().Err(ctx.Err()).Msg("context done")
		case resp := <-item.ch:
			if resp.err != nil {
				c.AbortWithError(http.StatusBadRequest, resp.err)
			} else {
				c.Data(http.StatusOK, resp.ct, resp.data)
			}
		}
	})

	engine.Run(":" + strconv.Itoa(port))
}

func main() {
	flag.IntVar(&port, "port", 80, "listen port")
	flag.Parse()
	run(port)
}
