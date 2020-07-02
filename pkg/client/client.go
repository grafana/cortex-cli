package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNoConfig         = errors.New("No config exists for this user")
	ErrResourceNotFound = errors.New("requested resource not found")
)

// Config is used to configure a Ruler Client
type Config struct {
	Key         string `yaml:"key"`
	Address     string `yaml:"address"`
	ID          string `yaml:"id"`
	TLScaFile   string `yaml:"tls_ca"`
	TLScertFile string `yaml:"tls_cert"`
	TLSkeyFile  string `yaml:"tls_key"`
}

// CortexClient is used to get and load rules into a cortex ruler
type CortexClient struct {
	key      string
	id       string
	endpoint *url.URL
	client   http.Client
}

// New returns a new Client
func New(cfg Config) (*CortexClient, error) {
	endpoint, err := url.Parse(cfg.Address)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"address": cfg.Address,
		"id":      cfg.ID,
	}).Debugln("New ruler client created")

	//if TLS flags, setup http client
	if cfg.TLScaFile != "" && cfg.TLScertFile != "" && cfg.TLSkeyFile != "" {

		// Load client cert
		cert, err := tls.LoadX509KeyPair(cfg.TLScertFile, cfg.TLSkeyFile)
		if err != nil {
			log.WithFields(log.Fields{
				"tls-cert": cfg.TLScertFile,
				"tls-key":  cfg.TLSkeyFile,
			}).Fatal(err)
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(cfg.TLScaFile)
		if err != nil {
			log.WithFields(log.Fields{
				"tls-ca": cfg.TLScaFile,
			}).Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		return &CortexClient{
			key:      cfg.Key,
			id:       cfg.ID,
			endpoint: endpoint,
			client:   http.Client{Transport: transport},
		}, nil
	}
	return &CortexClient{
		key:      cfg.Key,
		id:       cfg.ID,
		endpoint: endpoint,
		client:   http.Client{},
	}, nil
}

func (r *CortexClient) doRequest(path, method string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, r.endpoint.String()+path, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	if r.key != "" {
		req.SetBasicAuth(r.id, r.key)
	}

	req.Header.Add("X-Scope-OrgID", r.id)

	log.WithFields(log.Fields{
		"url":    req.URL.String(),
		"method": req.Method,
	}).Debugln("sending request to cortex api")

	resp, err := r.client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url":    req.URL.String(),
			"method": req.Method,
			"error":  err.Error(),
		}).Errorln("error during request to cortex api")
		return nil, err
	}

	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// checkResponse checks the API response for errors
func checkResponse(r *http.Response) error {
	log.WithFields(log.Fields{
		"status": r.Status,
	}).Debugln("checking response")
	if 200 <= r.StatusCode && r.StatusCode <= 299 {
		return nil
	}

	var msg string
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg = fmt.Sprintf("unable to decode body, %s", err.Error())
	} else {
		msg = fmt.Sprintf("request failed with response body %v", string(data))
	}

	if r.StatusCode == http.StatusNotFound {
		log.WithFields(log.Fields{
			"status": r.Status,
			"msg":    msg,
		}).Debugln("resource not found")
		return ErrResourceNotFound
	}

	log.WithFields(log.Fields{
		"status": r.Status,
		"msg":    msg,
	}).Errorln("requests failed")

	return errors.New("failed request to the cortex api")
}
