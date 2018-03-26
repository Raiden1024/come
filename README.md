# COME (COnnect ME)
## Synopsis
Come is a small tool to help sysadmin connecting users machine over ssh. 

It use the REST API plugin of an openfire server to display IP Address of users, verify online status or connect a user machine with SSH. 

## Versions
* 1.2 (current): Add interactive sessions list
* 1.1: Migration from Go language to Ruby language
* 1.0: Initial version

## Using "come"
This git repository contains an RPM Package (for Fedora & centOS: tested on CentOS 7.4). Install it and launch it in command line. 

## Available Options
* -l or --list - Display active users sessions list
* -m or --sessions-menu - Display interactive users sessions menu
* -i or --ip - Display User IP Address, ex: come -i <user>
* -c or --connect - SSH connect to a user machine, ex: come -c <user>
* -w or --waiting - Wait for user Online status, ex: come -w <user>
* -v or --version - Print version
* -h or --help - Print help

## Licence
3-Clauses BSD
