package procker

import (
	"strings"
	"testing"
)

func assertCommand(t *testing.T, p *Process, command string) {
	if p == nil {
		t.Fatalf("process not found")
	}
	assert(t, command, p.Command)
}

func TestParseProcfile(t *testing.T) {
	r := strings.NewReader(`web:     python ranking/manage.py runserver
db:      postgres -D /usr/local/var/postgres
redis:   redis-server /usr/local/etc/redis.conf`)

	p, _ := ParseProcfile(r)

	if len(p) != 3 {
		t.Fatalf("has length %d; want %d", len(p), 3)
	}

	assertCommand(t, p["web"], "python ranking/manage.py runserver")
	assertCommand(t, p["db"], "postgres -D /usr/local/var/postgres")
	assertCommand(t, p["redis"], "redis-server /usr/local/etc/redis.conf")
}

func TestParseProcfileIgnoreBlankLines(t *testing.T) {
	r := strings.NewReader(`web:     python ranking/manage.py runserver

db:      postgres -D /usr/local/var/postgres
     

redis:   redis-server /usr/local/etc/redis.conf

`)

	p, _ := ParseProcfile(r)

	if len(p) != 3 {
		t.Fatalf("has length %d; want %d", len(p), 3)
	}

	assertCommand(t, p["web"], "python ranking/manage.py runserver")
	assertCommand(t, p["db"], "postgres -D /usr/local/var/postgres")
	assertCommand(t, p["redis"], "redis-server /usr/local/etc/redis.conf")
}

func TestMustNotParseProcfileWithInvalidTypeNames(t *testing.T) {
	r := strings.NewReader(`web-9: bundle exec thin start
job: bundle exec rake jobs:work`)

	_, err := ParseProcfile(r)
	if err == nil {
		t.Fatalf("must not parse invalid lines")
	}
}
