# k8s-research

Repository for personal Kubernetes research and write-ups.

## Purpose of this project

Explore Kubernetes concepts by building a civilization simulation.

I will attempt to map Kubernetes concepts using a real-world analogies while building a runnable simulation.

## Project Structure

The project is structured as a series of entries that build on each other.  
Each entry will introduce new concepts that expand the civilization simulation analogy, usually by introducing new kubernetes resources and fitting them into the analogy.
There will be a separate directory for each entry that will be self-contained and have its own setup and cleanup scripts.

Each entry will have a `README.md` file that explains how to "run" the entry.  
The write-up for each entry will be in a separate markdown file in the parent `docs/` directory.

## Pre-requisites

- [KinD](https://kind.sigs.k8s.io/)
- [Helm](https://helm.sh/)
- [Skaffold](https://skaffold.dev/)