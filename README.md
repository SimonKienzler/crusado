# crusado

**C**urating **R**eusable **U**ser **S**tories for **A**zure **D**ev**O**ps

***

## Do I need this?

Absolutely you do, probably! `crusado` is for you if

* you have to **repeatedly create similar User Stories or Bugs** in Azure DevOps
* the Web UI's way of handling User Story templating is **not doing it for you**
* you work in the terminal a lot and **CLIs are your preferred way of doing things**
* the WYSIWYG editor of Azure DevOps scares you and **you prefer Markdown**
* you want to **put your User Story and Bug templates under version control**
* **automating the creation of many work items through scripts** is something you'd
  like to do

## A Word of Caution

`crusado` is a hobby project in very early stages of development. So it
definetly will be rough around the edges and there might be issues here and
there. If you find a bug, please open an issue!

## Getting Started

To get started with `crusado`, you need a couple of things:

1. Installed `crusado` CLI
1. Azure Personal Access Token (PAT) with correct permissions
1. A few environment variables set
1. Markdown files with Frontmatter that contain your User Story and Bug templates

Let's take care of these things one by one.

### 1 Install `crusado` CLI

Head to the [latest release](https://github.com/SimonKienzler/crusado/releases)
and download the version suitable for your system. Be sure to put the binary at
a file location that is in your `$PATH`. Verify the installation was successful
by running:

```sh
crusado version
```

Alternatively, you can build from source after cloning the repository:

```sh
make build
```

Then invoke the CLI using `./bin/crusado`.

### 2 Create Azure PAT

`crusado` uses an Azure Personal Access Token (PAT) to create work items on your
behalf. Follow the [Azure Docs on PAT
creation](https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate).
`crusado` needs the following scopes:

* `Work Items (Read & Write)`

That's it!

> If you encounter issues with permissions later on when using `crusado`, please
> open an issue.

Store the secret PAT somewhere safe.

### 3 Set Environment Variables

You need the following environment variables set to use `crusado`:

```sh
# the URL of the organization you want to operate in
export CRUSADO_AZURE_ORG_URL=https://dev.azure.com/<your organization name>

# an organization usually has many projects, specify your desired one here
export CRUSADO_AZURE_PROJECT_NAME=<your project name>

# the Azure PAT you created before. Please be careful with pasting sensitive
# data like this to your terminal, it might end up in the history of your shell
export CRUSADO_AZURE_PAT=<your PAT>

# crusado does not yet have a well-known, default config path. For now,
# you have to explicitly set the path to your profile (which we'll create
# in the next step). Recommended value: ~/.crusado/<your project name>.yaml
export CRUSADO_TEMPLATES_DIR=./example/profile.yaml
```

### 4 Create Your `crusado` Template Files

In the last step, you set the `CRUSADO_TEMPLATES_DIR` environment variable.
Create as many Markdown files as you like at the location you specified as that
variable's value. You can use these files to manage your User Story and Bug
templates with a combination of Markdown text and some Frontmatter.

Take a look at this minimal example:

```md
---
name: example-story
summary: This is a user story template.
type: UserStory

title: Try out crusado
tasks:
  - title: Download crusado
  - title: Test crusado
  - title: Document test results
    description: Optional description of a task
---

## Your User Story Description

This will end up in the **description field** of the user story.

[Markdown](https://en.m.wikipedia.org/wiki/Markdown) is supported!
```

For a more elaborate example, see the [example directory](./example/).

<details>
  <summary>More information about the available Frontmatter fields (click to toggle)</summary>

  * `name`: The name of the template within the context of `crusado`. This is
    the name you call in the `crusado template` subcommands to address this
    template. So chose a short and concise one! This is _not_ the name of the
    resulting User Story/Bug (that would be `title`).
  * `summary`: A short summary of the template within the context of
    `crusado`. This summary will not end up in the resulting User Story/Bug,
    but instead is used during `crusado template list` to give you a little
    more context on what the template contains. Use this field in whatever way
    best supports your workflow.
  * `type`: One of [`UserStory`, `Bug`]. Available options might be extended
    in the future.
  * `title`: This is the title of the resulting User Story/Bug in Azure DevOps
    once the template is applied.
  * `tasks`: The tasks to create as children of the User Story/Bug. Can be
    left empty if your template doesn't need subtasks.
    * `title`: Like the higher-level `title`, the title of the resulting task in
      Azure Devops.
    * `description`: The description of the resulting task in Azure Devops. You
      can leave this empty.
</details>

All Markdown content below the Frontmatter will be interpreted by `crusado` as
the content/description of the UserStory/Bug. (No guarantee that Azure DevOps
will accept all resulting HTML, but in my tests, most standard Markdown worked.)

That's a complete setup for `crusado`! Now continue with how to put it to use.

## Usage

> All command examples assume you put `crusado` in your path and are able to
> call it globally.

### General Help

For `crusado` and all of its subcommands, you can run

```sh
crusado <optional subcommand> --help
```

to get more information about the usage.

### Working with Templates

The `crusado template` subcommand is responsible for all interactions with your
templates. The alias `t` is defined, so you can substitute all calls to
`crusado template` with `crusado t`.

**Listing all Available Templates**

```sh
crusado template list
```

This will display a list of all available templates in the current profile
(remember that you can use different profiles and switch between them by
changing the `CRUSADO_TEMPLATES_DIR` environment variable).

Similar to `kubectl` and many other CLIs, `crusado` supports multiple output
formats via the `--output`/`-o` flag. E.g., call `crusado template list -ojson`
to get the templates in JSON format.

**Showing a Specific Template**

```sh
crusado template show <template name>
```

In case you want to take a closer look at the template you want to apply, use
this command. The `<template name>` is the one specified in the Frontmatter as
`name`. The name is also displayed in the `NAME` column when you `list` the
templates.

This command also supports multiple output formats via the `--output`/`-o` flag.

**Applying a Template**

This is where the fun actually begins! To apply any of your prepared templates,
execute:

```sh
crusado template apply <template name>
```

In the default configuration, like shown here, this command will show you the
work items to be created and ask for your confirmation. If confirmed, it will
create them in the _next_ iteration (meaning, most likely, the upcoming sprint).
However, you can tweak this command to suit your needs with some flags:

* `--dry-run`/`-d`: If you want to check that the Azure DevOps API would accept
  the work items and their titles/descriptions, add this bool flag to perform a
  validating dry-run. No actual work items will be created.
* `--yes`/`-y`: Skip the confirmation step and immediately apply the work items
  with this bool flag. Useful when you use `crusado` in automation scripts.
* `--iteration-offset=<int>`/`-i=<int>`: By default, `crusado` creates the work
  items in the next iteration of your project. This default was chosen because I
  think `crusado` will most likely be used to create User Stories in preparation
  for the next sprint.

  You can use this integer flag to override the default behavior. It's designed
  as an offset relative to the _current iteration_. So, in order to apply the
  template in the current iteration, pass `-i=0`. The default therefore is
  basically `-i=1` (one iteration in the future). You can even specify negative
  values to create User Stories/Bugs in past iterations, but I don't see how
  this would be useful.

  `crusado` will always show you the complete iteration path when you run the
  `apply` command. Thus, if you haven't disabled the confirmation step, you'll
  be able to double-check the iteration is correct before anything is applied.

## Anything Missing?

If you find deficencies in this documentation, please don't hesitate to open an
issue.

## Helpful Resources During Development

* ["Create Azure DevOps Work Item Action", by Colin Dembovsky](https://colinsalmcorner.com/azdo-create-work-item-action/)
* ["Reference guide for link types used in Azure DevOps and Azure Boards", Azure Docs](https://learn.microsoft.com/en-us/azure/devops/boards/queries/link-type-reference?view=azure-devops)
