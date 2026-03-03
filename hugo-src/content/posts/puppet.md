+++
date = '2023-09-20T10:00:00-06:00'
draft = false
title = 'Puppet: Systems Theory and Configuration Management'
+++

> “From a drop of water…a logician could infer the possibility of an Atlantic or a Niagara without having seen or heard of one or the other."
>
> — <cite>Sherlock Holmes, A Study in Scarlet</cite>

In a similar fashion, any astute inquirer may start with a single piece of creative output and find their bearings in an entirely new field of thinking. Today, that inquirer looks at **Puppet**, a software *system*—or, in its own words, a "system configuration management tool."

## Puppet Biography
Started by **Luke Kanies** in 2004, Puppet was born from a sysadmin who "learned to be a programmer" while creating it. Licensed as "open-core" (Apache core with proprietary extensions), it's now owned by Perforce.

### Lineage and History
The lineage of Puppet traces back to **CFEngine (Configuration Engine)**, created by Mark Burgess circa 1993. Before that, tools like IsConf v3 were being written to manage the infrastructure of the early 90s (Novell, Sun Microsystems).

Crucially, the tools came before the theory was formalized. In 1998, Steve Traugott published "[Bootstrapping an Infrastructure](https://www.usenix.org/legacy/publications/library/proceedings/lisa98/full_papers/traugott/traugott.pdf)" based on his work at NASA. His big idea: the infrastructure is one coherent machine—an "enterprise virtual machine."

## Systems Theory
Systems administration has found a market rationale. It is as much a science as computer science and undeniably a branch of engineering. Like any operational role, it often only comes to the fore in an organization when things go wrong (think [NotPetya](https://www.wired.com/story/notpetya-cyberattack-ukraine-russia-code-crashed-the-world/)).

### The Ideal World: Three Koans
When things go right, three truths of the field become apparent:

1.  **Policy:** Aims and wishes that are both human-readable and machine-readable. It is what we want and how things *should* be. (e.g., "Deny undergrads their unalienable right to free printing.")
2.  **Predictable:** The ability to reason about the state of the system. A defined "good" state defines reliability.
3.  **Scalable:** Systems that remain predictable as they grow in size.

### Core Definitions
*   **System:** A systematic structure where certain stimuli produce certain expected responses.
*   **Deterministic:** The property of having those certain expected responses.
*   **Nested:** Subsystems within larger systems—a *holonarchy*. Managing nodes of computers is easy; nodes of people (see [BOFH](https://en.wikipedia.org/wiki/Bastard_Operator_From_Hell)) is not.
*   **Uniformity vs. Variety:** Nodes must balance these. More uniformity leads to better static predictability, while more variety minimizes loss potential ("hedging your bets" against Murphy’s Law).

## Models for Applying Change
How do we move a system toward our desired state?
*   **Push vs. Pull:** Does the central server push the config, or does the node pull it?
*   **Divergence, Confluence, and Congruence:**
    *   **Divergence:** Moving away from the goal state.
    *   **Confluence:** Multiple paths leading to the goal state.
    *   **Congruence:** The goal state matches the actual state.

## Why Order Matters
In 2002, Steve Traugott explored why his methods were consistently successful. His intuition: the ordered, deterministic nature of the toolset granted success. Contra many of their peers, those who used deterministic tools could *prove* their superiority through reliability and fast rollouts.

## Today and Beyond
Puppet paved the way for modern infrastructure. Today, we see:
*   **NixOS:** A truly congruent model for system configuration.
*   **Databases:** Is confluence enough for stateful systems?
*   **Regulators & Audits:** How do we prove compliance through code?
*   **Cloud Infrastructure:** From Amazon AMIs to immutable infrastructure.

## Presentation Resources
The core principles of policy, predictability, and scalability remain the bedrock of systems engineering.

Below you can find the original sources for this presentation:
- [Presentation Notes (Markdown)](/downloads/puppet-presentation-notes.md)
- [Slideshow (OpenDocument Presentation)](/downloads/puppet-slideshow.odp)
