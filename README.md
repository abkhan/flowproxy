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
