## HOW TO INSTALL

Cryptopump can be used on Windows or Linux (with a MySQL or MariaDB database) or in a self-contained Docker environment. 

### DOCKER:

#### - CryptoPump is now available as a self-contained Docker container set for linux/amd64 and linux/arm/v7 (Raspberry Pi). Check it out at https://hub.docker.com/repository/docker/andreleibovici/cryptopump

This is the easiest way to get CryptoPump up and running. The Docker Compose file provides all the necessary components to run the system.

1 - Create a directory for the project.

2 - Copy the file docker-compose.yml to cryptopump directory, and edit the file, replacing <password> for "MYSQL_ROOT_PASSWORD" and "DB_PASS" with a password of your choice. Save the file.

3 - Execute docker-compose up -d

Docker Compose will download all the necessary images, including the database, and start Cryptopump. A persistent volume named "cryptopump_db_data" will be created and host the database files.

Cryptopump can be accessed by pointing your browser to http://localhost:8080.

### WINDOWS:

On Windows MYSQL can be used with Docker Desktop, use port 3306 and cryptopump as the database name. 
User should be root and you can set a password. 

Start the command prompt and do:
```
$ docker run -d --name cryptopump -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password mysql:latest
```

Now with your MYSQL server running on docker find the docker container ID
```
$ docker ps
```

and now add the cryptopump .sql file that's inside mysql folder of cryptopump to the database
```
$ docker exec -i <DOCKERID> mysql -uroot -p<ROOTPASSWORD> cryptopump < c:\path\to\cryptopump\mysql\cryptopump.sql
```

Now export the environment so the cryptopump executable is able to connect to the MYSQL server.

Using windows powershell
```
$env:DB_TCP_HOST="127.0.0.1"
$env:DB_PORT="3306"
$env:DB_USER="root"
$env:DB_PASS="<DB_PASSWORD>"
$env:DB_NAME="cryptopump"
$env:PORT="8090"
```

Now set your BINANCE API key under cryptopump configuration file cryptopump\config\config_default.yaml
Set it with apikey: and the secret under secretkey:

Withing the same shell start cryptopump and if everything worked ok you should see a browser window with cryptopump running.