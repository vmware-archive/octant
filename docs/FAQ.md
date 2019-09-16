# Frequently Asked Questions

## Q: How is Octant different from the _"official"_ Kubernetes Dashboard?

<details><summary> Answer: </summary> 

Kubernetes Dashboard is described as a **_"general purpose_**, web-based UI for Kubernetes clusters.", whereas **Octant** is designed to be a _"tool for **developers** to understand how applications run on a Kubernetes cluster."_ 

Octant provides more detail, is more extensible, uses newer technology, and is under more active development.

More specifically:
- Octant does not run in a cluster (by default). Instead, it runs locally on your workstation and uses your kubeconfig files
- Octant has a resource viewer which links related objects to better describe their relationship within the cluster
- Octant supports Custom Resource Definitions (CRDs)
- The dashboard functionality of Octant is _not_ the #1 priority. The tool was created to help give users of Kubernetes more information  in an easier fashion than _kubectl get_ or _kubectl describe_
- Octant can be extended with plugins <link to plugin guide>
- Octant is being actively developed
- Octant is based on newer web technologies. The Kubernetes dashboard is based on AngularJS that has been superseded by Angular. 

</details>

## Q: How do I install Octant?
<details><summary> Answer: </summary>

## Installation


### Package (Linux only)

1. Download the `.deb` or `.rpm` from the [releases page](https://github.com/vmware/octant/releases).

2. Install with either `dpkg -i` or `rpm -i` respectively.

###  Windows

#### Chocolatey

1. Install using chocolatey with the following one-liner:

   ```sh
   choco install octant --confirm
   ```

#### Scoop

1. Add the [extras](https://github.com/lukesampson/scoop-extras) bucket.

   ```sh
   scoop bucket add extras
   ```

 2. Install using scoop.

   ```sh
   scoop install octant
   ```

### macOS

#### Homebrew

1. Install using Homebrew with the following one-liner:

   ```sh
   brew install octant
   ```

### Download a Pre-built Binary (Linux, macOS, Windows)

1. Open the [releases page](https://github.com/vmware/octant/releases) from a browser and download the latest tarball or zip file.

2. Extract the tarball or zip where `X.Y` is the release version:

    ```sh
    $ tar -xzvf ~/Downloads/octant_0.X.Y_Linux-64bit.tar.gz
    octant_0.X.Y_Linux-64bit/README.md
    octant_0.X.Y_Linux-64bit/octant
    ```

3. Verify it runs:

    ```sh
    $ ./octant_0.X.Y_Linux-64bit/octant version
    ```

</details>

## Q: How can I contribute to Octant?
<details><summary> Answer: </summary>

New contributors will need to sign a CLA (contributor license agreement). We also ask that a changelog entry is included with your pull request. Details are described in our [contributing](CONTRIBUTING.md) documentation.

See our [hacking](../HACKING.md) guide for getting your development environment setup.

See our [roadmap](../ROADMAP.md) for tentative features in a 1.0 release.

</details>

## Q: Why doesn't Octant support -feature X-?

<details><summary> Answer: </summary>

Octant is a community driven project with contributions from volunteers around the world.

</details>

## Q: When will Octant get Feature X?

<details><summary> Answer: </summary>

See our [roadmap](../ROADMAP.md) for tentative features in a 1.0 release.

</details>

## Q: What are the requirements to run Octant?

<details><summary> Answer: </summary>

...

</details>

## Q: How can I configure Octant to run in my Cluster?

<details><summary> Answer: </summary>

While Octant is designed to run on a developers desktop or laptop, it is possible to configure Octant to run inside a Kubernetes Cluster. Instructions are located here: https://github.com/vmware/octant/tree/master/examples/in-cluster

</details>

## Q: Where is the Octant community?

<details><summary> Answer: </summary>

Community info...

</details>

