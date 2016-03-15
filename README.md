# AcidBath

## Name 
> This project is titled after Twisted Individual's "Acid Bath" 
> Check out his other drum and bass tracks here:

> [Spotify] (https://play.spotify.com/artist/1yJkKFH5pD56kG4ryEWk2H)

> [SoundCloud] (https://soundcloud.com/twisted-individual) 
 
## Description
> In short, this project is infrastructure code to connect with brokerages (currently only TD Ameritrade) to stream quotes and act on those streaming quotes. There's a whole lot of real time event processing between go routines.
> This code base is very much in an alpha phase. It has tons of places where it needs to improve, but it is in a functional state, since my main objective has been to get something working and crank out features. Currently the only brokerage firm supported is TD Ameritrade; however, the code does try to abstract brokers. The API isn't solid yet and is very subject to change, as I only created it suit my immediate needs.


> Application is basically architected in these packages

    
                   Web Layer                   Interface
      +---------------+------------+       +-----------------------+---------------------------------------+
      |               |            | +---> |    /broker/generic    |   /broker/tdapi                       |
      |  /web         |            |       |                       |   /broker/<implementation>            |
      |               | /handlers  |       +---------------------------------------------------------------+
      |               |            |       |   /eventproc/generic  |   /eventproc/reference                |
      |               |            | +---> |                       |   /eventproc/<your implementation>    |
      +---------------+------------+       +-----------------------+---------------------------------------+

                             +                                +                        +
                             |                                |                        |
                             |                                |                        |
                             v                                v                        v
                                              Common Stuff
                      +----------------------------+----------------------------------------------------------+
                      |                       /lib + misc reusable stuff                                      |
                      |                       /dm - data model components                                     |
                      +---------------------------+-----------------------------------------------------------+



> At a very high level, the application works like the following diagram.

  
                                                       Acid Bath
    +---------------------------------------------------------------------------------------------------------------------------------+   +----------+
    | All the HTTP Handlers execute                                                                                                   |   |          |
    | in their own go routine                                               +------------------------------------------------------------>+TD        |
    |                                                                       |                                                         |   |          |
    |                  +----------+ +-----------> +-----------------------+ |                               This code runs in own     |   |          |
    |                  | /login   |               | tdapi.Login()         | +                               go routine until          |   |          |
    |                  +----------+ <-----------+ +-----------------------+                                 end signal arrives        |   |          |
    |                                                                                                      +---------------------+    |   |          |
    |                  +----------+ +-----------> +-----------------------+                                |                     |    |   |          |
    |                  | /reqData |               | tdapi.Stream()        | +---->  go streamData() +----> |  parseData(stream)  | <------+TD Stream |
    |                  +----------+ <-----------+ +-----------------------+ <----+                         |                     |    |   |          |
    |                                                                                                +---> +---------+-----------+    |   |          |
    |                  +----------+ +-----------> +-----------------------+   1                      |               |                |   |          |
    |                  | /logout  |               | tdapi.Logout()        | +---->  endChan <+ true  +               | execute for    |   |          |
    |                  +----------+ <-----------+ +-----------------------+ |                                        v each "packet"  |   |          |
    |                                                                       |                                +-------+---------+      |   |          |
    |                                                                       |                                | event processor |      |   |          |
    |                                                                       | 2                              +-----------------+      |   |          |
    |                                                                       |                                                         |   |          |
    |                                                                       +------------------------------------------------------------>+TD        |
    |                                                                                                                                 |   |          |
    |                                                                                                                                 |   |          |
    +---------------------------------------------------------------------------------------------------------------------------------+   +----------+
 
> The API aims to be safe for concurrent access. So you can call it from HTTP handlers (or anywhere you like).
> After you request to stream quotes, there will be a brand new go routine created for the sole purpose of parsing the binary data stream coming from
> TD. As each "packet" of data is parsed, the parsing function executes a callback function to notify all the interested parties that some data
> has arrived. This is where people can stick their own functionality in, and do something interesting like display the new quote, or even do something
> wild like send an order to the broker
 
## Install/Running 
> * You need a TDA account and access to their API.
> * Get the code 
    
    go get github.com/marklaczynski/acidbath
> * Generate some certificates for https access
    
    go run /usr/local/go/src/crypto/tls/generate_cert.go --host localhost
    cp *.pem github.com/marklaczynski/acidbath/web/certificates
> * Populate the source id in the td config file, which is provided by TDA

    github.com/marklaczynski/acidbath/broker/tdapi/config/tdconfig.json 
    {
	"sourceid": "<sourceid here>",
	"version": "1"
    }

> * Run the applicaiton

    go run acidbath.go 
> or
    
    go build acidbath.go
    ./acidbath
> * point browser to 

    https://localhost:1111

## Thanks 
> This project uses the following external projects

> 1. [Gorilla MUX] (https://github.com/gorilla/mux)

> 2. [Go Statistics] (https://github.com/grd/statistics)

> 3. [AngularJS] (https://angularjs.org/)

> 4. [Bootstrap] (http://getbootstrap.com/)

