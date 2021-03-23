---
title: "Octant 2020 Community Survey Results"
image: /img/posts/2020/11/19/survey-results-banner.png
excerpt: After launching our first Octant community survey, we want to share some of the insights we’ve gathered from your responses
author: Isha Bagha
author_name: Isha Bagha
author_avatar: /img/contributors/isha-bagha.jpg
categories: ['kubernetes']
tags: ['survey', 'Isha Bagha']
date: 2020-11-19
slug: octant-2020-community-survey-summary
---

Not too long ago, we set out to learn more about the Octant community. To do so, we shared the first [Octant Community survey](/octant-community-survey) and received 13 responses. It’s not a vast sample size, but we were able to find patterns on how Octant is used and how it could be improved. Here, we want to share our findings.

#### Methodology

The survey mentioned above collected a mix of qualitative and quantitative information about user profiles, workflows, and behaviors. Aside from that, we also had one-on-one interviews. Overall, we collected data from:

 * 13 survey respondents
 * 1 interview with a plugin author
 * 1 interview with an internal stakeholder at VMware

#### User Profiles

Most respondents work at organizations that have 1-50 developers, though we’ve interacted with many users who work at larger organizations (1000+ developers). Almost all survey respondents have deployed in production with Kubernetes, which tells us about their organizational maturity and level of Kuberenetes knowledge.

![](/img/posts/2020/11/19/org-size.png)

![](/img/posts/2020/11/19/k8s-prod.png)

### Infrastructure

Operating systems are a mix of Windows, Mac, and Linux, with most skewing towards Linux. We found that there aren’t many users running Octant on the private cloud.

![](/img/posts/2020/11/19/infrastructure.png)

### Why (or why not) Octant?

Most respondents use Octant for its UI and troubleshooting capabilities.

![](/img/posts/2020/11/19/value.png)

Through the qualitative feedback, we learned that users also use Octant for:

 * Teaching and demos
 * Avoiding switching between terminal windows
 * Ease of using kubeconfig for configuration and authentication
 * Viewing CRD’s
 * Resource viewer
 * Viewing logs
 * YAML edits
 * Port forwarding

#### Adjacent tooling

The following are tools as alternatives or in conjunction with Octant:

![](/img/posts/2020/11/19/alternative.png)

Users prefer these tools or use them in conjunction with Octant because of:

 * **CLI to avoid context switching** between terminal and browser (Lens, kubectl)
 * **More detailed metrics** and dashboards (Prometheus)
 * **Creating app deployments** and services (k8s dashboard)
 * **Speed and performance** (Lens, k9s)

#### Local Hosting vs SaaS

We asked users how their workflows might change if Octant were available as a web service. Most users responded that they prefer Octant as a local tool because it addresses the security concerns that come with web-based applications. Octant was built to run locally to address these security concerns; we validated that users adopt Octant for these reasons.

Though Octant is intended to be run locally, there are some users who would want to use Octant as a web service, and find workarounds to running Octant as a web-based application by deploying it in-cluster and adding a layer of user management.

#### Plugins

Octant’s superpower lies in its extensibility via plugins. Users can extend and customize the Octant UI through installing published plugins or creating their own plugins.

### Improve Messaging Around Octant’s Extensibility

From a plugin author and stakeholder perspective, there’s potential for stronger messaging around how powerful this functionality is, and to make building and using plugins more accessible.

### Plugin Author Experience

A few insights and areas for improvement we took away the plugin author experience include:

 * **How do authors learn to build plugins?** Some users turn to documentation for learning how to build a plugin, some reverse engineer and reuse patterns from existing plugins
 * **More examples with implementation of components**: It could be helpful not only to have more examples, but also provide examples with a broader usage of UI components
 * **Testing & GitOps**: Better CI and testing for pushing new versions of a published plugin
 * **Discoverability**: It could be helpful to have central catalog of plugins for publishing and discoverability

Overall, there’s room to enhance the marketing and messaging around Octant’s extensibility, and improve the UX for plugin discoverability and publishing.

#### Closing thoughts

Thank you all once again for participating in the survey! The data we collected either helped us to validate assumptions or gave us ideas on how we can shape Octant moving forward.

For updates on our roadmap, keep an eye out our [Octant issues](https://github.com/vmware-tanzu/octant/issues) and our [roadmap](https://github.com/vmware-tanzu/octant/blob/master/ROADMAP.md). If you have additional comments or feedback reach out to us on Slack via [#octant](https://kubernetes.slack.com/archives/CM37M9FCG).


