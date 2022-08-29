package couch

import (
	"context"
	"fmt"
	_ "github.com/go-kivik/couchdb/v4"
	"github.com/go-kivik/kivik/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	"weeb_bot/internal/storage"
)

type Client struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	client   *kivik.Client
	Enabled  bool
}

func New() *Client {
	return &Client{
		Username: os.Getenv("COUCH_USER"),
		Password: os.Getenv("COUCH_PASS"),
		Host:     os.Getenv("COUCH_HOST"),
		Port:     os.Getenv("COUCH_PORT"),
		Database: os.Getenv("COUCH_DATABASE"),
		Enabled:  true,
	}
}

func (c *Client) Connection() (*kivik.Client, error) {
	if !c.Enabled {
		return nil, storage.ErrNoConnection
	}
	if c.client == nil {
		uri := fmt.Sprintf("http://%s:%s@%s:%s", c.Username, c.Password, c.Host, c.Port)
		client, err := kivik.New("couch", uri)
		if err != nil {
			c.Enabled = false
			log.Errorf("Failed to connect to couchdb: %v", err)
			return nil, storage.ErrNoConnection
		}
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		_, err = client.Ping(ctx)
		if err != nil {
			c.Enabled = false
			log.Errorf("Failed to verify connection to couchdb: %v", err)
			return nil, storage.ErrNoConnection
		}
		c.client = client
	}
	return c.client, nil
}

func (c *Client) DB() (*kivik.DB, error) {
	conn, err := c.Connection()
	if err != nil {
		return nil, err
	}
	db := conn.DB(c.Database)
	err = db.Err()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Client) Init(ctx context.Context) error {
	db, err := c.DB()
	if err != nil {
		return err
	}
	docs := db.DesignDocs(ctx)
	if err != nil {
		return err
	}
	defer docs.Close()
	// TODO Actually check if the design docs are the same
	if docs.Next() {
		return nil
	}

	_, err = db.Put(ctx, "_design/matches", map[string]interface{}{
		"_id": "_design/matches",
		"views": map[string]interface{}{
			"player_average": map[string]interface{}{
				"map": `function(doc) {
    if (!!doc.info && !!doc.info.participants && Array.isArray(doc.info.participants)) {
        const emObject = function(prefix, arg) {
        	const keys = Object.keys(arg)
            for (const k of keys) {
                const pre = [...prefix,k];
                if (arg[k] == undefined){
                    continue;
                }
                if (typeof arg[k] == 'object' && !Array.isArray(arg[k])) {
                    emObject(pre, arg[k]);
                } else if (typeof arg[k] == 'number') {
                    emit(pre, arg[k]);
                }
            }
        }
        for (const d of doc.info.participants) {
            emObject([d.puuid, doc.info.queueId], d);
        }
    }
}`,
				"reduce": "_stats",
			},
		},
	})

	return err
}
