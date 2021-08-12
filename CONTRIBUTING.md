# Contributing

## Communication

We prefer asynchronous to synchronous communication when ever possible. This ensures that all interested parties are able
to find relevant information even when ideas and information are being exchanged across many timezones. When ever
synchronous communication is happening, we encourage folks to make best efforts to transfer that information to an
issue if relevant or send it to the project-octant group. This could be as plain text notes or a link to a recording
or markdown document for example.

* [project-octant](https://groups.google.com/forum/#!forum/project-octant)
* [#octant Slack](https://kubernetes.slack.com/app_redirect?channel=CM37M9FCG)

### Weekly Community Meeting
We have a weekly Octant community meeting that is held live (recordings uploaded later). We highly encourage folks
who are interested in contributing to Octant to attend these meetings. More details on the
 [Octant community page](https://octant.dev/community/).

## Tools

The Octant project uses GitHub issues and pull requests as the primary method for communicating about what work needs
to be done, what work is currently being done, who work is assigned to, and the current state of that work.
Issues are classified under Projects as part of a larger scoped feature or a release milestone. The Backlog tracks new
issues as they are added and the current triage status.

### Github Projects

A project for a given feature or milestone contain the following lanes:

- To do: issues queued in order of priority
- In progress: issues with work in progress
- Review in progress: issues pending review or follow up to feedback
- Reviewer approved: issues pending merge
- Done: completed work

### Github Discussions

We use Github Discussions to archive old community meeting nodes and track longer conversations that are not suitable
for Slack. A discussion can be tagged with a category to increase visibility and searchability.

## Proposals

If you wish to add a larger feature or make a major refactor to Octant we encourage folks to write up a proposal document.
Though our process is not formal, the convention is to create a PR against Octant with your proposal as markdown in the
proposals folder. Proposal reviews and feedback will happen on the PR with the proposal.

    octant/proposals       <- folder
        YYYYMMDD-title.md  <- file

In your proposal it is good to consider and include some of the following:
* Goals
* Non-Goals
* Does this break any public APIs?
* Does this fix/solve any outstanding issues?

## Create a Pull Request

Contributors making a code change should fork the repository.

When creating a new pull request, fill in the template to reference an existing issue, if any. Run tests and linters locally
to check if the changes will pass CI. Stale pull requests should be updated with recent changes especially if there
are merge conflicts.

Commit messages should be concise and convey what changes are made. Additional context for reviewers is added to the
pull request along with comments to navigate the changes. Draft pull requests can also be used to solicit comments for
ongoing work in progress.

## CHANGELOG

Authors are expected to include a changelog file with their pull requests. The changelog file
should be a new file created in the `changelogs/unreleased` folder. The file should follow the
naming convention of `pr-username` and the contents of the file should be your text for the
changelog.

    octant/changelogs/unreleased   <- folder
        000-username               <- file

## Review Process

Pull requests are generally in one of the following states:

- Pending review: changes that have not been reviewed by one or more Octant maintainers
- CI failure: some combination of failing tests, style linting, merge conflicts, or unsigned DCO
- Pending changes from a feedback cycle: additional changes are required for merging
- Blocked: changes are dependent on another pull request
- Stale: changes require further discussion or abandoned by author
- Ready to merge: changes are approved

Reviews for a new pull request are targeted within one business day in order to avoid going stale. Large
changes are encouraged to be broken into smaller commits or multiple pull requests to facilitate a fast review cycle.

## DCO Sign off

All authors to the project retain copyright to their work. However, to ensure
that they are only submitting work that they have rights to, we are requiring
everyone to acknowledge this by signing their work.

Any copyright notices in this repo should specify the authors as "the Octant contributors".

To sign your work, just add a line like this at the end of your commit message:

```
Signed-off-by: Wayne Witzel III <wayne@riotousliving.com>
```

This can easily be done with the `--signoff` option to `git commit`.

By doing this you state that you can certify the following (from https://developercertificate.org/):

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.


Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```
