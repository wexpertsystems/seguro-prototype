# Seguro Bounty

# Overview

Seguro improves the Urbit binary's dependability by automatically replicating
its event log.

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
in service of an Urbit that can scale the globe.

# Project Requirements

> An articulation of what constitutes complete work. User stories, Figma
> designs, interface specifications, and other technical constraints are all
> examples of requirements.

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

## Event Flow

## Optimistic Run-Ahead

## Event Log Batching

## Event Log Truncation

## Previous Work

[Some effort](https://github.com/urbit/urbit/commit/cfeb35e37be63f96bb50fe1f60e2f59e35c07258)
towards a goal similar to Seguro's has already been performed by Tlon
engineering in the past. Adaption of dqlite/Raft for Urbit was conceived and
researched as far back as 2013. The concept needs a responsible organization to
champion it and provide resources and management so it can be finished. Rhe core
architecture is sound and, given sufficient development resources, Seguro ought
to be a reasonably achievable within a year's time.

# Worker Requirements

> A description of the skills and/or qualifications required of a prospective
> worker. Technologies or skills known, years of experience, work schedule (e.g.
> time availability), and demonstrable accomplishments are all examples of this.

# Milestones

> Logical segments of work that can be considered done on their own. Breaking
> the project up into smaller pieces will help a new worker pick up the work
> should that be necessary, keeps the worker motivated with incremental
> achievement and remuneration, and is generally a sign of a clear
> specification. Milestones don't always make sense; sometimes there's only one.

## Milestone 1 - Specification

Expected Completion: 3-6 months

Payment: $50,000

Write a detailed technical specification and have it reviewed by the champion
(and perhaps additional technical Urbit authorities), accordingly revised, and
ultimately approved for implementation.

Deliverables:

- Publication of a set of Markdown documents

## Milestone 2 - Implementation

Expected Completion: 3-6 months

Payment: $50,000

Deliverables:

- A merged PR in the `urbit/urbit` repository

## Milestone 3 - Testing, Integration, Communication

Expected Completion: 3 months

Payment: $50,000

Write a comprehensive test suite that guarantee operability of Seguro's
specified features and use cases.

Integrate Seguro features in a live hosting environment.

Deliverables:

- A merged PR in the `urbit/urbit` repository
- Overview of Seguro and demonstration of working system on Urbit Developer Call

# Timeline

> If you're trying hit a certain deadline, make sure to specify absolute dates
> on the milestones. Otherwise, use relative dates (e.g. two months) to give the
> worker an idea of how long each milestone will take. This helps form an
> agreement between benefactor and worker on volume of work, and grounds to seek
> other arrangements should schedules be missed.
