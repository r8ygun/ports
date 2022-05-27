# ports
This application will read a series of ports from an input file and write them to a database.
## Running
You can run the application by executing 
```bash
./run.sh
```
however, this does require that `docker compose` is installed.
## Next
Given more time I would:
- Allow the `ports.json` file to be stored externally and mounted as a volume 
- Cover more test cases, especially integration tests. I'd like to run a suite of tests against the running application using `docker compose`
- Possibly model this problem as microservices where one service read the file, wrote each object to a queue from which another service would pick the up and persist them.
