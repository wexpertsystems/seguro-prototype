# Seguro

Seguro improves the Urbit binary's dependability by automatically replicating
its event log.

# Background

Seguro, an improvement to `vere` (the Urbit runtime binary), replicates a ship's
event log, providing resiliency and redundancy to a currently fault-intolerant
Urbit program. When running Seguro, a hardware failure will result in some
minimal amount of network downtime but should not cause data loss and therefore
should never require a breach in order to gracefully recover.

# Design

![Seguro Architecture](https://user-images.githubusercontent.com/91502660/153295790-6eef34ff-9136-4bc2-8927-2b432525c07d.png)

Seguro will use [dqlite](https://dqlite.io) to store and replicate a ship's
event log across a user-configurable number of replicas. In practice, the Urbit
binary will receive the addition of the following runtime command-line options:

1. `-m machines` (or `-r replicas`) the comma-delimited list of IP:port
   addresses to replicate the event log across
2. ???

## Cluster

A Seguro cluster will consist of an elected master and a set of replication
slaves. Urbit network UDP events will be received by an IP-level load balancer
which will then replicate them to all the members of the cluster. Cluster
members will all process these events and record them in their local event log
but only the Seguro master will actually emit their side effects back to the
network.

## Event processing flow

## Optimistic run-ahead

## Event log batching

## Event log truncation

## Previous work

[Some effort](https://github.com/urbit/urbit/commit/cfeb35e37be63f96bb50fe1f60e2f59e35c07258)
towards a goal similar to Seguro's has already been performed by Tlon
engineering in the past. Adaption of dqlite/Raft for Urbit was conceived and
researched as far back as 2013. The concept needs a responsible organization to
champion it and provide resources and management so it can be finished. Rhe core
architecture is sound and, given sufficient development resources, Seguro ought
to be a reasonably achievable within a year's time.

# About Us

## TODO: add more info about Pax and I

## TODO: add job description and/or desired qualities in candidate?

[Wexpert Systems](https://wexpert.systems) is an Urbit software and services
company which offers full-service premium Urbit hosting on physical servers,
with minimal dependence on Big Tech cloud services. Wexpert Systems sells Urbit
planets to the public and offers bespoke hosting for stars and galaxies and
their moons. Shop our planet store or contact us for hosting or other questions.
