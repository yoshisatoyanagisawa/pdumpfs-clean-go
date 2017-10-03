# pdumpfs-clean-go
Go implementation of pdumpfs-clean.

## how to use?
```
% git checkout https://github.com/yoshisatoyanagisawa/pdumpfs-clean-go.git
% cd pdumpfs-clean-go
% go build pdumpfs-clean.go
```
then, you get `pdumpfs-clean` binary.

## supported options
* --dryrun: dry run
* --keep: specifies the rule to keep directories.
* --remove-empty: removes empty directory (note that original pdumpfs-clean did not remove empty directory).
* --verbose: verbose.

## limitations
Unlike original implementation, this pdumpfs-clean does not have --force option.
