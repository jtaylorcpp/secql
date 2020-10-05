SECQL
========

SecQL is an integration layer between a number of cloud provider SDKs and security tooling/agents. The goal is to be able to access the information about the entiretly of a polyglot stack while using a standard and well understood vocabulary and schema.

### Supported Cloud Providers

#### AWS (EC2)
SecQL supports EC2 instance metadata collection, AMI metadata colllection, and the use of EC2 Instance Connect to do interactive/shell-out style calls to supported tooling. SecQL, when using EC2 Instance Connect, will generate a unique key-pair per instance per server. All credentials are stored in memory and upon restarting the server new Connect sessions will be started begining with generating new SSH keys.


### Supported Secuirty Products

#### OSQuery
[OSQuery](https://osquery.io/) is an open source security agent which allows all actions, processes, etc. on a machine to be queried as if it were a collection of SQLite tables.

SecQL has 2 ways of interacting with OSQuery. First, by creating an interactive ssh session and directly calling `osqueryi`. Or, second, by using the `secqld` agent which uses the OSQuery configuration to dynamically generate endpoints which can be scraped. 

## Getting Started