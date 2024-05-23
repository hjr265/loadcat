// Copyright 2015 The Loadcat Authors. All rights reserved.

package nginx

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/coreos/go-systemd/v22/dbus"

	"github.com/hjr265/loadcat/cfg"
	"github.com/hjr265/loadcat/data"
	"github.com/hjr265/loadcat/feline"
)

var TplNginxConf = template.Must(template.New("").Parse(`
upstream {{.Balancer.Id.Hex}} {
	{{if eq .Balancer.Settings.Algorithm "least-connections"}}
		least_conn;
	{{else if eq .Balancer.Settings.Algorithm "source-ip"}}
		ip_hash;
	{{end}}

	{{if eq (len .Balancer.Servers) 0}}
		server  127.0.0.1:80 down;
	{{end}}
	{{range $srv := .Balancer.Servers}}
		server  {{$srv.Settings.Address}} weight={{$srv.Settings.Weight}} {{if eq $srv.Settings.Availability "available"}}{{else if eq $srv.Settings.Availability "backup"}}backup{{else if eq $srv.Settings.Availability "unavailable"}}down{{end}};
	{{end}}
}

server {
	{{if eq .Balancer.Settings.Protocol "http"}}
		listen  {{.Balancer.Settings.Port}};
	{{else if eq .Balancer.Settings.Protocol "https"}}
		listen  {{.Balancer.Settings.Port}} ssl;
	{{end}}
	server_name  {{.Balancer.Settings.Hostname}};

	{{if eq .Balancer.Settings.Protocol "https"}}

		ssl_certificate      /var/lib/loadcat/out/{{.Balancer.Id.Hex}}/server.crt;
		ssl_certificate_key  /var/lib/loadcat/out/{{.Balancer.Id.Hex}}/server.key;

	{{end}}

	{{if eq .Balancer.Settings.SSLOptions.SSLVerify "on"}}
		ssl_client_certificate /var/lib/loadcat/out/{{.Balancer.Id.Hex}}/ca.crt;
		ssl_verify_client on;
	{{end}}

	location / {
		proxy_set_header  Host $host;
		proxy_set_header  X-Real-IP $remote_addr;
		proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header  X-Forwarded-Proto $scheme;

		proxy_pass  http://{{.Balancer.Id.Hex}};

		proxy_http_version  1.1;

		proxy_set_header  Upgrade $http_upgrade;
		proxy_set_header  Connection 'upgrade';
	}
}
`))

type Nginx struct {
	sync.Mutex

	Systemd *dbus.Conn
	Cmd     *exec.Cmd
}

func (n *Nginx) Generate(dir string, bal *data.Balancer) error {
	n.Lock()
	defer n.Unlock()

	f, err := os.Create(filepath.Join(dir, "nginx.conf"))
	if err != nil {
		return err
	}
	err = TplNginxConf.Execute(f, struct {
		Dir      string
		Balancer *data.Balancer
	}{
		Dir:      dir,
		Balancer: bal,
	})
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	if bal.Settings.Protocol == "https" {
		err = ioutil.WriteFile(filepath.Join(dir, "server.crt"), bal.Settings.SSLOptions.Certificate, 0666)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(dir, "server.key"), bal.Settings.SSLOptions.PrivateKey, 0666)
		if err != nil {
			return err
		}
	}

	if bal.Settings.SSLOptions.SSLVerify == "on" {
		err = ioutil.WriteFile(filepath.Join(dir, "ca.crt"), bal.Settings.SSLVerifyClient.ClientCertificate, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Nginx) Start() error {
	n.Lock()
	defer n.Unlock()

	switch cfg.Current.Nginx.Mode {
	case "systemd":
		return nil

	case "exec":
		n.Cmd = exec.Command("nginx")
		n.Cmd.Stdout = os.Stdout
		n.Cmd.Stderr = os.Stderr
		return n.Cmd.Start()

	default:
		return errors.New("unknown Nginx mode")
	}
}

func (n *Nginx) Reload() error {
	n.Lock()
	defer n.Unlock()

	switch cfg.Current.Nginx.Mode {
	case "systemd":
		if n.Systemd == nil {
			c, err := dbus.NewSystemdConnection()
			if err != nil {
				return err
			}
			n.Systemd = c
		}

		ch := make(chan string)
		_, err := n.Systemd.ReloadUnit(cfg.Current.Nginx.Systemd.Service, "replace", ch)
		if err != nil {
			return err
		}
		<-ch

		return nil

	case "exec":
		cmd := exec.Command("nginx", "-s", "reload")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()

	default:
		return errors.New("unknown Nginx mode")
	}
}

func init() {
	feline.Register("nginx", &Nginx{})
}
