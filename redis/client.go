package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	raw  *redis.Client
	pipe redis.Pipeliner

	//properties
	addr         string
	username     string
	password     string
	db           int
	protocol     int
	poolSize     int
	minIdleConns int
	readTimeout  time.Duration
	writeTimeout time.Duration
	enableTls    bool
	caCertPath   string
	servername   string
}

// Option pattern use
type Option func(*Client)

func WithAddr(addr string) Option {
	return func(c *Client) {
		c.addr = addr
	}
}
func WithUsername(username string) Option {
	return func(c *Client) {
		c.username = username
	}
}
func WithPassword(password string) Option {
	return func(c *Client) {
		c.password = password
	}
}
func WithDB(db int) Option {
	return func(c *Client) {
		c.db = db
	}
}
func WithProtocol(protocol int) Option {
	return func(c *Client) {
		c.protocol = protocol
	}
}
func WithPoolSize(poolSize int) Option {
	return func(c *Client) {
		c.poolSize = poolSize
	}
}
func WithMinIdleConns(minIdleConns int) Option {
	return func(c *Client) {
		c.minIdleConns = minIdleConns
	}
}
func WithReadTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.readTimeout = timeout
	}
}
func WithWriteTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.writeTimeout = timeout
	}
}
func WithEnableTls(enable bool) Option {
	return func(c *Client) {
		c.enableTls = enable
	}
}
func WithCACertPath(caCertPath string) Option {
	return func(c *Client) {
		c.caCertPath = caCertPath
	}
}
func WithServerName(servername string) Option {
	return func(c *Client) {
		c.servername = servername
	}
}

func New(ctx context.Context, opts ...Option) (*Client, error) {
	client := &Client{
		addr:         "127.0.0.1:6379",
		db:           0,
		protocol:     2,
		poolSize:     10,
		minIdleConns: 3,
		readTimeout:  2 * time.Second,
		writeTimeout: 2 * time.Second,
		enableTls:    false,
	}

	for _, opt := range opts {
		opt(client)
	}

	rsOpts := &redis.Options{
		Addr:         client.addr,
		DB:           client.db,
		Protocol:     client.protocol,
		PoolSize:     client.poolSize,
		MinIdleConns: client.minIdleConns,
		ReadTimeout:  client.readTimeout,
		WriteTimeout: client.writeTimeout,
	}
	if client.username != "" {
		rsOpts.Username = client.username
	}
	if client.password != "" {
		rsOpts.Password = client.password
	}
	if client.enableTls {
		if client.caCertPath == "" || client.servername == "" {
			return nil, fmt.Errorf("enableTls requires caCertPath %s and servername %s", client.caCertPath, client.servername)
		}
		tlsCfg, err := buildTLSConfig(client.caCertPath, client.servername)
		if err != nil {
			return nil, err
		}
		rsOpts.TLSConfig = tlsCfg
	}
	rsClient := redis.NewClient(rsOpts)
	if err := rsClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	client.raw = rsClient
	client.pipe = client.raw.Pipeline()
	return client, nil
}

func buildTLSConfig(caCertPath string, servername string) (*tls.Config, error) {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate to pool")
	}
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    caPool,
		ServerName: servername,
	}, nil
}

func (c *Client) Raw() *redis.Client {
	return c.raw
}

func (c *Client) Pipe() redis.Pipeliner {
	return c.pipe
}

func (c *Client) Ping(ctx context.Context) string {
	return c.raw.Ping(ctx).Val()
}

func (c *Client) Close() error {
	if c.raw != nil {
		return c.raw.Close()
	}
	return nil
}
