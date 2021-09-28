# Contributing to Monoskope

:+1::tada: First of, thanks for help improving Monoskope!

The following is a set of guidelines for contributing to Monoskope and its packages, which are hosted in the finleap connect [organization](https://github.com/finleap-connect) on GitHub. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Code of Conduct

This project and everyone participating in it is governed by the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [fcloud-connect@finleap.com](mailto:fcloud-connect@finleap.com?subject=[m8]%20COD%20Violation).

## Developer Certificate of Origin

The DCO is a declaration attached to every contribution made by every developer. In the commit message of the contribution, the developer simply adds a Signed-off-by statement and thereby agrees to the [DCO](DCO).

Each commit must include a DCO which looks like this

```Signed-off-by: Jane Doe <jane.doe@email.com>```

You may type this line on your own when writing your commit messages. However, if your `user.name` and `user.email` are set in your git configs, you can use `-s` or `--signoff` to add the Signed-off-by line to the end of the commit message.

## Developing with Monoskope

See the detailed documentation at [docs/development/README.md](docs/development/README.md)

## Submitting a Pull Request

Do you have an improvement?

1. Submit an [issue][issue] describing your proposed change.
2. We will try to respond to your issue shortly.
3. Fork this repo, develop and test your code changes. See the project's
   [README](README.md) for further information about working in this repository.
4. Submit a pull request against this repo's `main` branch.
    - Include instructions on how to test your changes.
5. Your branch may be merged once all configured checks pass, including:
    - The branch has passed tests in CI.
    - A review from appropriate maintainers (see
      [MAINTAINERS.md](MAINTAINERS.md))

## Committing

We prefer squash or rebase commits so that all changes from a branch are
committed to main as a single commit. All pull requests are squashed when
merged, but rebasing prior to merge gives you better control over the commit
message.

### Commit messages

Finalized commit messages should be in the following format:

```txt
Subject

Problem

Solution

Validation

Fixes #[GitHub issue ID]
```

#### Subject

- one line, <= 50 characters
- describe what is done; not the result
- use the active voice
- capitalize first word and proper nouns
- do not end in a period — this is a title/subject
- reference the GitHub issue by number

##### Examples

```txt
bad: server disconnects should cause dst client disconnects.
good: Propagate disconnects from source to destination
```

```txt
bad: support tls servers
good: Introduce support for server-side TLS (#347)
```

#### Problem

Explain the context and why you're making that change.  What is the problem
you're trying to solve? In some cases there is not a problem and this can be
thought of as being the motivation for your change.

#### Solution

Describe the modifications you've made.

If this PR changes a behavior, it is helpful to describe the difference between
the old behavior and the new behavior. Provide before and after screenshots,
example CLI output, or changed YAML where applicable.

Describe any implementation changes which are particularly complex or
unintuitive.

List any follow-up work that will need to be done in a future PR and link to any
relevant GitHub issues.

#### Validation

Describe the testing you've done to validate your change.  Give instructions for
reviewers to replicate your tests.  Performance-related changes should include
before- and after- benchmark results.

[issue]: https://github.com/finleap-connect/monoskope/issues/new
