# Seguro Bounty

# Overview

Seguro improves the Urbit binary's dependability by automatically replicating
its event log across a set of machines in a cluster.

## Problem

The [Urbit runtime](https://github.com/urbit/urbit/tree/master/pkg/urbit/vere)
was designed to host just a single Urbit instance running in a Unix process,
with a single file volume attached for its event log and checkpoints. For
individuals running only one or a handful of ships, this architecture is
satisfactory. For providing Urbit to the world as a new, decentralized,
peer-to-peer network of personal servers, the current implementation has
inarguably succeeded.

However, for a quickly maturing platform which needs to scale to meet enormous
demand, current Urbit technology is not suitable. Its low dependability is one
element which continues to prevent the Urbit community from building truly
resilient and scalable hosting services.

## Background

**Seguro** is an upgrade to the Urbit binary which replicates a ship's event
log, providing resiliency and redundancy to a currently fault-intolerant Urbit
program. When running Seguro, a hardware failure will result in some minimal
amount of network downtime but should not cause data loss and therefore should
never require a breach in order to gracefully recover. Seguro is one component
in service of an Urbit that can scale the globe with dependability high enough
to trust Urbit with even the most critical workloads.

# Project Requirements

## Potential Architecture Diagram

![Seguro Architecture](https://user-images.githubusercontent.com/91502660/153295790-6eef34ff-9136-4bc2-8927-2b432525c07d.png)

Seguro will use [dqlite](https://dqlite.io) to store and replicate a ship's
event log across a user-configurable number of replicas. The Urbit binary will
receive the addition of the following runtime command-line interface options:

1. `-m machines` (or `-r replicas`, or `-s slaves`) the comma-delimited list of
   IP:port addresses to replicate the event log across
2. `-o optimism` the number [0..m] of replications to wait for before emitting
   an event's side effects

## Cluster

A Seguro cluster will consist of an elected master and `m` replication slaves
(user-specified). Urbit network UDP events will be received by an IP-level load
balancer (running standalone or on the master?) which will then replicate them
to all the members of the cluster. Slaves will process these events and record
them in their local event log but only the Seguro master will actually emit
their side effects back to the Urbit network.

## Performance

### Configuration

Users will have the option to configure Seguro's degree of optimism. the number
of healthy replicas must process an event before emitting its side effect from
the master. For example, if set to 0, the master will receive, process and emit
side effects from an event without waiting for any slaves in the cluster to
acknowledge their successful reception and processing of said event. If set to
N, where N is the number of slaves in the cluster, the master will wait for all
of the slaves to process the event before emitting its side effects. 0 is most
performant, N is most durable. In other words, setting a value of 0 tells Seguro
to be "optimistic" about event replications, and N tells Seguro to be the
opposite.

### Event Log Batching

An event log batching system could be designed and implemented to improve the
performance of Seguro by way of reducing per-event processing and replication
overhead (if there is any).

## Features

## Snapshots & Event Log Truncation

Seguro's support of optimistic run-ahead as defined in the Configuration section
depends on coupling segments of the log to binary versions (for error handling
and crash recovery), which depends on "epochs" and/or event log truncation.
There is a [PR currently under review](https://github.com/urbit/urbit/pull/5701)
which implements this feature. This PR assumes existence of a local log that can
be subdivided into epochs where each epoch is coupled to a particular snapshot.
If the log is moved off of the module, and epochs/truncation are supported,
there will need to be some way to make sure that the relevant snapshots are
persisted with the same durability on slave machines.

## Previous Work

[Some effort](https://github.com/urbit/urbit/commit/cfeb35e37be63f96bb50fe1f60e2f59e35c07258)
towards a goal similar to Seguro's has already been performed by Tlon
engineering in the past. Adaption of dqlite/Raft for Urbit was conceived and
researched as far back as 2013. The concept needs a responsible organization to
champion it and provide resources and management so it can be finished. Rhe core
architecture is sound and, given sufficient development resources, Seguro ought
to be a reasonably achievable within a year's time.

## Open Questions

1. Should the ordering of (compute event, commit event, release effect) change?
   Why? How?

2. What is the nature of dqlite integration? Are we making a new kind of pipe
   where we control both sides and the remote talks to an arbitrary database, or
   does the database have a client SDK that gets pushed into the runtime?

3. Should we "just" write a Nock hypervisor (have the runtime use a Nock core as
   a hypervisor, rather than operating on the Arvo core)?

4. What happens in case of failed nodes, master or slaves?
5. How does automatic failover work?
6. What are the recommended default settings for running Seguro in hosting
   environments?

# Worker Requirements

Prospective candidates should have intermediate to advanced proficiency in the C
programming language, a healthy appetite for Martian software (although no
previous experienced with it is required), some previous experience with
distributed systems, expertise in writing well-defined technical requirements
and specifications, and strong coding and sytle habits. 3+ years of experience
and a full-time commitment are also required. While no previous experience with
Urbit is necessary, a candidate without it should be able to quickly demonstrate
knowledge of the basics of the Urbit OS and its runtime (i.e., after reading the
[whitepaper](https://media.urbit.org/whitepaper.pdf) or the relevant sections in
the [docs](https://urbit.org/docs)).

# Milestones

A worker who can confidently commit to completion of the entire Seguro project,
with all 3 of its milestones, is highly preferable. Guidance and leadership will
be provided by ~mastyr-bottec of [Wexpert Systems](https://wexpert.systems).

A salary of $125,000 will be paid to the work on a bi-weekly basis throughout
the duration of the bounty (one year). If the work is completed before the year
is complete, ???. If completed after, ???...

## Milestone 1 - Specification

Write a detailed technical specification and have it reviewed by stakeholders
(and perhaps additional technical Urbit authorities), accordingly revised, and
ultimately approved for implementation.

Expected Completion: 3 months

Deliverables:

- Publication of a set of Markdown documents

## Milestone 2 - Implementation

Clean, legible code which clearly implements the specification. Where
implementation of the specification as written is impossible or impractical, the
specification should be updated accordingly after sufficient consideration and
approval of the stakeholders.

Expected Completion: 3 months

Deliverables:

- A merged PR in the `urbit/urbit` repository

## Milestone 3 - Testing, Integration, Communication

Write a comprehensive test suite that guarantees operability of Seguro's
specified features and use cases. Integrate Seguro features in a live hosting
environment (maybe delete this one???).

Expected Completion: 3 months

Deliverables:

- A merged PR in the `urbit/urbit` repository
- Overview of Seguro and demonstration of working system on Urbit Developer Call

# Timeline

This bounty should take one year to complete.
