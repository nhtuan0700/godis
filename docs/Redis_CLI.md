1. git clone https://github.com/redis/redis.git
2. cd redis
3. git checkout 6.2
3. make
4. ./src/redis-cli -p 3000

Or use docker
1. docker run --name test-redis -d redis:6.2
