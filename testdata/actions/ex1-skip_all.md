
# Inputs

## benchdiff_version

default: `0.4.1`

Version of benchdiff to use (exclude "v" from the front of the release name)

## benchdiff_dir

__Deprecated__ - Let's pretend this is deprecated

default: `${{ runner.temp }}/benchdiff`

Where benchdiff will be installed

## install_only

default: `false`

Whether to stop after installing. Any value other than "false" is interpreted as true.

## report_status

default: `true`

Whether to report the status. Any value other than "true" is interpreted as false.

## status_name

default: `benchdiff`

Name to use in reporting status.

## github_token

__Required__

Token to use for reporting status.

## benchdiff_args

default: `--base-ref=$default_base_ref`

Arguments for the benchdiff command.
All instances of $default_base_ref will be replaced with this repo's default branch.

See https://github.com/willabides/benchdiff for usage


# Outputs

## benchdiff_bin

path to the benchdiff executable

## benchstat_output

output from benchstat

## bench_command

command used to run benchmarks

## head_sha

git revision benchstat used as head

## base_sha

git revision benchstat used as base

## degraded_result

whether any part of the results is degraded
