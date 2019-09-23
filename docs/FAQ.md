# Frequently Asked Questions

## Q: How is Octant different from the _"official"_ Kubernetes Dashboard?

<details><summary> Answer: </summary> 

Kubernetes Dashboard is described as a **_"general purpose_**, web-based UI for Kubernetes clusters.", whereas **Octant** is designed to be a _"tool for **developers** to understand how applications run on a Kubernetes cluster."_ 

Octant provides more detail, is more extensible, uses newer technology, and is under more active development.

More specifically:
- Octant does not run in a cluster (by default). Instead, it runs locally on your workstation and uses your kubeconfig files
- Octant has a **resource viewer** which links related objects to better describe their relationship within the cluster
- Octant supports Custom Resource Definitions (CRDs)
- The dashboard functionality of Octant is _not_ the #1 priority. The tool was created to help give users of Kubernetes more information  in an easier fashion than _kubectl get_ or _kubectl describe_
- Octant can be extended with plugins 
    - Plugin docs here: [docs/plugins](https://github.com/vmware/octant/tree/master/docs/plugins)
- Octant is being very actively developed, with [major releases only weeks apart.](https://github.com/vmware/octant/releases)
- Octant is based on newer web technologies. The Kubernetes dashboard is based on "AngularJS" which has been superseded by "Angular". 

</details>

## Q: How do I install or Update Octant?
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

### Upgrading

The process of upgrading Octant will depend on how you installed it. Generally, you can use Update or Upgrade functions of the package manager you used to install Octant. (e.g. brew upgrade octant)

If you downloaded a pre-built binary, you could download the new version and replace the old one manually.

</details>

## Q: How can I contribute to Octant?
<details><summary> Answer: </summary>

Octant is a community-driven project, and as such welcomes new contributors from the community. 

Ways you can contribute with a Pull Request:
- Documentation
    - See something wrong or missing from our docs? 
    - Do you have a unique use-case not documented?
- Octant core
    - Octant is written mostly in Golang and Angular. Our hacking guide can be found [here](https://github.com/vmware/octant/blob/master/HACKING.md)
- Plugins
    - Octant has a very extensible plugin model designed to let contributors add functionality. A plugin can read objects, and allows users to add components to Octant's views.  
    - A sample plugin is available [here](https://github.com/vmware/octant/blob/master/cmd/octant-sample-plugin)
    - A list of community plugins for Octant will be assembled soon

New contributors will need to sign a CLA (contributor license agreement). We also ask that a changelog entry is included with your pull request. Details are described in our [contributing](CONTRIBUTING.md) documentation.

See our [hacking](../HACKING.md) guide for getting your development environment setup.

See our [roadmap](../ROADMAP.md) for tentative features in a 1.0 release.

**Ways to contribute without a Pull Request?**
- Share the love on social media with the hashtag #octant
- Participate in Octant community meetings
- 


</details>

## Q: Is Octant stable?

<details><summary> Answer: </summary>

Octant is under active development, but each release is considered stable. 

Release information can be found here:
- [Releases](https://github.com/vmware/octant/releases)

Open Issues can be found here: 
- [Open Issues](https://github.com/vmware/octant/issues)

</details>

## Q: Why doesn't Octant support Feature X?

<details><summary> Answer: </summary>

Octant is a community driven project with contributions from volunteers around the world. 

If a feature you want is not already on our [Roadmap](https://github.com/vmware/octant/blob/master/ROADMAP.md), please feel free to [file an issue](https://github.com/vmware/octant/issues/new) and request it, or submit a Pull Request with your feature to be reviewed and [merged](https://github.com/vmware/octant/blob/master/CONTRIBUTING.md).

</details>

## Q: When will Octant get Feature X?

<details><summary> Answer: </summary>

See our [roadmap](../ROADMAP.md) for tentative features in a 1.0 release.

</details>

## Q: What are the system requirements to run Octant?

<details><summary> Answer: </summary>

Octant supports running on macOS, Windows and Linux using either [pre-built binaries](https://github.com/vmware/octant/releases) or [building directly from source.](https://github.com/vmware/octant/blob/master/HACKING.md)

Octant requires an active KUBECONFIG (i.e. kubectl configured and working).

Octant does not reqiure special permissions within your cluster because it uses your local kubeconfig information.

</details>

## Q: How can I configure Octant to run in my Cluster?

<details><summary> Answer: </summary>

While Octant is designed to run on a developers desktop or laptop, it is possible to configure Octant to run inside a Kubernetes Cluster. Instructions are located here: https://github.com/vmware/octant/tree/master/examples/in-cluster

</details>

## Q: Where can I get help with Octant?

<details><summary>Answer: </summary>

The best way to get help is to file an issue on GitHub:
- https://github.com/vmware/octant/issues 


You can also reach out in our communities:

- On Slack 
    - slack.k8s.io #octant
- On Twitter 
    - [@ProjectOctant](https://twitter.com/projectoctant) 
    - Hashtag:  [#octant](https://twitter.com/search?q=%23octant)
- In Google Groups
    - https://groups.google.com/forum/#!forum/project-octant/

</details>

## Q: Where is the Octant community?

<details><summary> Answer: </summary>

We welcome community engagement in the following places:

- On Slack 
    - slack.k8s.io #octant
- On Twitter 
    - [@ProjectOctant](https://twitter.com/projectoctant) 
    - Hashtag:  [#octant](https://twitter.com/search?q=%23octant)
- In Google Groups
    - https://groups.google.com/forum/#!forum/project-octant/


</details>




