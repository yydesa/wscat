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
