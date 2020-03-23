# Changes

### v1.6.0

- refactored kiya so that it can be used as a library

### v1.5.0

- add "env" function for template command that reads OS environment values.

### v1.4.3

- fixes exit (1) on error (thanks to Frank Schroder)

### v1.4.1

- more logging when moving secrets from one to another profile

### v1.4.0

- add filter for list operation (thanks Tom Geurtsen)

### v1.3.5

- default generate character set is made URL encoding free
- after generate password copy it to clipboard
- do not log secrets if a command fails
- return with exit code 1 if kiya is aborted 
