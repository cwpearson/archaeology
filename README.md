# archaeology

A block-level file backup utility written in Go.

## Setup

Archaeology looks for `config.json` in `/etc/archaeology/`, `$HOME/.archaeology`, or the execution directory (`.`).
The configuration file specifies what files Archaeology should include in the backup,
what files should be ignored (which overrides include), and where the backups should be stored.
An example is below.  

    {
        "ignore": [
            ".git",
            "**/*.ignoreme"

        ]
    }

## Wishlist

* Command for block-level edit distance
* Track block-level changes