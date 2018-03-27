# COME (COnnect ME)
## Synopsis
Come is a small tool to help sysadmin connecting users machine over ssh. 

It use the REST API plugin of an openfire server to display IP Address of users, verify online status or connect a user machine with SSH. 

## Versions
* 1.2 (current): Add interactive sessions list
* 1.1: Migration from Go language to Ruby language
* 1.0: Initial version

## Installation
RPM & Deb packages are provides in the repo.
* RPM package tested with CentOS 7.4 & Fedora 27
* Deb package tested with Ubuntu 16.04 LTS & Ubuntu 17.10

WARNING: Come will only work actually if you have a correct CA for SSL connection on openfire Server (no self-signed). 

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
