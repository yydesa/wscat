# wscat

A netcat for websockets. Basically, take bytestreams
and stuff them into binary websocket frames;
then re-join at the other end and forward.

The clients behave just like netcat, the server
is a websocket server that forwards incoming
connections to configured TCP addresses and ports.

One obvious use is to start the server somewhere,
configure your web server to forward the proper
websocket connection to the websocket server, and
then put 
```
Host yourhost
  ProxyCommand wscat --addr yourhost.com
```
into `.ssh/config`. Result: You can ssh into
your machine, via a websocket connection on
port 443. SSH within SSL websockets, and it
works from everywhere https:// works; no need
to have port 22 open in any firewall on the way.

The default path `wscat` uses is `/webssh`,
because of the scenario above.

For example, setting up forwarding in nginx
is like this, in a server clause (with `wscatd`
listening on port 8080 which is its default):
```
  location /webssh {
    proxy_pass http://127.0.0.1:8080;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "Upgrade";

    proxy_read_timeout 600s;
  }
```
You will also need a keepalive in the ssh client.

In putty, you need to enter the full `wscat` command
at 'Telnet command' and select 'Local', both in
the Proxy panel.

## wscat

By default `wscat` connects to `localhost:8080`
(which you probably don't want).

`--addr host:port` changes the host and port
number to connect to (the `:port` can be omitted
if it is the normal 443 or 80).

`--path asdf` changes the URL path to connect to,
with the given path it would connect to
`wss://host/asdf`. Default is `webssh`.

`--suffix suf` allows to add a suffix to the
base URL path, e.g. `--suffix /sub` will cause
a connection to `wss://host/websssh/sub`. This
is useful only when you want to keep the default
base path, otherwise it is easier all put into
`--path`.

`--proxy host:port` allows to specify a proxy
to use, according to `http.ProxyURL` and `url.Parse`
of golang, e.g. `--proxy http://localhost:3128`
(socks(5) seems not to be supported?)

`--insecure` causes wscat to use `ws://` instead
of `wss://`, losing the SSL security properties.
(With SSL you get verification that you talk to
the correct host, and can more lightheartedly
accept the host key at the first connection.)

## wscatd

The server accepts as arguments a list of endpoints
and places to forward them. The default is
`webssh=localhost:22`, meaning that `ws://host/webssh`
will be forwarded to the local ssh daemon. If arguments
are given the default does not apply and you must specify
it if you want it, like
```
wscatd webssh=localhost:22 webssh/other=other:22
```
By default `wscatd` listens on port 8080, and only
on `localhost`, not being reachable from the outside.
This can be changed with the option `--addr` which takes
a single argument with the hostname/address and port number
to listen on, e.g. `--addr myhost:1234`.

`wscatd` itself does not provide SSL termination;
you need to do that in the reverse-proxying webserver.
