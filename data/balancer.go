// Copyright 2015 The Loadcat Authors. All rights reserved.

package data

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"os"
	"github.com/boltdb/bolt"
	"gopkg.in/mgo.v2/bson"
)

type Balancer struct {
	Id       bson.ObjectId
	Label    string
	Settings BalancerSettings
}

type BalancerSettings struct {
	Hostname   string
	Port       int
	Protocol   Protocol
	Algorithm  Algorithm
	SSLOptions SSLOptions
	SSLVerifyClient SSLVerifyClient
}

type SSLOptions struct {
	CipherSuite CipherSuite
	Certificate []byte
	PrivateKey  []byte
	SSLVerify	SSLVerify
	DNSNames    []string
	Fingerprint []byte
}

type SSLVerifyClient struct {

	ClientCertificate   []byte
}

func ListBalancers() ([]Balancer, error) {
	bals := []Balancer{}
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("balancers"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			bal := Balancer{}
			err := bson.Unmarshal(v, &bal)
			if err != nil {
				return err
			}
			bals = append(bals, bal)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bals, nil
}

func GetBalancer(id bson.ObjectId) (*Balancer, error) {
	bal := &Balancer{}
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("balancers"))
		v := b.Get([]byte(id.Hex()))
		if v == nil {
			bal = nil
			return nil
		}
		err := bson.Unmarshal(v, bal)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bal, nil
}

func (b *Balancer) Delete() error {
    servers, err := b.Servers()
    if err != nil {
        return err
    }
    for _, server := range servers {
        err := server.Delete()
        if err != nil {
            return err
        }
    }

    dirPath := "/var/lib/loadcat/out/" + b.Id.Hex()
    err = os.RemoveAll(dirPath)
    if err != nil {
        return err
    }

    return DB.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte("balancers"))
        return bucket.Delete([]byte(b.Id.Hex()))
    })
}

func (l *Balancer) Servers() ([]Server, error) {
	return ListServersByBalancer(l)
}

func (l *Balancer) Put() error {
	if !l.Id.Valid() {
		l.Id = bson.NewObjectId()
	}
	if l.Label == "" {
		l.Label = "Unlabelled"
	}
	if l.Settings.Protocol == "https" {
		buf := []byte{}
		raw := l.Settings.SSLOptions.Certificate
		for {
			p, rest := pem.Decode(raw)
			raw = rest
			if p == nil {
				break
			}
			buf = append(buf, p.Bytes...)
		}
		certs, err := x509.ParseCertificates(buf)
		if err != nil {
			return err
		}
		l.Settings.SSLOptions.DNSNames = certs[0].DNSNames
		sum := sha1.Sum(certs[0].Raw)
		l.Settings.SSLOptions.Fingerprint = sum[:]
	} else {
		l.Settings.SSLOptions.CipherSuite = ""
		l.Settings.SSLOptions.Certificate = nil
		l.Settings.SSLOptions.PrivateKey = nil
		l.Settings.SSLOptions.DNSNames = nil
		l.Settings.SSLOptions.Fingerprint = nil
		l.Settings.SSLVerifyClient.ClientCertificate = nil
		l.Settings.SSLOptions.SSLVerify = "off"
	}
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("balancers"))
		p, err := bson.Marshal(l)
		if err != nil {
			return err
		}
		return b.Put([]byte(l.Id.Hex()), p)
	})
}
