# Hemar

Named for the oldest container in the world, this is a simple container engine, inspired by the [Containers From Scratch](https://github.com/lizrice/containers-from-scratch) talks by Liz Rice and heavily influenced by [Vessel](https://github.com/0xc0d/vessel) by Ali Josie

## Usage
Right now, `hemar pull [image]` and `hemar run [image] [command]` work. There is full outbound networking with NAT, but it's pretty brittle. It requires that the outbound network device be named `eth0`, for example.

## Licensing

This is a learning tool based heavily on work by other people. That work is dual-licensed with Apache 2 and MIT, so this is as well.