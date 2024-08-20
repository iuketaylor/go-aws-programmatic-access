# go-ibl-wormhole
cli tool for gaining programatic access to BBC (ibl) AWS accounts 

# Usgae
* Add .env file with the following variables 

```
WORMHOLE_PATH=<path to ibl-scripts wormhole script>
AWS_CREDENTIALS_PATH=<path to your aws_credentials dotfile>
DEV_ID=<dev account ID>
TOOLING_ID=<tooling account ID>
PROD_ID=<prod account ID>
ICAT_ID=<icat account ID>
INNOVATION_ID=<innovation account ID>

```

* run `go.build` to build the binary 
* add the binary to your PATH env variable $HOME/<path to binary>
* optional add an alias eg: `alias wh='go-ibl-wormhole'`
* run `wh` then select which ever option you want
