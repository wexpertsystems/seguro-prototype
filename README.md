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

# About Us

## TODO: add more info about Pax and I

[Wexpert Systems](https://wexpert.systems) is an Urbit software and services
company which offers full-service premium Urbit hosting on physical servers,
with minimal dependence on Big Tech cloud services. Wexpert Systems sells Urbit
planets to the public and offers bespoke hosting for stars and galaxies and
their moons. Shop our planet store or contact us for hosting or other questions.

# Whitepaper Paste

Seguro replicates a ship's event log, providing resiliency and redundancy to a
currently fault-intolerant Urbit program. When running Seguro, a hardware
failure will result in some minimal amount of network downtime but crucially
should not cause data loss and therefore should never require a breach in order
to gracefully recover.

![Seguro Architecture](https://user-images.githubusercontent.com/91502660/153295790-6eef34ff-9136-4bc2-8927-2b432525c07d.png)

Seguro will be based on a dqlite and Raft architecture to provide data
replication, high availability and failover for the Urbit platform. A Seguro
cluster will consist of an elected master and two or more replication slaves.
Urbit network UDP events will be received by an IP-level load balancer which
will then replicate them to all the members of the cluster. Cluster members will
all process these events and record them in their local event log but only the
Seguro master will actually emit their side effects back to the network.

If the Seguro master crashes or is otherwise not available to the cluster a Raft
election will take place and a new elected master will be promoted and pick up
where the system left off. Seguro will be independent from Armada in cases where
only high availability is desired, such as hosts of popular groups or any other
situation that demands better availability than a standalone Urbit ship can
provide. It will unavoidably add some complexity to the intial setup and support
of the Urbit ship but, being entirely optional, that complexity ought to only be
only a factor for those who choose to use it.

Some effort towards this goal has already been performed and adaption of
dqlite/Raft for Urbit was conceived and researched as far back as 2013. The
concept needs some responsible organization to champion it and provide resources
and management so it can be finished, but the core architecture is sound and,
given sufficient development resources, Seguro ought to be a reasonably
achievable improvement at this point in Urbit's development.

Challenges to Seguro are in managing the event loop and recording events. An
implementation of optimistic run-ahead logic and batching of event logs writes
will be necessary for performance, but this should be safe to implement given
the added reliability of a clustered architecture. Ideally we may also want
aggressive event log truncation to reduce the storage hardware requirements but
this may not be necessary for a successful initial launch.
