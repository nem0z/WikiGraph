import os
import time
import mysql.connector
import redis
import pika

def get_mysql_conn(retry: int = 10):
    mysql_host = os.environ.get("MYSQL_HOST")
    mysql_user = os.environ.get("MYSQL_USER")
    mysql_password = os.environ.get("MYSQL_PASSWORD")
    mysql_database = os.environ.get("MYSQL_DATABASE")
    
    if not (mysql_host and mysql_user and mysql_password and mysql_database):
        raise Exception("You need to set all the required env variable for MySQL")
    
        
    for i in range(retry):
        try:
            return mysql.connector.connect(
                host=mysql_host,
                user=mysql_user,
                password=mysql_password,
                database=mysql_database
            )

        except Exception as e:
            print(f"MySQL is unavailable - retry {i+1}/{retry}.\n Error: {e}")
            time.sleep(5)
            
    raise Exception("Could not establish a connection with MySQL")

def get_redis_conn(retry: int = 10):
    redis_host = os.environ.get("REDIS_HOST")

    if not redis_host:
        raise Exception("REDIS_HOST env variable need to be set")
    
    for i in range(retry):
        try:
            r = redis.Redis(host=redis_host)
            r.ping()
            return r
            
        except Exception as e:
            print(f"Redis is unavailable - retry {i+1}/{retry}.\n Error: {e}")
            time.sleep(5)
            
    raise Exception("Could not establish a connection with Redis")

def get_rabbitmq_connection(retry: int = 10):
    host = os.environ.get("RABBITMQ_HOST")
    user = os.environ.get("RABBITMQ_DEFAULT_USER")
    password = os.environ.get("RABBITMQ_DEFAULT_PASS")
    
    if not (host and user and password):
        raise Exception("You need to set all the required env variable for RabbitMQ")

    credentials = pika.PlainCredentials(user, password)

    for i in range(retry):
        try:
            return pika.BlockingConnection(pika.ConnectionParameters(host=host, credentials=credentials))
        
        except Exception as e:
            print(f"Redis is unavailable - retry {i+1}/{retry}.\n Error: {e}")
            time.sleep(5)

    raise Exception("Could not establish a connection with RabbitMQ")


