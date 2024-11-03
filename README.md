# Getting started

You will need three terminals to run the code in this repository.
1. One for running the Cadence server via Docker Compose.
2. One for running the Cadence CLI.
3. One for running the Cadence worker contained in this repository.

## Running the Cadence server

In one terminal, run the following command from the repository root:
```bash
docker-compose up
```

You will see a lot of output as the Cadence server starts up. 
The cadence web UI will be available at http://localhost:8088.

## Running the Cadence CLI

First lets setup an alias for the Cadence CLI.
### Alias for Cadence CLI
An easy way to run the Cadence CLI is to use Docker. We will
create an alias for the command to make it easier to run. If you do not 
want to create an alias, you can run the command directly.

After setting up the alias you can run the Cadence CLI by running `cadence` in a terminal.

#### Mac:
Run the following command in a terminal:
```bash
alias cadence="docker run -i -t ubercadence/cli:master --address host.docker.internal:7933"
```

#### Linux:
Run the following command in a terminal:
```bash
alias cadence="docker run -i -t --network=host ubercadence/cli:master --address 127.0.0.1:7933"
```

#### Windows:
Create a file called `cadence.bat` with the following content:
```bash
docker run -i -t ubercadence/cli:master --address host.docker.internal:7933 %*
```

### Creating a domain for the worker
Before we can start the worker, we need to create a domain for it.
Run the following command in a terminal:
```bash
cadence --domain cadence-workshop domain register
```

You should now be able to see the domain in the Cadence web UI here:
http://localhost:8088/domains/cadence-workshop

## Running the Cadence worker
To run the cadence worker, simply run the following from the repository root:
```bash
go run main.go
```

## Starting a workflow
To start the hello world workflow, run the following command:
```bash
cadence --domain cadence-workshop wf start --tasklist tasklist --execution_timeout 10 --workflow_type HelloWorld --input '{"message":"Cadence"}'
```

The workflow should start and you should see the workflow in the web UI:
http://localhost:8088/domains/cadence-workshop
