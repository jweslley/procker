package procker

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestPrefixedWriter(t *testing.T) {
	b := &bytes.Buffer{}

	out := NewPrefixedWriter(b, "test> ")
	fmt.Fprint(out, "wintermute")
	assert(t, "test> wintermute", b.String())
}

func TestParseProcfile(t *testing.T) {
	r := strings.NewReader(`web:     python ranking/manage.py runserver
db:      postgres -D /usr/local/var/postgres
redis:   redis-server /usr/local/etc/redis.conf`)

	p, _ := ParseProcfile(r)

	if len(p) != 3 {
		t.Fatalf("has length %d; want %d", len(p), 3)
	}

	assert(t, "python ranking/manage.py runserver", p["web"])
	assert(t, "postgres -D /usr/local/var/postgres", p["db"])
	assert(t, "redis-server /usr/local/etc/redis.conf", p["redis"])
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

	assert(t, "python ranking/manage.py runserver", p["web"])
	assert(t, "postgres -D /usr/local/var/postgres", p["db"])
	assert(t, "redis-server /usr/local/etc/redis.conf", p["redis"])
}

func TestMustNotParseProcfileWithInvalidTypeNames(t *testing.T) {
	r := strings.NewReader(`web-9: bundle exec thin start
job: bundle exec rake jobs:work`)

	_, err := ParseProcfile(r)
	if err == nil {
		t.Fatalf("must not parse invalid lines")
	}
}

func TestParseEnv(t *testing.T) {
	r := strings.NewReader(`RAILS_ENV=production
QUEUE=system
PYTHONUNBUFFERED=True`)

	env, _ := ParseEnv(r)

	if len(env) != 3 {
		t.Fatalf("has length %d; want %d", len(env), 3)
	}

	assert(t, "RAILS_ENV=production", env[0])
	assert(t, "QUEUE=system", env[1])
	assert(t, "PYTHONUNBUFFERED=True", env[2])
}

func TestParseEnvExpandValues(t *testing.T) {
	path := os.Getenv("GOPATH")
	r := strings.NewReader(`RAILS_ENV=production
QUEUE=system
PYTHONUNBUFFERED=True
PROCKER_APP_ROOT=$GOPATH/xpto
PROCKER_APP_TMP=$PROCKER_APP_ROOT/tmp`)

	env, _ := ParseEnv(r)

	if len(env) != 5 {
		t.Fatalf("has length %d; want %d", len(env), 5)
	}

	assert(t, "RAILS_ENV=production", env[0])
	assert(t, "QUEUE=system", env[1])
	assert(t, "PYTHONUNBUFFERED=True", env[2])
	assert(t, fmt.Sprintf("PROCKER_APP_ROOT=%s/xpto", path), env[3])
	assert(t, fmt.Sprintf("PROCKER_APP_TMP=%s/xpto/tmp", path), env[4])
}

func TestShellCommand(t *testing.T) {
	cmd, err := NewShellCommand("echo 1 2 3")

	var env []string
	assert(t, nil, err)
	assert(t, "echo", cmd.Args[0])
	assert(t, []string{"1", "2", "3"}, cmd.Args[1:])
	assert(t, env, cmd.Env)
}

func TestShellCommandWithQuotedArgs(t *testing.T) {
	cmd, err := NewShellCommand("echo 'hello world' 42 'ping pong'")

	var env []string
	assert(t, nil, err)
	assert(t, "echo", cmd.Args[0])
	assert(t, []string{"hello world", "42", "ping pong"}, cmd.Args[1:])
	assert(t, env, cmd.Env)
}

func TestShellCommandWithCustomEnv(t *testing.T) {
	cmd, err := NewShellCommand("ANSWER=42 LOG=file.log echo 'hello world' 42 'ping pong'")

	assert(t, nil, err)
	assert(t, "echo", cmd.Args[0])
	assert(t, []string{"hello world", "42", "ping pong"}, cmd.Args[1:])
	assert(t, []string{"ANSWER=42", "LOG=file.log"}, cmd.Env)
}
