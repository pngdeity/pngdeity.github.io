+++
date = '2026-03-02T10:00:00-06:00'
draft = false
title = 'Threaded Fasteners for Computer Science'
+++

# Threaded Fasteners for Computer Science
*The first and last presentation you’ll ever need on the topic.*

In the world of computer science, we often focus on bits, bytes, and abstract algorithms. But when it comes to the hardware that runs our code, we eventually run into the physical world—and in that world, everything is held together by threaded fasteners.

## General Character
A fastener is more than just "a screw." It has a specific anatomy designed for engineering precision:

*   **Drive:** The interface where the tool (screwdriver) meets the fastener.
*   **Head:** The top part that provides torque and prevents the screw from passing through the material.
*   **Body/Shank:** The cylindrical portion, which may be partially or fully threaded.
*   **Threads:** The helical structure that provides the mechanical advantage.
*   **Point:** The very end of the fastener.

Key dimensions include:
*   **Nominal Length:** The length from the underside of the head to the point.
*   **Grip Length:** The unthreaded portion of the shank.
*   **Thread Length:** The portion covered in threads.

## Thread Geometry
Threads follow a standardized geometry. Most modern fasteners use a **60° thread angle**. 

Important diameters to know:
*   **Major Diameter (Dmaj):** The largest diameter of the thread.
*   **Minor Diameter (Dmin):** The smallest diameter (the "root" of the thread).
*   **Pitch Diameter (Dp):** An intermediate diameter where the thread width and the gap width are equal.
*   **Pitch (P):** The distance from one thread peak to the next.

## Standards
There are two primary standards you'll encounter in computing hardware:

### 1. Unified Thread Standard (UTS)
Regulated by **ASME/ANSI B1.2-1983**. For sizes smaller than 1/4 inch, they are designated by a number (e.g., #4, #6). 
*   **Example:** `#6-32 x 1/2"`
    *   `#6`: The major diameter.
    *   `32`: The number of threads per inch (TPI).
    *   `1/2"`: The length in inches.

### 2. Metric/ISO
Defined by **ISO 261** and **ISO 262**.
*   **Example:** `M8-1.25x30`
    *   `M8`: 8mm major diameter.
    *   `1.25`: The pitch in millimeters.
    *   `30`: The length in millimeters.

## Notable Examples in Computing
If you've ever opened a PC, you've used these:

| Screw Type | Common Application |
| :--- | :--- |
| **#6-32** | 3.5” HDDs, Case panels, Power supply mounting |
| **M3-0.5** | 2.5” SSD/HDD, 5.25” Optical drives, 3.5” Floppy drives |
| **M2-0.4** | M.2 (NGFF) SSDs, mini PCIe cards |
| **#4-40** | DVI, VGA, and the D-sub family of connectors |

## Epilogue: A Note on Pedantry
I'd just like to interject for a moment. What you're referring to as a "Philips head" is, in fact, a **Philips drive**, or as I've recently taken to calling it, **Ansi Type I Cross Recess**. 

## Original Presentation
This blog post is based on a lightning talk. You can view the original presentation below:

<embed src="/downloads/fasteners-lightning-talk.pdf" type="application/pdf" width="100%" height="600px" />

[Download the original PDF](/downloads/fasteners-lightning-talk.pdf)
