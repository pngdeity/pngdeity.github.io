+++
date = '2026-09-16T10:00:00-05:00'
draft = true
title = 'Open Source on Offense'
+++

> "Linux is a cancer."
>
> <cite>Steve Ballmer, Microsoft CEO, 2001</cite>

Why on Earth would Microsoft fund the creation and maintenance of an entirely free and open source operating system?

The answer is not that anyone at Microsoft read the GNU Manifesto and decided they liked it. It's that open source can be a weapon, and Microsoft has gotten very good at wielding it. The company that once called Linux a cancer[^ballmer] now treats it as raw material for the oldest trick in technology economics: [commoditize your complement](https://www.joelonsoftware.com/2002/06/12/strategy-letter-v/). Give away the layer you don't need, and keep the one you do.

[^ballmer]: The fuller quote, from a 2001 interview: "Linux is a cancer that attaches itself in an intellectual property sense to everything it touches." Notice what he objected to: the GPL's copyleft, not the quality of the kernel.

<!-- TODO(author): optional personal on-ramp. If you want to say how you first ran into this story, or why it gets under your skin, this is the place. Otherwise the argument stands on its own. -->

Once you see the pattern, a lot of corporate open source stops looking selfless, and a lot of otherwise baffling decisions start to make sense. Enter SONiC: a network operating system Microsoft built, gave away, and now leans on to keep the switching market a buyer's market.

## What SONiC is

**SONiC** — "Software for Open Networking in the Cloud" — is a [free and open source network operating system](https://en.wikipedia.org/wiki/SONiC_(operating_system)) built on Debian Linux, and it runs on the switches inside a data center. Microsoft developed it and [open-sourced it in 2016](https://azure.microsoft.com/en-us/blog/microsoft-showcases-software-for-open-networking-in-the-cloud-sonic/), building on an earlier release called the [Azure Cloud Switch](https://azure.microsoft.com/en-us/blog/microsoft-showcases-the-azure-cloud-switch-acs/), and on the Open Compute Project's [Switch Abstraction Interface](https://github.com/opencomputeproject/SAI) (SAI).

Two decisions do most of the work:

- **The software is decoupled from the hardware.** SONiC talks to the switch through SAI, so the same network OS runs on [more than 100 different switches](https://sonicfoundation.dev/) from a pile of different vendors and ASIC families.
- **The system is containerized.** Network functions run as microservices instead of one monolithic image, which lets operators swap and upgrade pieces without touching the rest.

In April 2022, Microsoft [handed governance to the Linux Foundation](https://www.linuxfoundation.org/press/press-release/software-for-open-networking-in-the-cloud-sonic-moves-to-the-linux-foundation). "We created it as open source so the entire networking ecosystem would grow stronger," said Dave Maltz, a corporate vice president for Azure networking. He added that SONiC already runs on millions of ports in the networks of cloud scalers, enterprises, and fintechs.

Here's the part I can't get over: the project's premier members today include Microsoft, Google, Alibaba, Broadcom, Dell, and NVIDIA, alongside **Cisco, Arista, and Nokia** — the incumbents whose proprietary switch software SONiC was built to displace.

## Why switches?

Of all the layers Microsoft could have gone after, why this one? Because a commoditization play only works when three things line up at once, and data-center switching was one of the few places where they did.

First, it is a thing Microsoft buys in enormous volume and never sells. The Azure network alone ran [more than 180,000 switches](https://www.nextplatform.com/2021/05/12/microsoft-does-the-math-on-azure-datacenter-switch-failures/) across 130 locations as of 2021, and SONiC now runs on "millions of ports." Every datacenter build buys the whole layer again, forever.

Second, a supplier was extracting rent from it. Cisco's product gross margins have sat in the 60s for a decade, and it sold the box and the software as a single unit, so the premium was never itemized. You could not buy the switch and decline the lock-in.

Third, the technical path to breaking that lock had just opened. Broadcom's merchant silicon put a commodity switch ASIC on the open market, which meant the hardware was no longer the moat. What remained was the network OS — and a network OS is software, which is a thing Microsoft already knew how to write at scale.

Notice what the third condition rules out. Microsoft cannot commoditize top-end CPUs; there is no merchant alternative, and RISC-V is not close. And it will never commoditize operating systems in general, because it sells one. The sweet spot is a layer where it is a giant buyer, a rentier sits on the other side, and a technical shift has just cleared the way. Switches checked every box. It is the mirror image of the PC play: Microsoft used to sell the OS and commoditize the box; in the datacenter it buys the box and commoditizes the OS.

To put a number on the rent, take the same generation of switch — 3.2 Tbps of capacity, 32 100-gigabit ports — from each side:[^prices]

| | Cisco Nexus 9332C (list) | Edgecore AS7712-32X (bare metal) |
| :--- | ---: | ---: |
| Capacity | 3.2 Tbps | 3.2 Tbps |
| Price | ~$51,400 | ~$8,340 |
| Per 100G port | ~$1,606 | ~$261 |
| Per gigabit of capacity | ~$16.06 | ~$2.61 |

[^prices]: Cisco list price for the N9K-C9332C from the Cisco global price list; Edgecore AS7712-32X price from Colfax Direct. Both as of September 2026. The Edgecore unit is bare metal with no network OS; the Cisco price includes NX-OS.

Call it a sixfold spread at list. Cisco's street prices are far below list, so the defensible claim is "a large multiple," not "six times" — but even at 60 or 70 percent off, the incumbent still lands several times above the white-box price. Spread a delta of that size across 180,000 switches and the arithmetic reaches into the billions, which is why this was worth Microsoft's engineering time before a single line of SONiC was written.

The market shares tell the same story with fewer caveats. In 2015 Cisco took about [60 percent of all Ethernet switch revenue](https://www.fierce-network.com/cloud/cisco-grabs-60-percent-ethernet-switch-revenue-2015-delloro-group). By 2022, Omdia had Cisco at 37 percent of data-center Ethernet switch revenue, Arista at 18, and white-box vendors at 14. The best bookend of all: in 2015 Cisco published a post arguing that [white-box switches were "no bargain"](https://blogs.cisco.com/news/myth-busting-white-box-switches-are-no-bargain) — 20 to 30 percent *more* expensive once support was counted — and by 2026 it was [selling its Silicon One silicon to the white-box builders and SONiC operators](https://hyperframeresearch.com/2026/05/24/ciscos-silicon-one-an-8-12bn-business/). When the incumbent starts arming the commodity ecosystem, the commodity has won.

Two honest caveats. List price overstates the gap, because nobody pays list. And the white-box figure excludes optics, NOS support, and integration, while the Cisco figure includes NX-OS and a support contract. SONiC also did not set these prices by itself: merchant silicon did the heavy lifting, and SONiC removed the last lock-in. The fair claim is narrower, and more interesting, than "SONiC cut switch prices sixfold." The software was the final brick in the wall.

<!-- TODO(author): numbers to own or replace. The prices are public list figures, not quotes you would actually get. If you have real procurement numbers — or would rather cut the table and let the share data carry the point — say the word. -->

## Commoditizing the complement

There's a name for this. In 2002, Joel Spolsky gave it one: [commoditize your complement](https://www.joelonsoftware.com/2002/06/12/strategy-letter-v/).

A *complement* is a product people buy alongside yours. Demand for your product rises when the price of its complements falls, so your strategic interest is to drive those prices down toward marginal cost — and there is no cheaper price than free.

Spolsky's examples are the familiar ones: IBM documenting the PC so the add-in market would commoditize, and Microsoft licensing MS-DOS to every clone-maker so the PC itself became a commodity. In both cases the point was never the layer being commoditized; it was the layer that stayed scarce. For IBM that layer was the PC, briefly. For Microsoft it was the operating system, for a long time.

SONiC is that play, run on switches. What Microsoft gave away was the network OS; what it bought itself was a future in which a switch is an interchangeable box, bought on price from whichever white-box vendor is cheapest.

## Giving away the standard

Commoditizing the software was phase one. Phase two was to stop owning it.

A project controlled by Microsoft would always be suspect, especially to the vendors it was commoditizing. By moving SONiC to the Linux Foundation, Microsoft made it neutral enough that competitors could adopt it and contribute to it without having to trust Redmond. That is how a Microsoft project became an industry standard: not by winning the market, but by declining to own the thing.

This isn't new either. Eric S. Raymond described the same maneuver in 1999, in a section of *The Cathedral and the Bazaar* called "[Open Source as a Strategic Weapon](http://www.catb.org/~esr/writings/cathedral-bazaar/magic-cauldron/ar01s11.html)": DEC funding the X Window System to "reset the competition" against Sun, and vendors funding Apache so that none of them had to beat Microsoft alone.

## The license is a means, not a value

Here is the uncomfortable part for anyone who came to free software for the freedom. The license is a mechanism, not a motive. Copyleft and permissive licensing are levers with different properties: copyleft stops a layer from being enclosed by anyone, and permissive licensing maximizes adoption. SONiC ships under a [mix of GPL and Apache](https://en.wikipedia.org/wiki/SONiC_(operating_system)), chosen for what it accomplishes rather than for what it professes.

The test of whether openness is a value or a means is simple: does it survive when it stops paying? Usually it does not. MongoDB moved to the SSPL in 2018, Elastic followed in 2021 before adding the AGPL back in 2024, HashiCorp relicensed Terraform in 2023 and got forked as OpenTofu, Redis left the BSD license in 2024, and Wizards of the Coast tried to revoke the open D&D license in 2023 once it had done its job. The modern move is to build the escape hatch in advance: a permissive license plus a contributor license agreement, so the company keeps the option to relicense later.

None of which means the contributors are cynical, or that the code isn't genuinely free. Both things can be true at once: the software is free, and the freedom is on the company's terms.

<!-- TODO(author): your paragraph. Say plainly what you think the uncanny part is — that "open" here is a tactic whose value is instrumental. -->

## Further reading

- Joel Spolsky, [Strategy Letter V: The Economics of Open Source](https://www.joelonsoftware.com/2002/06/12/strategy-letter-v/) (2002)
- Eric S. Raymond, [Open Source as a Strategic Weapon](http://www.catb.org/~esr/writings/cathedral-bazaar/magic-cauldron/ar01s11.html), in *The Cathedral and the Bazaar* (1999)
- Gwern Branwen, [Laws of Tech: Commoditize Your Complement](https://gwern.net/complement) (2018–2022)
- Linux Foundation, [SONiC Moves to the Linux Foundation](https://www.linuxfoundation.org/press/press-release/software-for-open-networking-in-the-cloud-sonic-moves-to-the-linux-foundation) (2022)
- [SONiC Foundation](https://sonicfoundation.dev/) and [SONiC on Wikipedia](https://en.wikipedia.org/wiki/SONiC_(operating_system))

<!-- TODO(author): still needed — your personal stake in the topic, and a closing line. The middle is factual scaffolding you can rewrite in your voice. -->
