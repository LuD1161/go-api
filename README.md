# TNBT

## Setting up the Project

1. `git clone https://github.com/LuD1161/restructuring-tnbt/`
2. `cd restructuring-tnbt`
3. Create a `.env` file or rename `.env.sample` to `.env`
4. `docker-compose up` or if you want it to be running in the background `docker-compose up -d`
5. Go to `http://localhost:3000/swagger/index.html` and you can play with the API.

### FAQ
1. How do I interact with the API ( like create user, get JWToken, login etc ) ?
- The `http://localhost:3000/swagger/index.html` is an interactive swagger UI to play with the API.

2. Where can I access the PGAdmin console ( yes there's that too :p ) ?
- It can be accessed at `http://localhost:5050/`. It's username and password is what you've in your `.env` file.

3. How to query the database ?
- Create a new server in PGAdmin dashboard. Find the IP Address of your postgres database in the docker swarm using `docker inspect full_db_postgres | grep -i 'ipaddress'`
### Output of docker inspect full_db_postgres | grep -i 'ipaddress': 
<img width="454" alt="image" src="https://user-images.githubusercontent.com/17861054/74604183-67597100-50e1-11ea-963d-70879fe30774.png">

### PGAdmin Dashboard - Creating a new server : 
<img width="702" alt="image" src="https://user-images.githubusercontent.com/17861054/74604139-22353f00-50e1-11ea-89c3-1630f364c927.png">
