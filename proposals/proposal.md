# DNS Corrector

The DNS corrector is a tool for bare metal homelabs that run off a home connection
where the ISP is prone to swtiching the IP address frequently enough that having
a mechanism to update IPs on that cycle is necessary.

## Process

Every time period DNS Corrector looks up (a) the network address provided by the ISP and
(b) the network address for a specified address.

If they match, we sleep. If we fail to resolve either (a) or (b) we mark a failure but
still sleep.

If they don't match, then we trigger a DNS update.

The DNS update will look at the provided list of hostnames, and using the provided keys will
update our DNS settings. We keep trying with backoff until the DNS settings are updated, and
then pause for the TTL time before trying again.

## Tasks

1. Bring up dns-corrector in cluster
1. Resolve the current ip address and log it, and write to metric
1. Resolve the ip address of the anchor hostname and log it, and write to metric
1. Log if there's a difference between the two settings
1. Run this in a loop
1. Write code to update the DNS through cloudflare
1. Trigger update code when the DNS changes
1. Metrics on trigger time, and updte duration
1. Dashboard that captures the two resolved IPs
1. Dashboard captures last update time
1. Dashboard captures last update duration
1. Dashboard captures when dns-corrector is the process of correction
