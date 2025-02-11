## GO-DOT-SERVER
a simple server to test godot games exported to the web locally - useful for testing/learning

using this to learn how to self host a godot game (locally during development and in prod), networking in golang/web, multiplayer in godot and websockets - so this is a learning kind of project

the goal here is to write a program to self-host a godot export on the web using go - first locally and eventually in prod and learn how godot handles networking.

## Project Versions
godot - Version 4.2.1
Go - Version 1.2.3

## Setup
Place the godot export into the ./public folder

Generate a Self-Signed TLS certificate using openSSL - place the certificates into the ./certs folder. 

Open https://127.0.0.1:8080
