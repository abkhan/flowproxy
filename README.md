# Flow Proxy

## Requirement

Targets (where to send the traffic to)

    Learning load balancer


Default for all new sources

Old Sources stick to old destinations.


## Flow

Config has all the learned targets,
and must have a current default.

All traffic coming on learned targets just follow the same.

New sender, goes to default.

Once default server is at capacity, update the default.

Old traffic goes to the same old.

New goes to new default.

## Running FlowProxy

### Command Line Options
  -d string
    	IP or FQDN of default destination
  -dlist string
    	comma separated IP or FQDN of destination(s)
  -f string
    	config file name (default "./config.yml")
  -init
    	create new config, must provide default destination
  -p int
    	UDP listen port (default 9995)

### Start with Init
    $ flowproxy -init -f ./config.yml -d 192.168.2.3

As the flowproxy would learn, it would update the config file.

### Update the default proxy while keeping other stuff the same
    $flowproxy -f ./config.yml -d 192.168.2.4

Flowproxy would start and read all previously learned routes from config file, and use the new default destination for all new sources.
