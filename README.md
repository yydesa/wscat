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

For example, setting up forwarding in nginx
is like this, in a server clause (with `wscatd`
listening on port 7070):
```
  location /webssh {
    proxy_pass http://127.0.0.1:7070;
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
