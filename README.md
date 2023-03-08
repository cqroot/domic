<div align="center">
  <h1>Doter</h1>

  <p><i>Manage your <b>dotfiles</b> more easily.</i></p>

  <p>
    <a href="https://github.com/cqroot/doter/actions">
      <img src="https://github.com/cqroot/doter/workflows/test/badge.svg" alt="Action Status" />
    </a>
    <a href="https://github.com/cqroot/doter/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/cqroot/doter" />
    </a>
    <a href="https://github.com/cqroot/doter/issues">
      <img src="https://img.shields.io/github/issues/cqroot/doter" />
    </a>
  </p>
</div>

## Quick Start

If you are using doter for the first time, you can run the following command:

```bash
doter init
```

It will detect which dotfiles in your user directory can be managed with doter, and print some commands to help you create the dotfiles repository.

After you have moved the dotfiles to be managed to the repository according to the guidance, run the `doter` command to view the current status.
All dotfiles should now be in an unapplied state.

Running the `doter apply` command will generate the dotfiles in the repository to the correct location through symbolic links.
