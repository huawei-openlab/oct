In this session, we define how the case should be organized.

The case should be listed like this:

```
casedir
  |__ group one
  |        |_____ case one
  |        |_____ case two
  |
  |__ group two
  |__ group three
	   |____  case three
			|___ config.json
			|___ source/container-file
			|___ source/container-script one
			|___ report.md

```

The `unit-test.sh` is used to check whether the 'tcserver' works as expect
