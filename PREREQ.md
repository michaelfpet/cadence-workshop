## Prerequisites to run the Cadence server locally

You need to install prerequisites which are necessary to run Cadence.  
You need `git`, `docker` and `docker-compose`.  
If you're missing any of them, follow the next session for your OS for step-by-step guidance:

## MacOS

### Install Homebrew (if you don't have it yet):
1. Open **Terminal**
2. Run the following command to install Homebrew:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### Install Git:
Once Homebrew is installed, run the following command to install Git:

```bash
brew install git
````
After installation, verify that Git was installed correctly by checking its version:
```bash
git --version
````
You should see something like git version 2.x.x.

### Install Docker and docker-compose:
1. Go to the official Docker website: [download Docker Desktop for Mac](https://docs.docker.com/desktop/setup/install/mac-install/).
2. Install Docker Desktop:
    - Open the .dmg file you just downloaded.
    - Drag the Docker icon to the Applications folder to install it.
3. Start Docker Desktop:
    - Open the Applications folder and click on the Docker app to launch it.
    - Docker may prompt you to enter your macOS password to complete the installation.
4. Verify Docker Installation:
    - After Docker Desktop starts, you should see the Docker icon in the macOS menu bar.
    - To verify the installation, open a terminal and run:
       ```bash
       docker --version
       ```
5. Verify Docker Compose Installation: Docker Compose is bundled with Docker Desktop, so it should already be installed. To verify, run:
    ```bash
    docker-compose --version
    ````
This will display the installed version of Docker Compose.

## Linux

We assume you're using Ubuntu distribution. If you're using other, you most probably need to change `apt-get` to the package manager used in the distro.

### Install Git:

1. Open **Terminal**
2. Run the following command to install `git`:
```bash
apt-get install git
````
After installation, verify that Git was installed correctly by checking its version:
```bash
git --version
````
You should see something like git version 2.x.x.

TODO - write about docker-compose
