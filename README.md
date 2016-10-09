puber
=====

Is a queryable store for public keys. Its intended use is for automatic
generation of SSH `authorized_keys` files.

It currently supports [redis](https://redis.io) and memory back ends
for storing data.

Access to add and remove public keys is controlled via whitelists and
yubikey authentication.

If a host is in the whitelist, it can send add / remove requests
assuming the yubikey authentication is successful.

## Usage

### Starting

First you must acquire an API Key from Yubico via their [web
portal](https://upgrade.yubico.com/getapikey/). It requires a valid
email address and a yubikey.

Once you have the Client ID and API Key, you can pass them to puber on
startup:

```
puber -ycid CLIENTID -yskey APIKEY
```

This will start puber in memory store mode with a whitelist that only
allows requests to `/add` and `/rm` from `localhost`. More hosts can
be specified with the `-wl` flag.

### Adding keys

### Querying keys

Keys can be pulled per user or all at once for all users. This uses
`/user/$username` and `/all` respectively.

If you are running in debug mode (-debug flag), three fake keys will
be added so that you can test querying:

To get all keys for user `debug`:
```
curl http://localhost:8081/user/debug
```

To get all keys stored in puber:
```
curl http://localhost:8081/all
```

### Querying keys via Ansible

Obviously if you are using puber for dishing out your public keys via
Ansible, you don't want to use the temporary memory store. Using
something like redis is highly recommended!

Ansible can query a HTTP endpoint to get a users public keys. Here is
the example from the [authorized key module](https://docs.ansible.com/ansible/authorized_key_module.html):

```
# Using github url as key source
- authorized_key: user=charlie key=https://github.com/charlie.keys
```

puber would look like this:

```
- authorized_key: user=charlie key=https://puberurl/charlie
```

More users can be specified by simply specifying them via a ','
seperated list:

```
- authorized_keys: user=root key=https://puberurl/charlie,fred,jane
```

## TODO

* [ ] Make yubikey auth optional.
* [ ] Allow for arbitrary yubi client id and API keys to be used.
* [ ] Add a BoltDB backend.
