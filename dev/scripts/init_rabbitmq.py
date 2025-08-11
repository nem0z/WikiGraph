from utils import get_mysql_conn, get_rabbitmq_connection

UNPROCESSED_URL_QUEUE   = "unprocessed_articles"
ARTICLES_QUEUE          = "articles"
RELATIONS_QUEUE         = "relations"

channel = None
rabbit_mq_conn = None
mysql_conn = None

try:
    rabbit_mq_conn = get_rabbitmq_connection()
    channel = rabbit_mq_conn.channel()
    channel.queue_declare(queue=UNPROCESSED_URL_QUEUE, durable=True)
    channel.queue_declare(queue=ARTICLES_QUEUE, durable=True)
    channel.queue_declare(queue=RELATIONS_QUEUE, durable=True)
    
    mysql_conn = get_mysql_conn()
    cursor = mysql_conn.cursor()
    cursor.execute("SELECT link FROM articles WHERE processed = 0")
    links = cursor.fetchall()
    if not links:
        raise Exception("No unprocessed links to load from MySQL")
    
    for (link, ) in links:
        channel.basic_publish(
        exchange="",
        routing_key=UNPROCESSED_URL_QUEUE,
        body=link.encode("utf-8"), # type: ignore
    )
    
finally:
    if channel is not None:
        channel.close()
    if rabbit_mq_conn is not None:
        rabbit_mq_conn.close()
    if mysql_conn is not None:
        mysql_conn.close()