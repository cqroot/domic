<div align="center">
  <h1>Domic</h1>

  <p><i>Manage your <b>dotfiles</b> more easily.</i></p>

  <p>
    <a href="https://github.com/cqroot/domic/actions">
      <img src="https://github.com/cqroot/domic/workflows/test/badge.svg" alt="Action Status" />
    </a>
    <a href="https://github.com/cqroot/domic/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/cqroot/domic" />
    </a>
    <a href="https://github.com/cqroot/domic/issues">
      <img src="https://img.shields.io/github/issues/cqroot/domic" />
    </a>
  </p>
</div>

## Quick Start

If you are using domic for the first time, you can run the following command:

```bash
domic init
```

It will detect which dotfiles in your user directory can be managed with domic, and print some commands to help you create the dotfiles repository.

After you have moved the dotfiles to be managed to the repository according to the guidance, run the `domic` command to view the current status.
All dotfiles should now be in an unapplied state.

Running the `domic apply` command will generate the dotfiles in the repository to the correct location through symbolic links.
