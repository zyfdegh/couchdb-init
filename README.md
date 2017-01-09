# couchdb-init
Create design document in CouchDB, insert user and acl documents. Query data by view.

# Prerequisites
Docker (1.12+ recommended)
Go (1.7+ recommended)

# Build
```sh
mkdir -p $GOPATH/src/github.com/zyfdegh
cd $GOPATH/src/github.com/zyfdegh
git clone https://github.com/zyfdegh/couchdb-init
cd couchdb-init
./build.sh
```

# Run
Start CouchDB in container
```sh
docker run -p 5984:5984 -d --restart=always --name couch couchdb:1.6.1
```
This start a couchdb, with port 5984 mapped to host machine.

Execute program, this will insert design documents, user documents, acl documents into CouchDB,
after that, it queries one acl and one user from CouchDB.
```sh
$./couchdb-init
2017/01/06 11:01:46 &{Welcome 8bf3c8073da62d44f5affd3ca66a6ca1 {1.6.1 The Apache Software Foundation} 1.6.1}
2017/01/06 11:01:46 >>> Creating db acl
2017/01/06 11:01:47 >>> Adding design doc _design/acl
2017/01/06 11:01:47 &{true bf00d35e1214161a22e78c2b00000076 1-6722270fc0288e04673d17366318cbba}
2017/01/06 11:01:47 &{true bf00d35e1214161a22e78c2b00000882 1-85463bdc682b280c68e8917b9a7d6557}
2017/01/06 11:01:47 &{true bf00d35e1214161a22e78c2b0000180e 1-1dfa68f27e94527b02a779c857e0a7d7}
2017/01/06 11:01:47 &{true bf00d35e1214161a22e78c2b00002052 1-f12db59f814d38802ab9cb7fa49bda3f}
2017/01/06 11:01:47 &{true bf00d35e1214161a22e78c2b000024e4 1-20fcc4e701e2eaf1c07509c700b66ff1}
2017/01/06 11:01:47 &{true bf00d35e1214161a22e78c2b00002be7 1-e4b308a910931a0c82d49d2b3ca70dc4}
2017/01/06 11:01:47 >>> Query db acl
2017/01/06 11:01:47 &{6 [] 6 0}
2017/01/06 11:01:47 >>> Creating db user
2017/01/06 11:01:47 >>> Adding design doc _design/user
2017/01/06 11:01:47 >>> Adding users
2017/01/06 11:01:47 &{true bf00d35e1214161a22e78c2b000033be 1-4e9b44eafd0fbc6bd4a04c0f73868627}
2017/01/06 11:01:47 &{true bf00d35e1214161a22e78c2b00003445 1-86b93aecd0c17b328c66120eec550f87}
2017/01/06 11:01:47 &{true bf00d35e1214161a22e78c2b000041ae 1-c7c60f7d0d2c7c7f3895333d0580e23a}
2017/01/06 11:01:47 >>> Query db user
2017/01/06 11:01:47 &{0 [{bf00d35e1214161a22e78c2b000033be admin@email.com map[_id:bf00d35e1214161a22e78c2b000033be _rev:1-4e9b44eafd0fbc6bd4a04c0f73868627 email:admin@email.com password:secret123 username:admin] map[]}] 3 0}
2017/01/06 11:01:47 map[email:admin@email.com password:secret123 username:admin _id:bf00d35e1214161a22e78c2b000033be _rev:1-4e9b44eafd0fbc6bd4a04c0f73868627]
```

Open CouchDB dashboard in web browser to view and edit databases and documents. 
```sh
http://127.0.0.1:5984/utils
```

# Full workspace
To see how couchdb-init work with docker registry, docker_auth and couchdb to build a registry with authz & authn.
Redirect to [zyfdegh/dockerauth-workspace](https://github.com/zyfdegh/dockerauth-workspace).

The workspace contains configs, certs and a compose file, you can deploy all containers with a single command.

