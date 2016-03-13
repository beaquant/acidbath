# AcidBath

## Name 
> This project is titled after Twisted Individual's "Acid Bath" 
> Check out his other drum and bass track here:

> [Spotify] (https://play.spotify.com/artist/1yJkKFH5pD56kG4ryEWk2H)

> [SoundCloud] (https://soundcloud.com/twisted-individual) 
 
## Description
> In short, this project is infrastructure code to connect with brokerages (currently only TD Ameritrade) to stream quotes and act on those streaming quotes.
> This code base is very much in an alpha phase. It has tons of places where it needs to improve, but it is in a functional state, since my main objective has been to get something working and crank out features. Currently the only brokerage firm supported is TD Ameritrade; however, the code does try to abstract brokers. The API isn't solid yet and is very subject to change, as I only created it suit my immediate needs.

## Install/Running 
> * You need a TDA account and access to their API.
> * Get the code 
    
    go get github.com/marklaczynski/acidbath
> * Generate some certificates for https access
    
    go run /usr/local/go/src/crypto/tls/generate_cert.go --host localhost
    cp *.pem github.com/marklaczynski/acidbath/web/certificates
> * Populate the source id in the td config file, which is provided by TDA

    github.com/marklaczynski/acidbath/broker/tdapi/config/tdconfig.json 
> * Run the applicaiton

    go run acidbath.go 
> or
    
    go build acidbath.go
    ./acidbath
> * point browser to 

    https://localhost:1111

## Thanks 
> This project uses the following github projects

> 1. Gorilla MUX "github.com/gorilla/mux"

> 2. Go Statistics https://github.com/grd/statistics

